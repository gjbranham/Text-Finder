package concurrency

import "sync"

type FileInfo struct {
	Key, File string
	LineNum   int
}

type MatchInformation struct {
	mu      sync.Mutex
	Count   int
	Matches []FileInfo
}

func (c *MatchInformation) CounterInc(count int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Count = c.Count + count
}

func (c *MatchInformation) AddMatch(info ...FileInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Matches = append(c.Matches, info...)
}
