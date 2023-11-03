package application

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/gjbranham/Text-Finder/internal/args"
	c "github.com/gjbranham/Text-Finder/internal/concurrency"
)

var testDir = "/tmp/text-finder"

func TestMain(m *testing.M) {
	directorySetup()

	exitCode := m.Run()

	directoryTeardown()

	os.Exit(exitCode)
}

func TestSimpleRootLevelSearch(t *testing.T) {
	fileName := "TestSimpleRootLevelSearch"
	key := "foo"
	writeTestFile(fileName, key)

	os.Args = []string{"myapp", "-d", fmt.Sprintf("%v", testDir), key}
	args, _, _ := args.ProcessArgs(os.Args[0], os.Args[1:])

	app := TextFinder{Args: args, MatchInfo: new(c.MatchInformation)}

	app.FindFiles(testDir)

	if len(app.MatchInfo.Matches) != 1 {
		t.Errorf("Wrong number of file matches\nExpected: 1\nGot:      %v\n", len(app.MatchInfo.Matches))
	}

	matchInfo := app.MatchInfo.Matches[0]

	if matchInfo.File != filepath.Join(testDir, fileName) {
		t.Errorf("File found does not match expected\nExpected: %v\nGot:      %v\n", testDir, matchInfo.File)
	}

	if matchInfo.Key != key {
		t.Errorf("Search terms do not match expected\nExpected: %v\nGot:      %v\n", key, matchInfo.Key)
	}

	if matchInfo.LineNum != 1 {
		t.Errorf("Wrong line number saved for match\nExpected: %v\nGot:      %v\n", 1, matchInfo.LineNum)
	}

	removeTestFile(fileName)
}

func TestRecursiveSearch(t *testing.T) {
	recursiveDir := "recursiveDir"
	if err := os.Mkdir(path.Join(testDir, recursiveDir), 0777); err != nil {
		t.Errorf("Failed to create subdirectory for recursive search test: %v\n", err)
	}
	fileName := "TestRecursiveSearch"
	key := "foo"
	writeTestFile(recursiveDir+"/"+fileName, key)

	os.Args = []string{"myapp", "-r", "-d", fmt.Sprintf("%v", testDir), key}
	args, _, _ := args.ProcessArgs(os.Args[0], os.Args[1:])

	app := TextFinder{Args: args, MatchInfo: new(c.MatchInformation)}

	app.FindFiles(testDir)

	if len(app.MatchInfo.Matches) != 1 {
		t.Errorf("Wrong number of file matches\nExpected: 1\nGot:      %v\n", len(app.MatchInfo.Matches))
	}

	matchInfo := app.MatchInfo.Matches[0]

	if matchInfo.File != filepath.Join(testDir+"/"+recursiveDir, fileName) {
		t.Errorf("File found does not match expected\nExpected: %v\nGot:      %v\n", testDir, matchInfo.File)
	}

	if matchInfo.Key != key {
		t.Errorf("Search terms do not match expected\nExpected: %v\nGot:      %v\n", key, matchInfo.Key)
	}

	if matchInfo.LineNum != 1 {
		t.Errorf("Wrong line number saved for match\nExpected: %v\nGot:      %v\n", 1, matchInfo.LineNum)
	}

	removeTestFile(recursiveDir + "/" + fileName)
}

func writeTestFile(name string, keyWord string) {
	path := path.Join(testDir, name)
	if err := os.WriteFile(path, []byte(keyWord), 0777); err != nil {
		log.Fatalf("Failed to write file: %v\n", err)
	}
}

func removeTestFile(name string) {
	path := path.Join(testDir, name)
	if err := os.Remove(path); err != nil {
		log.Fatalf("Failed to remove file: %v\n", err)
	}
}

func directorySetup() {
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		if err == nil {
			if err := os.RemoveAll(testDir); err != nil {
				log.Fatalf("Test directory already exists but could not be removed: %v\n", err)
			}
		} else {
			log.Fatalf("Fatal error while attempting to retrieve information for '%v': %v\n", testDir, err)
		}
	}

	if err := os.Mkdir(testDir, 0777); err != nil {
		log.Fatalf("Failed to create testing directory: %v\n", err)
	}
}

func directoryTeardown() {
	if err := os.RemoveAll(testDir); err != nil {
		log.Fatalf("Failed to remove testing directory: %v\n", err)
	}
}
