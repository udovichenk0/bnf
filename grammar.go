package main

import (
	"fmt"
	"math/rand"
)

type Rule struct {
	Head string
	Body Expr
}

type Rules map[string]Expr

type Grammar struct {
	rules       Rules
	startSymbol string
}

func NewGrammar(startSymbol string) Grammar {
	return Grammar{
		rules:       make(Rules),
		startSymbol: startSymbol,
	}
}

func (g *Grammar) Generate(expr Expr) (string, error) {
	switch expr := expr.(type) {
	case EqualExpr:
		chosenAlternative := rand.Intn(len(expr))
		sequece := expr[chosenAlternative]
		return g.Generate(sequece)
	case SequenceExpr:
		var res string
		for i := 0; i < len(expr); i++ {
			r, err := g.Generate(expr[i])
			if err != nil {
				return "", err
			}
			res += r
		}
		return res, nil
	case NonTerminalExpr:
		exp, err := g.GetExpressionFromGrammar(expr.Text)
		if err != nil {
			return "", err
		}
		str, err := g.Generate(exp)
		if err != nil {
			return "", err
		}
		return str, nil

	case StringExpr:
		return expr.Text, nil
	default:
		return "", fmt.Errorf("unknown Expression")
	}
}

func (g *Grammar) AddRule(rule Rule) {
	g.rules[rule.Head] = rule.Body
}
func (g *Grammar) GetExpressionFromGrammar(head string) (Expr, error) {

	expr, ok := g.rules[head]
	if !ok {
		return nil, fmt.Errorf("unknown symbol: %s", head)
	}
	return expr, nil
}
