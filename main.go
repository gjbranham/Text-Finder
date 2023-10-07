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

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	processArgs()
}

func main() {
	log.Printf("Starting search at: %v\n", rootPath)
	log.Printf("Search terms: %v\n", searchTerms)
	log.Printf("Recursive: %v\n", *recursiveSearch)
	log.Print("-----Results-----\n")

	start := time.Now()

	info, err := os.Stat(rootPath)
	if err != nil {
		log.Fatalf("Fatal error: could not get info for path '%v'\n", rootPath)
	}

	if info.IsDir() {
		findFiles()
	} else {
		checkFileForMatch(rootPath)
	}

	log.Printf("-----------------\n")
	log.Printf("Found %v matching files in %v", matchCounter.count, time.Since(start))
}
