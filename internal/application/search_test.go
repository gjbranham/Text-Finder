package application

import (
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/gjbranham/Text-Finder/internal/args"
	c "github.com/gjbranham/Text-Finder/internal/concurrency"
	o "github.com/gjbranham/Text-Finder/internal/output"
)

var testDir = "/tmp/text-finder"

func TestMain(m *testing.M) {
	fileSetup()

	o.SetPrinter(&o.Stdout{})

	exitCode := m.Run()

	fileTeardown()

	os.Exit(exitCode)
}

// func TestSimpleRootLevelSearch(t *testing.T) {
// 	fileName := "TestSimpleRootLevelSearch"
// 	key := "foo"
// 	writeTestFile(fileName, key)

// 	os.Args = []string{"myapp", "-d", testDir, key}
// 	args, _, _ := args.ProcessArgs(os.Args[0], os.Args[1:])

// 	app := TextFinder{Args: args, MatchInfo: new(c.MatchInformation)}

// 	app.FindFiles(testDir)

// 	if len(app.MatchInfo.Matches) != 1 {
// 		t.Errorf("Wrong number of file matches\nExpected: 1\nGot:      %v\n", len(app.MatchInfo.Matches))
// 	}

// 	matchInfo := app.MatchInfo.Matches[0]

// 	if matchInfo.File != filepath.Join(testDir, fileName) {
// 		t.Errorf("File found does not match expected\nExpected: %v\nGot:      %v\n", testDir, matchInfo.File)
// 	}

// 	if matchInfo.Key != key {
// 		t.Errorf("Search terms do not match expected\nExpected: %v\nGot:      %v\n", key, matchInfo.Key)
// 	}

// 	if matchInfo.LineNum != 1 {
// 		t.Errorf("Wrong line number saved for match\nExpected: %v\nGot:      %v\n", 1, matchInfo.LineNum)
// 	}

// 	removeTestFile(fileName)
// }

// func TestRecursiveSearch(t *testing.T) {
// 	recursiveDir := "recursiveDir"
// 	if err := os.Mkdir(path.Join(testDir, recursiveDir), 0777); err != nil {
// 		t.Errorf("Failed to create subdirectory for recursive search test: %v\n", err)
// 	}
// 	fileName := "TestRecursiveSearch"
// 	key := "foo"
// 	writeTestFile(recursiveDir+"/"+fileName, key)

// 	os.Args = []string{"myapp", "-r", "-d", testDir, key}
// 	args, _, _ := args.ProcessArgs(os.Args[0], os.Args[1:])

// 	app := TextFinder{Args: args, MatchInfo: new(c.MatchInformation)}

// 	app.FindFiles(testDir)

// 	if len(app.MatchInfo.Matches) != 1 {
// 		t.Errorf("Wrong number of file matches\nExpected: 1\nGot:      %v\n", len(app.MatchInfo.Matches))
// 	}

// 	matchInfo := app.MatchInfo.Matches[0]

// 	if matchInfo.File != filepath.Join(testDir+"/"+recursiveDir, fileName) {
// 		t.Errorf("File found does not match expected\nExpected: %v\nGot:      %v\n", testDir, matchInfo.File)
// 	}

// 	if matchInfo.Key != key {
// 		t.Errorf("Search terms do not match expected\nExpected: %v\nGot:      %v\n", key, matchInfo.Key)
// 	}

// 	if matchInfo.LineNum != 1 {
// 		t.Errorf("Wrong line number saved for match\nExpected: %v\nGot:      %v\n", 1, matchInfo.LineNum)
// 	}

// 	removeTestFile(recursiveDir + "/" + fileName)
// }

func TestTable(t *testing.T) {
	type test struct {
		args         []string
		matches      int
		matchingFile string
		matchingKey  string
		matchingLine int
	}

	tests := []test{
		{
			args: []string{"myapp", "-d", testDir, "foo"}, matches: 1, matchingFile: "simpleSearch", matchingKey: "foo", matchingLine: 1,
		},
		{
			args: []string{"myapp", "-d", testDir, "foo", "bar"}, matches: 2, matchingFile: "multiItemSearch", matchingKey: "foo", matchingLine: 1,
		},
		{
			args: []string{"myapp", "-i", "-d", testDir, "Foo"}, matches: 1, matchingFile: "caseInsensitive", matchingKey: "Foo", matchingLine: 1,
		},
		{
			args: []string{"myapp", "-r", "-d", testDir, "foo"}, matches: 1, matchingFile: "recursiveDir/recursiveSearch", matchingKey: "foo", matchingLine: 1,
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			writeTestFile(tt.matchingFile, "foo bar")
			os.Args = tt.args
			args, _, _ := args.ProcessArgs(os.Args[0], os.Args[1:])
			app := TextFinder{Args: args, MatchInfo: new(c.MatchInformation)}
			app.FindFiles(testDir)

			if len(app.MatchInfo.Matches) != tt.matches {
				t.Errorf("Wrong number of file matches\nExpected: %v\nGot:      %v\n", tt.matches, len(app.MatchInfo.Matches))
			}

			matchInfo := app.MatchInfo.Matches[0]

			if matchInfo.File != path.Join(testDir, tt.matchingFile) {
				t.Errorf("File found does not match expected\nExpected: %v\nGot:      %v\n", tt.matchingFile, matchInfo.File)
			}
			if matchInfo.Key != tt.matchingKey {
				t.Errorf("Search terms do not match expected\nExpected: %v\nGot:      %v\n", tt.matchingKey, matchInfo.Key)
			}
			if matchInfo.LineNum != tt.matchingLine {
				t.Errorf("Wrong line number saved for match\nExpected: %v\nGot:      %v\n", tt.matchingLine, matchInfo.LineNum)
			}

			removeTestFile(tt.matchingFile)
		})
	}
}

func writeTestFile(name string, content string) {
	path := path.Join(testDir, name)
	if err := os.WriteFile(path, []byte(content), 0777); err != nil {
		log.Fatalf("Failed to write file: %v\n", err)
	}
}

func removeTestFile(name string) {
	path := path.Join(testDir, name)
	if err := os.Remove(path); err != nil {
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
