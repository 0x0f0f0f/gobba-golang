package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// This file contains definitions for synthesization rules

// TODO SynthesizesTo
func (c Context) SynthesizesTo(exp ast.Expression) (ast.TypeValue, Context, error) {
	c.debugSection("synth", exp.String())
	switch ve := exp.(type) {
	case *ast.UnitLiteral: // Rule 1I=>
		c.debugRuleOut("1I=>")
		return &ast.UnitType{}, c, nil
	case *ast.IntegerLiteral: // Rule intI=>
		c.debugRuleOut("intI=>")
		return ast.NewVariableType("int"), c, nil
	case *ast.FloatLiteral: // Rule floatI=>
		c.debugRuleOut("floatI=>")
		return ast.NewVariableType("float"), c, nil

	case *ast.ComplexLiteral: // Rule complexI=>
		c.debugRuleOut("complexI=>")
		return ast.NewVariableType("complex"), c, nil

	case *ast.BoolLiteral: // Rule boolI=>
		c.debugRuleOut("boolI=>")
		return ast.NewVariableType("bool"), c, nil

	case *ast.StringLiteral: // Rule stringI=>
		c.debugRuleOut("stringI=>")
		return ast.NewVariableType("string"), c, nil

	case *ast.RuneLiteral: // Rule runeI=>
		c.debugRuleOut("runeI=>")
		return ast.NewVariableType("rune"), c, nil

	case *ast.IdentifierExpr:
		// Rule Var
		c.debugRuleOut("Var")
		annot := c.GetAnnotation(ve.Identifier)
		if annot == nil {
			c.debugRuleFail("Var")
			return nil, c, c.notInContextError(ve.Identifier)
		}
		c.debugRuleOut("Var")
		return *annot, c, nil
	case *ast.IfExpression:
		// Rules ifthen<:else=> and ifelse<:then=> share the first
		// 3 premises
		c.debugRule("ifthen<:else=> or ifelse<:then=>")

		gamma1, err := c.CheckAgainst(ve.Condition, ast.NewVariableType("bool"))
		if err != nil {
			c.debugRuleFail("ifthen<:else=> or ifelse<:then=>")
			return nil, c, err
		}
		thent, theta, err := gamma1.SynthesizesTo(ve.Consequence)
		if err != nil {
			c.debugRuleFail("ifthen<:else=> or ifelse<:then=>")
			return nil, c, err
		}
		elset, theta1, err := theta.SynthesizesTo(ve.Alternative)
		if err != nil {
			c.debugRuleFail("ifthen<:else=> or ifelse<:then=>")
			return nil, c, err
		}

		// Try to see if thent <: elset
		var delta Context
		delta, err = theta1.Subtype(thent, elset)
		if err != nil {
			// Try other case where elset <: thent
			delta, err = theta1.Subtype(elset, thent)
			if err != nil {
				c.debugRuleFail("ifthen<:else=> or ifelse<:then=>")
				return nil, c, c.expectedSameTypeIfBranches(thent, elset)
			}
			// Rule ifelse<:then=>
			// thent is a supertype of elset
			delta.debugRuleOut("ifelse<:then=>")

			return thent, delta, nil
		}
		// Rule ifthen<:else=>
		// elset is a supertype of thent
		delta.debugRuleOut("ifthen<:else=>")
		return elset, delta, nil

	case *ast.InfixExpression:
		return c.synthInfixExpr(ve)
	case *ast.FunctionLiteral:
		// Rule ->l=>
		c.debugRule("->I=>")

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
			c.debugRuleFail("->I=>")
			return nil, c, err
		}

		funtype := &ast.LambdaType{Domain: alphaext, Codomain: betaext}
		deltadrop := delta.Drop(annot)
		deltadrop.debugRuleOut("->I=>")

		return funtype, deltadrop, nil
	case *ast.FixExpr:
		// Rule fixI=>
		c.debugRule("fixI=>")

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
			c.debugRuleFail("fixI=>")
			return nil, c, err
		}
		// funtype := &ast.LambdaType{Domain: alphaext, Codomain: betaext}
		deltadrop := delta.Drop(annot)
		deltadrop.debugRuleOut("fixI=>")

		return betaext, deltadrop, nil
	case *ast.ApplyExpr:
		// Rule ->E
		c.debugRule("->E")

		a, theta, err := c.SynthesizesTo(ve.Function)
		if err != nil {
			return nil, c, err
		}
		theta.debugRuleOut("->E")
		return theta.ApplicationSynthesizesTo(theta.Apply(a), ve.Arg)
		//TODO Rule Anno
	case *ast.AnnotExpr:
		if c.IsWellFormed(ve.Type) {
			// Rule Anno
			c.debugRule("Anno")
			delta, err := c.CheckAgainst(ve.Body, ve.Type)
			if err != nil {
				c.debugRuleFail("Anno")
				return nil, c, err
			}

			delta.debugRuleOut("Anno")
			return ve.Type, delta, nil

		}
	}
	return nil, c, c.synthError(exp)
}

// TODO document
func (c Context) ApplicationSynthesizesTo(
	ty ast.TypeValue,
	exp ast.Expression) (ast.TypeValue, Context, error) {

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
			c.debugRuleFail("α^App")
			return nil, c, err
		}

		delta.debugRuleOut("α^App")
		return alpha2ext, delta, nil
	case *ast.ForAllType:
		// Rule ∀App
		c.debugRule("∀App")

		alpha := ast.GenUID("α")
		alphaexv := &ExistentialVariable{Identifier: alpha}
		alphaext := &ast.ExistsType{Identifier: alpha}
		gamma := c.InsertHead(alphaexv)
		sub_a := Substitution(vty.Type, alphaext, vty.Identifier)

		gamma.debugRuleOut("∀App")
		return gamma.ApplicationSynthesizesTo(sub_a, exp)
	case *ast.LambdaType:
		// Rule ->App
		c.debugRule("->App")

		delta, err := c.CheckAgainst(exp, vty.Domain)
		if err != nil {
			c.debugRuleFail("->App")
			return nil, c, err
		}
		return vty.Codomain, delta, nil
	}

	return nil, c, c.synthError(exp)
}

func (c Context) SynthExpr(exp ast.Expression) (ast.TypeValue, error) {
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
