package main

import (
	"context"
	"fmt"
	"os/exec"
)

func Create() *LocalCmd{
	return &LocalCmd{defaultShell: "bash"}
}

// Create the local command struct
type LocalCmd struct {
	defaultShell string
}

func (lc *LocalCmd) Create() *LocalCmd {
//	l := LocalCmd{defaultShell: "bash"}
	return &LocalCmd{defaultShell: "bash"}
}
//
//	Exec(ctx context.Context, cmd string) ([]byte, error)
//	SetShell(s string) error
//	PrintShell() string

func (lc *LocalCmd) SetShell(s string) error {
	// TODO configure 'supported' shells e.g. bash, sh
	// TODO check that shells are in the path
	lc.defaultShell = s
	return nil
}

func (lc *LocalCmd) PrintShell() string {
	return lc.defaultShell
}

func (lc *LocalCmd) Exec(ctx context.Context, cmd string) ([]byte, error) {
	// Wrap the command in the configured shell e.g. "sh -c 'command sent to method'"
	fmt.Printf("defaultShell in exec is %s\n\n", lc.defaultShell)
	fmt.Printf("LocalCmd is %v\n\n", lc)
	wrappedCmd := exec.CommandContext(ctx, "bash", "-c", cmd)
	//wrappedCmd := exec.CommandContext(ctx, "bash", "-c", "ls")
	output, err := wrappedCmd.CombinedOutput()

	return output, err
}
