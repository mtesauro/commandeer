package main

import (
	"fmt"
	"os"
	"time"

	c "github.com/mtesauro/commandeer/pkg"
)

func main() {
	fmt.Println("Comanndeer!")

	// Create a command pkg
	testPkg := c.NewPkg("POC")

	// Set the location to run the commands - in the local terminal
	// Default is LocalTerm so the command below isn't required
	testPkg.SetLocation(c.LocalTerm{})

	// Set some targets
	testPkg.AddTarget("Ubuntu:21.04", "Ubuntu", "21.04", "Linux", "bash")
	testPkg.AddTarget("Ubuntu:20.04", "Ubuntu", "20.04", "Linux", "bash")
	testPkg.AddTarget("CentOS:7", "CentOS", "7", "Linux", "bash")
	testPkg.AddTarget("RHEL:8", "RHEL", "8", "Linux", "bash")

	// Add a single command
	testPkg.AddCmd("ls ./", "ls command failed", false, 0, "Ubuntu:21.04")

	// Load multiple commands
	cmdList := []c.SingleCmd{
		{ // List the files
			Cmd:     "free -m",
			Errmsg:  "free failed",
			Hard:    false,
			Timeout: (3 * time.Second),
		},
		{ // Check uptime
			Cmd:     "uptime",
			Errmsg:  "uptime failed",
			Hard:    false,
			Timeout: 0,
		},
		{ // Check uptime
			Cmd:     "pstree | head",
			Errmsg:  "pstree failed",
			Hard:    false,
			Timeout: 0,
		},
	}
	err := testPkg.LoadCmds(cmdList, "Ubuntu:21.04")
	if err != nil {
		fmt.Println("Unable to load commands to the target")
		fmt.Printf("Error was:\n%v\n", err)
		os.Exit(1)
	}

	// Exec the commands
	out, err := testPkg.ExecPkgCombined("Ubuntu:21.04")
	if err != nil {
		fmt.Println("An error has occured:")
		fmt.Printf("\t%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Output:\n%s\n", out)

}
