package main

import (
	// "fmt"

	"flag"
	"log"
	"os"
	"path/filepath"
	"time"
)

var matchInfo *matchInformation = new(matchInformation)
var globalArgs *arguments = new(arguments)

func main() {
	args, out, err := processArgs(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		log.Print(out)
		os.Exit(2)
	} else if err != nil {
		log.Printf("Failed to parse command-line arguments: %v", err)
		log.Printf("Info: %v\n", out)
		os.Exit(1)
	}
	globalArgs = args

	absPath, err := filepath.Abs(globalArgs.rootPath)
	if err != nil {
		log.Fatalf("Fatal error: could not resolve absolute path for '%v'\n", globalArgs.rootPath)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		log.Fatalf("Fatal error: could not get info for path '%v'\n", absPath)
	}

	start := time.Now()

	if info.IsDir() {
		findFiles(absPath)
	} else {
		checkFileForMatch(absPath)
	}

	printResults(start)
}
