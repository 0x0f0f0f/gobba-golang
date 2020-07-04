package typecheck

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// This file contains definitions for type checker errors

type TypeError struct {
	Msg string
}

func (te TypeError) Error() string {
	s := fmt.Sprintf("type error: ")
	s += fmt.Sprintf("%s", te.Msg)
	return s
}

// TODO how to handle errors
func (c *Context) malformedError(t ast.TypeValue) *TypeError {
	return &TypeError{fmt.Sprintf("type %s is not well formed", t)}
}

func (c *Context) subtypeError(a, b ast.TypeValue) *TypeError {
	return &TypeError{fmt.Sprintf("expected %s to be a subtype of %s", a, b)}
}

func (c *Context) synthError(expr ast.Expression) *TypeError {
	return &TypeError{fmt.Sprintf("failed to infer type for %s", expr)}
}

func (c *Context) notInContextError(id ast.UniqueIdentifier) *TypeError {
	return &TypeError{fmt.Sprintf("identifier %s not in context", id)}
}

func (c *Context) unexpectedType(expected, found ast.TypeValue) *TypeError {
	return &TypeError{
		fmt.Sprintf("unexpected value of type %s, expected a value of type %s", found, expected),
	}
}

func (c *Context) expectedSameTypeIfBranches(tt, ft ast.TypeValue) *TypeError {
	return &TypeError{
		fmt.Sprintf("type mismatch in if expression. then branch "+
			"has type %s while else branch has type %s", tt, ft),
	}

}
