package typecheck

import (
	// "fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// Returns true if a type variable identifier occurs in a given type
func OccursIn(alpha ast.UniqueIdentifier, a ast.TypeValue) bool {
	switch va := a.(type) {
	case *ast.TyUnVar:
		return va.Identifier == alpha
	case *ast.TyExVar:
		return va.Identifier == alpha
	case *ast.TyLambda:
		return OccursIn(alpha, va.Domain) || OccursIn(alpha, va.Codomain)
	case *ast.TyForAll:
		return va.Identifier == alpha || OccursIn(alpha, va.Type)
	default:
		// Type variables do not occur in monotypes
		return false
	}

}

// TODO document
func Substitution(a, b ast.TypeValue, alpha ast.UniqueIdentifier) ast.TypeValue {
	switch va := a.(type) {
	case *ast.TyUnVar:
		if va.Identifier == alpha {
			return b
		} else {
			return a
		}
	case *ast.TyExVar:
		if va.Identifier == alpha {
			return b
		} else {
			return a
		}
	case *ast.TyForAll:
		if va.Identifier == alpha {
			return &ast.TyForAll{
				Identifier: va.Identifier,
				Sort:       va.Sort,
				Type:       b,
			}
		} else {
			return &ast.TyForAll{
				Identifier: va.Identifier,
				Type:       Substitution(va.Type, b, alpha),
			}
		}
	case *ast.TyLambda:
		return &ast.TyLambda{
			Domain:   Substitution(va.Domain, b, alpha),
			Codomain: Substitution(va.Codomain, b, alpha),
		}
	default:
		return a

	}

}
