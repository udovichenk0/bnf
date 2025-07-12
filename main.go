package main

import (
	"bytes"
	"log"
	"os"
)

func main() {
	var startSymbol = []rune("full-name")
	file, err := os.ReadFile("./b.bnf")
	if err != nil {
		log.Fatal(err)
	}
	b := bytes.Split(file, []byte("\n"))
	grammar := NewGrammar(startSymbol)
	for _, line := range b {

		scanner := NewScanner(bytes.Runes(line))
		tokens, err := scanner.Scan()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		if tokens != nil {
			head := tokens[0]
			body := tokens[2:]
			parser := NewParser(body, bytes.Runes(line))
			expr, err := parser.Parse()
			if err != nil {
				log.Fatalf("Error: %v", err)
				break
			}
			grammar.AddRule(Rule{
				Head: head.Text,
				Body: expr,
			})
		}
	}
	startExpr, err := grammar.GetExpressionFromGrammar(startSymbol)
	if err != nil {
		log.Fatal(err)
	}
	result, err := grammar.Generate(startExpr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(result))
}
