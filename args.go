package main

import (
	"bytes"
	"flag"
	"log"
)

type arguments struct {
	rootPath        string
	recursiveSearch bool
	searchTerms     []string
}

func processArgs(exeName string, sysArgs []string) (parsedArgs *arguments, output string, err error) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flags := flag.NewFlagSet(exeName, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var args arguments
	flags.StringVar(&args.rootPath, "d", "./", "Root directory to start searching for matches")
	flags.BoolVar(&args.recursiveSearch, "r", false, "Search recursively starting at the root directory")

	if err = flags.Parse(sysArgs); err != nil {
		return nil, buf.String(), err
	}

	args.searchTerms = flags.Args()

	return &args, output, nil
}
