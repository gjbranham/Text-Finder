package main

import (
	"log"
	"time"
)

func main() {
	startDir, _ := processArgs()

	start := time.Now()

	fileList := getAllFiles(startDir)

	log.Printf("Found %v files in %v", len(fileList), time.Since(start))

	divisor := 64
	chunk := len(fileList) / divisor
	remainder := len(fileList) - (divisor * chunk)

	fileMatchCh := make(chan string, len(fileList))
	matchCountCh := make(chan int)

	start = time.Now()

	for i := 0; i < divisor; i++ {
		go checkFilesForMatch(fileList[i*chunk:(i+1)*chunk], fileMatchCh, matchCountCh)
	}
	// Need to check the remainder N files as well
	go checkFilesForMatch(fileList[len(fileList)-remainder:], fileMatchCh, matchCountCh)

	matches := 0
	for i := 0; i < divisor+1; i++ {
		m := <-matchCountCh
		matches += m
	}

	matchedFiles := make([]string, 0)
	for i := 0; i < matches; i++ {
		f := <-fileMatchCh
		matchedFiles = append(matchedFiles, f)
	}
	log.Printf("Found %v matching files in %v", len(matchedFiles), time.Since(start))
}
