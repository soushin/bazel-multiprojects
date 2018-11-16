package util

import "testing"

func TestAdd(t *testing.T) {
	actual := Add("test")
	expected := "test - built by Bazel!"
	if actual != expected {
		t.Errorf("invalid text: got %s want %s", actual, expected)
	}
}
