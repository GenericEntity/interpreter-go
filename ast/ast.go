package ast

import "github.com/GenericEntity/interpreter-go/monkey/token"

type Node interface {
	TokenLiteral() string // used for debugging and testing
}

type Statement interface {
	Node
	// dummy method so a struct can implement Statement without necessarily implementing Expression
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

// An Identifier is an expression because for simplicity we do not distinguish identifiers
// used in the LHS of an assignment from identifiers used in expressions.
// We could create two identifier types to distinguish their usages.
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
