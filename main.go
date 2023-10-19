package main

import (
	// "fmt"

	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"golang.org/x/exp/slices"
)

type fileInfo struct {
	key, file string
	lineNum   int
}

type matchInformation struct {
	mu      sync.Mutex
	count   int
	matches []fileInfo
}

func (c *matchInformation) counterInc(count int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count = c.count + count
}

func (c *matchInformation) addMatch(info ...fileInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.matches = append(c.matches, info...)
}

var matchInfo *matchInformation = new(matchInformation)
var globalArgs *arguments
var padding int = 0

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
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

	// sort for nice looking output
	matchInfoCopy := copyMatchInfo()
	sortMatchInfoByKeyThenFile(matchInfoCopy)

	matchCount := printAndCountMatches(matchInfoCopy)

	log.Printf("Found %v matches in %v files in %v", matchInfoCopy.count, matchCount, time.Since(start))
}

func calcPadding() int {
	for _, t := range globalArgs.searchTerms {
		if len(t) > padding {
			padding = len(t)
		}
	}
	return padding
}

func sortMatchInfoByKeyThenFile(matchInfo *matchInformation) {
	sort.Slice(matchInfo.matches, func(i, j int) bool {
		if matchInfo.matches[i].key == matchInfo.matches[j].key {
			return matchInfo.matches[i].file < matchInfo.matches[j].file
		}
		return matchInfo.matches[i].key < matchInfo.matches[j].key
	})
}

func printAndCountMatches(matchInfo *matchInformation) int {
	padding = calcPadding()
	uniqFiles := make([]string, 0)
	customFmt := fmt.Sprintf("%%-%ds: %%s line %%v", padding)
	for _, item := range matchInfo.matches {
		log.Printf(customFmt, item.key, item.file, item.lineNum)
		if !slices.Contains(uniqFiles, item.file) {
			uniqFiles = append(uniqFiles, item.file)
		}
	}
	return len(uniqFiles)
}

func copyMatchInfo() *matchInformation {
	var matchInfoCopy matchInformation
	matchInfoCopy.count = matchInfo.count
	matchInfoCopy.matches = matchInfo.matches
	return &matchInfoCopy
}
