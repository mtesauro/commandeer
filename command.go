package main

import "context"

// Cmd is an interface for running OS commands
type Command interface {
	Exec(ctx context.Context, cmd string) ([]byte, error)
	SetShell(s string) error
	PrintShell() string
}
