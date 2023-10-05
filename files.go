package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/exp/slices"
)

func getAllFiles(root string) {
	var wg sync.WaitGroup

	ex, e := os.Executable()
	if e != nil {
		log.Fatal("Fatal error: could not determine path name for running executable")
	}
	exPath := filepath.Dir(ex)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		info, e := os.Lstat(path)
		if e != nil {
			log.Fatalf("Fatal error: could not retrieve file info for file '%v'", path)
		}

		if d.IsDir() || strings.Contains(path, exPath) || (info.Mode()&os.ModeSymlink) == os.ModeSymlink || slices.Contains(getIgnoreExts(), filepath.Ext(path)) {
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
		log.Fatalf("Fatal error occurred while walking directories: %v", err)
	}
	wg.Wait()
}

func checkFileForMatch(file string) {
	chunk := 1024 * 1024
	buf := make([]byte, chunk)

	mflag := 0

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
				matchCounter.inc()
				mflag = 1
				break
			}
		}
		if mflag == 1 {
			break
		}
	}
}
