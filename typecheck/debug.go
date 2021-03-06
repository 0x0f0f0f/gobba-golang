package typecheck

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"os"
)

// Global flag to print typechecking debug statements
var DebugTypeCheck bool = false

func (c Context) debugRule(name string) {
	if DebugTypeCheck {
		fmt.Fprintln(os.Stderr, "\tApplying rule", name, c)
	}
}

func (c Context) debugRuleOut(name string) {
	if DebugTypeCheck {
		fmt.Fprintln(os.Stderr, "\tSUCCESSfully applied rule", name, c)
	}
}

func (c Context) debugRuleFail(name string) {
	if DebugTypeCheck {
		fmt.Fprintln(os.Stderr, "\tFAILed to apply rule", name)
	}
}

func (c Context) debugErr(err error) {
	if DebugTypeCheck {
		fmt.Fprintln(os.Stderr, "Type Error: ", err, c)
	}
}

func (c Context) debugSection(name string, rest ...string) {
	if DebugTypeCheck {
		fmt.Fprintln(os.Stderr, name, rest)
	}
}

func (c Context) debugSynth(exp ast.Expression, t ast.TypeValue, printctx bool) {
	if DebugTypeCheck {
		if printctx {
			fmt.Fprintln(os.Stderr, exp, "=>", t.FullString(), "in", c)
		} else {
			fmt.Fprintln(os.Stderr, exp, "=>", t.FullString())
		}

	}
}
