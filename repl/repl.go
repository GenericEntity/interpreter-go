package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/GenericEntity/interpreter-go/monkey/lexer"
	"github.com/GenericEntity/interpreter-go/monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Yikes! Watch what you type!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
