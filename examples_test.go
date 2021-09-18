//hershey fonts drawing tests

package hershey

import (
	"testing"
)

func TestDrawAllFontImage(t *testing.T) {
	err := DrawAllFontImage()
	if err != nil {
		t.Fatalf(`Got %v, want <nil>`, err)
	}
}

func TestDrawAllFontStringImage(t *testing.T) {
	err := DrawAllFontStringImage()
	if err != nil {
		t.Fatalf(`Got %v, want <nil>`, err)
	}
}

func TestDrawAStringLines(t *testing.T) {
	err := DrawAStringLines()
	if err != nil {
		t.Fatalf(`Got %v, want <nil>`, err)
	}
}
