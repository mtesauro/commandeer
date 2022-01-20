package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func interfaceFun(c Command) {
	fmt.Println("In interfaceFun")
	fmt.Println(c.PrintShell())
}

func main() {
	fmt.Println("Comanndeer!")

	//Setup a context for long running command cancellation
	//Below is a silly example of a 2 minute timout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var cmd Command
	cmd = Create()
	interfaceFun(cmd)
	cmd.SetShell("sh")
	out, err := cmd.Exec(ctx, "ls")
	if err != nil {
		log.Fatalf("Error exec'ing command was:\n\t%v\n", err)
	}
	fmt.Println(string(out))

}
