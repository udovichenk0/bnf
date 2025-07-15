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
		if s.IsSymbol(c) {
			s.GetLiteralToken()
			continue
		} else {
			l := fmt.Sprintf("unknown symbol: %s in line:", string(c))
			return nil, fmt.Errorf("%s \n%s\n%s^", l, string(s.line), strings.Repeat("-", s.current))
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
		log.Fatalf("unknown literal token '%s' in line: \n%s\n%s^", string(lit), string(s.line), strings.Repeat("-", s.start))
		return
	}
	s.tokens = append(s.tokens, Token{
		Text:      lit,
		TokenType: e,
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
