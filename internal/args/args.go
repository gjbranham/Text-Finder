package args

import (
	"bytes"
	"flag"
	"log"
)

type Arguments struct {
	RootPath        string
	RecursiveSearch bool
	SearchTerms     []string
}

func ProcessArgs(exeName string, sysArgs []string) (*Arguments, string, error) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flags := flag.NewFlagSet(exeName, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var args Arguments
	flags.StringVar(&args.RootPath, "d", "./", "Root directory to start searching for matches")
	flags.BoolVar(&args.RecursiveSearch, "r", false, "Search recursively starting at the root directory")

	if err := flags.Parse(sysArgs); err != nil {
		return nil, buf.String(), err
	}

	args.SearchTerms = flags.Args()

	return &args, buf.String(), nil
}
