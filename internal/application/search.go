package application

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	c "github.com/gjbranham/Text-Finder/internal/concurrency"
	o "github.com/gjbranham/Text-Finder/internal/output"
)

func (a *TextFinder) FindFiles(rootPath string) {
	var wg sync.WaitGroup

	files, err := os.ReadDir(rootPath)
	if err != nil {
		o.Print(fmt.Sprintf("Fatal error occurred while walking root dir: %v\n", err))
	}

	for _, fo := range files {
		path := path.Join(rootPath, fo.Name())

		if fo.IsDir() && a.Args.RecursiveSearch {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				a.FindFiles(path)
			}(path)
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			a.CheckFileForMatch(path)
		}(path)
	}
	wg.Wait()
}

func (a *TextFinder) CheckFileForMatch(file string) {
	fileObj, err := os.Open(file)
	if err != nil {
		o.Print(fmt.Sprintf("Failed to open file '%v': %v\n", file, err))
		return
	}
	defer fileObj.Close()

	lineNum := 1
	localMatchCnt := 0
	localMatchList := []c.FileInfo{}

	r := bufio.NewScanner(fileObj)
	for r.Scan() {
		line := r.Text()
		for _, key := range a.Args.SearchTerms {
			if strings.Contains(line, "\x00") {
				o.Print(fmt.Sprintf("Ignoring binary file %v", file))
				return
			} else if strings.Contains(line, strings.TrimSpace(key)) {
				localMatchCnt++
				localMatchList = append(localMatchList, c.FileInfo{Key: key, File: file, LineNum: lineNum})
			}
		}
		lineNum++
	}
	a.MatchInfo.CounterInc(localMatchCnt)
	a.MatchInfo.AddMatch(localMatchList...)
}
