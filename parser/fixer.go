package parser

import (
	"github.com/extrame/otto/file"
	"github.com/extrame/otto/token"
)

type Fixture struct {
	Str string
}

type Fixer interface {
	Fix(string, file.Idx, token.Token) (*Fixture, error)
}
