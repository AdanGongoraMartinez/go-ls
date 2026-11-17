package main

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestRenderLong(t *testing.T) {
	when := time.Date(2026, 11, 17, 7, 6, 0, 0, time.Local)
	entries := []Entry{
		{Name: "a.txt", Mode: 0o644, Nlink: 1, Owner: "adan", Group: "adan", Size: 3, Mtime: when},
		{
			Name: "link", Mode: os.ModeSymlink | 0o777, Nlink: 1, Owner: "adan", Group: "adan",
			Size: 10, Mtime: when, Link: "a.txt",
		},
	}

	var buf bytes.Buffer
	renderLong(&buf, entries)

	want := "-rw-r--r-- 1 adan adan  3 Nov 17 07:06 a.txt\n" +
		"lrwxrwxrwx 1 adan adan 10 Nov 17 07:06 link -> a.txt\n"
	if got := buf.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
