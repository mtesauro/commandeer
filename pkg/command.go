package commandeer

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

/////////////////////////////
// Command Package Structs //
/////////////////////////////

// Holds a named package of commands for one or more targets
type CmdPkg struct {
	Label    string   // Holds a human friendly label for the command package
	Targets  []Target // Hold the target(s) the commands where the commands run
	Location Terminal // Holds the command-line interface for local or remote command invocation
}

// TODO: Decide if this is needed
// Collection of multiple command packages
type CmdCollection struct {
	Label      string   // Holds a human friendly label for the collection
	Collection []CmdPkg // Allows commands for multiple targets to be stored together
}

// Set the Location - either local aka on the host or remove aka over SSH
func (cp *CmdPkg) SetLocation(cl Terminal) {
	cp.Location = cl
}

// Set the target for the collection of commands
func (cp *CmdPkg) AddTarget(id string, dist string, rel string, os string, sh string) {
	tg := Target{
		ID:      id,
		Release: rel,
		Distro:  dist,
		OS:      os,
		Shell:   sh,
		PkgCmds: []SingleCmd{},
	}
	cp.Targets = append(cp.Targets, tg)
}

// Add a single command to a command package target
func (cp *CmdPkg) AddCmd(c string, e string, h bool, d time.Duration, t string) error {
	// Check that the target exists
	tg, err := FindTarget(cp, t)
	if err != nil {
		return err
	}

	// Add the command
	cmd := SingleCmd{
		Cmd:     c,
		Errmsg:  e,
		Hard:    h,
		Timeout: d,
	}
	tg.PkgCmds = append(tg.PkgCmds, cmd)

	return errors.New(fmt.Sprintf("Cannot add command to non-existent target %s", t))
}

// Add a slice of SingleCmd to a command package target
func (cp *CmdPkg) LoadCmds(c []SingleCmd, t string) error {
	// Check that the target exists
	tg, err := FindTarget(cp, t)
	if err != nil {
		return err
	}

	// Add the provided commands to the command package
	for i := range c {
		// Iterate through sent commands c and append to the target's collection
		tg.PkgCmds = append(tg.PkgCmds, c[i])
	}

	return nil
}

// Execute the commands for the provided target t returning a slice of bytes
// representing stdout and stderr combined for the commands run. An error is
// returned if the target isn't found in the command package or an error
// occurs during running the commands.
func (cp *CmdPkg) ExecPkgCombined(t string) ([]byte, error) {
	// Check that the target exists
	tg, err := FindTarget(cp, t)
	if err != nil {
		return nil, err
	}

	// Setup to run multiple commands
	var fullOut []byte
	for k := range tg.PkgCmds {
		// Set a default contenxt
		ctx := context.Background()

		// Does thei command have a timeout
		if tg.PkgCmds[k].Timeout != 0 {
			// Set a timeout with context
			new, cancel := context.WithTimeout(context.Background(), tg.PkgCmds[k].Timeout)
			ctx = new
			defer cancel()

		}
		out, err := cp.Location.ExecCombined(ctx, tg.PkgCmds[k].Cmd, tg.Shell)
		if err != nil {
			return nil, err
		}
		fullOut = append(fullOut, out...)
	}

	return fullOut, nil
}

