package rand

import (
	"testing"

	"github.com/ReconfigureIO/fixed"
)

func TestNormals(t *testing.T) {
	r := New(42)
	out := make(chan fixed.Int26_6)

	go r.Normals(out)

	var s fixed.Int26_6
	for i := 0; i < 1024; i++ {
		o := <-out
		s += o
	}
	// The mean should be 0
	if s.Floor()/1024 != 0 {
		t.Fail()
	}
}
