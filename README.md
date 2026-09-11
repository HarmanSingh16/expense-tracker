# Expense Tracker

A simple command-line expense tracker written in Go. Persists expenses to a local `expenses.json` file.

## Build

```bash
go build -o expense-tracker
```

## Usage

```
expense-tracker <command> [flags]
```

Commands: `add`, `list`, `summary`, `delete`, `update`

## Commands

### Add an expense

```bash
expense-tracker add -description="Coffee" -amount=5
```

Flags:
- `-description` (required) — Description of the expense
- `-amount` (required) — Amount in whole dollars (>= 0)

### List all expenses

```bash
expense-tracker list
```

Outputs a table with columns: `ID`, `Date`, `Description`, `Amount`.

### Summary

Total of all expenses:

```bash
expense-tracker summary
```

Total for a specific month (1–12):

```bash
expense-tracker summary -month=9
```

### Delete an expense

```bash
expense-tracker delete -id=3
```

### Update an expense

Only the fields you provide will be changed.

```bash
expense-tracker update -id=3 -description="Tea" -amount=3
```

Flags:
- `-id` (required) — ID of the expense to update
- `-description` — New description
- `-amount` — New amount (>= 0)

## Data

Expenses are stored as JSON in `expenses.json` in the current directory. The tool reads from / writes to this file on every command.

Each expense has the shape:

```json
{
  "id": 1,
  "date": "2026-09-11T14:26:44.757585+05:30",
  "description": "Bottle",
  "amount": 20
}
```

## Notes

- Amounts are integers (whole units of currency). The display prefixes them with `$`.
- IDs are assigned automatically and monotonically.
- If `expenses.json` does not exist yet, it will be created on the first `add` or `update` / `delete`.