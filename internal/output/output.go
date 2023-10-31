package output

import (
	"fmt"
	"sort"
	"time"

	"golang.org/x/exp/slices"

	c "github.com/gjbranham/Text-Finder/internal/concurrency"
)

func PrintResults(start time.Time, searchTerms []string, matchInfo *c.MatchInformation) {
	// copy matchInfo then sort it for nice looking output
	matchInfoCopy := copyMatchInfo(matchInfo)
	sortMatchInfoByKeyThenFile(matchInfoCopy)

	matchCount := printAndCountMatches(matchInfoCopy, searchTerms)

	Print(fmt.Sprintf("Found %v matches in %v files in %v", matchInfoCopy.Count, matchCount, time.Since(start)))
}

func copyMatchInfo(matchInfo *c.MatchInformation) *c.MatchInformation {
	var matchInfoCopy c.MatchInformation
	matchInfoCopy.Count = matchInfo.Count
	matchInfoCopy.Matches = matchInfo.Matches
	return &matchInfoCopy
}

func sortMatchInfoByKeyThenFile(matchInfo *c.MatchInformation) {
	sort.Slice(matchInfo.Matches, func(i, j int) bool {
		if matchInfo.Matches[i].Key == matchInfo.Matches[j].Key {
			return matchInfo.Matches[i].File < matchInfo.Matches[j].File
		}
		return matchInfo.Matches[i].Key < matchInfo.Matches[j].Key
	})
}

func printAndCountMatches(matchInfo *c.MatchInformation, searchTerms []string) int {
	padding := calcPadding(searchTerms)
	uniqFiles := make([]string, 0)
	customFmt := fmt.Sprintf("%%-%ds: %%s line %%v", padding)
	for _, item := range matchInfo.Matches {
		Print(fmt.Sprintf(customFmt, item.Key, item.File, item.LineNum))
		if !slices.Contains(uniqFiles, item.File) {
			uniqFiles = append(uniqFiles, item.File)
		}
	}
	return len(uniqFiles)
}

func calcPadding(searchTerms []string) int {
	padding := -1
	for _, t := range searchTerms {
		if len(t) > padding {
			padding = len(t)
		}
	}
	return padding
}
