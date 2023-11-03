package args

import (
	"bytes"
	"flag"
)

type Arguments struct {
	CaseInsensitive bool
	RecursiveSearch bool
	RootPath        string
	SearchTerms     []string
}

func ProcessArgs(exeName string, sysArgs []string) (*Arguments, string, error) {
	flags := flag.NewFlagSet(exeName, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var args Arguments
	flags.BoolVar(&args.CaseInsensitive, "i", false, "Turns on case insensitivity")
	flags.BoolVar(&args.RecursiveSearch, "r", false, "Turns on recursive search. Will traverse all sub-directories of root directory")
	flags.StringVar(&args.RootPath, "d", ".", "Root directory to start searching for matches")

	if err := flags.Parse(sysArgs); err != nil {
		return nil, buf.String(), err
	}

	args.SearchTerms = flags.Args()

	return &args, buf.String(), nil
}
