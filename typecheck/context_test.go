package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
	// "github.com/0x0f0f0f/gobba-golang/lexer"
	// "github.com/0x0f0f0f/gobba-golang/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

var alphaid ast.UniqueIdentifier = ast.GenUID("alpha")
var betaid ast.UniqueIdentifier = ast.GenUID("beta")
var gammaid ast.UniqueIdentifier = ast.GenUID("gamma")
var deltaid ast.UniqueIdentifier = ast.GenUID("delta")
var epsilonid ast.UniqueIdentifier = ast.GenUID("epsilon")

var alphaext ContextValue = &ExistentialVariable{
	Identifier: alphaid,
}
var betaext ContextValue = &ExistentialVariable{
	Identifier: betaid,
}
var gammaext ContextValue = &ExistentialVariable{
	Identifier: gammaid,
}
var deltaext ContextValue = &ExistentialVariable{
	Identifier: deltaid,
}
var epsilonext ContextValue = &ExistentialVariable{
	Identifier: epsilonid,
}

var alphauniv ContextValue = &UniversalVariable{
	Identifier: alphaid,
}
var betauniv ContextValue = &UniversalVariable{
	Identifier: betaid,
}
var gammauniv ContextValue = &UniversalVariable{
	Identifier: gammaid,
}
var deltauniv ContextValue = &UniversalVariable{
	Identifier: deltaid,
}
var epsilonuniv ContextValue = &UniversalVariable{
	Identifier: epsilonid,
}

var SampleContext1 Context = Context{
	Contents: []ContextValue{
		alphaext,
		betaext,
		gammauniv,
	},
}

func TestInsert(t *testing.T) {
	tests := []struct {
		Cont    Context
		El      ContextValue
		Inserts []ContextValue
		Result  Context
	}{
		{SampleContext1, betaext, []ContextValue{epsilonext, deltaext},
			Context{
				Contents: []ContextValue{alphaext, epsilonext, deltaext, gammauniv},
			}},
		{SampleContext1, gammaext, []ContextValue{epsilonext, deltaext},
			Context{
				Contents: []ContextValue{alphaext, betaext, gammauniv, epsilonext, deltaext},
			}},
	}

	for _, tt := range tests {
		nc := tt.Cont.Insert(tt.El, tt.Inserts)

		assert.Equal(t, tt.Result, nc)
	}
}

func TestDrop(t *testing.T) {
	tests := []struct {
		Cont   Context
		El     ContextValue
		Result Context
	}{
		{SampleContext1, betaext,
			Context{
				Contents: []ContextValue{alphaext, gammauniv},
			}},
		{SampleContext1, gammauniv,
			Context{
				Contents: []ContextValue{alphaext, betaext},
			}},
		{SampleContext1, deltaext,
			Context{
				Contents: []ContextValue{alphaext, betaext, gammauniv},
			}},
	}

	for _, tt := range tests {
		nc := tt.Cont.Drop(tt.El)

		assert.Equal(t, tt.Result, nc)
	}
}
