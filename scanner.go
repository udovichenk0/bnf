package main

import (
	"fmt"
	"log"
	"strings"
)

type Scanner struct {
	current int
	start   int
	line    []rune
	tokens  []Token
}

type Token struct {
	Text      []rune
	TokenType TokenType
	Loc       int
}

type TokenType int

const (
	NonTerminalSym TokenType = iota + 1
	String
	Digit
	Equal
	Choice
	OpenParen
	CloseParen
	OpenCurlyBrace
	CloseCurlyBrace
	Asterisk
)

func NewScanner(line []rune) Scanner {
	return Scanner{line: line}
}

var LiteralTokens = map[string]TokenType{
	"|":   Choice,
	"/":   Choice,
	"::=": Equal,
	":=":  Equal,
	"=":   Equal,
	"(":   OpenParen,
	")":   CloseParen,
	"{":   OpenCurlyBrace,
	"}":   CloseCurlyBrace,
	"*":   Asterisk,
}

func TokenToString(tokenType TokenType) (string, error) {
	for lit, tt := range LiteralTokens {
		if tokenType == tt {
			return lit, nil
		}
	}
	return "", fmt.Errorf("unknown tokenType: %d", tokenType)
}

func (s Scanner) Scan() ([]Token, error) {
	s.current = 0
	s.start = 0
	s.tokens = nil
	if len(s.line) == 0 {
		return nil, nil
	}

	for !s.IsAtEnd() {
		c := s.Peek()
		switch c {
		case '<':
			s.GetVariableToken()
			continue
		case '"':
			s.GetStringToken()
			continue
		case ' ':
			s.Next()
			continue
		}
		if s.IsDigit(c) {
			s.GetDigitToken()
			continue
		}
		if s.IsSymbol(c) {
			s.GetLiteralToken()
			continue
		} else {
			s.start = s.current
			return nil, fmt.Errorf("unknown symbol: %s in line:%s", string(c), s.PointToLoc())
		}
	}
	return s.tokens, nil
}

func (s *Scanner) GetVariableToken() {
	s.current++
	s.start = s.current
	for !s.IsAtEnd() && !s.Match(('>')) {
		s.Next()
	}
	if s.IsAtEnd() {
		log.Fatalf("Error: Expected '%c' symbol at the end of the line: '%s'", '>', string(s.line))
		return
	}

	token := Token{
		Text:      s.line[s.start:s.current],
		TokenType: NonTerminalSym,
		Loc:       s.start,
	}
	s.current++
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) GetStringToken() {
	s.current++
	s.start = s.current
	for !s.IsAtEnd() && !s.Match(('"')) {
		s.Next()
	}
	if s.IsAtEnd() {
		log.Fatalf("Error: Expected '%c' symbol at the end of the line: '%s'", '"', string(s.line))
		return
	}

	token := Token{
		Text:      s.line[s.start:s.current],
		TokenType: String,
		Loc:       s.start,
	}
	s.current++
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) GetLiteralToken() {
	s.start = s.current
	var lit []rune
	for !s.IsAtEnd() && s.IsSymbol(s.Peek()) {
		lit = append(lit, s.Next())
	}

	e, ok := LiteralTokens[string(lit)]
	if !ok {
		log.Fatalf("unknown literal token '%s' in line:%s", string(lit), s.PointToLoc())
		return
	}
	s.tokens = append(s.tokens, Token{
		Text:      lit,
		TokenType: e,
		Loc:       s.start,
	})
}

func (s *Scanner) GetDigitToken() {
	s.start = s.current
	var lit []rune
	for !s.IsAtEnd() && s.IsDigit(s.Peek()) {
		lit = append(lit, s.Next())
	}

	s.tokens = append(s.tokens, Token{
		Text:      lit,
		TokenType: Digit,
		Loc:       s.start,
	})
}

func (s *Scanner) Peek() rune {
	if s.IsAtEnd() {
		return '\x00'
	}
	return s.line[s.current]
}

func (s *Scanner) IsSymbol(b rune) bool {
	switch b {
	case '|', '/', ':', '=', '(', ')', '{', '}', '*':
		return true
	}
	return false
}

func (s *Scanner) IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) IsAtEnd() bool {
	return s.current >= len(s.line) || s.line[s.current] == '\r' || s.line[s.current] == '\n'
}

func (s *Scanner) Match(b rune) bool {
	return s.line[s.current] == b
}

func (s *Scanner) Next() rune {
	if s.IsAtEnd() {
		return '\x00'
	}
	s.current++
	return s.line[s.current-1]
}

func (s *Scanner) PointToLoc() string {
	count := int((s.current-s.start)/2) + s.start
	return fmt.Sprintf("\n%s\n%s^\n", string(s.line), strings.Repeat("-", count))
}

func (t Token) OfType(tokType TokenType) bool {
	return t.TokenType == tokType
}
