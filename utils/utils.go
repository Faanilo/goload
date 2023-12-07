package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

func GetFilePathFromArgs(args []string) string {
	if len(args) < 2 {
		return ""
	}
	return args[1]
}

func SetupWatcher(fileToRun string) (*fsnotify.Watcher, error) {
	absPath, err := filepath.Abs(fileToRun)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(absPath)
	if err := watcher.Add(dir); err != nil {
		return nil, err
	}

	return watcher, nil
}

func WatchChange(fileToRun string, watcher *fsnotify.Watcher, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write && strings.HasSuffix(event.Name, ".go") {
				fmt.Println("Detected change in", event.Name)
				fmt.Println("Restarting the application...")
				if err := RestartApp(fileToRun); err != nil {
					fmt.Println("Error restarting application:", err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error watching:", err)
		}
	}
}

func RestartApp(file string) error {
	os.Exit(0)
	// Uncomment the lines below if you want to spawn a new process instead of exiting the current one
	// runCmd := exec.Command("go", "run", file)
	// runCmd.Stdout = os.Stdout
	// runCmd.Stderr = os.Stderr
	// return runCmd.Start()
	return nil
}

func StartServer(file string) error {
	runCmd := exec.Command("go", "run", file)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	return runCmd.Run()
}
