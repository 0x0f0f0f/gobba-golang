package repl

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/0x0f0f0f/gobba-golang/parser"
	"github.com/0x0f0f0f/gobba-golang/token"
	"github.com/alecthomas/repr"
	"github.com/c-bata/go-prompt"
)

// TODO: configurable
const PROMPT = "> "

type Repl struct {
	promptString string
	prompt       *prompt.Prompt
}

func executor(line string) {
	l := lexer.New(line)
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Printf("%+v\n", tok)
	}
	l = lexer.New(line)
	p := parser.New(l)
	repr.Println(p)
	program := p.ParseProgram()
	repr.Println(p)
	repr.Println(program)

	for _, err := range p.Errors() {
		fmt.Printf("%s\n", err)
	}

	for _, stmt := range program.Statements {
		fmt.Printf("%#v\n", stmt)
	}
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func New() *Repl {
	return &Repl{
		promptString: "> ",
		prompt:       prompt.New(executor, completer),
	}
}

func (r *Repl) Start() {
	r.prompt.Run()
}
