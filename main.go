package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

const jsonFile = "tasks.json"

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func loadTasks() ([]Task, error) {
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		return []Task{}, nil
	}
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return []Task{}, nil // Return empty slice if JSON is invalid
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(jsonFile, data, 0644)
}

func getNextID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

func addTask(description string) {
	tasks, _ := loadTasks()
	now := time.Now().Format(time.RFC3339)
	newTask := Task{
		ID:          getNextID(tasks),
		Description: description,
		Status:      "todo",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tasks = append(tasks, newTask)
	saveTasks(tasks)
	fmt.Printf("Task added with ID: %d\n", newTask.ID)
}

func updateTask(id int, description string) {
	tasks, _ := loadTasks()
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			saveTasks(tasks)
			fmt.Printf("Task %d updated\n", id)
			return
		}
	}
	fmt.Printf("Task %d not found\n", id)
}

func deleteTask(id int) {
	tasks, _ := loadTasks()
	initialLen := len(tasks)
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
	if len(tasks) < initialLen {
		saveTasks(tasks)
		fmt.Printf("Task %d deleted\n", id)
	} else {
		fmt.Printf("Task %d not found\n", id)
	}
}

func updateStatus(id int, status string) {
	tasks, _ := loadTasks()
	validStatuses := []string{"todo", "in-progress", "done"}
	valid := false
	for _, s := range validStatuses {
		if s == status {
			valid = true
			break
		}
	}
	if !valid {
		fmt.Printf("Invalid status. Use: todo, in-progress, done\n")
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			saveTasks(tasks)
			fmt.Printf("Task %d status updated to %s\n", id, status)
			return
		}
	}
	fmt.Printf("Task %d not found\n", id)
}

func listTasks(statusFilter string) {
	tasks, _ := loadTasks()
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}
	filtered := tasks
	if statusFilter != "" {
		filtered = []Task{}
		for _, task := range tasks {
			if task.Status == statusFilter {
				filtered = append(filtered, task)
			}
		}
	}
	for _, task := range filtered {
		fmt.Printf("ID: %d, Description: %s, Status: %s, Created: %s, Updated: %s\n",
			task.ID, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./task_manager <command> [args]")
		fmt.Println("Commands: add, update, delete, in-progress, done, list, list-done, list-todo, list-in-progress")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ./task_manager add <description>")
			return
		}
		description := ""
		for i := 2; i < len(os.Args); i++ {
			description += os.Args[i] + " "
		}
		addTask(description[:len(description)-1])

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: ./task_manager update <id> <description>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Task ID must be a number")
			return
		}
		description := ""
		for i := 3; i < len(os.Args); i++ {
			description += os.Args[i] + " "
		}
		updateTask(id, description[:len(description)-1])

	case "delete", "in-progress", "done":
		if len(os.Args) != 3 {
			fmt.Printf("Usage: ./task_manager %s <id>\n", command)
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Task ID must be a number")
			return
		}
		if command == "delete" {
			deleteTask(id)
		} else if command == "in-progress" {
			updateStatus(id, "in-progress")
		} else {
			updateStatus(id, "done")
		}

	case "list":
		listTasks("")

	case "list-done":
		listTasks("done")

	case "list-todo":
		listTasks("todo")

	case "list-in-progress":
		listTasks("in-progress")

	default:
		fmt.Println("Invalid command")
	}
}
