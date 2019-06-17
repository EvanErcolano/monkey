package ast

import "monkey/token"

// Every valid Monkey program is a sequence of statements contained
// within Program.Statements, a slice of AST nodes that impl Statement
// interface

// Program is the root node of our AST and impls the Node interface
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

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type LetStatement struct {
	Token token.Token // the token.Let token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type ReturnStatement struct {
	Token token.Token // the token.RETURN token
	Value Expression
}

func (rs *ReturnStatement) statementNode()      {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
