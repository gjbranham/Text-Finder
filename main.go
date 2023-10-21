package main

import (
	// "fmt"

	"flag"
	"log"
	"os"
	"time"
)

var matchInfo *matchInformation = new(matchInformation)
var globalArgs *arguments

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

	info, err := os.Stat(globalArgs.rootPath)
	if err != nil {
		log.Fatalf("Fatal error: could not get info for path '%v'\n", globalArgs.rootPath)
	}

	start := time.Now()

	if info.IsDir() {
		findFiles()
	} else {
		checkFileForMatch(globalArgs.rootPath)
	}

	printResults(start)
}
