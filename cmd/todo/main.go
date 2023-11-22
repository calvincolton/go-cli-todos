package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	todo "github.com/calvincolton/go-cli-todos"
)

var todoFileName = ".todo.json"

func main() {
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed by C$.\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2023\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}
	add := flag.Bool("add", false, "Add a task to the To-Do list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Mark a task as completed by passing its number")
	delete := flag.Int("delete", 0, "Delete a task by passing its number")
	listIncomplete := flag.Bool("list-incomplete", false, "List tasks that have not been marked as completed")

	flag.Parse()

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *add:
		// Get the task from standard input
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Add the task
		l.Add(t)

		// Save the To-Do list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *list:
		// List current To-Do items
		fmt.Print(l)
	case *listIncomplete:
		// List current To-Do items that have not been completed
		filteredList := todo.List{}
		for _, t := range *l {
			if !t.Done {
				filteredList = append(filteredList, t)
			}
		}
		fmt.Print(&filteredList)
	case *complete > 0:
		// Complete the To-Do item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the To-Do list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > 0:
		// Delete the To-do item
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the To-Do list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}

	return s.Text(), nil
}
