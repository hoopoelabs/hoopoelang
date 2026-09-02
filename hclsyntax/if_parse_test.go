package hclsyntax

import (
	"testing"
	"github.com/hoopoelabs/hoopoelang"
)

func TestParseIfBlock(t *testing.T) {
	src := []byte(`
if var.env == "prod" {
    instance_type = "large"
} else if var.env == "staging" {
    instance_type = "medium"
} else {
    instance_type = "small"
}
`)

	file, diags := ParseConfig(src, "test.hop", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		t.Fatalf("unexpected errors: %s", diags)
	}

	body, ok := file.Body.(*Body)
	if !ok {
		t.Fatalf("expected *Body, got %T", file.Body)
	}

	if len(body.IfBlocks) != 1 {
		t.Fatalf("expected 1 if block, got %d", len(body.IfBlocks))
	}

	ifBlock := body.IfBlocks[0]
	if ifBlock.Condition == nil {
		t.Fatalf("expected condition on if block")
	}

	if len(ifBlock.Body.Attributes) != 1 {
		t.Errorf("expected 1 attribute in if body, got %d", len(ifBlock.Body.Attributes))
	}

	if ifBlock.Else == nil {
		t.Fatalf("expected else block")
	}

	if len(ifBlock.Else.IfBlocks) != 1 {
		t.Fatalf("expected else-if block inside else body, got %d", len(ifBlock.Else.IfBlocks))
	}

	elseIfBlock := ifBlock.Else.IfBlocks[0]
	if len(elseIfBlock.Body.Attributes) != 1 {
		t.Errorf("expected 1 attribute in else-if body, got %d", len(elseIfBlock.Body.Attributes))
	}

	if elseIfBlock.Else == nil {
		t.Fatalf("expected final else block")
	}

	if len(elseIfBlock.Else.Attributes) != 1 {
		t.Errorf("expected 1 attribute in final else body, got %d", len(elseIfBlock.Else.Attributes))
	}
}
