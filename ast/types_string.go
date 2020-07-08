package ast

import (
	"fmt"
)

// This file contains string representation of type values

func (u *UnitType) String() string     { return "unit" }
func (u *ExistsType) String() string   { return "∃'" + u.Identifier.String() }
func (u *VariableType) String() string { return u.Identifier.String() }
func (u *ForAllType) String() string {
	return fmt.Sprintf("∀%s.%s", u.Identifier.String(), u.Type.String())
}

func (u *LambdaType) String() string {
	return fmt.Sprintf("%s -> %s", u.Domain.String(), u.Codomain.String())
}

func (u *UnitType) FullString() string     { return u.String() }
func (u *VariableType) FullString() string { return u.Identifier.FullString() }
func (u *ForAllType) FullString() string {
	return fmt.Sprintf("∀%s.%s", u.Identifier.FullString(), u.Type.String())
}
func (u *LambdaType) FullString() string {
	return fmt.Sprintf("%s -> %s", u.Domain.FullString(), u.Codomain.FullString())
}
func (u *ExistsType) FullString() string { return "∃'" + u.Identifier.FullString() }

// helper for generating fancy names
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

func (u *UnitType) FancyString(occ map[UniqueIdentifier]int) string { return "unit" }
func (u *ExistsType) FancyString(occ map[UniqueIdentifier]int) string {
	return "'" + genFancy(occ, u.Identifier)
}
func (u *VariableType) FancyString(occ map[UniqueIdentifier]int) string {
	return u.String()
}
func (u *ForAllType) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("∀%s.%s", genFancy(occ, u.Identifier), u.Type.FancyString(occ))
}
func (u *LambdaType) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("%s -> %s", u.Domain.FancyString(occ), u.Codomain.FancyString(occ))
}
