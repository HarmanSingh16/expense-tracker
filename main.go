package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type Expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
}

const filename = "expenses.json"

func nextID(expenses []Expense) int {
	maxID := 0

	for _, expense := range expenses {
		if expense.ID > maxID {
			maxID = expense.ID
		}
	}

	return maxID + 1
}

func loadExpenses() []Expense {
	var expenses []Expense

	data, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return expenses
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if err := json.Unmarshal(data, &expenses); err != nil {
		fmt.Println(err)
	}

	return expenses
}

func saveExpenses(expenses []Expense) {
	data, err := json.MarshalIndent(expenses, "", " ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.WriteFile(filename, data, 0644)
}

func addExpense(description string, amount int, expenses []Expense) []Expense {
	id := nextID(expenses)
	var expense = Expense{
		ID:          id,
		Date:        time.Now(),
		Description: description,
		Amount:      amount,
	}

	expenses = append(expenses, expense)
	fmt.Printf("Expense added successfully (ID: %d)\n", id)

	return expenses
}

func deleteExpense(id int, expenses []Expense) []Expense {
	before := len(expenses)
	expenses = slices.DeleteFunc(expenses, func(expense Expense) bool { return expense.ID == id })

	if len(expenses) == before {
		fmt.Printf("No expense found (ID: %d)\n", id)
	} else {
		fmt.Println("Expense deleted successfully")
	}
	return expenses
}

func main() {
	var expenses = loadExpenses()
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage: expense-tracker add|delete|summary|list")
		return
	}
	switch args[0] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		dscpPtr := addCmd.String("description", "", "Describe your expense")
		amountPtr := addCmd.Int("amount", -1, "Enter the amount of your expense")
		addCmd.Parse(os.Args[2:])

		description := strings.TrimSpace(*dscpPtr)
		amount := *amountPtr

		if description == "" {
			fmt.Println("Invalid Description")
			fmt.Println("Use flag: -description=<description>")
			return
		} else if amount < 0 {
			fmt.Println("Invalid Amount")
			fmt.Println("Use flag: -amount=<amount>, amount >= 0")
			return
		}

		expenses = addExpense(*dscpPtr, *amountPtr, expenses)
		saveExpenses(expenses)

	case "list":
		fmt.Println("ID\tDate\tDescription\tAmount")
		for _, expense := range expenses {
			fmt.Printf("%d\t%s\t%s\t$%d\n", expense.ID, expense.Date.Format("2006-01-02"), expense.Description, expense.Amount)
		}

	case "summary":
		var total int
		if len(args) < 2 {
			for _, expense := range expenses {
				total += expense.Amount
			}
			fmt.Printf("Total expenses: $%d\n", total)
		} else {
			summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
			monthPtr := summaryCmd.Int("month", -1, "Filter expenses based on month")
			summaryCmd.Parse(os.Args[2:])

			if *monthPtr < 1 || *monthPtr > 12 {
				fmt.Println("Invalid month")
				return
			}

			month := time.Month(*monthPtr)
			for _, expense := range expenses {
				if expense.Date.Month() == month {
					total += expense.Amount
				}
			}

			fmt.Printf("Total expenses for month %s: $%d\n",month, total)
		}

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		idPtr := deleteCmd.Int("id", -1, "Enter id of the task to delete")
		deleteCmd.Parse(os.Args[2:])

		id := *idPtr
		if id == -1 {
			fmt.Println("Usage: expense-tracker delete -id=<id>")
		} else {
			expenses = deleteExpense(*idPtr, expenses)
			saveExpenses(expenses)
		}
	
	default:
		fmt.Println("Invalid Argument")
		fmt.Println("Usage: expense-tracker add|delete|summary|list")
	}
}
