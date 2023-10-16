package main

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/exp/slices"
)

func findFiles() {
	var wg sync.WaitGroup

	if !args.recursiveSearch {
		files, err := os.ReadDir(args.rootPath)
		if err != nil {
			log.Fatalf("Fatal error occurred while walking root dir: %v\n", err)
		}
		absPath, err := filepath.Abs(args.rootPath)
		if err != nil {
			log.Fatalf("Fatal error occurred while obtaining absolute path for starting point: %v\n", err)
		}
		for _, fo := range files {
			if fo.IsDir() || slices.Contains(getIgnoreExts(), strings.ToLower(filepath.Ext(fo.Name()))) {
				continue
			}
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				checkFileForMatch(path)
			}(filepath.Join(absPath, fo.Name()))
		}
	} else {
		err := filepath.Walk(args.rootPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				log.Fatalf("Fatal error: could not retrieve file info for file '%v'\n", path)
			}
			if info.IsDir() || slices.Contains(getIgnoreExts(), strings.ToLower(filepath.Ext(path))) {
				return nil
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				checkFileForMatch(path)
			}()
			return nil
		})
		if err != nil {
			log.Fatalf("Fatal error occurred while walking directories: %v\n", err)
		}
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

	localMatchCnt := 0
	localMatchList := []fileInfo{}

	r := bufio.NewScanner(fileObj)
	for r.Scan() {
		line := r.Text()
		for _, key := range args.searchTerms {
			if strings.Contains(line, "\x00") {
				log.Printf("Ignoring binary file %v", file)
				return
			} else if strings.Contains(line, strings.TrimSpace(key)) {
				localMatchCnt++
				localMatchList = append(localMatchList, fileInfo{key: key, file: file})
			}
		}
	}
	matchInfo.counterInc(localMatchCnt)
	matchInfo.addMatch(localMatchList...)
}
