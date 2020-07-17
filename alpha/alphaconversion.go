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
	case *ast.ExprPrefix:
		nright, err := a.ExpressionAlphaConversion(ve.Right)
		if err != nil {
			return nil, err
		}
		var nexpr ast.ExprPrefix
		err = copier.Copy(&nexpr, ve)
		if err != nil {
			return nil, err
		}

		nexpr.Right = nright
		return &nexpr, nil
	case *ast.ExprInfix:
		nright, err := a.ExpressionAlphaConversion(ve.Right)
		if err != nil {
			return nil, err
		}
		nleft, err := a.ExpressionAlphaConversion(ve.Left)
		if err != nil {
			return nil, err
		}
		var nexpr ast.ExprInfix
		err = copier.Copy(&nexpr, ve)
		if err != nil {
			return nil, err
		}
		nexpr.Left = nleft
		nexpr.Right = nright
		return &nexpr, nil

	case *ast.ExprIdentifier:
		uid, err := a.Get(ve.Identifier.Value)
		if err != nil {
			return nil, err
		}
		var newexpr ast.ExprIdentifier
		err = copier.Copy(&newexpr, ve)
		if err != nil {
			return nil, err
		}
		newexpr.Identifier = uid
		return &newexpr, nil
	case *ast.ExprLambda:
		na := NewAlphaEnvironmentExtension(a)
		nid := na.IdentifierAlphaConversion(ve.Param.Identifier)

		nbody, err := na.ExpressionAlphaConversion(ve.Body)
		if err != nil {
			return nil, err
		}

		var newexpr ast.ExprLambda
		err = copier.Copy(&newexpr, ve)
		if err != nil {
			return nil, err
		}
		newexpr.Param.Identifier = nid
		newexpr.Body = nbody
		return &newexpr, nil
	case *ast.ExprRec:
		na := NewAlphaEnvironmentExtension(a)
		nid := na.IdentifierAlphaConversion(ve.Name.Identifier)

		nbody, err := na.ExpressionAlphaConversion(ve.Body)
		if err != nil {
			return nil, err
		}

		var newexpr ast.ExprRec
		err = copier.Copy(&newexpr, ve)
		if err != nil {
			return nil, err
		}
		newexpr.Name.Identifier = nid
		newexpr.Body = nbody
		return &newexpr, nil
	case *ast.ExprApplySpine:
		nfun, err := a.ExpressionAlphaConversion(ve.Function)
		if err != nil {
			return nil, err
		}

		nspine := []ast.Expression{}

		for _, arg := range ve.Spine {
			narg, err := a.ExpressionAlphaConversion(arg)
			if err != nil {
				return nil, err
			}
			nspine = append(nspine, narg)
		}

		var nexpr ast.ExprApplySpine
		err = copier.Copy(&nexpr, ve)
		if err != nil {
			return nil, err
		}
		nexpr.Function = nfun
		nexpr.Spine = nspine
		return &nexpr, nil
	case *ast.ExprIf:
		ncond, err := a.ExpressionAlphaConversion(ve.Condition)
		if err != nil {
			return nil, err
		}
		ntbr, err := a.ExpressionAlphaConversion(ve.Consequence)
		if err != nil {
			return nil, err
		}
		nfbr, err := a.ExpressionAlphaConversion(ve.Alternative)
		if err != nil {
			return nil, err
		}
		var nexpr ast.ExprIf
		err = copier.Copy(&nexpr, ve)
		if err != nil {
			return nil, err
		}
		nexpr.Condition = ncond
		nexpr.Consequence = ntbr
		nexpr.Alternative = nfbr

		return &nexpr, nil
	case *ast.ExprAnnot:
		nbody, err := a.ExpressionAlphaConversion(ve.Body)
		if err != nil {
			return nil, err
		}
		var nexpr ast.ExprAnnot
		err = copier.Copy(&nexpr, ve)
		if err != nil {
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
