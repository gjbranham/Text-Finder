package main

import (
	"bufio"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

func findFiles(rootPath string) {
	var wg sync.WaitGroup

	files, err := os.ReadDir(rootPath)
	if err != nil {
		log.Fatalf("Fatal error occurred while walking root dir: %v\n", err)
	}

	for _, fo := range files {
		path := path.Join(rootPath, fo.Name())

		if fo.IsDir() && globalArgs.recursiveSearch {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				findFiles(path)
			}(path)
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			checkFileForMatch(path)
		}(path)
	}
	wg.Wait()
}

func checkFileForMatch(file string) {
	fileObj, err := os.Open(file)
	if err != nil {
		log.Printf("Failed to open file '%v': %v\n", file, err)
		return
	}
	defer fileObj.Close()

	lineNum := 1
	localMatchCnt := 0
	localMatchList := []fileInfo{}

	r := bufio.NewScanner(fileObj)
	for r.Scan() {
		line := r.Text()
		for _, key := range globalArgs.searchTerms {
			if strings.Contains(line, "\x00") {
				log.Printf("Ignoring binary file %v", file)
				return
			} else if strings.Contains(line, strings.TrimSpace(key)) {
				localMatchCnt++
				localMatchList = append(localMatchList, fileInfo{key: key, file: file, lineNum: lineNum})
			}
		}
		lineNum++
	}
	matchInfo.counterInc(localMatchCnt)
	matchInfo.addMatch(localMatchList...)
}
