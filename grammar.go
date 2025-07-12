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
	rules       Rules
	startSymbol []rune
}

func NewGrammar(startSymbol []rune) Grammar {
	return Grammar{
		rules:       make(Rules),
		startSymbol: startSymbol,
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
		exp, err := g.GetExpressionFromGrammar(expr.Text)
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
func (g *Grammar) GetExpressionFromGrammar(head []rune) (Expr, error) {
	expr, ok := g.rules[string(head)]
	if !ok {
		return nil, fmt.Errorf("unknown symbol: %s", string(head))
	}
	return expr, nil
}
