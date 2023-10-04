package main

import (
	"log"
	"time"
)

func main() {
	startDir, _ := processArgs()
	// copyDir := os.Mkdir(path.Join(saveDir, "fileCopies"), 0644)

	start := time.Now()

	fileList := getAllFiles(startDir)

	log.Printf("Found %v files in %v", len(fileList), time.Since(start))

	divisor := 64
	chunk := len(fileList) / divisor
	remainder := len(fileList) - (divisor * chunk)

	fileMatchCh := make(chan string)
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

	file_list := make([]string, 0)
	for i := 0; i < matches; i++ {
		f := <-fileMatchCh
		file_list = append(file_list, f)
	}

	log.Printf("Found %v matching files in %v", matches, time.Since(start))
}
