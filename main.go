package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/mtesauro/commandeer/pkg/commandeer"
)

func interfaceFun(c Command) {
	fmt.Println("In interfaceFun")
	fmt.Println(c.PrintShell())
}

// osCmds - holds a group of commands to be executed
type osCmds struct {
	id     string   // Holds an optional description/label for the group of commands
	cmds   []string // Holds the OS commands
	errmsg []string // Holds the error message to log if the matching OS command fails
	hard   []bool   // Flag to know if an error on the matching command should be a hard/breaking error
}

type cmdList struct {
	list []singleCmd
}

type singleCmd struct {
	cmd    string
	errmsg string
	hard   bool
}

func newCmds(d string) *osCmds {
	var o osCmds
	o.id = d
	return &o
}

func addCmd(o *osCmds, cmd string, lerr string, hard bool) {
	// Append command to existing list
	o.cmds = append(o.cmds, cmd)
	o.errmsg = append(o.cmds, lerr)
	o.hard = append(o.hard, hard)
}

// TODO Pull out OS figuring out into it's own package from godojo
// add the methods to discovery.go

func main() {
	fmt.Println("Comanndeer!")

	// Setup a command
	var cmd Command
	// TODO Once a library, make this cmd.Create()
	cmd = Create()
	interfaceFun(cmd)

	// Add some commands to the osCmds struct
	fmt.Println("Setting up command data")
	rando := newCmds("Random Test")
	addCmd(rando, "ls", "ls failed", false)
	addCmd(rando, "uptime", "uptime failed", false)
	addCmd(rando, "pstree", "pstree failed", false)

	// Create a command package
	var testPkg cmdPkg
	testPkg.ID = "Ubuntu:21.04"
	testPkg.Distro = "Ubuntu"
	testPkg.Release = "21.04"
	testPkg.OS = "Linux"
	testPkg.Label = "POC"
	getCmdPackage(&testPkg)

	fmt.Println("Test package looks like:")
	fmt.Printf("%+v\n\n", testPkg)
	for k, v := range(*testPkg.PkgCmds) {
		fmt.Printf("Command #%v is %v\n", k, v)
	}

	//listCmd := cmdList {
	one := []singleCmd{
		{
			// List the files
			"ls",
			"ls failed",
			false,
		},
		{
			// Check uptime
			"uptime",
			"uptime failed",
			false,
		},
		{
			// Check uptime
			"pstree",
			"pstree failed",
			false,
		},
	}
	// TODO create a load commands method to take ^ and create a filled in osCmds struct - return *osCmds and an error

	fmt.Printf("Blah: %v\n\n", one)

	// Run the commands
	fmt.Println("Running command data")
	// Run without a timeout
	ctx := context.Background()
	cmd.TryCmds(ctx, rando)

	// Setup a context for long running command cancellation
	// Below is a silly example of a 2 minute timout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Test the change shell method
	fmt.Println("Testing setting the shell")
	cmd.SetShell("sh")
	out, err := cmd.Exec(ctx, "ls")
	if err != nil {
		log.Fatalf("Error exec'ing command was:\n\t%v\n", err)
	}
	fmt.Println(string(out))

}
