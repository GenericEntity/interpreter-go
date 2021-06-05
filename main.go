package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/GenericEntity/interpreter-go/monkey/repl"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is a generic REPL for the Monkey Programming Language!\n", u.Username)
	fmt.Printf("Please enjoy the infinite loop of interpretation! (CTRL-C to terminate)\n")
	repl.Start(os.Stdin, os.Stdout)
}
