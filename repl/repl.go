package repl

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/alpha"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/0x0f0f0f/gobba-golang/parser"
	"github.com/0x0f0f0f/gobba-golang/token"
	"github.com/0x0f0f0f/gobba-golang/typecheck"
	"github.com/alecthomas/repr"
	"github.com/c-bata/go-prompt"
	"os"
)

type ReplOptions struct {
	ShowAST      bool
	ShowTok      bool
	DebugParser  bool
	PromptString string
}

type Repl struct {
	Options *ReplOptions
	prompt  *prompt.Prompt
}

// TODO go-prompt live prefix

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

func (r *Repl) executor(line string) {
	if r.Options.ShowTok {
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}

	l := lexer.New(line)
	p := parser.New(l)
	p.TraceOnError = r.Options.DebugParser

	pri := repr.New(os.Stdout, repr.Hide(token.Token{}))
	program := p.ParseProgram()

	if r.Options.ShowAST {
		pri.Println(program)
	}

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		return
	}

	// Do alpha conversion on the program (generate unique identifiers)
	alphaconv_program, err := alpha.ProgramAlphaConversion(program)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if r.Options.ShowAST {
		pri.Println(alphaconv_program)
		fmt.Println((*alphaconv_program).String())
	}

	// Typecheck
	// TODO default context with primitives
	// TODO preserve context between statements in the repl
	ctx := typecheck.NewContext()
	ast.ResetUIDCounter()
	ty, err := ctx.SynthExpr(*alphaconv_program)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("- : %s\n", ty.FancyString(map[ast.UniqueIdentifier]int{}))

	// TODO evaluation
}

func (r *Repl) completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}
