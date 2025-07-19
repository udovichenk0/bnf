package main

import (
	"fmt"
	"strconv"
	"strings"
)

const defaultMaxRepetitions = 10

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

type DigitExpr struct {
	Text      []rune
	TokenType TokenType
	Loc       int
}

type RepetitionExpr struct {
	Min  int
	Max  int
	Expr Expr
}

type OptionalExpr struct {
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

	for !p.IsAtEnd() && p.Peek().OfType(Choice) {
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
	case Digit:
		p.Next()
		min, err := strconv.Atoi(string(token.Text))
		if err != nil {
			return nil, err
		}

		asterisk := p.Next()
		if asterisk.OfType(Asterisk) {
			max := defaultMaxRepetitions
			IsDefaultMax := true
			digit := p.Peek()
			if digit.OfType(Digit) {
				max, err = strconv.Atoi(string(digit.Text))
				if err != nil {
					return nil, err
				}
				IsDefaultMax = false
				p.Next()
			}
			expr, err := p.ParsePrimaryExpr()
			if err != nil {
				return nil, err
			}
			if IsDefaultMax {
				max = min + max
			}

			if max < min {
				return nil, fmt.Errorf("repetition expression error: max repeat count (%d) must be greater than min repeat count (%d):%s", max, min, p.PointToLoc(token))
			}
			return RepetitionExpr{
				Min:  min,
				Max:  max,
				Expr: expr,
			}, nil
		}

		return nil, nil
	case Asterisk:
		p.Next()
		digit := p.Peek()
		if digit.OfType(Digit) {
			p.Next()
			expr, err := p.ParsePrimaryExpr()
			if err != nil {
				return nil, err
			}
			max, err := strconv.Atoi(string(digit.Text))
			if err != nil {
				return nil, err
			}
			return RepetitionExpr{
				Expr: expr,
				Max:  max,
			}, nil
		}
		expr, err := p.ParsePrimaryExpr()
		if err != nil {
			return nil, err
		}
		return RepetitionExpr{
			Expr: expr,
			Max:  defaultMaxRepetitions,
		}, nil
	case OpenSquareBrace:
		p.Next()
		expr, err := p.Parse()
		if err != nil {
			return nil, err
		}
		err = p.Expect(CloseSquareBrace)
		if err != nil {
			return nil, err
		}
		return OptionalExpr{
			Expr: expr,
		}, nil
	case OpenCurlyBrace:
		p.Next()
		expr, err := p.Parse()
		if err != nil {
			return nil, err
		}
		err = p.Expect(CloseCurlyBrace)
		if err != nil {
			return nil, err
		}
		return RepetitionExpr{
			Max:  defaultMaxRepetitions,
			Expr: expr,
		}, nil
	default:
		return nil, fmt.Errorf("expect expression")
	}
}

func (p *Parser) CanStartExpr() bool {
	switch p.Peek().TokenType {
	case String, OpenParen, NonTerminalSym, OpenCurlyBrace, Asterisk, Digit, OpenSquareBrace:
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
	if token.OfType(expected) {
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

func (p *Parser) PointToLoc(token Token) string {
	count := int(len(token.Text)/2) + token.Loc
	return fmt.Sprintf("\n%s\n%s^\n", string(p.line), strings.Repeat("-", count))
}
