package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var rootPath string

// Init the package.
func Init() {
	// Root path is automatically determined from the calling function's source file location.
	// Catch the calling function's source file path.
	_, file, _, _ := runtime.Caller(1)
	// Save the directory alone.
	rootPath = filepath.Dir(file)
}

type Tracer interface {
	Trace() string
	Error() string
	JSON() []byte
}

type traceInfo struct {
	file string
	line int
	name string
}

// StorageError - error type containing cause and the stack trace.
type storageError struct {
	e     error
	trace []traceInfo
}

func (se storageError) Error() string {
	return se.e.Error()
}

// Trace - returns stack trace.
func (se storageError) Trace() string {
	var traceArr []string
	for _, info := range se.trace {
		traceArr = append(traceArr, fmt.Sprintf("%s:%d:%s",
			info.file, info.line, info.name))
	}
	return strings.Join(traceArr, " ")
}

func (se storageError) JSON() []byte {
	return nil
}

// NewStorageError - return new Error type.
func NewStorageError(e error) error {
	if e == nil {
		return nil
	}
	err := &storageError{}
	err.e = e

	stack := make([]uintptr, 40)
	length := runtime.Callers(2, stack)
	if length > len(stack) {
		length = len(stack)
	}
	stack = stack[:length]

	for _, pc := range stack {
		pc = pc - 1
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		name := fn.Name()
		file = strings.TrimPrefix(file, rootPath+string(os.PathSeparator))
		err.trace = append(err.trace, traceInfo{file, line, name})
	}

	return err
}

// XLError constructed at XL layer combining errors from the XL's underlying disks
type xlError struct {
	e     error
	trace []traceInfo
	// errors from the StorageAPI
	errs []error
}

func (xe xlError) Error() string {
	return xe.e.Error()
}

func (xe xlError) Trace() string {
	var traceArr []string
	for _, info := range xe.trace {
		traceArr = append(traceArr, fmt.Sprintf("%s:%d:%s",
			info.file, info.line, info.name))
	}
	return strings.Join(traceArr, " ")
}

// NewXLError - returns new XLError type.
func NewXLError(e error, errs ...error) error {
	if e == nil {
		return nil
	}

	err := &xlError{}
	err.e = e
	err.errs = errs

	stack := make([]uintptr, 40)
	length := runtime.Callers(2, stack)
	if length > len(stack) {
		length = len(stack)
	}
	stack = stack[:length]

	for _, pc := range stack {
		pc = pc - 1
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		name := fn.Name()
		file = strings.TrimPrefix(file, rootPath+string(os.PathSeparator))
		err.trace = append(err.trace, traceInfo{file, line, name})
	}

	return err
}
