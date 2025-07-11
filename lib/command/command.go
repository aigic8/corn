package command

import (
	"context"
	"fmt"
	"io"
	"time"

	"mvdan.cc/sh/v3/interp"
)

type RunCommandOpts struct {
	Cmd     *ParsedCommand
	Stdout  io.Writer
	Stderr  io.Writer
	Timeout time.Duration
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
