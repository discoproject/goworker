package jobutil

import (
	"testing"
)

func assert(key, expected string, t *testing.T) {
	if val := Setting(key); val != expected {
		t.Error("incorrect value: _", val, "_")
	}
}

func TestSingle(t *testing.T) {
	input := "hello = world\n"
	addLine(input)
	assert("hello", "world", t)
}

func TestComment(t *testing.T) {
	input := "# this is a comment"
	addLine(input)
}

func TestQuote(t *testing.T) {
	input := "this = \" # is not a comment\""
	addLine(input)
	assert("this", "# is not a comment", t)
}

func TestTrim(t *testing.T) {
	input := "empty = \" \""
	addLine(input)
	assert("empty", "", t)
}
