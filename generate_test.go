package hershey

import "testing"

func TestGenerate(t *testing.T) {
	err := Generate()
	if err != nil {
		t.Fatalf(`Got %v, want <nil>`, err)
	}
}

func TestGenerateHeights(t *testing.T) {
	err := GenerateHeights()
	if err != nil {
		t.Fatalf(`Got %v, want <nil>`, err)
	}
}
