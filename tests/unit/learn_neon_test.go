package tests

import (
	"testing"
	"vectorengine/src"
)

func TestAssemblyParameter(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8}
	v1 := src.TestAssemblyParameter(values)
	if v1 != 36 {
		t.Fatalf("expected 36, got %d", v1)
	}
}
