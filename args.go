package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func processArgs() (string, string) {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Print("Please provide a start dir and a save dir")
		os.Exit(1)
	}

	start, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatalf("Failed generating abspath: %v", err)
	}

	save, err := filepath.Abs(args[1])
	if err != nil {
		log.Fatalf("Failed generating abspath: %v", err)
	}

	if stat, err := os.Stat(start); err != nil {
		log.Fatalf("os.Stat failed: %v", err)
	} else if !stat.IsDir() {
		log.Fatalf("'%v' is not a directory", start)
	}

	if stat, err := os.Stat(save); err != nil {
		os.Mkdir(save, 0644)
	} else if !stat.IsDir() {
		log.Fatalf("%v is a file, not directory", start)
	}

	return start, save
}
