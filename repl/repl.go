package repl

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/0x0f0f0f/gobba-golang/parser"
	"github.com/0x0f0f0f/gobba-golang/token"
	"github.com/alecthomas/repr"
	"github.com/c-bata/go-prompt"
	"os"
)

// TODO: configurable
const PROMPT = "> "

type ReplOptions struct {
	ShowAST      bool
	ShowTok      bool
	PromptString string
}

type Repl struct {
	Options *ReplOptions
	prompt  *prompt.Prompt
}

func (r *Repl) executor(line string) {
	if r.Options.ShowTok {
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}

	l := lexer.New(line)
	p := parser.New(l)

	pri := repr.New(os.Stdout, repr.Hide(token.Token{}))
	program := p.ParseProgram()

	if r.Options.ShowAST {
		pri.Println(program)
	}

	for _, err := range p.Errors() {
		fmt.Printf("%s\n", err)
	}

}

func (r *Repl) completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func New(o *ReplOptions) *Repl {
	r := &Repl{}
	p := prompt.New(r.executor, r.completer)
	r.prompt = p
	r.Options = o

	return r
}

func (r *Repl) Start() {
	r.prompt.Run()
}