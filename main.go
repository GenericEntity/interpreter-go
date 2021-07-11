package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/GenericEntity/interpreter-go/monkey/interpreter"
	"github.com/GenericEntity/interpreter-go/monkey/repl"
)

var (
	flagScriptFile = flag.String("f", "", "path to file to interpret. if blank, opens a REPL")
)

func main() {
	flag.Parse()

	switch strings.TrimSpace(*flagScriptFile) {
	case "":
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		repl.Greet(os.Stdout, u.Username)
		repl.Start(os.Stdin, os.Stdout)

	default:
		contents, err := ioutil.ReadFile(*flagScriptFile)
		if err != nil {
			fmt.Printf("Error when reading file. %v", err)
			return
		}
		interpreter.Interpret(string(contents), os.Stdout)
	}
}
