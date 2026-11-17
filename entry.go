package main

import (
	"os"
	"time"
)

type entry struct {
	name  string
	mode  os.FileMode
	nlink uint64
	owner string
	group string
	size  int64
	mtime time.Time
	link  string
	err   error
}
