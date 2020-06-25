package repl

import (
    "fmt"
    "github.com/0x0f0f0f/gobba-golang/lexer"
    "github.com/0x0f0f0f/gobba-golang/token"
    "github.com/c-bata/go-prompt"
)

// TODO: configurable
const PROMPT = "> "

type Repl struct {
    promptString string
    prompt *prompt.Prompt
}

func executor(line string) {
    l := lexer.New(line)
    for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
        fmt.Printf("%+v\n", tok)
    }
}

func completer(t prompt.Document) []prompt.Suggest {
    return []prompt.Suggest{}
}

func New() *Repl {
    return &Repl{
        promptString: "> ",
        prompt: prompt.New(executor, completer),
    }
}

func (r *Repl) Start() {
    r.prompt.Run()
}
