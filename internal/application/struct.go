package application

import (
	"github.com/gjbranham/Text-Finder/internal/args"
	c "github.com/gjbranham/Text-Finder/internal/concurrency"
)

type TextFinder struct {
	Args      *args.Arguments
	MatchInfo *c.MatchInformation
}
