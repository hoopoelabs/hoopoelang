package hclsyntax

import (
	"testing"
	"github.com/hoopoelabs/hoopoelang"
)

func TestParseForBlock(t *testing.T) {
	src := []byte(`
for k, v in var.tags {
    tag {
        name = k
        value = v
    }
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

	if len(body.ForBlocks) != 1 {
		t.Fatalf("expected 1 for block, got %d", len(body.ForBlocks))
	}

	forBlock := body.ForBlocks[0]
	if forBlock.KeyVar != "k" {
		t.Errorf("expected key variable 'k', got %q", forBlock.KeyVar)
	}
	if forBlock.ValVar != "v" {
		t.Errorf("expected value variable 'v', got %q", forBlock.ValVar)
	}

	if forBlock.CollExpr == nil {
		t.Fatalf("expected collection expression on for block")
	}

	if len(forBlock.Body.Blocks) != 1 {
		t.Errorf("expected 1 block in for body, got %d", len(forBlock.Body.Blocks))
	}
}

func TestParseForBlockSingleVar(t *testing.T) {
	src := []byte(`
for item in var.items {
    id = item.id
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

	if len(body.ForBlocks) != 1 {
		t.Fatalf("expected 1 for block, got %d", len(body.ForBlocks))
	}

	forBlock := body.ForBlocks[0]
	if forBlock.KeyVar != "" {
		t.Errorf("expected empty key variable, got %q", forBlock.KeyVar)
	}
	if forBlock.ValVar != "item" {
		t.Errorf("expected value variable 'item', got %q", forBlock.ValVar)
	}

	if len(forBlock.Body.Attributes) != 1 {
		t.Errorf("expected 1 attribute in for body, got %d", len(forBlock.Body.Attributes))
	}
}
