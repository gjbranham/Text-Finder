package main

import (
	"fmt"
	"log"
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
	processArgs()
}

func main() {
	fmt.Printf("Starting search at root dir: %v\n", rootDir)
	fmt.Printf("Search terms: %v\n", terms)
	fmt.Printf("Recursive: %v\n", *recurseSearch)

	start := time.Now()

	findFiles()

	log.Printf("Found %v matching files in %v", matchCounter.count, time.Since(start))
}
