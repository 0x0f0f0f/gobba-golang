package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/alpha"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/0x0f0f0f/gobba-golang/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSynthExpr(t *testing.T) {
	tests := map[string]string{
		"();":                                    "unit",
		"true;":                                  "bool",
		"false;":                                 "bool",
		"4;":                                     "int",
		"4.5;":                                   "float",
		"4.5+3.2e-2i;":                           "complex",
		"fun (x) {x};":                           "'a -> 'a",
		"fun (x) {if x then 3 else 4}":           "bool -> int",
		"fun (x) {if x then x else x}(true)":     "bool",
		"fun (x) {if x then x else false}(true)": "bool",
		"fun (x) {if true then x else 4.5}(4)":   "float",
		"fun (x) {x}(2)":                         "int",
		"fun (x) {x}(2.2)":                       "float",
		"if false then 4.5+3i else 4.5":          "complex",
		"if true then 4 else 4.5":                "float",
		"if true then true else false":           "bool",
		"fun (x) {x()}(fun (y) {y})":             "unit",
		// Stressing inference
		"let swap = fun(x,y,f) {f(y,x)}; " +
			"let firstid = fun(a,b) {a};" +
			"swap(3,\"ciao\",firstid)": "string",
		// Fixed point combinator
		"let f = fun (x) {if x <= 1 then 1 else f(x)}; f":             "int -> int",
		"let id = fun(a){a}; let id1 = fun(b){b}; let f = id(id1); f": "'a -> 'a",
		"let x = 4 and y = 3.2 and f = fun(x,y) {x}; f(y)":            "'a -> float",
		// Arithmetic Operators
		"4.5 +. 4":     "float",
		"4 +: 3+3i":    "complex",
		"0 + 1":        "int",
		"4.5 +: 3+14i": "complex",
		// Left in binop is existential
		"fun (x) {x+1}(3)":              "int",
		"fun (x) {x +. 1.5}(3)":         "float",
		"fun (x) {x +. 1.5}(3.5)":       "float",
		"fun (x) {x +: 1.5+3i}(3)":      "complex",
		"fun (x) {x +: 1.5+3i}(3.5)":    "complex",
		"fun (x) {x +: 1.5+3i}(3.5+3i)": "complex",
		"fun (x) {x-1}(3)":              "int",
		"fun (x) {x -. 1.5}(3)":         "float",
		"fun (x) {x -. 1.5}(3.5)":       "float",
		"fun (x) {x -: 1.5-3i}(3)":      "complex",
		"fun (x) {x -: 1.5-3i}(3.5)":    "complex",
		"fun (x) {x -: 1.5-3i}(3.5-3i)": "complex",
		// Right in binop is existential
		"fun (x) {1 + x}(3)":            "int",
		"fun (x) {1.5 +. x}(3)":         "float",
		"fun (x) {1.5 +. x}(3.5)":       "float",
		"fun (x) {1.5+3i +: x}(3)":      "complex",
		"fun (x) {1.5+3i +: x}(3.5)":    "complex",
		"fun (x) {1.5+3i +: x}(3.5+3i)": "complex",

		// Unary operators
		"!true": "bool",

		// Type annotation functions
		"fun (x: int, y: int) { if x = 2 then y else 0}":                     "int -> int -> int",
		"fun (x: bool, y) {x = y}":                                           "bool -> bool -> bool",
		"let fib = fun(n) { if n < 2 then n else fib(n-1) + fib(n-2) }; fib": "int -> int",
	}

	for input, expected := range tests {
		t.Log("--- TEST CASE", input, "---")
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()
		// parser.CheckParseErrors(t, p)
		assert.Len(t, p.Errors(), 0)
		alphaconv_program, err := alpha.ProgramAlphaConversion(program)
		if err != nil {
			assert.Fail(t, "could not α-convert expression")
			return
		}

		ctx := NewContext()
		ast.ResetUIDCounter()
		ty, err := ctx.SynthExpr(*alphaconv_program)
		if assert.Nil(t, err) {
			assert.Equal(t, expected, ty.FancyString(map[ast.UniqueIdentifier]int{}), expected, input)
		}
	}
}

func TestSynthExprFail(t *testing.T) {
	tests := []string{
		// Nonsensical programs
		"2 + \"ciao\"",
		"fun(x) {x+1}(\"ciao\")",
		"fun(x) {1+x}(\"ciao\")",
		"fun(x,y) {x+y}(\"ciao\")",
		"fun (x) {x()}(fun (y) {y+1})",
		// Impredicativeness
		"fun (x) {x(x, ())}",
		// Arithmetical imprecision
		"fun (x) {x+1}(3.5)",
		"fun (x) {x+1}(3.5+3i)",
		"fun (x) {x+1.5}(3.5+3i)",
		"fun (x) {1+x}(3.5)",
		"fun (x) {1+x}(3.5+3i)",
		"fun (x) {1.5+x}(3.5+3i)",
	}

	for _, input := range tests {
		t.Log("--- TEST CASE", input, "---")
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()
		// parser.CheckParseErrors(t, p)
		assert.Len(t, p.Errors(), 0)
		alphaconv_program, err := alpha.ProgramAlphaConversion(program)
		if err != nil {
			assert.Fail(t, "could not α-convert expression")
			return
		}

		ctx := NewContext()
		ast.ResetUIDCounter()
		_, err = ctx.SynthExpr(*alphaconv_program)
		if err == nil {
			assert.Fail(t, "did not find any error")
		} else {
			t.Log(err)
		}
		assert.NotNil(t, err)
	}
}
