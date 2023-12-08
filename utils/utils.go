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

var pauseSignal = make(chan struct{})
var resumeSignal = make(chan struct{})

func GetTargetFilePath(args []string) string {
	if len(args) < 2 {
		return ""
	}
	return args[1]
}

func InitializeWatcher(targetFile string) (*fsnotify.Watcher, error) {
	absPath, err := filepath.Abs(targetFile)
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

func MonitorChanges(targetFile string, watcher *fsnotify.Watcher, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-pauseSignal:
			// Pausing watcher
			<-resumeSignal // Wait for resume signal
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write && strings.HasSuffix(event.Name, ".go") {
				fmt.Println("Detected change in", event.Name)
				fmt.Println("Restarting the application...")
				if err := RestartApplication(targetFile); err != nil {
					fmt.Println("Error restarting application:", err)
				} else {
					fmt.Println("Application restarted successfully")
					close(pauseSignal) // Stop the watcher
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error watching:", err)
			fmt.Println("Stopping the server...")
			if err := StopServerProcess(targetFile); err != nil {
				fmt.Println("Error stopping application:", err)
			}
			close(pauseSignal) // Stop the watcher
			return
		}
	}
}

func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RestartApplication(file string) error {
	// stop server before restarting
	if err := executeCommand("pkill", "-f", filepath.Base(file)); err != nil {
		fmt.Println("Error stopping the server:", err)
	}

	return executeCommand("go", "run", file)
}

func StartServerProcess(file string) error {
	return executeCommand("go", "run", file)
}
func StopServerProcess(file string) error {
	return executeCommand("pkill", "-f", filepath.Base(file))
}

func ResumeWatcher() {
	close(resumeSignal)
}
