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

	// Setup logging
	logLoc, err := c.LogToFile("./logs", "example.log")
	if err != nil {
		fmt.Printf("Error was:\n%v\n", err)
		os.Exit(1)
	}
	testPkg.SetLogging(logLoc)
	// Setup command logging
	logCmd, err := c.LogToFile("./logs", "commands.log")
	if err != nil {
		fmt.Printf("Error was:\n%v\n", err)
		os.Exit(1)
	}
	testPkg.SetCmdLog(logCmd)
	testPkg.TurnOnCmdLog()

	// Set the location to run the commands - in the local terminal
	// Default is LocalTerm so the command below isn't required
	testPkg.SetLocation(c.LocalTerm{})

	// Set some targets
	testPkg.AddTarget("Ubuntu:21.04", "Ubuntu", "21.04", "Linux", "bash")
	testPkg.AddTarget("Ubuntu:20.04", "Ubuntu", "20.04", "Linux", "bash")
	testPkg.AddTarget("CentOS:7", "CentOS", "7", "Linux", "bash")
	testPkg.AddTarget("RHEL:8", "RHEL", "8", "Linux", "bash")

	// Write some log entries
	testPkg.LogTrace("All Targets have been created - trace")
	testPkg.LogInfo("All Targets have been created - info")
	testPkg.LogWarn("All Targets have been created - warn")
	testPkg.LogError("All Targets have been created - error")

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
	err = testPkg.LoadCmds(cmdList, "Ubuntu:21.04")
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
