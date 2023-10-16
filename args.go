package main

import (
	"flag"
)

type arguments struct {
	rootPath        string
	recursiveSearch bool
	searchTerms     []string
}

var args arguments

func processArgs() {
	flag.StringVar(&args.rootPath, "d", "./", "Root directory to start searching for matches")
	flag.BoolVar(&args.recursiveSearch, "r", false, "Search recursively starting at the root directory")

	flag.Parse()

	args.searchTerms = flag.Args()
}
