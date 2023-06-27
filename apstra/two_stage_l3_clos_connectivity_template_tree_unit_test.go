package apstra

import (
	"testing"
)

func TestObjPolicyBatch(t *testing.T) {
	a := ObjPolicyBatchAttributes{Subpolicies: []ObjectId{"a", "b"}}

	etn := "batch"
	if etn != a.typeName() {
		t.Fatalf("expected %q got %q", etn, a.typeName())
	}

	raw, err := a.marshal()
	if err != nil {
		t.Fatal(err)
	}

	eraw := `{"subpolicies":["a","b"]}`
	if eraw != string(raw) {
		t.Fatalf("expected %q, got %q", eraw, string(raw))
	}
}

func TestObjPolicyPipeline(t *testing.T) {
	sp1 := ObjectId("c")
	sp2 := ObjectId("d")
	a := ObjPolicyPipelineAttributes{
		FirstSubpolicy:  &sp1,
		SecondSubpolicy: &sp2,
	}

	etn := "pipeline"
	if etn != a.typeName() {
		t.Fatalf("expected %q got %q", etn, a.typeName())
	}

	raw, err := a.marshal()
	if err != nil {
		t.Fatal(err)
	}

	eraw := `{"first_subpolicy":"c","second_subpolicy":"d"}`
	if eraw != string(raw) {
		t.Fatalf("expected %q, got %q", eraw, string(raw))
	}
}
