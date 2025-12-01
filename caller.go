package logger

import (
	"runtime"
	"strings"
	"sync"
)

var (

	// qualified package name, cached at first use
	goLoggerPackage string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for reporting the line-number inside the go-logger package, used for testing.
	reportCallerInGoLoggerPackage bool

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	maximumCallerDepth  int = 25
	knownGoLoggerFrames int = 4
)

// getCaller retrieves the name of the first non-go-logger calling function
func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				goLoggerPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownGoLoggerFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	var prev runtime.Frame
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != goLoggerPackage {
			if reportCallerInGoLoggerPackage {
				return &prev
			}
			return &f //nolint:scopelint
		}
		prev = f
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// getPackageName returns the package path part of a fully-qualified function name.
//
// It is intended for strings in the format produced by runtime.FuncForPC.Name(),
// where the function/type suffix is appended using '.' characters.
//
// Example mappings:
//
//	"github.com/myorg/my-repo.something.func1" → "github.com/myorg/my-repo"
//	"myrepo.myfile.MyFunc"                     → "myrepo"
//	"package/subpackage.Function"              → "package/subpackage"
//	"simplepkg"                                → "simplepkg"
//
// The rule is: take everything from the start of the string up to the first '.'
// that appears after the last '/' (if any). If there is no such '.', the whole
// string is treated as the package name.
func getPackageName(f string) string {
	// Start searching for '.' just after the last '/' (or at 0 if there is no '/').
	start := strings.LastIndexByte(f, '/')
	if start == -1 {
		start = 0
	} else {
		start++ // move past the '/'
	}

	// Find the first '.' after that point.
	if i := strings.IndexByte(f[start:], '.'); i != -1 {
		return f[:start+i]
	}

	// No '.' after the last '/', so treat the entire string as the package name.
	return f
}
