package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var rootDir string
var recurseSearch *bool
var terms []string

func processArgs() {

	flag.StringVar(&rootDir, "d", "./", "Root directory to start searching for matches")
	recurseSearch = flag.Bool("r", false, "Search recursively starting at the root directory")

	flag.Parse()

	terms = flag.Args()

	start, err := filepath.Abs(rootDir)
	if err != nil {
		log.Fatalf("Could not determine path for root dir: %v", err)
	}

	if stat, err := os.Stat(start); err != nil {
		log.Fatalf("os.Stat failed: %v", err)
	} else if !stat.IsDir() {
		log.Fatalf("'%v' is not a directory", start)
	}
}
