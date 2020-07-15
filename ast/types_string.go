package ast

import (
	"fmt"
)

// This file contains string representation of type values

func (u *TyUnit) String() string  { return "unit" }
func (u *TyExVar) String() string { return "∃'" + u.Identifier.String() }
func (u *TyUnVar) String() string { return u.Identifier.String() }
func (u *TyForAll) String() string {
	return fmt.Sprintf("∀%s.%s", u.Identifier.String(), u.Type.String())
}
func (u *TyLambda) String() string {
	return fmt.Sprintf("%s -> %s", u.Domain.String(), u.Codomain.String())
}
func (u *TySum) String() string {
	return fmt.Sprintf("%s -> %s", u.Left.String(), u.Right.String())
}
func (u *TyProduct) String() string {
	return fmt.Sprintf("%s -> %s", u.Left.String(), u.Right.String())
}

func (u *TyUnit) FullString() string  { return u.String() }
func (u *TyUnVar) FullString() string { return u.Identifier.FullString() }
func (u *TyExVar) FullString() string { return "∃'" + u.Identifier.FullString() }
func (u *TyForAll) FullString() string {
	return fmt.Sprintf("∀%s.%s", u.Identifier.FullString(), u.Type.String())
}
func (u *TyLambda) FullString() string {
	return fmt.Sprintf("%s -> %s", u.Domain.FullString(), u.Codomain.FullString())
}
func (u *TySum) FullString() string {
	return fmt.Sprintf("%s -> %s", u.Left.FullString(), u.Right.FullString())
}
func (u *TyProduct) FullString() string {
	return fmt.Sprintf("%s -> %s", u.Left.FullString(), u.Right.FullString())
}

// helper for generating fancy type names in OCaml style
func genFancy(occ map[UniqueIdentifier]int, id UniqueIdentifier) string {
	if num, ok := occ[id]; ok {
		return string(rune(num + 97))
	}

	// FIXME generate decent names
	max := -1
	for _, v := range occ {
		if v > max {
			max = v
		}
	}

	occ[id] = max + 1
	return string(rune(max + 1 + 97))

}

func (u *TyUnit) FancyString(occ map[UniqueIdentifier]int) string { return "unit" }
func (u *TyExVar) FancyString(occ map[UniqueIdentifier]int) string {
	return "'" + genFancy(occ, u.Identifier)
}
func (u *TyUnVar) FancyString(occ map[UniqueIdentifier]int) string {
	return u.String()
}
func (u *TyForAll) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("∀%s.%s", genFancy(occ, u.Identifier), u.Type.FancyString(occ))
}
func (u *TyLambda) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("%s -> %s", u.Domain.FancyString(occ), u.Codomain.FancyString(occ))
}
func (u *TySum) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("%s -> %s", u.Left.FancyString(occ), u.Right.FancyString(occ))
}
func (u *TyProduct) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("%s -> %s", u.Left.FancyString(occ), u.Right.FancyString(occ))
}
