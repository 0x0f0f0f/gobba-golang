// Contains data structures and method definitions for α-conversion
// of gobba programs and expressions. De-Brujin notation is used
package alpha

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/jinzhu/copier"
)

type AlphaConversionError struct {
	Msg string
}

func (ae AlphaConversionError) Error() string {
	s := fmt.Sprintf("alpha conversion error: ")
	s += fmt.Sprintf("%s", ae.Msg)
	return s
}

func unboundError(name string) *AlphaConversionError {
	return &AlphaConversionError{Msg: fmt.Sprintf("unbound identifier %s", name)}
}

// ======================================================================
// Alpha Environment definition and methods
// ======================================================================

// Contains mappings to integers for unique identifiers
type AlphaEnvironment struct {
	store map[string]int
	outer *AlphaEnvironment
}

// Create a new empty α environment for α-conversion
func NewAlphaEnvironment() *AlphaEnvironment {
	s := make(map[string]int)
	return &AlphaEnvironment{store: s, outer: nil}
}

func NewAlphaEnvironmentExtension(a *AlphaEnvironment) *AlphaEnvironment {
	n := NewAlphaEnvironment()
	n.outer = a
	return n
}

// Search for an identifier in the environment or return an error
func (a AlphaEnvironment) Get(name string) (ast.UniqueIdentifier, error) {
	uid, ok := a.store[name]
	if !ok {
		if a.outer != nil {
			return a.outer.Get(name)
		}
		return ast.UniqueIdentifier{Value: name, Id: 0}, unboundError(name)
	}
	return ast.UniqueIdentifier{Value: name, Id: uid}, nil
}

func (a *AlphaEnvironment) IdentifierAlphaConversion(uid ast.UniqueIdentifier) ast.UniqueIdentifier {
	nuid, err := a.Get(uid.Value)
	if err != nil {
		a.store[uid.Value] = 0

		return ast.UniqueIdentifier{Value: uid.Value, Id: 0}
	}

	a.store[uid.Value] = nuid.Id + 1
	return ast.UniqueIdentifier{Value: uid.Value, Id: nuid.Id + 1}
}

