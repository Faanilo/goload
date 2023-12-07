package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// Get the arguments passed to your program
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Please provide the path to the Go file to run.")
		return
	}
	fileToRun := args[0]

	// get path of the file to run
	absPath, err := filepath.Abs(fileToRun)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	// waiting for the change in code
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	// Add the directory containing the file to the watcher
	dir := filepath.Dir(absPath)
	err = watcher.Add(dir)
	if err != nil {
		fmt.Println("Error adding directory to watcher:", err)
		return
	}

	fmt.Println("Watching for changes in", dir)

	// Set up a goroutine to handle file change events
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Detected change in", event.Name)
					if strings.HasSuffix(event.Name, ".go") {
						fmt.Println("Restarting the application...")
						if err := restartApp(fileToRun); err != nil {
							fmt.Println("Error restarting application:", err)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("Error watching:", err)
			}
		}
	}()

	// Run the initial instance of the application
	fmt.Println("Running", fileToRun)
	if err := runServer(fileToRun); err != nil {
		fmt.Println("Error running server:", err)
		wg.Wait() // Wait for file changes despite the initial error
	} else {
		wg.Wait()
	}
}

func restartApp(file string) error {
	// Kill the current process
	os.Exit(0)
	// Uncomment the lines below if you want to spawn a new process instead of exiting the current one
	// runCmd := exec.Command("go", "run", file)
	// runCmd.Stdout = os.Stdout
	// runCmd.Stderr = os.Stderr
	// return runCmd.Start()
	return nil
}

func runServer(file string) error {
	runCmd := exec.Command("go", "run", file)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	return runCmd.Run()
}
