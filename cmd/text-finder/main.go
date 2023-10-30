package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	app "github.com/gjbranham/Text-Finder/internal/application"
	"github.com/gjbranham/Text-Finder/internal/args"
	c "github.com/gjbranham/Text-Finder/internal/concurrency"
	o "github.com/gjbranham/Text-Finder/internal/output"
)

func main() {
	args, out, err := args.ProcessArgs(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		log.Print(out)
		os.Exit(2)
	} else if err != nil {
		log.Printf("Failed to parse command-line arguments: %v", err)
		log.Printf("Info: %v\n", out)
		os.Exit(1)
	}
	app := app.TextFinder{Args: args, MatchInfo: new(c.MatchInformation)}

	absPath, err := filepath.Abs(app.Args.RootPath)
	if err != nil {
		log.Fatalf("Fatal error: could not resolve absolute path for '%v'\n", app.Args.RootPath)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		log.Fatalf("Fatal error: could not get info for path '%v'\n", absPath)
	}

	start := time.Now()

	if info.IsDir() {
		app.FindFiles(absPath)
	} else {
		app.CheckFileForMatch(absPath)
	}

	o.PrintResults(start, app.Args.SearchTerms, app.MatchInfo)
}
