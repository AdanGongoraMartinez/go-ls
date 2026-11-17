package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type options struct {
	all  bool // -a
	long bool // -l
}

func main() {
	path := "."
	opt := options{all: false, long: false}

	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			path = arg
			continue
		}
		// -l -la -a
		for _, char := range arg[1:] {
			switch char {
			case 'a':
				opt.all = true
			case 'l':
				opt.long = true
			default:
				fmt.Fprintf(os.Stderr, "usage: %s [-a] [path]\n", os.Args[0])
				os.Exit(1)
			}
		}
	}

	entries, err := readEntries(path, opt)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if opt.long {
		renderLong(w, entries)
	} else {
		renderShort(w, entries)
	}
}

func modeString(m os.FileMode) string {
	var b [10]byte
	b[0] = '-'
	switch {
	case m.IsDir():
		b[0] = 'd'
	case m&os.ModeSymlink != 0:
		b[0] = 'l'
	case m&os.ModeCharDevice != 0:
		b[0] = 'c'
	case m&os.ModeDevice != 0:
		b[0] = 'b'
	case m&os.ModeNamedPipe != 0:
		b[0] = 'p'
	case m&os.ModeSocket != 0:
		b[0] = 's'
	}

	const rwx = "rwxrwxrwx"
	for i := range 9 {
		if m&(1<<uint(8-i)) != 0 {
			b[i+1] = rwx[i]
		} else {
			b[i+1] = '-'
		}
	}

	if m&os.ModeSetuid != 0 {
		b[3] = pick(b[3], 'x', 's', 'S')
	}
	if m&os.ModeSetgid != 0 {
		b[6] = pick(b[6], 'x', 's', 'S')
	}
	if m&os.ModeSticky != 0 {
		b[9] = pick(b[9], 'x', 't', 'T')
	}

	return string(b[:])
}

// pick devuelve el carácter "con bit" si el permiso de ejecución estaba prendido,
// o el "sin bit" (mayúscula) si no lo estaba.
func pick(actual, conExec, conBit, sinBit byte) byte {
	if actual == conExec {
		return conBit
	}
	return sinBit
}

func newEntry(name string, info os.FileInfo, dir string, ids *idCache) Entry {
	e := Entry{Name: name, Mode: info.Mode(), Size: info.Size(), Mtime: info.ModTime()}
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		e.Nlink = st.Nlink
		e.Owner = ids.userName(st.Uid)
		e.Group = ids.groupName(st.Gid)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		if target, err := os.Readlink(filepath.Join(dir, name)); err == nil {
			e.Link = target
		}
	}
	return e
}

func readEntries(dir string, opts options) ([]Entry, error) {
	dirents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	ids := &idCache{users: map[uint32]string{}, groups: map[uint32]string{}}
	out := make([]Entry, 0, len(dirents)+2)

	// . and .. for -a
	if opts.all {
		for _, n := range []string{".", ".."} {
			if !opts.long {
				out = append(out, Entry{Name: n})
				continue
			}
			info, err := os.Lstat(filepath.Join(dir, n))
			if err != nil {
				continue
			}
			out = append(out, newEntry(n, info, dir, ids))
		}
	}

	for _, de := range dirents {
		if !opts.all && strings.HasPrefix(de.Name(), ".") {
			continue
		}
		e := Entry{Name: de.Name()}
		if !opts.long {
			out = append(out, e)
			continue
		}

		info, err := de.Info() // 1 solo Lstat por entrada, no 2
		if err != nil {
			e.Err = err
			out = append(out, e)
			continue
		}
		e.Mode = info.Mode()
		e.Size = info.Size()
		e.Mtime = info.ModTime()

		if st, ok := info.Sys().(*syscall.Stat_t); ok {
			e.Nlink = st.Nlink
			e.Owner = ids.userName(st.Uid)
			e.Group = ids.groupName(st.Gid)
		}

		if info.Mode()&os.ModeSymlink != 0 {
			if target, err := os.Readlink(filepath.Join(dir, de.Name())); err == nil {
				e.Link = target
			}
		}
		out = append(out, e)
	}
	return out, nil
}

func renderShort(w io.Writer, entries []Entry) {
	for _, e := range entries {
		if e.Err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", e.Name, e.Err)
			continue
		}

		fmt.Fprintf(w, "%s\n", e.Name)
	}
}

func renderLong(w io.Writer, entries []Entry) {
	var wNlink, wOwner, wGroup, wSize int
	for _, e := range entries {
		wNlink = max(wNlink, len(strconv.FormatUint(e.Nlink, 10)))
		wOwner = max(wOwner, len(e.Owner))
		wGroup = max(wGroup, len(e.Group))
		wSize = max(wSize, len(strconv.FormatInt(e.Size, 10)))
	}

	for _, e := range entries {
		if e.Err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", e.Name, e.Err)
			continue
		}
		name := e.Name
		if e.Link != "" {
			name += " -> " + e.Link
		}
		fmt.Fprintf(w, "%s %*d %-*s %-*s %*d %s %s\n",
			modeString(e.Mode),
			wNlink, e.Nlink,
			wOwner, e.Owner,
			wGroup, e.Group,
			wSize, e.Size,
			e.Mtime.Format("Jan _2 15:04"),
			name,
		)
	}
}
