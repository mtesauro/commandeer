package main

import (
	"fmt"

	c "github.com/mtesauro/commandeer/pkg"
)

// USE THIS TO TEST OUT HOW THE LIBRARY WILL BE WITH godojo

// TODO Update to return an error instead of os.Exit()
func GetCmdPackage(cPkg *c.CmdPkg) {
	// Switch on OS
	// FIXME for new command package structure
	//switch cPkg.OS {
	//case "linux":
	//	linuxPkgs(cPkg)
	//case "darwin":
	//	fmt.Println("OS X/Darwin is not a supported installation platform")
	//	os.Exit(1)
	//case "windows":
	//	fmt.Println("Windows is not a supported installation platform")
	//	os.Exit(1)
	//}

	pCmds := []c.SingleCmd{
		{
			// List the files
			Cmd:    "ls",
			Errmsg: "ls failed",
			Hard:   false,
		},
		{
			// Check uptime
			Cmd:    "uptime",
			Errmsg: "uptime failed",
			Hard:   false,
		},
		{
			// Check uptime
			Cmd:    "pstree",
			Errmsg: "pstree failed",
			Hard:   false,
		},
	}

	fmt.Println(pCmds)

	// FIXME
	//cPkg.PkgCmds = pCmds

}

// switch on Distro
// Update to return an error
func linuxPkgs(cPkg *c.CmdPkg) {
	// FIXME
	//switch strings.ToLower(cPkg.Distro) {
	//case "debian":
	//	fallthrough
	//case "ubuntu":
	//	//ubuntuPkgs(cPkg)
	//	// Updated with errors will be like below
	//	// return ubuntu(cPkgs)
	//	return
	//default:
	//	fmt.Println("Unsupported OS, quitting.")
	//	os.Exit(1)

	//}
}

// switch on Release for Distro
// Update to return an error

// Need a function that lets you set a function to call for a specific package label
// AKA when using as an external package, you write a function for the command package and the use this method to link the package name to your function.
// godojo defines a label "prep-install" and a function "ubuntuPrepCmds()" so that when commandeer finds "prep-install", it calls the matching function.
// Will also need a way to associate that function with supported OS/Distro/release as well
