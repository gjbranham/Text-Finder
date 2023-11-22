package application

import (
	"bytes"
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/gjbranham/Text-Finder/internal/args"
	c "github.com/gjbranham/Text-Finder/internal/concurrency"
	o "github.com/gjbranham/Text-Finder/internal/output"
)

/***
- Need to add a test to verify multiple search terms from a file
***/

var testDir = "/tmp/text-finder"
var buf *bytes.Buffer = new(bytes.Buffer)

func TestMain(m *testing.M) {
	fileSetup()

	o.SetPrinter(&o.Buffer{Buf: buf})

	exitCode := m.Run()

	fileTeardown()

	os.Exit(exitCode)
}

func TestSimpleCmdLineArguments(t *testing.T) {
	type test struct {
		args              []string
		fileContent       string
		expectedMatchInfo c.MatchInformation
	}

	tests := []test{
		{
			args:        []string{"myapp", "-d", testDir, "foo"},
			fileContent: "foo",
			expectedMatchInfo: c.MatchInformation{
				Count: 1,
				Matches: []c.FileInfo{
					{Key: "foo", File: path.Join(testDir, "simpleSearch"), LineNum: 1},
				},
			},
		},
		{
			args:        []string{"myapp", "-i", "-d", testDir, "Foo"},
			fileContent: "foo",
			expectedMatchInfo: c.MatchInformation{
				Count: 1,
				Matches: []c.FileInfo{
					{Key: "Foo", File: path.Join(testDir, "caseInsensitiveSearch"), LineNum: 1},
				},
			},
		},
		{
			args:        []string{"myapp", "-i", "-d", "", "foo"},
			fileContent: "foo",
			expectedMatchInfo: c.MatchInformation{
				Count: 1,
				Matches: []c.FileInfo{
					{Key: "foo", File: path.Join(testDir, "emptyRootDir"), LineNum: 1},
				},
			},
		},
		{
			args:        []string{"myapp", "-r", "-d", testDir, "foo"},
			fileContent: "foo",
			expectedMatchInfo: c.MatchInformation{
				Count: 1,
				Matches: []c.FileInfo{
					{Key: "foo", File: path.Join(testDir, "recursiveDir/recursiveSearch"), LineNum: 1},
				},
			},
		},
		{
			// no search terms
			args:        []string{"myapp", "-d", testDir},
			fileContent: "",
			expectedMatchInfo: c.MatchInformation{
				Count:   0,
				Matches: []c.FileInfo{},
			},
		},
		{
			// binary file
			args:        []string{"myapp", "-d", testDir},
			fileContent: "foo" + string([]byte{0}),
			expectedMatchInfo: c.MatchInformation{
				Count:   0,
				Matches: []c.FileInfo{},
			},
		},
		{
			args:        []string{"myapp", "-d", testDir, "foo"},
			fileContent: "foo",
			expectedMatchInfo: c.MatchInformation{
				Count: 2,
				Matches: []c.FileInfo{
					{Key: "foo", File: path.Join(testDir, "firstMatchingFile"), LineNum: 1},
					{Key: "foo", File: path.Join(testDir, "secondMatchingFile"), LineNum: 1},
				},
			},
		},
		{
			args:        []string{"myapp", "-d", testDir, "foo"},
			fileContent: "foo\n",
			expectedMatchInfo: c.MatchInformation{
				Count: 2,
				Matches: []c.FileInfo{
					{Key: "foo", File: path.Join(testDir, "multiLineMatches"), LineNum: 1},
					{Key: "foo", File: path.Join(testDir, "multiLineMatches"), LineNum: 2},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			for _, m := range tt.expectedMatchInfo.Matches {
				writeTestFile(m.File, tt.fileContent)
			}
			os.Args = tt.args
			args, _, _ := args.ProcessArgs(os.Args[0], os.Args[1:])
			app := TextFinder{Args: args, MatchInfo: new(c.MatchInformation)}
			app.FindFiles(testDir)

			if err := checkMatch(tt.expectedMatchInfo, *app.MatchInfo); err != nil {
				t.Errorf("Test failed: %v", err)
			}
			for _, m := range tt.expectedMatchInfo.Matches {
				removeTestFile(m.File)
			}
		})
	}
}

func checkMatch(expected, actual c.MatchInformation) error {
	if expected.Count != actual.Count {
		return errors.New("matching file count did not match")
	}
	if len(expected.Matches) == 0 && len(actual.Matches) == 0 { // for searches that return no results
		return nil
	}
	for _, a := range actual.Matches {
		for _, e := range expected.Matches {
			if reflect.DeepEqual(a, e) {
				return nil
			}
		}
	}
	return errors.New("file match not found")
}

func writeTestFile(filePath string, content string) {
	fo, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open file for writing: %v\n", err)
	}

	if _, err := fo.Write([]byte(content)); err != nil {
		log.Fatalf("Failed to write file: %v\n", err)
	}

	// if err := os.WriteFile(filePath, []byte(content), 0777); err != nil {
	// 	log.Fatalf("Failed to write file: %v\n", err)
	// }
}

func removeTestFile(filePath string) {
	if _, err := os.Stat(filePath); errors.Is(err, fs.ErrNotExist) {
		return
	}
	if err := os.Remove(filePath); err != nil {
		log.Fatalf("Failed to remove file: %v\n", err)
	}
}

func fileSetup() {
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		if err == nil {
			if err := os.RemoveAll(testDir); err != nil {
				log.Fatalf("Test directory already exists but could not be removed: %v\n", err)
			}
		} else {
			log.Fatalf("Fatal error while attempting to retrieve information for '%v': %v\n", testDir, err)
		}
	}

	if err := os.MkdirAll(path.Join(testDir, "recursiveDir"), 0777); err != nil {
		log.Fatalf("Failed to create testing directory: %v\n", err)
	}
}

func fileTeardown() {
	if err := os.RemoveAll(testDir); err != nil {
		log.Fatalf("Failed to remove testing directory: %v\n", err)
	}
}
