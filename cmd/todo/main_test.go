package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task := "test task number 1"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTaskFromArguments", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s\n  2: %s\n", task, task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead", expected, string(out))
		}
	})

	task3 := "Add a third task to be marked as completed"
	t.Run("CompleteTask", func(t *testing.T) {
		// Add the task
		cmd := exec.Command(cmdPath, "-add", task3)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		// Check all tasks are present
		expected := fmt.Sprintf("  1: %s\n  2: %s\n  3: %s\n", task, task2, task3)
		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead", expected, string(out))
		}

		// Mark the new task as completed
		cmd = exec.Command(cmdPath, "-complete", "3")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		// Test that the new task is present in the completed list
		cmd = exec.Command(cmdPath, "-list-complete")
		out, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected = fmt.Sprintf("X 1: %s\n", task3)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead", expected, string(out))
		}

		// Test that the new task is NOT present in the incomplete list
		cmd = exec.Command(cmdPath, "-list-incomplete")
		out, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected = fmt.Sprintf("  1: %s\n  2: %s\n", task, task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead", expected, string(out))
		}
	})
}
