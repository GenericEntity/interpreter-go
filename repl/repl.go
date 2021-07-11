package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/GenericEntity/interpreter-go/monkey/evaluator"
	"github.com/GenericEntity/interpreter-go/monkey/lexer"
	"github.com/GenericEntity/interpreter-go/monkey/object"
	"github.com/GenericEntity/interpreter-go/monkey/parser"
)

const PROMPT = ">> "

func Greet(out io.Writer, name string) {
	fmt.Fprintf(out, "Hello %s! This is a generic REPL for the Monkey Programming Language!\n", name)
	fmt.Fprintf(out, "Please enjoy the infinite loop of interpretation! (CTRL-C to terminate)\n")
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		// Read
		fmt.Fprint(out, PROMPT)
		isScanned := scanner.Scan()
		if !isScanned {
			break
		}

		// Eval
		line := scanner.Text()
		lex := lexer.New(line)
		p := parser.New(lex)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Yikes! Watch what you type!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
