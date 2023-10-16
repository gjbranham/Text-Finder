package main

import (
	"log"
	"os"
	"path"
	"testing"
)

var test_dir = "/tmp/text-finder"

func TestMain(m *testing.M) {
	if err := os.Mkdir(test_dir, 0777); err != nil {
		log.Fatalf("Failed to create testing directory '%v': %v\n", test_dir, err)
	}

	exitCode := m.Run()

	if err := os.RemoveAll(test_dir); err != nil {
		log.Fatalf("Failed to remove testing directory '%v': %v\n", test_dir, err)
	}

	os.Exit(exitCode)
}

func TestBasicUsage(t *testing.T) {
	WriteFile("TestBasicUsage", "foo")

}

func WriteFile(name string, content string) {
	path := path.Join(test_dir, name)
	if err := os.WriteFile(path, []byte(content), 0777); err != nil {
		log.Fatalf("Failed to write file '%v': %v\n", path, err)
	}
}
