package commandeer

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

func Create() *LocalCmd {
	return &LocalCmd{defaultShell: "bash"}
}

// Create the local command struct
type LocalCmd struct {
	defaultShell string
}

func (lc *LocalCmd) Create() *LocalCmd {
	return &LocalCmd{defaultShell: "bash"}
}

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
	wrappedCmd := exec.CommandContext(ctx, lc.defaultShell, "-c", cmd)
	//wrappedCmd := exec.CommandContext(ctx, "bash", "-c", "ls")
	output, err := wrappedCmd.CombinedOutput()

	return output, err
}

// TryCmd(cmd string, lerr string, hard bool) error
func (lc *LocalCmd) TryCmd(cmd string, lerr string, hard bool) error {
	// TODO Actually write this
	// TODO Use lerr to log the error
	fmt.Printf("====================\nCommand to run is %s\n====================\n", cmd)
	return nil
}

// TryCmds(c osCmds) error
func (lc *LocalCmd) TryCmds(ctx context.Context, c *osCmds) error {
	// Cycle through the provided commands, trying them one at at time
	for i := range c.cmds {
		err := lc.TryCmd(
			c.cmds[i],
			c.errmsg[i],
			c.hard[i],
		)

		if err != nil {
			return errors.New(c.errmsg[i])
		}
	}

	return nil
}
