package main

import (
	"log"
	"os"
	"sync"
	"time"
)

type mutexCounter struct {
	mu    sync.Mutex
	count int
}

func (c *mutexCounter) inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

var matchCounter *mutexCounter = new(mutexCounter)


func main() {
	// Can't set flags in init() because startup order
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	processArgs()

	log.Printf("Starting search at: %v\n", args.rootPath)
	log.Printf("Search terms: %v\n", args.searchTerms)
	log.Printf("Recursive: %v\n", args.recursiveSearch)
	log.Print("-----Results-----\n")

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

	log.Printf("-----------------\n")
	log.Printf("Found %v matching files in %v", matchCounter.count, time.Since(start))
}
