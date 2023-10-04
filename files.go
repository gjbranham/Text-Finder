package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

func getAllFiles(root string) []string {
	fileList := []string{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || strings.Contains(path, "go-text-finder") || strings.Contains(path, "Documents/Bitfloor") || slices.Contains(getIgnoreExts(), filepath.Ext(path)) {
			return nil
		}
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		log.Fatalf("Failure occurred while walking directories: %v", err)
	}
	return fileList
}

func checkFilesForMatch(fileList []string, result chan string, numMatches chan int) {
	chunk := 8 * 1024 * 1024
	buf := make([]byte, chunk)

	mflag := 0
	matches := 0

	for _, file := range fileList {
		fileObj, err := os.Open(file)
		if err != nil {
			log.Printf("Failed to open file '%v", file)
		}
		defer fileObj.Close()

		for {
			bytesRead, err := fileObj.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("Failed to read chunk from file: %v", err.Error())
				}
				break
			}
			for _, keyword := range getSearchStrings() {
				if strings.Contains(string(buf[:bytesRead]), keyword) {
					fmt.Printf("%v\n", file)
					matches++
					result <- file
					mflag = 1
					break
				}
			}
			if mflag == 1 {
				mflag = 0
				break
			}
		}
	}
	numMatches <- matches
}
