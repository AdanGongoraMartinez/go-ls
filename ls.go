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

func printLong(dir, name string) {
	fileInfo, err := os.Lstat(filepath.Join(dir, name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		return
	}
	fileStats := fileInfo.Sys().(*syscall.Stat_t)
	fmt.Println(
		fileInfo.Mode().String(),
		fileStats.Nlink,
		fileStats.Uid,
		fileStats.Gid,
		fileInfo.Size(),
		name)
}
