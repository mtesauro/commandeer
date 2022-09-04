package commandeer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

/////////////////////////////
// Command Package Structs //
/////////////////////////////

// Holds a named package of commands for one or more targets
type CmdPkg struct {
	Label        string                 // Holds a human friendly label for the command package
	Targets      []Target               // Hold the target(s) the commands where the commands run
	Location     Terminal               // Holds the command-line interface for local or remote command invocation
	Log          map[string]*log.Logger // Holds a pointer to a Logger using Go's log module from the stadard library
	Redact       bool                   // Determine if redaction should be turned on; defaults to true
	StrRedact    []string               // List of strings to redact from logging
	CmdLog       *log.Logger            // Optional 2nd logger to write out commands and output from them
	EnableCmdLog bool                   /// Determine if commands should be logged; defaults to false
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
	// Sanity checks
	if len(c) == 0 {
		return errors.New("Command was empty")
	}
	if len(c) == 0 {
		return errors.New("Error message was empty")
	}
	if d < 0 {
		return errors.New("Timeout cannot be negative and must be zero or greater")
	}


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

	// Sanity check length of PkgCmds to ensure there's at least 1
	if len(tg.PkgCmds) == 0 {
		return nil, errors.New("Cannot Exec a Package with no commands added, please add commands before Exec'ing")
	}

	// Setup to run multiple commands
	var fullOut []byte
	for k := range tg.PkgCmds {
		// Set a default contenxt
		ctx := context.Background()

		// Sanity check duration - it should be greater than or equal to zero
		if tg.PkgCmds[k].Timeout < 0 {
			return nil, errors.New("Timeout cannot be negative and must be zero or greater")
		}

		// Since duration is >= 0 and not zero, create a new context.WithTimeout for that duration
		if tg.PkgCmds[k].Timeout != 0 {
			// Set a timeout with context
			new, cancel := context.WithTimeout(context.Background(), tg.PkgCmds[k].Timeout)
			ctx = new
			defer cancel()

		}

		// Execute the command
		out, err := cp.Location.ExecCombined(ctx, tg.PkgCmds[k].Cmd, tg.Shell)
		if err != nil {
			return nil, err
		}

		// Log command if configured
		if cp.EnableCmdLog {
			cp.LogCmd(tg.PkgCmds[k].Cmd + "\n" + string(out))
		}

		// Gather output to return
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

	// Sanity check length of PkgCmds to ensure there's at least 1
	if len(tg.PkgCmds) == 0 {
		return errors.New("Cannot Exec a Package with no commands added, please add commands before Exec'ing")
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

		// Log command if configured
		if cp.EnableCmdLog {
			cp.LogCmd(tg.PkgCmds[k].Cmd)
		}

		// Execute the command
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

	// Sanity check length of PkgCmds to ensure there's at least 1
	if len(tg.PkgCmds) == 0 {
		return errors.New("Cannot Exec a Package with no commands added, please add commands before Exec'ing")
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

		// Log command if configured
		if cp.EnableCmdLog {
			cp.LogCmd(tg.PkgCmds[k].Cmd)
		}

		// Execute the command
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

	// Sanity check length of PkgCmds to ensure there's at least 1
	if len(tg.PkgCmds) == 0 {
		return nil, errors.New("Cannot Exec a Package with no commands added, please add commands before Exec'ing")
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

		// Execute the command
		out, err := cp.Location.ExecStdout(ctx, tg.PkgCmds[k].Cmd, tg.Shell)
		if err != nil {
			return nil, err
		}

		// Log command if configured
		if cp.EnableCmdLog {
			cp.LogCmd(tg.PkgCmds[k].Cmd + "\n" + string(out))
		}

		// Gather the output to return
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

	// Sanity check length of PkgCmds to ensure there's at least 1
	if len(tg.PkgCmds) == 0 {
		return nil, errors.New("Cannot Exec a Package with no commands added, please add commands before Exec'ing")
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

		// Execute the command
		out, err := cp.Location.ExecStderr(ctx, tg.PkgCmds[k].Cmd, tg.Shell)
		if err != nil {
			return nil, err
		}

		// Log command if configured
		if cp.EnableCmdLog {
			cp.LogCmd(tg.PkgCmds[k].Cmd + "\n" + string(out))
		}

		// Gather the output to return
		fullOut = append(fullOut, out...)
	}

	return fullOut, nil
}

// SetLogging takes an io.Writer and creates new log.loggers for the following
// logging levels: trace, info, warning, error. All log levels will be written
// to the same io.Writer with the level prepended to the log line.
func (cp *CmdPkg) SetLogging(logHandler io.Writer) {
	cp.Log["trace"] = log.New(logHandler, "TRACE:   ", log.Ldate|log.Ltime|log.Lmsgprefix)
	cp.Log["info"] = log.New(logHandler, "INFO:    ", log.Ldate|log.Ltime|log.Lmsgprefix)
	cp.Log["warning"] = log.New(logHandler, "WARNING: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	cp.Log["error"] = log.New(logHandler, "ERROR:   ", log.Ldate|log.Ltime|log.Lmsgprefix)
}

// Write a trace level message to the configured logger
func (cp *CmdPkg) LogTrace(msg string) {
	cp.Log["trace"].Println(cp.Redactatron(msg))
}

// Write a info level message to the configured logger
func (cp *CmdPkg) LogInfo(msg string) {
	cp.Log["info"].Println(cp.Redactatron(msg))
}

// Write a warning level message to the configured logger
func (cp *CmdPkg) LogWarn(msg string) {
	cp.Log["warning"].Println(cp.Redactatron(msg))
}

// Write a error level message to the configured logger
func (cp *CmdPkg) LogError(msg string) {
	cp.Log["error"].Println(cp.Redactatron(msg))
}

// Turn off redacting of log messages
func (cp *CmdPkg) TurnOffRedaction() {
	cp.Redact = false
}

// Add a single item to the list of items to redact from log messages
func (cp *CmdPkg) AddRedact(s string) {
	cp.StrRedact = append(cp.StrRedact, s)
}

// Add a slice of strings to the list of items to redact from logs messages
func (cp *CmdPkg) AddRedactSlice(s []string) {
	for i := range s {
		if len(s[i]) > 0 {
			cp.StrRedact = append(cp.StrRedact, s[i])
		}
	}
}

// Redactatron - redacts sensitive information from being written to the logs
// Redaction is configurable with Install's Redact boolean config.
// If true (the default), sensitive info will be redacted
func (cp *CmdPkg) Redactatron(l string) string {
	// Redact sensitive data if it's turned on
	if cp.Redact {
		// Redact sensitive info from the logger
		clean := l
		r := "[~REDACTED~]"
		for i := range cp.StrRedact {
			if strings.Contains(clean, cp.StrRedact[i]) {
				clean = strings.Replace(clean, cp.StrRedact[i], r, -1)
			}
		}
		return clean
	}
	return l
}

func (cp *CmdPkg) SetCmdLog(logHandler io.Writer) {
	cp.CmdLog = log.New(logHandler, "<commandeer> $ ", log.Ldate|log.Ltime|log.Lmsgprefix)
}

// Turn on logging of commands, assumes you have properly configured CmdLog
func (cp *CmdPkg) TurnOnCmdLog() {
	cp.EnableCmdLog = true
}

// Write a error level message to the configured logger
func (cp *CmdPkg) LogCmd(msg string) {
	cp.CmdLog.Println(msg)
}

///////////////////////////////////////
// Command Package Utility functions //
///////////////////////////////////////

// Create a new empty command package with the provided label l
func NewPkg(l string) *CmdPkg {
	lg := make(map[string]*log.Logger)
	str := make([]string, 0)
	return &CmdPkg{
		Label:        l,
		Targets:      []Target{},
		Location:     &LocalTerm{},
		Log:          lg,
		Redact:       true,
		StrRedact:    str,
		CmdLog:       &log.Logger{},
		EnableCmdLog: false,
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

// LogToFile takes a path and a file name and attempts to create the path and file name
// as the target for the CmdPkg.Log output by returning a io.Writer to the resulting
// path + file. For example, sending in "./logs" and "app.log" will return a io.Writer
// pointing at "./logs/app.log" relative to where the Go program is run
func LogToFile(fpath string, fname string) (io.Writer, error) {
	// Create the full path
	logPath := path.Join(fpath, fname)
	// Create the logs directory if it does not exist
	_, err := os.Stat(logPath)
	if err != nil {
		// logs directory doesn't exist
		err = os.MkdirAll(fpath, 0755)
		if err != nil {
			// Can't create logs directory for some reason, return sending the error
			return nil, err
		}
	}

	// Create log file
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return logFile, nil
}