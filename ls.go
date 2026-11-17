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

var nlinkWidth, ownerWidth, groupWidth, fsizeWidth, mtimeWidth int

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

	if longFormat {
		setWidths(path, names)
		for _, name := range names {
			printLong(path, name)
		}
	} else {
		for _, name := range names {
			fmt.Println(name)
		}
	}
}

func setWidths(dir string, names []string) {
	for _, name := range names {
		_, nlink, owner, group, fsize, mtime := getLong(dir, &name)
		nlinkWidth = max(nlinkWidth, len(strconv.Itoa(int(nlink))))
		ownerWidth = max(ownerWidth, len(owner))
		groupWidth = max(groupWidth, len(group))
		fsizeWidth = max(fsizeWidth, len(strconv.Itoa(int(fsize))))
		mtimeWidth = max(mtimeWidth, len(mtime))
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

func getLong(dir string, name *string) (modestr string, nlink uint64, owner, group string, fsize int64, mtime string) {
	fullpath := filepath.Join(dir, *name)
	fileInfo, err := os.Lstat(fullpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", *name, err)
		return
	}

	fileStats := fileInfo.Sys().(*syscall.Stat_t)

	owner, group = "?", "?"

	if u, err := user.LookupId(strconv.Itoa(int(fileStats.Uid))); err == nil {
		owner = u.Username
	}
	if g, err := user.LookupGroupId(strconv.Itoa(int(fileStats.Gid))); err == nil {
		group = g.Name
	}

	if fileStats.Mode&syscall.S_IFLNK == syscall.S_IFLNK {
		target, err := os.Readlink(fullpath)
		if err == nil {
			*name += " -> " + target
		}
	}

	return modeString(fileStats.Mode), fileStats.Nlink, owner, group, fileInfo.Size(), fileInfo.ModTime().Format("Jan _2 15:04")
}

func printLong(dir, name string) {
	modestr, nlink, owner, group, fsize, mtime := getLong(dir, &name)

	fmt.Printf("%s %*d %*s %*s %*d %*s %s\n",
		modestr,
		nlinkWidth, nlink,
		ownerWidth, owner,
		groupWidth, group,
		fsizeWidth, fsize,
		mtimeWidth, mtime,
		name,
	)
}
