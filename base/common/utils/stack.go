package utils

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"runtime"
	"strings"
)

// Frame holds information about a single frame in the call stack.
type Frame struct {
	// Unique, package path-qualified name for the function of this call
	// frame.
	Function string

	// File and line number of our location in the frame.
	//
	// Note that the line number does not refer to where the function was
	// defined but where in the function the next call was made.
	File string
	Line int
}

func (f Frame) String() string {
	// This takes the following forms.
	//  (path/to/file.go)
	//  (path/to/file.go:42)
	//  path/to/package.MyFunction
	//  path/to/package.MyFunction (path/to/file.go)
	//  path/to/package.MyFunction (path/to/file.go:42)

	var sb strings.Builder
	sb.WriteString(f.Function)
	if len(f.File) > 0 {
		if sb.Len() > 0 {
			sb.WriteRune(' ')
		}
		_, err := fmt.Fprintf(&sb, "(%v", f.File)
		PanicAndPrintIfNotNil(err)
		if f.Line > 0 {
			_, err = fmt.Fprintf(&sb, ":%d", f.Line)
			PanicAndPrintIfNotNil(err)
		}
		sb.WriteRune(')')
	}

	if sb.Len() == 0 {
		return "unknown"
	}

	return sb.String()
}

const _defaultCallersDepth = 8

// Stack is a stack of call frames.
//
// Formatted with %v, the output is in a single-line, in the form,
//
//	foo/bar.Baz() (path/to/foo.go:42); bar/baz.Qux() (bar/baz/qux.go:12); ...
//
// Formatted with %+v, the output is in the form,
//
//	foo/bar.Baz()
//		path/to/foo.go:42
//	bar/baz.Qux()
//		bar/baz/qux.go:12
type Stack []Frame

// Returns a single-line, semi-colon representation of a Stack. For a
// multi-line representation, use %+v.
func (fs Stack) String() string {
	items := make([]string, len(fs))
	for i, f := range fs {
		items[i] = f.String()
	}
	return strings.Join(items, "; ")
}

// Format implements fmt.Formatter to handle "%+v".
func (fs Stack) Format(w fmt.State, c rune) {
	if !w.Flag('+') {
		// Without %+v, fall back to String().
		_, err := io.WriteString(w, fs.String())
		PanicAndPrintIfNotNil(err)
		return
	}

	for _, f := range fs {
		_, err := fmt.Fprintln(w, f.Function)
		PanicAndPrintIfNotNil(err)

		_, err = fmt.Fprintf(w, "\t%v:%v\n", f.File, f.Line)
		PanicAndPrintIfNotNil(err)
	}
}

// CallerName returns the name of the first caller in this stack that isn't
// owned by the Fx library.
func (fs Stack) CallerName() string {
	for _, f := range fs {
		return f.Function
	}
	return "n/a"
}

func GetParentCallerMethodName() string {
	stack := CallerStack(1, 2)
	callerName := stack.CallerName()
	index := strings.LastIndex(callerName, ".")
	if index == -1 || index == len(callerName)-1 {
		return callerName
	}
	return callerName[index+1:]
}

// CallerStack returns the call stack for the calling function, up to depth frames
// deep, skipping the provided number of frames, not including Callers itself.
//
// If zero, depth defaults to 8.
func CallerStack(skip, depth int) Stack {
	if depth <= 0 {
		depth = _defaultCallersDepth
	}

	pcs := make([]uintptr, depth)

	// +2 to skip this frame and runtime.Callers.
	n := runtime.Callers(skip+2, pcs)
	pcs = pcs[:n] // truncate to number of frames actually read

	result := make([]Frame, 0, n)
	frames := runtime.CallersFrames(pcs)
	for f, more := frames.Next(); more; f, more = frames.Next() {
		result = append(result, Frame{
			Function: sanitize(f.Function),
			File:     f.File,
			Line:     f.Line,
		})
	}
	return result
}

// Match from beginning of the line until the first `vendor/` (non-greedy)
var vendorRe = regexp.MustCompile("^.*?/vendor/")

// sanitize makes the function name suitable for logging display. It removes
// url-encoded elements from the `dot.git` package names and shortens the
// vendored paths.
func sanitize(function string) string {
	// Use the stdlib to un-escape any package import paths which can happen
	// in the case of the "dot-git" postfix. Seems like a bug in stdlib =/
	if unescaped, err := url.QueryUnescape(function); err == nil {
		function = unescaped
	}

	// strip everything prior to the vendor
	return vendorRe.ReplaceAllString(function, "vendor/")
}
