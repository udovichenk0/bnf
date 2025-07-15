package main

import (
	"fmt"
)

type Parser struct {
	pos    int
	tokens []Token
	line   []rune
}

type Expr interface{}

type SequenceExpr []Expr

type ChoiceExpr []Expr

type NonTerminalExpr struct {
	Text      []rune
	TokenType TokenType
	Loc       int
}

type StringExpr struct {
	Text      []rune
	TokenType TokenType
	Loc       int
}

type RepetitionExpr struct {
	Min  int
	Max  int
	Expr Expr
}

func NewParser(tokens []Token, line []rune) Parser {
	return Parser{tokens: tokens, line: line}
}

func (p *Parser) ParseSequenceExpr() (SequenceExpr, error) {
	primary, err := p.ParsePrimaryExpr()

	if err != nil {
		return nil, err
	}

	sequence := SequenceExpr{
		primary,
	}

	for !p.IsAtEnd() && p.CanStartExpr() {
		primary, err := p.ParsePrimaryExpr()
		if err != nil {
			return nil, err
		}
		sequence = append(sequence, primary)
	}
	return sequence, nil
}

func (p *Parser) ParseChoiceExpr() (Expr, error) {
	sequenceExpr, err := p.ParseSequenceExpr()
	if err != nil {
		return nil, err
	}

	if len(sequenceExpr) == 0 {
		return nil, nil
	}

	var expr ChoiceExpr
	expr = append(expr, sequenceExpr)

	for !p.IsAtEnd() && p.Peek().TokenType == Choice {
		if err := p.Expect(Choice); err != nil {
			return nil, err
		}
		sequenceExpr, err := p.ParseSequenceExpr()
		if err != nil {
			return nil, err
		}
		expr = append(expr, sequenceExpr)
	}
	return expr, nil
}

func (p *Parser) ParsePrimaryExpr() (Expr, error) {
	token := p.Peek()
	switch token.TokenType {
	case NonTerminalSym:
		p.Next()
		return NonTerminalExpr(token), nil
	case String:
		p.Next()
		return StringExpr(token), nil
	case OpenParen:
		p.Next()
		expr, _ := p.Parse()
		err := p.Expect(CloseParen)
		if err != nil {
			return nil, err
		}
		return expr, nil
	case Asterisk:
		p.Next()
		expr, err := p.ParsePrimaryExpr()
		if err != nil {
			return nil, err
		}
		return RepetitionExpr{
			Expr: expr,
			Min:  0,
			Max:  10,
		}, nil
	case OpenCurlyBrace:
		p.Next()
		expr, _ := p.Parse()
		err := p.Expect(CloseCurlyBrace)
		if err != nil {
			return nil, err
		}
		return RepetitionExpr{
			Min:  0,
			Max:  10,
			Expr: expr,
		}, nil
	default:
		return nil, fmt.Errorf("expect expression")
	}
}

func (p *Parser) CanStartExpr() bool {
	switch p.Peek().TokenType {
	case String, OpenParen, NonTerminalSym, OpenCurlyBrace, Asterisk:
		return true
	default:
		return false
	}
}

func (p *Parser) IsAtEnd() bool {
	return p.pos >= len(p.tokens)
}

func (p *Parser) Peek() Token {
	if p.IsAtEnd() {
		return Token{}
	}
	return p.tokens[p.pos]
}

func (p *Parser) Next() Token {
	if p.IsAtEnd() {
		return Token{}
	}
	p.pos++
	return p.tokens[p.pos-1]
}

func (p *Parser) Expect(expected TokenType) error {
	token := p.Peek()
	if token.TokenType == expected {
		p.Next()
		return nil
	}
	tokenStr, err := TokenToString(expected)
	if err != nil {
		return err
	}
	return fmt.Errorf("expected token: %s, got: %s", tokenStr, string(token.Text))
}

func (p *Parser) Parse() (Expr, error) {
	expr, err := p.ParseChoiceExpr()
	if err != nil {
		return nil, err
	}
	return expr, nil
}
