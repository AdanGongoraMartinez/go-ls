package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	path := "."
	showAll := false
	longFormat := false

	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			path = arg
			continue
		}
		// -l -la -a
		for _, char := range arg[1:] {
			switch char {
			case 'a':
				showAll = true
			case 'l':
				longFormat = true
			default:
				fmt.Fprintf(os.Stderr, "usage: %s [-a] [path]\n", os.Args[0])
				os.Exit(1)
			}
		}
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	var names []string
	if showAll {
		names = append(names, ".", "..")
	}

	for _, entry := range entries {
		if !showAll && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		names = append(names, entry.Name())
	}

	for _, name := range names {
		if longFormat {
			printLong(path, name)
		} else {
			fmt.Println(name)
		}
	}
}

func modeString(mode uint32) string {
	var s [10]byte
	switch mode & syscall.S_IFMT {
	case syscall.S_IFDIR:
		s[0] = 'd'
	case syscall.S_IFLNK:
		s[0] = 'l'
	case syscall.S_IFCHR:
		s[0] = 'c'
	case syscall.S_IFBLK:
		s[0] = 'b'
	case syscall.S_IFIFO:
		s[0] = 'p'
	case syscall.S_IFSOCK:
		s[0] = 's'
	default:
		s[0] = '-'
	}
	return string(s[:1]) + os.FileMode(mode).Perm().String()[1:]
}

func printLong(dir, name string) {
	fileInfo, err := os.Lstat(filepath.Join(dir, name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		return
	}
	fileStats := fileInfo.Sys().(*syscall.Stat_t)
	fmt.Println(
		modeString(fileStats.Mode),
		fileStats.Nlink,
		fileStats.Uid,
		fileStats.Gid,
		fileInfo.Size(),
		name)
}
