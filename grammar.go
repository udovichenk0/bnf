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
	case ChoiceExpr:
		chosenAlternative := rand.Intn(len(expr))
		sequence := expr[chosenAlternative]
		return g.Generate(sequence)
	case OptionalExpr:
		if rand.Intn(2) == 1 {
			return g.Generate(expr.Expr)
		}
		return nil, nil
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
	case DigitExpr:
		return expr.Text, nil
	case RepetitionExpr:
		var result []rune
		repetitionCount := rand.Intn(expr.Max-expr.Min+1) + expr.Min
		for i := 0; i < repetitionCount; i++ {
			r, err := g.Generate(expr.Expr)
			if err != nil {
				return nil, fmt.Errorf("unknown expression")
			}
			result = append(result, r...)
		}

		return result, nil
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
