package main

import (
	"bytes"
	"log"
	"os"
)

const (
	startSymbolFlag string = "-s"
	fileFlag        string = "-f"
)

var flags = map[string]string{
	startSymbolFlag: startSymbolFlag,
	fileFlag:        fileFlag,
}

func FlagExist(flag string) bool {
	_, ok := flags[flag]
	return ok
}

func main() {
	args := os.Args[1:]
	var startSymbol string
	var fileName string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case startSymbolFlag:
			i++
			if i >= len(args) || FlagExist(args[i]) {
				log.Fatalf("%s flag should have a start symbol", startSymbolFlag)
			}

			startSymbol = args[i]
		case fileFlag:
			i++
			if i >= len(args) || FlagExist(args[i]) {
				log.Fatalf("%s flag should have a file name", fileFlag)
			}
			fileName = args[i]
		default:
			log.Fatalf("unnknown flag: %s", arg)
		}
	}

	if fileName == "" {
		log.Fatalf("%s flag should be specified", fileFlag)
	}

	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	b := bytes.Split(file, []byte("\r\n"))
	grammar := NewGrammar()
	for _, line := range b {

		scanner := NewScanner(bytes.Runes(line))
		tokens, err := scanner.Scan()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		if tokens != nil {
			head := tokens[0]
			body := tokens[2:]
			if startSymbol == "" && len(head.Text) > 0 {
				startSymbol = string(head.Text)
			}
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

func IsFlag(val string) {

}
