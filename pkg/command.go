package commandeer

import "context"

// Cmd is an interface for running OS commands
type Command interface {
	Exec(ctx context.Context, cmd string) ([]byte, error)
	SetShell(s string) error
	PrintShell() string
	TryCmd(cmd string, lerr string, hard bool) error
	TryCmds(ctx context.Context, c *osCmds) error
	//LoadCmds( SOMETHING) error
	// validate commands call to ensure they are in the path?
}
