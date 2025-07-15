package common

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type ParsedCmd = syntax.File

type RunCommandOpts struct {
	Cmd     *ParsedCmd
	Stdout  io.Writer
	Stderr  io.Writer
	Timeout time.Duration
}

func ParseCommand(cmd string) (*ParsedCmd, error) {
	parser := syntax.NewParser()
	// TODO: maybe the program (the second parameter) is the shell runner (bash/zsh/fish)?
	// check if it is and add the custom shell functionality
	return parser.Parse(strings.NewReader(cmd), "")
}

func RunCommand(o *RunCommandOpts) error {
	runner, err := interp.New(interp.StdIO(nil, o.Stdout, o.Stderr))
	if err != nil {
		return fmt.Errorf("creating a runner: %w", err)
	}
	ctx := context.Background()
	var cancel context.CancelFunc
	if o.Timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, o.Timeout)
		defer cancel()
	}
	if err := runner.Run(ctx, o.Cmd); err != nil {
		return fmt.Errorf("running command: %w", err)
	}
	return nil
}

func JoinCommandParts(command string, args []string) string {
	return command + " " + strings.Join(args, " ")
}
