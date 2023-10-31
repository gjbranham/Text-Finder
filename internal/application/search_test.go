package application

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/gjbranham/Text-Finder/internal/args"
	c "github.com/gjbranham/Text-Finder/internal/concurrency"
	o "github.com/gjbranham/Text-Finder/internal/output"
)

var test_dir = "/tmp/text-finder"
var buf *bytes.Buffer = new(bytes.Buffer)

func TestMain(m *testing.M) {
	if err := os.Mkdir(test_dir, 0777); err != nil {
		log.Fatalf("Failed to create testing directory '%v': %v\n", test_dir, err)
	}

	o.SetPrinter(&o.Buffer{Buf: buf})

	exitCode := m.Run()

	if err := os.RemoveAll(test_dir); err != nil {
		log.Fatalf("Failed to remove testing directory '%v': %v\n", test_dir, err)
	}

	os.Exit(exitCode)
}

func TestBasicUsage(t *testing.T) {
	fileName := "TestBasicUsage"
	key := "foo"
	WriteFile(fileName, key)

	os.Args = []string{"myapp", "-d", fmt.Sprintf("%v", test_dir), key}
	args, _, _ := args.ProcessArgs(os.Args[0], os.Args[1:])

	app := TextFinder{Args: args, MatchInfo: new(c.MatchInformation)}

	app.FindFiles(test_dir)

	if len(app.MatchInfo.Matches) != 1 {
		t.Errorf("Wrong number of file matches\nExpected: 1\nGot:      %v\n", len(app.MatchInfo.Matches))
	}

	matchInfo := app.MatchInfo.Matches[0]

	if matchInfo.File != filepath.Join(test_dir, fileName) {
		t.Errorf("File found does not match expected\nExpected: %v\nGot:      %v\n", test_dir, matchInfo.File)
	}

	if matchInfo.Key != key {
		t.Errorf("Search terms do not match expected\nExpected: %v\nGot:      %v\n", key, matchInfo.Key)
	}

	if matchInfo.LineNum != 1 {
		t.Errorf("Wrong line number saved for match\nExpected: %v\nGot:      %v\n", 1, matchInfo.LineNum)
	}

}

func WriteFile(name string, content string) {
	path := path.Join(test_dir, name)
	if err := os.WriteFile(path, []byte(content), 0777); err != nil {
		log.Fatalf("Failed to write file '%v': %v\n", path, err)
	}
}
