package minerrors

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

// Error - error type containing cause and the stack trace.
type Error struct {
	cause error
	stack []string
}

func (e Error) Error() string {
	return e.cause.Error()
}

// Cause - cause of the Error.
func (e Error) Cause() error {
	return e.cause
}

// NewError - return new Error type.
func NewError(e error) error {
	if e == nil {
		return nil
	}

	var stackStrs []string

	stack := make([]uintptr, 40)
	length := runtime.Callers(2, stack)
	stack = stack[:length]

	for _, pc := range stack {
		pc = pc - 1
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		name := fn.Name()
		if strings.HasPrefix(name, "runtime") {
			break
		}
		file = strings.TrimPrefix(file, rootPath+string(os.PathSeparator))
		stackStrs = append(stackStrs, fmt.Sprintf("%s:%d:%s()", file, line, name))
	}

	return &Error{e, stackStrs}
}

// Stack - returns top 'n' stack frames. if n is not specified all the frames are returned.
func (e Error) Stack(n ...int) string {
	length := len(e.stack)
	if len(n) != 0 {
		length = n[0]
	}
	return strings.Join(e.stack[:length], " ")
}

// Is - check if e and original are of the same type.
func Is(e error, original error) bool {
	if e == original {
		return true
	}

	if e, ok := e.(*Error); ok {
		return Is(e.cause, original)
	}

	if original, ok := original.(*Error); ok {
		return Is(e, original.cause)
	}

	return false
}

// XLError constructed at XL layer combining errors from the XL's underlying disks
type XLError struct {
	cause error
	stack []string
	// errors from the StorageAPI
	errs []error
}

func (xlErr XLError) Error() string {
	return xlErr.cause.Error()
}

// NewXLError - returns new XLError type.
func NewXLError(err error, errs ...error) error {
	if err == nil {
		return nil
	}

	var stackStrs []string

	stack := make([]uintptr, 40)
	length := runtime.Callers(2, stack)
	stack = stack[:length]

	for _, pc := range stack {
		pc = pc - 1
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		name := fn.Name()
		if strings.HasPrefix(name, "runtime") {
			break
		}
		file = strings.TrimPrefix(file, rootPath+string(os.PathSeparator))
		stackStrs = append(stackStrs, fmt.Sprintf("%s:%d:%s()", file, line, name))
	}

	return &XLError{err, stackStrs, errs}
}
