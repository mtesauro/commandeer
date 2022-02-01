package commandeer

import (
	"fmt"
	"os"
	"strings"
)

type cmdPkg struct {
	Label   string       // Holds a human friendly labe for the command package
	ID      string       // Holds the Distro ID e.g. Ubuntu:21.10
	Release string       // Holds the Release e.g. 21.10
	Distro  string       // Holds the Distro name e.g. Ubuntu
	OS      string       // Holds the operating system e.g. Linux
	PkgCmds *[]singleCmd // Holds the commands in this command package
}

// TODO Update to return an error instead of os.Exit()
func getCmdPackage(cPkg *cmdPkg) {
	// Switch on OS
	switch cPkg.OS {
	case "linux":
		linuxPkgs(cPkg)
	case "darwin":
		fmt.Println("OS X/Darwin is not a supported installation platform")
		os.Exit(1)
	case "windows":
		fmt.Println("Windows is not a supported installation platform")
		os.Exit(1)
	}

	pCmds := []singleCmd{
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

	cPkg.PkgCmds = &pCmds

}

// switch on Distro
// Update to return an error
func linuxPkgs(cPkg *cmdPkg) {
	switch strings.ToLower(cPkg.Distro) {
	case "debian":
		fallthrough
	case "ubuntu":
		//ubuntuPkgs(cPkg)
		// Updated with errors will be like below
		// return ubuntu(cPkgs)
		return
	default:
		fmt.Println("Unsupported OS, quitting.")
		os.Exit(1)

	}
}

// switch on Release for Distro
// Update to return an error

// Need a function that lets you set a function to call for a specific package label
// AKA when using as an external package, you write a function for the command package and the use this method to link the package name to your function.
// godojo defines a label "prep-install" and a function "ubuntuPrepCmds()" so that when commandeer finds "prep-install", it calls the matching function.
// Will also need a way to associate that function with supported OS/Distro/release as well
