package main

import (
	"io"
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

	ex, e := os.Executable()
	if e != nil {
		log.Fatal("Fatal error: could not determine path name for running executable")
	}
	exPath := filepath.Dir(ex)

	if !*recursiveSearch {
		files, err := os.ReadDir(rootPath)
		if err != nil {
			log.Fatalf("Fatal error occurred while walking root dir: %v\n", err)
		}
		absPath, err := filepath.Abs(rootPath)
		if err != nil {
			log.Fatalf("Fatal error occurred while obtaining absolute path for starting point: %v\n", err)
		}
		for _, fo := range files {
			if fo.IsDir() || strings.Contains(fo.Name(), exPath) || slices.Contains(getIgnoreExts(), strings.ToLower(filepath.Ext(fo.Name()))) {
				continue
			}
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				checkFileForMatch(path)
			}(filepath.Join(absPath, fo.Name()))
		}
	} else {
		err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
			if e != nil {
				log.Fatalf("Fatal error: could not retrieve file info for file '%v'", path)
			}
			if info.IsDir() || strings.Contains(path, exPath) || slices.Contains(getIgnoreExts(), strings.ToLower(filepath.Ext(path))) {
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
	}
	wg.Wait()
}

func checkFileForMatch(file string) {
	chunk := 1024 * 1024
	buf := make([]byte, chunk)

	mflag := 0

	fileObj, err := os.Open(file)
	if err != nil {
		log.Printf("Failed to open file '%v': %v", file, err)
		return
	}
	defer fileObj.Close()

	for {
		bytesRead, err := fileObj.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read chunk from file '%v': %v", file, err.Error())
			}
			break
		}
		for _, keyword := range searchTerms {
			if strings.Contains(string(buf[:bytesRead]), keyword) {
				log.Printf("%v\n", file)
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
