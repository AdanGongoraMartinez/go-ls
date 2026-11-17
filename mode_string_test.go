package main

import (
	"os"
	"testing"
)

func TestModeString(t *testing.T) {
	tests := []struct {
		mode os.FileMode
		want string
	}{
		{0o644, "-rw-r--r--"},
		{os.ModeSetuid | 0o755, "-rwsr-xr-x"}, // el bug de §1.2
		{os.ModeSetuid | 0o644, "-rwSr--r--"}, // setuid sin 'x'
		{os.ModeDir | 0o755, "drwxr-xr-x"},
		{os.ModeSymlink | 0o777, "lrwxrwxrwx"},
		{os.ModeSticky | 0o777, "-rwxrwxrwt"},
		{os.ModeNamedPipe | 0o644, "prw-r--r--"},
	}
	for _, tt := range tests {
		if got := modeString(tt.mode); got != tt.want {
			t.Errorf("modeString(%o) = %q, want %q", tt.mode, got, tt.want)
		}
	}
}
