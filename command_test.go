package commandeer_test

import (
	"errors"
	"strings"
	"testing"

	c "github.com/mtesauro/commandeer"
)

func TestPoorlyLocalCmdLs(t *testing.T) {
	// Setup a Command Package with a target
	testPkg := c.NewPkg("Test local command using ls")
	testPkg.AddTarget("Ubuntu:21.04", "Ubuntu", "21.04", "Linux", "bash")

	// Add the ls command
	err := testPkg.AddCmd("ls ./", "ls command failed", false, 0, "Ubuntu:21.04")
	if err != nil {
		t.Fatal(err)
	}

	// Execute the command package
	out, err := testPkg.ExecPkgCombined("Ubuntu:21.04")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(out), "command_test.go") {
		t.Errorf("expected output to contain a current go file, got:\n%s", string(out))
	}
}

func TestExecPkgCombinedSuccess(t *testing.T) {
	// Setup a mock terminal
	m := c.MockTerm{
		Out:         []byte("command.go  command_test.go  discovery.go  localTerm.go  mockTerm.go  targets.go  terminal.go"),
		Err:         nil,
		LastCommand: "",
	}

	// Setup a Command Package with a target
	testPkg := c.NewPkg("Test local command using ls")

	// Set Terminal to mock terminal
	testPkg.SetLocation(&m)

	// Create a target for the commands
	testPkg.AddTarget("Ubuntu:21.04", "Ubuntu", "21.04", "Linux", "bash")

	// Add at least 1 command even though Terminal is mocked
	err := testPkg.AddCmd("ls ./", "ls command failed", false, 0, "Ubuntu:21.04")
	if err != nil {
		t.Fatal(err)
	}

	// Execute the command package
	out, err := testPkg.ExecPkgCombined("Ubuntu:21.04")
	if err != nil {
		t.Fatal(err)
	}

	// Output should match the output we mocked
	if !strings.Contains(string(out), string(m.Out)) {
		t.Errorf("Expected output to match what was mocked (expected), got:\n%s\nexpected:\n%s", string(out), string(m.Out))
	}
}

// TODO Decide if this is an actually useful test
func TestExecPkgCombinedFailed(t *testing.T) {
	// Setup a mock terminal
	m := c.MockTerm{
		Out:         []byte("command.go  command_test.go  discovery.go  localTerm.go  mockTerm.go  targets.go  terminal.go"),
		Err:         nil,
		LastCommand: "",
	}

	// Setup a Command Package with a target
	testPkg := c.NewPkg("Test local command using ls")

	// Set Terminal to mock terminal
	testPkg.SetLocation(&m)

	// Create a target for the commands
	testPkg.AddTarget("Ubuntu:21.04", "Ubuntu", "21.04", "Linux", "bash")

	// Add at least 1 command even though Terminal is mocked
	err := testPkg.AddCmd("ls ./", "ls command failed", false, 0, "Ubuntu:21.04")
	if err != nil {
		t.Fatal(err)
	}

	// Execute the command package
	_, err = testPkg.ExecPkgCombined("Ubuntu:21.04")
	if err != nil {
		t.Fatal(err)
	}

	// TEST: Force output from ExecPkgCombined to be different than what was mocked above
	// Force output to be different than mocked output
	out := []byte("DIFFERENT")

	// Output should NOT match the output we mocked for this test of failure
	if strings.Contains(string(out), string(m.Out)) {
		t.Errorf("Expected output to be different than what was mocked (expected), got:\n%s\nexpected:\n%s", string(out), string(m.Out))
	}
}

func TestExecPkgCombinedErrored(t *testing.T) {
	// Setup a mock terminal
	m := c.MockTerm{
		Out:         []byte("command.go  command_test.go  discovery.go  localTerm.go  mockTerm.go  targets.go  terminal.go"),
		Err:         errors.New("Forcing an error for testing"),
		LastCommand: "",
	}

	// Setup a Command Package with a target
	testPkg := c.NewPkg("Test local command using ls")

	// Set Terminal to mock terminal
	testPkg.SetLocation(&m)

	// Create a target for the commands
	testPkg.AddTarget("Ubuntu:21.04", "Ubuntu", "21.04", "Linux", "bash")

	// Add at least 1 command even though Terminal is mocked
	err := testPkg.AddCmd("ls ./", "ls command failed", false, 0, "Ubuntu:21.04")
	if err != nil {
		t.Fatal(err)
	}

	// Execute the command package
	_, err = testPkg.ExecPkgCombined("Ubuntu:21.04")

	// TEST: Mock forces an error to be returned from ExecPkgCombined so error should never be nil
	if err == nil {
		t.Fatal(errors.New("Error was mocked but no error returned from ExecPkgCombined()"))
	}

}

func TestExecPkgCombinedBadTimeout(t *testing.T) {
	// Setup a mock terminal
	m := c.MockTerm{
		Out:         []byte("command.go  command_test.go  discovery.go  localTerm.go  mockTerm.go  targets.go  terminal.go"),
		Err:         nil,
		LastCommand: "",
	}

	// Setup a Command Package with a target
	testPkg := c.NewPkg("Test local command using ls")

	// Set Terminal to mock terminal
	testPkg.SetLocation(&m)

	// Create a target for the commands
	testPkg.AddTarget("Ubuntu:21.04", "Ubuntu", "21.04", "Linux", "bash")

	// Add at least 1 command even though Terminal is mocked
	err := testPkg.AddCmd("ls ./", "ls command failed", false, -7, "Ubuntu:21.04")
	if err != nil {
		t.Fatal(err)
	}

	// Execute the command package
	_, err = testPkg.ExecPkgCombined("Ubuntu:21.04")

	// TEST: Command added with negative duration so error should never be nil
	if err == nil {
		t.Fatal(errors.New("Error was mocked but no error returned from ExecPkgCombined()"))
	}

}
