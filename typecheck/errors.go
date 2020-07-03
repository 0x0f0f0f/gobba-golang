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
