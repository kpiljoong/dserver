package server

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

func WatchConfigFile(filePath string, reloadConfigFunc func() error, verbose bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v\n", err)
	}
	defer watcher.Close()

	addFileToWatcher(watcher, filePath, verbose)

	for {
		select {
		case event := <-watcher.Events:
			if verbose {
				log.Println("Config file changed. Reloading...")
			}

			// Handle write and create events
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				time.Sleep(100 * time.Millisecond) // Allow file stabilization
				if err := reloadConfigFunc(); err != nil {
					log.Printf("Failed to reload config: %v\n", err)
				} else if verbose {
					log.Println("Config reloaded successfully")
				}
			}

			if event.Op&fsnotify.Rename != 0 {
				if verbose {
					log.Println("Config file renamed. Re-adding to watcher...")
				}
				time.Sleep(100 * time.Millisecond)           // Allow file stabilization
				addFileToWatcher(watcher, filePath, verbose) // Re-attach watcher
				if err := reloadConfigFunc(); err != nil {
					log.Printf("Failed to reload config after rename: %v\n", err)
				}
			}
		case err := <-watcher.Errors:
			log.Printf("Watcher error: %v\n", err)
		}
	}
}

func addFileToWatcher(watcher *fsnotify.Watcher, path string, verbose bool) {
	if err := watcher.Add(path); err != nil {
		log.Printf("Failed to add file to watcher: %v\n", err)
	} else if verbose {
		log.Printf("Watching file: %s", path)
	}
}
