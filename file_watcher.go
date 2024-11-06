package main

import (
	"bufio"
	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Regex for finding step number in build status file
var regexStarted = regexp.MustCompile("^STEP\\s([0-9]+):\\sstarted")
var regexFinished = regexp.MustCompile("^STEP\\s([0-9]+):\\sfinished")

// Initialize a new watcher to watch for changes in the build status file
func initializeWatcher() {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// Close watcher after everything is done
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Fatal("Error closing watcher", err)
		}
	}(watcher)

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Only trigger if the file is the one specified in the config
				if filepath.Base(event.Name) != Configuration.Backend.StatusFile {
					continue
				}

				// New line written? Then check if we should display a notification
				if event.Has(fsnotify.Write) {
					// Read from status file
					var line string
					f, _ := os.Open(Configuration.Backend.StatusFile)
					scanner := bufio.NewScanner(f)
					for scanner.Scan() {
						line = scanner.Text()
					}
					err := f.Close()
					if err != nil {
						log.Fatal("Error closing file: ", err)
					}

					// Notify with last read line and client gotten from path
					// But only notify, if the line is one of the Step(s) that have been configured!
					stepStarted := regexStarted.FindStringSubmatch(line)
					stepFinished := regexFinished.FindStringSubmatch(line)
					stepStartedID := 0
					stepFinishedID := 0
					if stepStarted != nil && len(stepStarted) > 0 {
						stepStartedID, _ = strconv.Atoi(stepStarted[1])
					}
					if stepFinished != nil && len(stepFinished) > 0 {
						stepFinishedID, _ = strconv.Atoi(stepFinished[1])
					}

					if slices.Contains(Configuration.Userspace.NotifyOnStepsStart, stepStartedID) || slices.Contains(Configuration.Userspace.NotifyOnStepsFinish, stepFinishedID) {
						pushNotif(Client, line)
					}

					// Oh oh, the build failed!
					if strings.Contains(line, "ERROR") {
						pushNotif(Client+": Build failed", line)
					}
				}
				// New file = new build
				if event.Has(fsnotify.Create) {
					pushNotif(Client, "Started new Build for Client: "+Client)
				}
				// File has gone away = build finished
				if event.Has(fsnotify.Remove) {
					pushNotif(Client, "Finished Build for Client: "+Client)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Fatal("error:", err)
			}
		}
	}()

	// Watch the current folder
	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
	}

	// Keep watching
	<-make(chan struct{})
}

// Push a notification to the os
func pushNotif(title string, message string) {
	err := beeep.Notify(title, message, Configuration.Backend.LogoFile)
	if err != nil {
		log.Fatal("Couldn't start GUI-Notification.", err)
	}
}
