package main

import (
	// "fmt"
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
}

type matchInformation struct {
	mu      sync.Mutex
	count   int
	matches []fileInfo
}

func (c *matchInformation) counterInc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *matchInformation) addMatch(info fileInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.matches = append(c.matches, info)
}

var matchInfo *matchInformation = new(matchInformation)

var padding int = 0

func main() {
	// can't set flags in init() because of startup order on tests
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	processArgs()

	for _, t := range args.searchTerms {
		if len(t) > padding {
			padding = len(t)
		}
	}

	start := time.Now()

	info, err := os.Stat(args.rootPath)
	if err != nil {
		log.Fatalf("Fatal error: could not get info for path '%v'\n", args.rootPath)
	}

	if info.IsDir() {
		findFiles()
	} else {
		checkFileForMatch(args.rootPath)
	}

	// sort for nice looking output
	sort.Slice(matchInfo.matches, func(i, j int) bool {
		return matchInfo.matches[i].key < matchInfo.matches[j].key && matchInfo.matches[i].file < matchInfo.matches[j].file
	})

	uniqFiles := make([]string, 0)

	customFmt := fmt.Sprintf("%%-%ds: %%s\n", padding)
	for _, item := range matchInfo.matches {
		log.Printf(customFmt, item.key, item.file)
		if !slices.Contains(uniqFiles, item.file) {
			uniqFiles = append(uniqFiles, item.file)
		}
	}

	log.Printf("Found %v matches in %v files in %v", matchInfo.count, len(uniqFiles), time.Since(start))
}
