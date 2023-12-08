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

func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RestartApp(file string) error {
	// stop server before restarting
	if err := executeCommand("pkill", "-f", filepath.Base(file)); err != nil {
		fmt.Println("Error stopping the server:", err)
	}

	return executeCommand("go", "run", file)
}

func StartServer(file string) error {
	return executeCommand("go", "run", file)
}