func (a *AlphaEnvironment) ExpressionAlphaConversion(exp ast.Expression) (ast.Expression, error) {
	// fmt.Println("uniquifying", exp, "in", a.store)

	switch ve := exp.(type) {
	case *ast.UnitLiteral:
		return ve, nil
	case *ast.IntegerLiteral:
		return ve, nil
	case *ast.FloatLiteral:
		return ve, nil
	case *ast.ComplexLiteral:
		return ve, nil
	case *ast.BoolLiteral:
		return ve, nil
	case *ast.StringLiteral:
		return ve, nil
	case *ast.RuneLiteral:
		return ve, nil
	case *ast.PrefixExpression:
		nright, err := a.ExpressionAlphaConversion(ve.Right)
		if err != nil {
			panic(err)
			return nil, err
		}
		var nexpr ast.PrefixExpression
		err = copier.Copy(&nexpr, ve)
		if err != nil {
			panic(err)
			return nil, err
		}

		nexpr.Right = nright
		return &nexpr, nil
	case *ast.InfixExpression:
		nright, err := a.ExpressionAlphaConversion(ve.Right)
		if err != nil {
			panic(err)
			return nil, err
		}
		nleft, err := a.ExpressionAlphaConversion(ve.Left)
		if err != nil {
			panic(err)
			return nil, err
		}
		var nexpr ast.InfixExpression
		err = copier.Copy(&nexpr, ve)
		if err != nil {
			panic(err)
			return nil, err
		}
		nexpr.Left = nleft
		nexpr.Right = nright
		return &nexpr, nil

	case *ast.IdentifierExpr:
		uid, err := a.Get(ve.Identifier.Value)
		if err != nil {
			panic(err)
			return nil, err
		}
		var newexpr ast.IdentifierExpr
		err = copier.Copy(&newexpr, ve)
		if err != nil {
			panic(err)
			return nil, err
		}
		newexpr.Identifier = uid
		return &newexpr, nil
	case *ast.FunctionLiteral:
		na := NewAlphaEnvironmentExtension(a)
		nid := na.IdentifierAlphaConversion(ve.Param.Identifier)

		nbody, err := na.ExpressionAlphaConversion(ve.Body)
		if err != nil {
			panic(err)
			return nil, err
		}

		var newexpr ast.FunctionLiteral
		err = copier.Copy(&newexpr, ve)
		if err != nil {
			panic(err)
			return nil, err
		}
		newexpr.Param.Identifier = nid
		newexpr.Body = nbody
		return &newexpr, nil
	case *ast.FixExpr:
		na := NewAlphaEnvironmentExtension(a)
		nid := na.IdentifierAlphaConversion(ve.Param.Identifier)

		nbody, err := na.ExpressionAlphaConversion(ve.Body)
		if err != nil {
			panic(err)
			return nil, err
		}

		var newexpr ast.FixExpr
		err = copier.Copy(&newexpr, ve)
		if err != nil {
			panic(err)
			return nil, err
		}
		newexpr.Param.Identifier = nid
		newexpr.Body = nbody
		return &newexpr, nil
	case *ast.ApplyExpr:
		nfun, err := a.ExpressionAlphaConversion(ve.Function)
		if err != nil {
			panic(err)
			return nil, err
		}
		narg, err := a.ExpressionAlphaConversion(ve.Arg)
		if err != nil {
			panic(err)
			return nil, err
		}

		var nexpr ast.ApplyExpr
		err = copier.Copy(&nexpr, ve)
		if err != nil {
			panic(err)
			return nil, err
		}
		nexpr.Function = nfun
		nexpr.Arg = narg
		return &nexpr, nil
	case *ast.IfExpression:
		ncond, err := a.ExpressionAlphaConversion(ve.Condition)
		if err != nil {
			panic(err)
			return nil, err
		}
		ntbr, err := a.ExpressionAlphaConversion(ve.Consequence)
		if err != nil {
			panic(err)
			return nil, err
		}
		nfbr, err := a.ExpressionAlphaConversion(ve.Alternative)
		if err != nil {
			panic(err)
			return nil, err
		}
		var nexpr ast.IfExpression
		err = copier.Copy(&nexpr, ve)
		if err != nil {
			panic(err)
			return nil, err
		}
		nexpr.Condition = ncond
		nexpr.Consequence = ntbr
		nexpr.Alternative = nfbr

		return &nexpr, nil
	case *ast.AnnotExpr:
		nbody, err := a.ExpressionAlphaConversion(ve.Body)
		if err != nil {
			panic(err)
			return nil, err
		}
		var nexpr ast.AnnotExpr
		err = copier.Copy(&nexpr, ve)
		if err != nil {
			panic(err)
			return nil, err
		}
		nexpr.Body = nbody
		return &nexpr, nil

	default: // TODO other expressions
		panic(fmt.Sprintf("alpha conversion not implemented yet for expression of type %T", ve))
	}
}

// TODO include primitives
var default_alpha_environment = NewAlphaEnvironment()

// // Apply α-conversion on a program
// func ProgramAlphaConversion(p *ast.Program) (*ast.Program, error) {
// 	np := &ast.Program{}
// 	np.Statements = make([]ast.Statement, 0)

// 	env := NewAlphaEnvironment()

// 	for _, stmt := range p.Statements {
// 		newstmt, err := env.StatementAlphaConversion(stmt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		np.Statements = append(np.Statements, newstmt)
// 	}

// 	return np, nil
// }

func ProgramAlphaConversion(p ast.Expression) (*ast.Expression, error) {
	env := NewAlphaEnvironment()
	np, err := env.ExpressionAlphaConversion(p)
	if err != nil {
		return nil, err
	}

	return &np, nil
}
