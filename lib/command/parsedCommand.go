package command

import (
	"io"

	"mvdan.cc/sh/v3/syntax"
)

type (
	CommandParser struct {
		p *syntax.Parser
	}
)

type ParsedCommand = syntax.File

func NewCommandParser() *CommandParser {
	return &CommandParser{p: syntax.NewParser()}
}

func (p *CommandParser) Parse(r io.Reader) (*ParsedCommand, error) {
	return p.p.Parse(r, "")
}
