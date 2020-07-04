package ast

import (
	"fmt"
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
	// outer *AlphaEnvironment
}

// Create a new empty α environment for α-conversion
func NewAlphaEnvironment() *AlphaEnvironment {
	s := make(map[string]int)
	return &AlphaEnvironment{store: s}
}

// Create a new empty α environment for α-conversion, enclosed in another α environment
func NewEnclosedAlphaEnvironment(outer *AlphaEnvironment) *AlphaEnvironment {
	env := NewAlphaEnvironment()
	return env
}

// Search for an identifier in the environment or return an error
func (a *AlphaEnvironment) Get(name string) (UniqueIdentifier, error) {
	uid, ok := a.store[name]
	if !ok {
		return UniqueIdentifier{name, 0}, unboundError(name)
	}
	return UniqueIdentifier{Value: name, Id: uid}, nil
}

func (a *AlphaEnvironment) IdentifierAlphaConversion(uid UniqueIdentifier) UniqueIdentifier {
	uid, err := a.Get(uid.Value)
	if err != nil {
		a.store[uid.Value] = 0
		return UniqueIdentifier{uid.Value, 0}
	}

	a.store[uid.Value] = uid.Id + 1
	uid.Id += 1
	return uid
}

func (a *AlphaEnvironment) ExpressionAlphaConversion(exp Expression) (Expression, error) {
	switch ve := exp.(type) {
	case *UnitLiteral:
		return ve, nil
	case *IntegerLiteral:
		return ve, nil
	case *FloatLiteral:
		return ve, nil
	case *ComplexLiteral:
		return ve, nil
	case *BoolLiteral:
		return ve, nil
	case *StringLiteral:
		return ve, nil
	case *RuneLiteral:
		return ve, nil
	case *PrefixExpression:
		nright, err := a.ExpressionAlphaConversion(ve.Right)
		if err != nil {
			return nil, err
		}
		var nexpr PrefixExpression
		err = copier.Copy(&nexpr, ve)
		if err != nil {
			return nil, err
		}

		nexpr.Right = nright
		return &nexpr, nil
	case *InfixExpression:
		nright, err := a.ExpressionAlphaConversion(ve.Right)
		if err != nil {
			return nil, err
		}
		nleft, err := a.ExpressionAlphaConversion(ve.Left)
		if err != nil {
			return nil, err
		}
		var nexpr InfixExpression
		err = copier.Copy(&nexpr, ve)
		if err != nil {
			return nil, err
		}
		nexpr.Left = nleft
		nexpr.Right = nright
		return &nexpr, nil

	case *IdentifierExpr:
		uid, err := a.Get(ve.Identifier.Value)
		if err != nil {
			return nil, err
		}
		var newexpr *IdentifierExpr
		err = copier.Copy(newexpr, ve)
		if err != nil {
			return nil, err
		}
		newexpr.Identifier = uid
		return newexpr, nil
	case *FunctionLiteral:
		nid := a.IdentifierAlphaConversion(ve.Param.Identifier)
		var newexpr FunctionLiteral
		err := copier.Copy(&newexpr, ve)
		if err != nil {
			return nil, err
		}
		newexpr.Param.Identifier = nid
		return &newexpr, nil
	case *ApplyExpr:
		nfun, err := a.ExpressionAlphaConversion(ve.Function)
		if err != nil {
			return nil, err
		}
		narg, err := a.ExpressionAlphaConversion(ve.Arg)
		if err != nil {
			return nil, err
		}

		var nexpr ApplyExpr
		copier.Copy(&nexpr, ve)
		if err != nil {
			return nil, err
		}
		nexpr.Function = nfun
		nexpr.Arg = narg
		return &nexpr, nil
	// TODO other expressions

	default:
		panic(fmt.Sprintf("alpha conversion not implemented yet for expression of type %T", ve))
	}
}

// TODO include primitives
var default_alpha_environment = NewAlphaEnvironment()

// Apply α-conversion on a statement
func (a *AlphaEnvironment) StatementAlphaConversion(stmt Statement) (Statement, error) {
	switch vs := stmt.(type) {
	case *ExpressionStatement:
		// The new statement to return
		var ns ExpressionStatement
		// The new expression to return
		nexpr, err := a.ExpressionAlphaConversion(vs.Expression)
		if err != nil {
			return nil, err
		}
		err = copier.Copy(&ns, vs)
		if err != nil {
			panic(err)
			// return nil, err
		}
		ns.Expression = nexpr
		return &ns, nil
	// TODO other statements
	default:
		panic(fmt.Sprintf("alpha conversion not implemented yet for statement of type %T", vs))
	}

}

// Apply α-conversion on a program
func ProgramAlphaConversion(p *Program) (*Program, error) {
	np := &Program{}
	np.Statements = make([]Statement, 0)

	env := NewAlphaEnvironment()

	for _, stmt := range p.Statements {
		newstmt, err := env.StatementAlphaConversion(stmt)
		if err != nil {
			return nil, err
		}
		np.Statements = append(np.Statements, newstmt)
	}

	return np, nil
}
