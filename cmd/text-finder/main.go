package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	app "github.com/gjbranham/Text-Finder/internal/application"
	"github.com/gjbranham/Text-Finder/internal/args"
	c "github.com/gjbranham/Text-Finder/internal/concurrency"
	o "github.com/gjbranham/Text-Finder/internal/output"
)

func main() {
	o.SetPrinter(&o.Stdout{})
	args, out, err := args.ProcessArgs(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		o.P.Printer.Print(out)
		os.Exit(2)
	} else if err != nil {
		o.Print(fmt.Sprintf("Failed to parse command-line arguments: %v", err))
		o.Print(fmt.Sprintf("Info: %v\n", out))
		os.Exit(1)
	}
	app := app.TextFinder{Args: args, MatchInfo: new(c.MatchInformation)}

	absPath, err := filepath.Abs(app.Args.RootPath)
	if err != nil {
		o.Print(fmt.Sprintf("Fatal error: could not resolve absolute path for '%v'\n", app.Args.RootPath))
	}

	info, err := os.Stat(absPath)
	if err != nil {
		o.Print(fmt.Sprintf("Fatal error: could not get info for path '%v'\n", absPath))
	}

	start := time.Now()

	if info.IsDir() {
		app.FindFiles(absPath)
	} else {
		app.CheckFileForMatch(absPath)
	}

	o.PrintResults(start, app.Args.SearchTerms, app.MatchInfo)
}
