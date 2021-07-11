package interpreter

import (
	"io"

	"github.com/GenericEntity/interpreter-go/monkey/evaluator"
	"github.com/GenericEntity/interpreter-go/monkey/lexer"
	"github.com/GenericEntity/interpreter-go/monkey/object"
	"github.com/GenericEntity/interpreter-go/monkey/parser"
)

func Interpret(code string, out io.Writer) {
	env := object.NewEnvironment()
	lex := lexer.New(code)
	p := parser.New(lex)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
