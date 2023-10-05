package main

import (
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

func main() {
	startDir, _ := processArgs()

	start := time.Now()

	getAllFiles(startDir)

	log.Printf("Found %v matching files in %v", matchCounter.count, time.Since(start))
}