// Execute the commands for the provided target t returning only a Go
// error  if the target isn't found in the command package or an error
// occurs during running the commands.  Stdout and Stderr are silently
// dropped.
func (cp *CmdPkg) ExecPkgError(t string) error {
	// Check that the target exists
	tg, err := FindTarget(cp, t)
	if err != nil {
		return err
	}

	// Setup to run multiple commands
	for k := range tg.PkgCmds {
		// Set a default contenxt
		ctx := context.Background()

		// Does thei command have a timeout
		if tg.PkgCmds[k].Timeout != 0 {
			// Set a timeout with context
			new, cancel := context.WithTimeout(context.Background(), tg.PkgCmds[k].Timeout)
			ctx = new
			defer cancel()

		}
		err := cp.Location.ExecError(ctx, tg.PkgCmds[k].Cmd, tg.Shell)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute the commands for the provided target t returning only a Go
// error  if the target isn't found in the command package. Stdout and
// Stderr are silently dropped.
func (cp *CmdPkg) ExecPkgOnly(t string) error {
	// Check that the target exists
	tg, err := FindTarget(cp, t)
	if err != nil {
		return err
	}

	// Setup to run multiple commands
	for k := range tg.PkgCmds {
		// Set a default contenxt
		ctx := context.Background()

		// Does thei command have a timeout
		if tg.PkgCmds[k].Timeout != 0 {
			// Set a timeout with context
			new, cancel := context.WithTimeout(context.Background(), tg.PkgCmds[k].Timeout)
			ctx = new
			defer cancel()

		}
		cp.Location.ExecOnly(ctx, tg.PkgCmds[k].Cmd, tg.Shell)
	}

	return nil
}

// Execute the commands for the provided target t returning a slice of bytes
// representing only stdout for the commands run. An error is returned if the
// target isn't found in the command package or an error occurs during running
// the commands.
func (cp *CmdPkg) ExecPkgStdout(t string) ([]byte, error) {
	// Check that the target exists
	tg, err := FindTarget(cp, t)
	if err != nil {
		return nil, err
	}

	// Setup to run multiple commands
	var fullOut []byte
	for k := range tg.PkgCmds {
		// Set a default contenxt
		ctx := context.Background()

		// Does thei command have a timeout
		if tg.PkgCmds[k].Timeout != 0 {
			// Set a timeout with context
			new, cancel := context.WithTimeout(context.Background(), tg.PkgCmds[k].Timeout)
			ctx = new
			defer cancel()

		}
		out, err := cp.Location.ExecStdout(ctx, tg.PkgCmds[k].Cmd, tg.Shell)
		if err != nil {
			return nil, err
		}
		fullOut = append(fullOut, out...)
	}

	return fullOut, nil
}

// Execute the commands for the provided target t returning a slice of bytes
// representing only stdout for the commands run. An error is returned if the
// target isn't found in the command package or an error occurs during running
// the commands.
// ExecStderr(ctx context.Context, cmd string, shell string) ([]byte, error)
func (cp *CmdPkg) ExecPkgStderr(t string) ([]byte, error) {
	// Check that the target exists
	tg, err := FindTarget(cp, t)
	if err != nil {
		return nil, err
	}

	// Setup to run multiple commands
	var fullOut []byte
	for k := range tg.PkgCmds {
		// Set a default contenxt
		ctx := context.Background()

		// Does thei command have a timeout
		if tg.PkgCmds[k].Timeout != 0 {
			// Set a timeout with context
			new, cancel := context.WithTimeout(context.Background(), tg.PkgCmds[k].Timeout)
			ctx = new
			defer cancel()

		}
		out, err := cp.Location.ExecStderr(ctx, tg.PkgCmds[k].Cmd, tg.Shell)
		if err != nil {
			return nil, err
		}
		fullOut = append(fullOut, out...)
	}

	return fullOut, nil
}

///////////////////////////////////////
// Command Package Utility functions //
///////////////////////////////////////

// Create a new empty command package with the provided label l
func NewPkg(l string) *CmdPkg {
	return &CmdPkg{
		Label:    l,
		Targets:  []Target{},
		Location: LocalTerm{},
	}
}

// Look for the provided target t in a command package and either
// return a pointer to the Target struct or an error if the target
// cannot be found
func FindTarget(cp *CmdPkg, t string) (*Target, error) {
	for k := range cp.Targets {
		if strings.Compare(cp.Targets[k].ID, t) == 0 {
			// Found a matching target
			return &cp.Targets[k], nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Command Package does not support target %s", t))
}
