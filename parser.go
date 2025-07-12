package main

import (
	"fmt"
)

type Parser struct {
	pos    int
	tokens []Token
	expr   Expr
	line   []rune
}

type Expr interface{}

type SequenceExpr []Expr

type EqualExpr []Expr

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

func NewParser(tokens []Token, line []rune) Parser {
	return Parser{tokens: tokens, line: line}
}

func (p *Parser) ParseSequenceExpr() SequenceExpr {
	primary := p.ParsePrimaryExpr()
	if primary == nil {
		return SequenceExpr{}
	}

	sequence := SequenceExpr{
		primary,
	}
	choiceBar := p.Peek()
	if choiceBar.TokenType == ChoiceSym {
		return sequence
	}
	for primary != nil && !p.IsAtEnd() {
		choiceBar := p.Peek()
		if choiceBar.TokenType != ChoiceSym {
			primary = p.ParsePrimaryExpr()
			if primary != nil {
				sequence = append(sequence, primary)
			}
		} else {
			break
		}
	}
	return sequence
}

func (p *Parser) ParseEqualExpr() error {
	sequenceExpr := p.ParseSequenceExpr()
	if len(sequenceExpr) == 0 {
		return fmt.Errorf("syntax error: %s", string(p.line))
	}

	equalExpr, isEqualExpr := p.expr.(EqualExpr)
	if isEqualExpr {
		p.expr = append(equalExpr, sequenceExpr)
	}

	if !p.IsAtEnd() {
		err := p.Expect(ChoiceSym)
		if err != nil {
			return err
		}
		if !isEqualExpr {
			p.expr = EqualExpr{sequenceExpr}
		}
		err = p.ParseEqualExpr()
		if err != nil {
			return err
		}
	} else {
		if !isEqualExpr {
			p.expr = sequenceExpr
		}
	}
	return nil
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
	token := p.Next()
	if token.TokenType == expected {
		return nil
	}
	tokenStr, err := TokenToString(expected)
	if err != nil {
		return err
	}
	return fmt.Errorf("expected token: %s, got: %s", tokenStr, string(token.Text))
}

func (p *Parser) ParsePrimaryExpr() Expr {
	token := p.Next()
	switch token.TokenType {
	case NonTerminalSym:
		return NonTerminalExpr(token)
	case StringSym:
		return StringExpr(token)
	default:
		return nil
	}
}

func (p *Parser) IsAtEnd() bool {
	return p.pos >= len(p.tokens)
}

func (p *Parser) Parse() (Expr, error) {
	err := p.ParseEqualExpr()
	if err != nil {
		return nil, err
	}
	return p.expr, nil
}
