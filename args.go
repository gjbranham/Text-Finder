package main

import (
	"flag"
)

var rootPath string
var recursiveSearch *bool
var searchTerms []string

func processArgs() {
	flag.StringVar(&rootPath, "d", "./", "Root directory to start searching for matches")
	recursiveSearch = flag.Bool("r", false, "Search recursively starting at the root directory")

	flag.Parse()

	searchTerms = flag.Args()
}
