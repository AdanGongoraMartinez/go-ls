package main

import (
	"bytes"
	"testing"
)

func TestRenderShort(t *testing.T) {
	entries := []Entry{{Name: "a"}, {Name: "b"}}
	var buf bytes.Buffer
	renderShort(&buf, entries)

	if got, want := buf.String(), "a\nb\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
