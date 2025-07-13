package main

import (
	"fmt"
	"math/rand"
)

type Rule struct {
	Head []rune
	Body Expr
}

type Rules map[string]Expr

type Grammar struct {
	rules Rules
}

func NewGrammar() Grammar {
	return Grammar{
		rules: make(Rules),
	}
}

func (g *Grammar) Generate(expr Expr) ([]rune, error) {
	switch expr := expr.(type) {
	case EqualExpr:
		chosenAlternative := rand.Intn(len(expr))
		sequence := expr[chosenAlternative]
		return g.Generate(sequence)
	case SequenceExpr:
		var sequenceString []rune
		for i := 0; i < len(expr); i++ {
			r, err := g.Generate(expr[i])
			if err != nil {
				return nil, err
			}
			sequenceString = append(sequenceString, r...)
		}
		return sequenceString, nil
	case NonTerminalExpr:
		exp, err := g.GetExpressionFromGrammar(string(expr.Text))
		if err != nil {
			return nil, err
		}
		str, err := g.Generate(exp)
		if err != nil {
			return nil, err
		}
		return str, nil

	case StringExpr:
		return expr.Text, nil
	default:
		return nil, fmt.Errorf("unknown Expression")
	}
}

func (g *Grammar) AddRule(rule Rule) {
	g.rules[string(rule.Head)] = rule.Body
}
func (g *Grammar) GetExpressionFromGrammar(head string) (Expr, error) {
	expr, ok := g.rules[head]
	if !ok {
		return nil, fmt.Errorf("unknown symbol: %s", head)
	}
	return expr, nil
}
