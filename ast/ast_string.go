package ast

import (
	"bytes"
)

func (a *Assignment) String() string {
	var b bytes.Buffer

	b.WriteString(a.Name.String() + " = ")
	b.WriteString(a.Value.String())

	return b.String()
}

func (ls *LetStatement) String() string {
	var b bytes.Buffer

	b.WriteString(ls.TokenLiteral() + " ")
	for i, ass := range ls.Assignments {
		b.WriteString(ass.String())
		if i < len(ls.Assignments)-1 {
			b.WriteString(" and ")
		}
	}
	b.WriteString(";")

	return b.String()
}
func (le *ExprLet) String() string {
	var b bytes.Buffer

	b.WriteString("(let ")
	b.WriteString(le.Assignment.String())
	b.WriteString("; ")
	b.WriteString(le.Body.String())
	b.WriteString(")")

	return b.String()
}
func (i *ExprIf) String() string {
	var b bytes.Buffer

	b.WriteString("(if ")
	b.WriteString(i.Condition.String())
	b.WriteString(" then ")
	b.WriteString(i.Consequence.String())
	b.WriteString(" else ")
	b.WriteString(i.Alternative.String())
	b.WriteString(")")

	return b.String()
}

func (f *ExprLambda) String() string {
	var b bytes.Buffer

	b.WriteString("(λ ")
	b.WriteString(f.Param.String())
	b.WriteString(" . ")
	b.WriteString(f.Body.String())
	b.WriteString(")")

	return b.String()
}

func (f *ExprApply) String() string {
	var b bytes.Buffer

	b.WriteString(f.Function.String())
	b.WriteString("(")
	b.WriteString(f.Arg.String())
	b.WriteString(")")

	return b.String()
}
func (p *ExprPrefix) String() string {
	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(p.Operator)
	b.WriteString(p.Right.String())
	b.WriteString(")")
	return b.String()
}
func (p *ExprInfix) String() string {
	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(p.Left.String())
	b.WriteString(" " + p.Operator + " ")
	b.WriteString(p.Right.String())
	b.WriteString(")")
	return b.String()
}
func (i *ExprIdentifier) String() string {
	return i.Identifier.String()
}

func (i *ExprAnnot) String() string {
	if i.Type == nil {
		return i.Body.String()
	}
	return "(" + i.Body.String() + ": " + i.Type.String() + ")"
}

func (i *ExprFix) String() string {
	return "(fix" + i.Param.String() + " . " + i.Body.String() + ")"
}
