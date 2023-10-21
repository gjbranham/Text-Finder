package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"golang.org/x/exp/slices"
)

func printResults(start time.Time) {
	// copy the global slice then sort it for nice looking output
	matchInfoCopy := copyMatchInfo()
	sortMatchInfoByKeyThenFile(matchInfoCopy)

	matchCount := printAndCountMatches(matchInfoCopy)

	log.Printf("Found %v matches in %v files in %v", matchInfoCopy.count, matchCount, time.Since(start))
}

func sortMatchInfoByKeyThenFile(matchInfo *matchInformation) {
	sort.Slice(matchInfo.matches, func(i, j int) bool {
		if matchInfo.matches[i].key == matchInfo.matches[j].key {
			return matchInfo.matches[i].file < matchInfo.matches[j].file
		}
		return matchInfo.matches[i].key < matchInfo.matches[j].key
	})
}

func printAndCountMatches(matchInfo *matchInformation) int {
	padding := calcPadding()
	uniqFiles := make([]string, 0)
	customFmt := fmt.Sprintf("%%-%ds: %%s line %%v", padding)
	for _, item := range matchInfo.matches {
		log.Printf(customFmt, item.key, item.file, item.lineNum)
		if !slices.Contains(uniqFiles, item.file) {
			uniqFiles = append(uniqFiles, item.file)
		}
	}
	return len(uniqFiles)
}

func calcPadding() int {
	padding := -1
	for _, t := range globalArgs.searchTerms {
		if len(t) > padding {
			padding = len(t)
		}
	}
	return padding
}

func copyMatchInfo() *matchInformation {
	var matchInfoCopy matchInformation
	matchInfoCopy.count = matchInfo.count
	matchInfoCopy.matches = matchInfo.matches
	return &matchInfoCopy
}
