package parser

import (
	"errors"
	"strings"
)

// Parser represents parser interface for other
type Parser interface {
	Parse([]byte) (string, string)
}

var NotFound = errors.New("Parser not found")

// New created new parser
func New(name string) (Parser, error) {
	switch {
	case strings.HasSuffix(name, ".md"):
		return NewMdParser(), nil
	}

	return nil, NotFound
}
