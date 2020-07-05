package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// This file contains definitions for synthesization rules

// TODO SynthesizesTo
func (c Context) SynthesizesTo(exp ast.Expression) (ast.TypeValue, Context, *TypeError) {
	c.debugSection("synth", exp.String())
	switch ve := exp.(type) {
	case *ast.UnitLiteral:
		// Rule 1l=>
		c.debugRule("1l=>")

		return &ast.UnitType{}, c, nil
	case *ast.IntegerLiteral:
		return &ast.IntegerType{}, c, nil
	case *ast.FloatLiteral:
		return &ast.FloatType{}, c, nil
	case *ast.ComplexLiteral:
		return &ast.ComplexType{}, c, nil
	case *ast.BoolLiteral:
		return &ast.BoolType{}, c, nil
	case *ast.StringLiteral:
		return &ast.StringType{}, c, nil
	case *ast.RuneLiteral:
		return &ast.RuneType{}, c, nil
	case *ast.IdentifierExpr:
		// Rule Var
		c.debugRule("Var")

		annot := c.GetAnnotation(ve.Identifier)
		if annot == nil {
			return nil, c, c.notInContextError(ve.Identifier)
		}
		return *annot, c, nil
	case *ast.IfExpression:
		// Rule ifthenelse=>
		c.debugRule("ifthenelse=>")

		gamma1, err := c.CheckAgainst(ve.Condition, &ast.BoolType{})
		if err != nil {
			return nil, c, err
		}
		tbrancht, theta, err := gamma1.SynthesizesTo(ve.Consequence)
		if err != nil {
			return nil, c, err
		}
		fbrancht, delta, err := theta.SynthesizesTo(ve.Alternative)
		if err != nil {
			return nil, c, err
		}

		if tbrancht != fbrancht {
			//FIXME use subtyping
			return nil, c, c.expectedSameTypeIfBranches(tbrancht, fbrancht)
		}
		return tbrancht, delta, nil

	// TODO case Binary operators
	// TODO case hastype
	case *ast.FunctionLiteral:
		// Rule ->l=>
		c.debugRule("->l=>")

		alpha := ast.GenUID("α")
		beta := ast.GenUID("β")
		alphaext := &ast.ExistsType{
			Identifier: alpha,
		}
		betaext := &ast.ExistsType{
			Identifier: beta,
		}
		alphaexv := &ExistentialVariable{
			Identifier: alpha,
		}
		betaexv := &ExistentialVariable{
			Identifier: beta,
		}
		annot := &TypeAnnotation{
			Identifier: ve.Param.Identifier,
			Value:      alphaext,
		}
		gamma := c.InsertHead(annot).InsertHead(betaexv).InsertHead(alphaexv)
		delta, err := gamma.CheckAgainst(ve.Body, betaext)
		if err != nil {
			return nil, c, err
		}

		funtype := &ast.LambdaType{Domain: alphaext, Codomain: betaext}
		deltadrop := delta.Drop(annot)
		return funtype, deltadrop, nil
	case *ast.ApplyExpr:
		// Rule ->E
		c.debugRule("->E")

		a, theta, err := c.SynthesizesTo(ve.Function)
		if err != nil {
			return nil, c, err
		}
		return theta.ApplicationSynthesizesTo(theta.Apply(a), ve.Arg)
		//TODO Rule Anno

	}
	return nil, c, c.synthError(exp)
}

// TODO document
func (c Context) ApplicationSynthesizesTo(
	ty ast.TypeValue,
	exp ast.Expression) (ast.TypeValue, Context, *TypeError) {

	switch vty := ty.(type) {
	case *ast.ExistsType:
		// Rule α^App
		c.debugRule("α^App")

		idexv := &ExistentialVariable{Identifier: vty.Identifier}
		alpha1 := ast.GenUID("α")
		alpha2 := ast.GenUID("α")
		alpha1exv := &ExistentialVariable{Identifier: alpha1}
		alpha2exv := &ExistentialVariable{Identifier: alpha2}
		alpha1ext := &ast.ExistsType{Identifier: alpha1}
		alpha2ext := &ast.ExistsType{Identifier: alpha2}

		var funt ast.TypeValue = &ast.LambdaType{
			Domain:   alpha1ext,
			Codomain: alpha2ext,
		}
		solvedexv := &ExistentialVariable{
			Identifier: vty.Identifier,
			Value:      &funt,
		}

		gamma := c.Insert(idexv, []ContextValue{
			alpha2exv,
			alpha1exv,
			solvedexv,
		})

		delta, err := gamma.CheckAgainst(exp, alpha1ext)
		if err != nil {
			return nil, c, err
		}

		return alpha2ext, delta, nil
	case *ast.ForAllType:
		// Rule ∀App
		c.debugRule("∀App")

		alpha := ast.GenUID("α")
		alphaexv := &ExistentialVariable{Identifier: alpha}
		alphaext := &ast.ExistsType{Identifier: alpha}
		gamma := c.InsertHead(alphaexv)
		sub_a := Substitution(vty.Type, alphaext, vty.Identifier)
		return gamma.ApplicationSynthesizesTo(sub_a, exp)
	case *ast.LambdaType:
		// Rule ->App
		c.debugRule("->App")

		delta, err := c.CheckAgainst(exp, vty.Domain)
		if err != nil {
			return nil, c, err
		}
		return vty.Codomain, delta, nil
	}

	return nil, c, c.synthError(exp)
}

func (c Context) SynthExpr(exp ast.Expression) (ast.TypeValue, *TypeError) {
	t, nc, err := c.SynthesizesTo(exp)
	if err != nil {
		c.debugErr(err)
		return nil, err
	}
	nc.debugSynth(exp, t, true)

	t = nc.Apply(t)
	nc.debugSynth(exp, t, false)
	return t, nil
}
