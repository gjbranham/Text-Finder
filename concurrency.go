package main

import "sync"

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
