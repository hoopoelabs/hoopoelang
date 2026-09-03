// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hclsyntax

import (
	"testing"

	"github.com/hoopoelabs/hoopoelang"
	"github.com/zclconf/go-cty/cty"
)

func TestNullCoalesce(t *testing.T) {
	tests := []struct {
		input     string
		ctx       *hcl.EvalContext
		want      cty.Value
		diagCount int
	}{
		{
			`null ?? "disabled"`,
			nil,
			cty.StringVal("disabled"),
			0,
		},
		{
			`"enabled" ?? "disabled"`,
			nil,
			cty.StringVal("enabled"),
			0,
		},
		{
			`0 ?? 1`,
			nil,
			cty.NumberIntVal(0),
			0,
		},
		{
			`false ?? true`,
			nil,
			cty.False,
			0,
		},
		{
			`null ?? null ?? "x"`,
			nil,
			cty.StringVal("x"),
			0,
		},
		{
			`true ? null ?? "a" : "b"`,
			nil,
			cty.StringVal("a"),
			0,
		},
		{
			`null ?? true ? "t" : "f"`,
			nil,
			cty.StringVal("t"),
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			expr, parseDiags := ParseExpression([]byte(test.input), "", hcl.Pos{Line: 1, Column: 1})
			if len(parseDiags) != 0 {
				t.Fatalf("unexpected parse diagnostics: %s", parseDiags.Error())
			}
			got, diags := expr.Value(test.ctx)
			if len(diags) != test.diagCount {
				t.Fatalf("wrong number of diagnostics %d; want %d\n%s", len(diags), test.diagCount, diags.Error())
			}
			if test.diagCount > 0 {
				return
			}
			if !got.RawEquals(test.want) {
				t.Fatalf("wrong result\ngot:  %#v\nwant: %#v", got, test.want)
			}
		})
	}
}

func TestOptionalChain(t *testing.T) {
	obj := cty.ObjectVal(map[string]cty.Value{
		"id": cty.StringVal("acl-123"),
		"nested": cty.ObjectVal(map[string]cty.Value{
			"port": cty.NumberIntVal(80),
		}),
		"tags": cty.TupleVal([]cty.Value{
			cty.StringVal("a"),
			cty.StringVal("b"),
		}),
	})

	tests := []struct {
		input     string
		ctx       *hcl.EvalContext
		want      cty.Value
		diagCount int
	}{
		{
			`obj?.id`,
			&hcl.EvalContext{Variables: map[string]cty.Value{"obj": obj}},
			cty.StringVal("acl-123"),
			0,
		},
		{
			`obj?.id ?? "disabled"`,
			&hcl.EvalContext{Variables: map[string]cty.Value{"obj": obj}},
			cty.StringVal("acl-123"),
			0,
		},
		{
			`obj?.id ?? "disabled"`,
			&hcl.EvalContext{Variables: map[string]cty.Value{"obj": cty.NullVal(cty.DynamicPseudoType)}},
			cty.StringVal("disabled"),
			0,
		},
		{
			`aws_waf_web_acl.main?.id ?? "disabled"`,
			&hcl.EvalContext{Variables: map[string]cty.Value{
				"aws_waf_web_acl": cty.ObjectVal(map[string]cty.Value{
					"main": cty.NullVal(cty.DynamicPseudoType),
				}),
			}},
			cty.StringVal("disabled"),
			0,
		},
		{
			`aws_waf_web_acl.main?.id ?? "disabled"`,
			&hcl.EvalContext{Variables: map[string]cty.Value{
				"aws_waf_web_acl": cty.ObjectVal(map[string]cty.Value{"main": obj}),
			}},
			cty.StringVal("acl-123"),
			0,
		},
		{
			`obj?.nested.port`,
			&hcl.EvalContext{Variables: map[string]cty.Value{"obj": cty.NullVal(cty.DynamicPseudoType)}},
			cty.NullVal(cty.DynamicPseudoType),
			0,
		},
		{
			`obj?.nested.port`,
			&hcl.EvalContext{Variables: map[string]cty.Value{"obj": obj}},
			cty.NumberIntVal(80),
			0,
		},
		{
			`obj?.nested?.port ?? 0`,
			&hcl.EvalContext{Variables: map[string]cty.Value{
				"obj": cty.ObjectVal(map[string]cty.Value{
					"nested": cty.NullVal(cty.DynamicPseudoType),
				}),
			}},
			cty.NumberIntVal(0),
			0,
		},
		{
			`obj?.tags?.[1]`,
			&hcl.EvalContext{Variables: map[string]cty.Value{"obj": obj}},
			cty.StringVal("b"),
			0,
		},
		{
			`obj?.tags?.[1]`,
			&hcl.EvalContext{Variables: map[string]cty.Value{"obj": cty.NullVal(cty.DynamicPseudoType)}},
			cty.NullVal(cty.DynamicPseudoType),
			0,
		},
		{
			// Missing variable is still an error; ?. does not suppress it.
			`missing?.id ?? "disabled"`,
			&hcl.EvalContext{},
			cty.DynamicVal,
			1,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			expr, parseDiags := ParseExpression([]byte(test.input), "", hcl.Pos{Line: 1, Column: 1})
			if len(parseDiags) != 0 {
				t.Fatalf("unexpected parse diagnostics: %s", parseDiags.Error())
			}
			got, diags := expr.Value(test.ctx)
			if len(diags) != test.diagCount {
				t.Fatalf("wrong number of diagnostics %d; want %d\n%s", len(diags), test.diagCount, diags.Error())
			}
			if test.diagCount > 0 {
				return
			}
			if !got.RawEquals(test.want) {
				t.Fatalf("wrong result\ngot:  %#v\nwant: %#v", got, test.want)
			}
		})
	}
}

func TestOptionalChainAndNullCoalesceScan(t *testing.T) {
	tokens := scanTokens([]byte(`a?.b ?? c`), "", hcl.Pos{Line: 1, Column: 1}, scanNormal)
	wantTypes := []TokenType{
		TokenIdent,
		TokenQuestionDot,
		TokenIdent,
		TokenQuestionQuestion,
		TokenIdent,
		TokenEOF,
	}
	if len(tokens) != len(wantTypes) {
		t.Fatalf("wrong token count %d; want %d\n%#v", len(tokens), len(wantTypes), tokens)
	}
	for i, want := range wantTypes {
		if tokens[i].Type != want {
			t.Errorf("token %d is %s; want %s", i, tokens[i].Type, want)
		}
	}
}
