package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
}

const dataFile = "expenses.json"

func loadExpenses() ([]Expense, error) {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return []Expense{}, nil
	}
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	var expenses []Expense
	err = json.Unmarshal(data, &expenses)
	return expenses, err
}

func saveExpenses(expenses []Expense) error {
	data, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func main() {
	// Define subcommands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDesc := addCmd.String("description", "", "Description of the expense")
	addAmt := addCmd.Float64("amount", 0, "Amount of the expense")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID := deleteCmd.Int("id", 0, "ID of the expense to delete")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	summaryMonth := summaryCmd.Int("month", 0, "Month number (1-12) for summary")

	// Check if a subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: expense-tracker <command> [options]")
		fmt.Println("Commands: add, delete, list, summary")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addDesc == "" || *addAmt <= 0 {
			fmt.Println("Error: description and positive amount are required")
			os.Exit(1)
		}
		expenses, err := loadExpenses()
		if err != nil {
			fmt.Printf("Error loading expenses: %v\n", err)
			os.Exit(1)
		}
		id := len(expenses) + 1
		expense := Expense{
			ID:          id,
			Date:        time.Now(),
			Description: *addDesc,
			Amount:      *addAmt,
		}
		expenses = append(expenses, expense)
		if err := saveExpenses(expenses); err != nil {
			fmt.Printf("Error saving expenses: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Expense added successfully (ID: %d)\n", id)

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		expenses, err := loadExpenses()
		if err != nil {
			fmt.Printf("Error loading expenses: %v\n", err)
			os.Exit(1)
		}
		found := false
		for i, e := range expenses {
			if e.ID == *deleteID {
				expenses = append(expenses[:i], expenses[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Error: expense ID not found")
			os.Exit(1)
		}
		if err := saveExpenses(expenses); err != nil {
			fmt.Printf("Error saving expenses: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Expense deleted successfully")

	case "list":
		listCmd.Parse(os.Args[2:])
		expenses, err := loadExpenses()
		if err != nil {
			fmt.Printf("Error loading expenses: %v\n", err)
			os.Exit(1)
		}
		if len(expenses) == 0 {
			fmt.Println("No expenses found")
			return
		}
		fmt.Println("ID  Date       Description  Amount")
		for _, e := range expenses {
			fmt.Printf("%d   %s  %-12s $%.2f\n", e.ID, e.Date.Format("2006-01-02"), e.Description, e.Amount)
		}

	case "summary":
		summaryCmd.Parse(os.Args[2:])
		expenses, err := loadExpenses()
		if err != nil {
			fmt.Printf("Error loading expenses: %v\n", err)
			os.Exit(1)
		}
		total := 0.0
		if *summaryMonth > 0 {
			if *summaryMonth < 1 || *summaryMonth > 12 {
				fmt.Println("Error: month must be between 1 and 12")
				os.Exit(1)
			}
			for _, e := range expenses {
				if e.Date.Year() == time.Now().Year() && int(e.Date.Month()) == *summaryMonth {
					total += e.Amount
				}
			}
			fmt.Printf("Total expenses for %s: $%.2f\n", time.Month(*summaryMonth), total)
		} else {
			for _, e := range expenses {
				total += e.Amount
			}
			fmt.Printf("Total expenses: $%.2f\n", total)
		}

	default:
		fmt.Println("Unknown command. Use: add, delete, list, summary")
		os.Exit(1)
	}
}
