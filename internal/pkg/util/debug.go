package util

import (
	"os"
	"runtime"
	"strings"
)

var filePrefix = String.Concat(os.Getenv("PWD"), "/")

// Current :
var Current = debugBlock{1}

// Prev :
var Prev = debugBlock{2}

// PPrev :
var PPrev = debugBlock{3}

type debugBlock struct {
	skip int
}

// File : like __FILE__
func (d debugBlock) File() string {
	_, file, _, _ := runtime.Caller(d.skip)
	return strings.TrimPrefix(file, filePrefix)
}

// Line : like __LINE__
func (d debugBlock) Line() int {
	_, _, line, _ := runtime.Caller(d.skip)
	return line
}

// Func : like __func__
func (d debugBlock) Func() string {
	pc, _, _, _ := runtime.Caller(d.skip)
	return strings.TrimPrefix(strings.TrimSuffix(runtime.FuncForPC(pc).Name(), ".0"), "")
}
