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
	// "github.com/c-bata/go-prompt"
	"github.com/peterh/liner"
	"os"
	"path/filepath"
)

type ReplOptions struct {
	ShowAST         bool
	ShowTok         bool
	DebugParser     bool
	PromptString    string
	HistoryFilename string
}

type Repl struct {
	Options *ReplOptions
	// prompt  *prompt.Prompt
	line *liner.State
}

// TODO go-prompt live prefix

func New(o *ReplOptions) *Repl {
	r := &Repl{}
	// p := prompt.New(r.executor, r.completer)
	line := liner.NewLiner()

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	o.HistoryFilename = filepath.Join(home, ".gobba_history")

	if f, err := os.Open(o.HistoryFilename); err == nil {
		line.ReadHistory(f)
		f.Close()
	} else if os.IsNotExist(err) {
		if f, err := os.Create(o.HistoryFilename); err == nil {
			f.Close()
		} else {
			panic(err)

		}
	} else {
		panic(err)
	}

	line.SetCtrlCAborts(true)

	r.Options = o
	r.line = line
	return r
}

func (r *Repl) Start() {
	run := true
	for run {
		if l, err := r.line.Prompt("> "); err == nil {
			r.line.AppendHistory(l)
			r.executor(l)
		} else if err == liner.ErrPromptAborted {
			fmt.Fprintln(os.Stderr, "Aborted")
		} else if err.Error() == "EOF" {
			run = false
		} else {
			repr.Println(err)
			fmt.Fprintln(os.Stderr, "Error reading line: ", err)
			run = false
		}

		if f, err := os.Create(r.Options.HistoryFilename); err != nil {
			fmt.Fprintln(os.Stderr, "Error writing history file: ", err)
			run = false
		} else {
			r.line.WriteHistory(f)
			f.Close()
		}

	}
	r.line.Close()
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

// func (r *Repl) completer(t prompt.Document) []prompt.Suggest {
// return []prompt.Suggest{}
// }
