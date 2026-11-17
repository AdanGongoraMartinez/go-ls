package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
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
	modes := map[uint32]rune{
		syscall.S_IFDIR:  'd',
		syscall.S_IFLNK:  'l',
		syscall.S_IFCHR:  'c',
		syscall.S_IFBLK:  'b',
		syscall.S_IFIFO:  'p',
		syscall.S_IFSOCK: 's',
	}

	if r, ok := modes[mode&syscall.S_IFMT]; ok {
		s[0] = byte(r)
	} else {
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

	owner, group := "?", "?"
	u, err := user.LookupId(strconv.Itoa(int(fileStats.Uid)))
	if err == nil {
		owner = u.Username
	}
	g, err := user.LookupGroupId(strconv.Itoa(int(fileStats.Gid)))
	if err == nil {
		group = g.Name
	}

	fmt.Printf("%s %d %s %s %d %s %s\n",
		modeString(fileStats.Mode),
		fileStats.Nlink,
		owner,
		group,
		fileInfo.Size(),
		fileInfo.ModTime().Format("Jan _2 15:04"),
		name)
}
