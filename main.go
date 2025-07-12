package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.ReadFile("./b.bnf")
	if err != nil {
		fmt.Println(err)
	}
	b := bytes.Split(file, []byte("\n"))
	grammar := NewGrammar("full-name")
	for _, line := range b {
		scanner := NewScanner(line)
		tokens, err := scanner.Scan()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		if tokens != nil {
			head := tokens[0]
			body := tokens[2:]
			parser := NewParser(body)
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
	const startSymbol = "full-name"
	startExpr, err := grammar.GetExpressionFromGrammar(startSymbol)
	if err != nil {
		log.Fatal(err)
	}
	result, err := grammar.Generate(startExpr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
}
