package main

import (
	"os"
	"time"
)

type Entry struct {
	Name  string
	Mode  os.FileMode
	Nlink uint64
	Owner string
	Group string
	Size  int64
	Mtime time.Time
	Link  string
	Err   error
}
