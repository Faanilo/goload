package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Faanilo/goload/utils"
)

func main() {
	fileToRun := utils.GetFilePathFromArgs(os.Args)
	if fileToRun == "" {
		fmt.Println("Please provide the path to the Go file to run.")
		return
	}

	watcher, err := utils.SetupWatcher(fileToRun)
	if err != nil {
		fmt.Println("Error setting up watcher:", err)
		return
	}
	defer watcher.Close()

	fmt.Println("Watching for changes in", filepath.Dir(fileToRun))

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go utils.WatchChange(fileToRun, watcher, wg)

	fmt.Println("Running", fileToRun)
	if err := utils.StartServer(fileToRun); err != nil {
		fmt.Println("Error running server:", err)
		wg.Wait()
	} else {
		wg.Wait()
	}
}
