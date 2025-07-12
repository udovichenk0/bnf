package main

import (
	"fmt"
	"log"
)

type Scanner struct {
	current int
	start   int
	line    []byte
	tokens  []Token
}

type Token struct {
	Text      string
	TokenType TokenType
}

type TokenType int

const (
	NonTerminalSym TokenType = iota + 1
	StringSym
	EqualSym
	ChoiceSym
)

func NewScanner(line []byte) Scanner {
	return Scanner{line: line}
}

var LiteralTokens = map[string]TokenType{
	"|":   ChoiceSym,
	"::=": EqualSym,
}

func (s Scanner) Scan() ([]Token, error) {
	s.current = 0
	s.start = 0
	s.tokens = nil

	if len(s.line) == 0 {
		return nil, nil
	}
	for !s.IsEnd() {
		c := s.Peek()
		switch c {
		case '<':
			s.GetVariableToken()
			continue
		case '"':
			s.GetStringToken()
			continue
		case ' ':
			s.GetChar()
			continue
		}
		if s.IsSymbol(c) {
			s.GetLiteralToken()
			continue
		} else {
			return nil, fmt.Errorf("unknown symbol: %s", string(c))
		}
	}
	return s.tokens, nil
}

func (s *Scanner) GetVariableToken() {
	s.current++
	s.start = s.current
	for !s.IsEnd() && !s.Match(('>')) {
		s.GetChar()
	}
	if s.IsEnd() {
		log.Fatalf("Error: Expected '%c' symbol at the end of the line: '%s'", '>', s.line)
		return
	}

	token := Token{
		Text:      string(s.line[s.start:s.current]),
		TokenType: NonTerminalSym,
	}
	s.current++
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) GetStringToken() {
	s.current++
	s.start = s.current
	for !s.IsEnd() && !s.Match(('"')) {
		s.GetChar()
	}
	if s.IsEnd() {
		log.Fatalf("Error: Expected '%c' symbol at the end of the line: '%s'", '"', s.line)
		return
	}

	token := Token{
		Text:      string(s.line[s.start:s.current]),
		TokenType: StringSym,
	}
	s.current++
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) GetLiteralToken() {
	s.start = s.current
	var lit []byte
	for !s.IsEnd() && s.IsSymbol(s.Peek()) {
		lit = append(lit, s.GetChar())
	}

	e, ok := LiteralTokens[string(lit)]
	if !ok {
		log.Fatalf("Error: Unknown literal token '%s'.", lit)
		return
	}
	s.tokens = append(s.tokens, Token{
		Text:      string(lit),
		TokenType: e,
	})
}

func (s *Scanner) Peek() byte {
	if s.IsEnd() {
		return 0
	}
	return s.line[s.current]
}

func (s *Scanner) IsSymbol(b byte) bool {
	var isSymbol bool
	switch b {
	case '|', ':', '=', '<', '>', '-', '_':
		isSymbol = true
	}
	return isSymbol
}

func (s *Scanner) IsEnd() bool {
	return s.current >= len(s.line) || s.line[s.current] == '\r' || s.line[s.current] == '\n'
}

func (s *Scanner) Match(b byte) bool {
	return s.line[s.current] == b
}

func (s *Scanner) GetChar() byte {
	s.current++
	return s.line[s.current-1]
}
