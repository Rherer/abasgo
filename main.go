package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var Configuration *Config
var Abasgui *exec.Cmd
var Client string

func init() {
	Client = getClient()
	// Load configuration from file
	Configuration = InitializeConfiguration()
	// Get the executable from the read config
	Abasgui = exec.Command("./"+Configuration.Backend.ExePath, Configuration.Backend.ExeArgs)
	// Start a parallel file watcher
	go initializeWatcher()
}

func main() {
	startExe()
}

// Get client name from network share path
func getClient() string {
	// Get client name from executable path
	// Get the Path of the Executable, then get everything after win
	ex, err := os.Executable()
	if err != nil {
		log.Panicln("Could not get executable path:", err)
	}
	dir := filepath.Dir(ex)
	dir = filepath.Base(dir)
	client := dir[strings.LastIndex(dir, "win")+3:]
	log.Println("Build Watcher started for client:", client)

	return client
}

// Start executable defined in config
func startExe() {
	// When the child is finished, release its memory and shutdown
	defer func(Abasgui *exec.Cmd) {
		err := Abasgui.Wait()
		if err != nil {
			log.Fatal("Couldn't release memory for Executable.", err)
		}
		os.Exit(0)
	}(Abasgui)
	// Start executable
	err := Abasgui.Start()
	if err != nil {
		log.Panicln("PANIC! Child Executable failed to start.", err)
	}
}
