// +build ignore

package main

import (
	"context"
	_flag "flag"
	_fmt "fmt"
	_ioutil "io/ioutil"
	_log "log"
	"os"
	"os/signal"
	_filepath "path/filepath"
	_sort "sort"
	"strconv"
	_strings "strings"
	"syscall"
	_tabwriter "text/tabwriter"
	"time"
	golib_mageimport "github.com/coopnorge/mage/targets/golib"
	
)

func main() {
	// Use local types and functions in order to avoid name conflicts with additional magefiles.
	type arguments struct {
		Verbose       bool          // print out log statements
		List          bool          // print out a list of targets
		Help          bool          // print out help for a specific target
		Timeout       time.Duration // set a timeout to running the targets
		Args          []string      // args contain the non-flag command-line arguments
	}

	parseBool := func(env string) bool {
		val := os.Getenv(env)
		if val == "" {
			return false
		}		
		b, err := strconv.ParseBool(val)
		if err != nil {
			_log.Printf("warning: environment variable %s is not a valid bool value: %v", env, val)
			return false
		}
		return b
	}

	parseDuration := func(env string) time.Duration {
		val := os.Getenv(env)
		if val == "" {
			return 0
		}		
		d, err := time.ParseDuration(val)
		if err != nil {
			_log.Printf("warning: environment variable %s is not a valid duration value: %v", env, val)
			return 0
		}
		return d
	}
	args := arguments{}
	fs := _flag.FlagSet{}
	fs.SetOutput(os.Stdout)

	// default flag set with ExitOnError and auto generated PrintDefaults should be sufficient
	fs.BoolVar(&args.Verbose, "v", parseBool("MAGEFILE_VERBOSE"), "show verbose output when running targets")
	fs.BoolVar(&args.List, "l", parseBool("MAGEFILE_LIST"), "list targets for this binary")
	fs.BoolVar(&args.Help, "h", parseBool("MAGEFILE_HELP"), "print out help for a specific target")
	fs.DurationVar(&args.Timeout, "t", parseDuration("MAGEFILE_TIMEOUT"), "timeout in duration parsable format (e.g. 5m30s)")
	fs.Usage = func() {
		_fmt.Fprintf(os.Stdout, `
%s [options] [target]

Commands:
  -l    list targets in this binary
  -h    show this help

Options:
  -h    show description of a target
  -t <string>
        timeout in duration parsable format (e.g. 5m30s)
  -v    show verbose output when running targets
 `[1:], _filepath.Base(os.Args[0]))
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		// flag will have printed out an error already.
		return
	}
	args.Args = fs.Args()
	if args.Help && len(args.Args) == 0 {
		fs.Usage()
		return
	}
		
	// color is ANSI color type
	type color int

	// If you add/change/remove any items in this constant,
	// you will need to run "stringer -type=color" in this directory again.
	// NOTE: Please keep the list in an alphabetical order.
	const (
		black color = iota
		red
		green
		yellow
		blue
		magenta
		cyan
		white
		brightblack
		brightred
		brightgreen
		brightyellow
		brightblue
		brightmagenta
		brightcyan
		brightwhite
	)

	// AnsiColor are ANSI color codes for supported terminal colors.
	var ansiColor = map[color]string{
		black:         "\u001b[30m",
		red:           "\u001b[31m",
		green:         "\u001b[32m",
		yellow:        "\u001b[33m",
		blue:          "\u001b[34m",
		magenta:       "\u001b[35m",
		cyan:          "\u001b[36m",
		white:         "\u001b[37m",
		brightblack:   "\u001b[30;1m",
		brightred:     "\u001b[31;1m",
		brightgreen:   "\u001b[32;1m",
		brightyellow:  "\u001b[33;1m",
		brightblue:    "\u001b[34;1m",
		brightmagenta: "\u001b[35;1m",
		brightcyan:    "\u001b[36;1m",
		brightwhite:   "\u001b[37;1m",
	}
	
	const _color_name = "blackredgreenyellowbluemagentacyanwhitebrightblackbrightredbrightgreenbrightyellowbrightbluebrightmagentabrightcyanbrightwhite"

	var _color_index = [...]uint8{0, 5, 8, 13, 19, 23, 30, 34, 39, 50, 59, 70, 82, 92, 105, 115, 126}

	colorToLowerString := func (i color) string {
		if i < 0 || i >= color(len(_color_index)-1) {
			return "color(" + strconv.FormatInt(int64(i), 10) + ")"
		}
		return _color_name[_color_index[i]:_color_index[i+1]]
	}

	// ansiColorReset is an ANSI color code to reset the terminal color.
	const ansiColorReset = "\033[0m"

	// defaultTargetAnsiColor is a default ANSI color for colorizing targets.
	// It is set to Cyan as an arbitrary color, because it has a neutral meaning
	var defaultTargetAnsiColor = ansiColor[cyan]

	getAnsiColor := func(color string) (string, bool) {
		colorLower := _strings.ToLower(color)
		for k, v := range ansiColor {
			colorConstLower := colorToLowerString(k)
			if colorConstLower == colorLower {
				return v, true
			}
		}
		return "", false
	}

	// Terminals which  don't support color:
	// 	TERM=vt100
	// 	TERM=cygwin
	// 	TERM=xterm-mono
    var noColorTerms = map[string]bool{
		"vt100":      false,
		"cygwin":     false,
		"xterm-mono": false,
	}

	// terminalSupportsColor checks if the current console supports color output
	//
	// Supported:
	// 	linux, mac, or windows's ConEmu, Cmder, putty, git-bash.exe, pwsh.exe
	// Not supported:
	// 	windows cmd.exe, powerShell.exe
	terminalSupportsColor := func() bool {
		envTerm := os.Getenv("TERM")
		if _, ok := noColorTerms[envTerm]; ok {
			return false
		}
		return true
	}

	// enableColor reports whether the user has requested to enable a color output.
	enableColor := func() bool {
		b, _ := strconv.ParseBool(os.Getenv("MAGEFILE_ENABLE_COLOR"))
		return b
	}

	// targetColor returns the ANSI color which should be used to colorize targets.
	targetColor := func() string {
		s, exists := os.LookupEnv("MAGEFILE_TARGET_COLOR")
		if exists == true {
			if c, ok := getAnsiColor(s); ok == true {
				return c
			}
		}
		return defaultTargetAnsiColor
	}

	// store the color terminal variables, so that the detection isn't repeated for each target
	var enableColorValue = enableColor() && terminalSupportsColor()
	var targetColorValue = targetColor()

	printName := func(str string) string {
		if enableColorValue {
			return _fmt.Sprintf("%s%s%s", targetColorValue, str, ansiColorReset)
		} else {
			return str
		}
	}

	list := func() error {
		
		targets := map[string]string{
			"catalogInfo:changes": "returns the string true or false depending on the fact that the current branch contains changes compared to the main branch.",
			"catalogInfo:validate": "validates all terraform projects",
			"go:changes": "returns the string true or false depending on the fact that the current branch contains changes compared to the main branch.",
			"go:downloadModules": "download the go modules",
			"go:fetchConfigs": "syncs all configuration files into the repository.",
			"go:fetchGolangCILintConfig": "writes the golangci-lint configuration file provided path relative to root if it doesn't already exist.",
			"go:fix": "runs auto fixes on the Go source code in the repository.",
			"go:generate": "runs commands described by directives within existing files with the intent to generate Go code.",
			"go:lint": "checks all Go source code for issues.",
			"go:lintFix": "fixes found issues (if it's supported by the linters)",
			"go:test": "automates testing the packages named by the import paths, see also: go test.",
			"go:validate": "runs validation check on the Go source code in the repository.",
			"pallets:changes": "returns the string true or false depending on the fact that the current branch contains changes compared to the main branch.",
			"pallets:validate": "validates all terraform projects",
			"policyBotConfig:changes": "returns the string true or false depending on the fact that the current branch contains changes compared to the main branch.",
			"policyBotConfig:validate": "validates all terraform projects",
			"build": "runs validate",
			"clean": "removes validate and build output.",
			"fix": "fixes found issues (if it's supported by the linters)",
			"generate": "runs commands described by directives within existing files with the intent to generate Go code.",
			"validate": "runs validation check on the source code in the repository.",
		}

		keys := make([]string, 0, len(targets))
		for name := range targets {
			keys = append(keys, name)
		}
		_sort.Strings(keys)

		_fmt.Println("Targets:")
		w := _tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
		for _, name := range keys {
			_fmt.Fprintf(w, "  %v\t%v\n", printName(name), targets[name])
		}
		err := w.Flush()
		return err
	}

	var ctx context.Context
	ctxCancel := func(){}

	// by deferring in a closure, we let the cancel function get replaced
	// by the getContext function.
	defer func() {
		ctxCancel()
	}()

	getContext := func() (context.Context, func()) {
		if ctx == nil {
			if args.Timeout != 0 {
				ctx, ctxCancel = context.WithTimeout(context.Background(), args.Timeout)
			} else {
				ctx, ctxCancel = context.WithCancel(context.Background())
			}
		}

		return ctx, ctxCancel
	}

	runTarget := func(logger *_log.Logger, fn func(context.Context) error) interface{} {
		var err interface{}
		ctx, cancel := getContext()
		d := make(chan interface{})
		go func() {
			defer func() {
				err := recover()
				d <- err
			}()
			err := fn(ctx)
			d <- err
		}()
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT)
		select {
		case <-sigCh:
			logger.Println("cancelling mage targets, waiting up to 5 seconds for cleanup...")
			cancel()
			cleanupCh := time.After(5 * time.Second)

			select {
			// target exited by itself
			case err = <-d:
				return err
			// cleanup timeout exceeded
			case <-cleanupCh:
				return _fmt.Errorf("cleanup timeout exceeded")
			// second SIGINT received
			case <-sigCh:
				logger.Println("exiting mage")
				return _fmt.Errorf("exit forced")
			}
		case <-ctx.Done():
			cancel()
			e := ctx.Err()
			_fmt.Printf("ctx err: %v\n", e)
			return e
		case err = <-d:
			// we intentionally don't cancel the context here, because
			// the next target will need to run with the same context.
			return err
		}
	}
	// This is necessary in case there aren't any targets, to avoid an unused
	// variable error.
	_ = runTarget

	handleError := func(logger *_log.Logger, err interface{}) {
		if err != nil {
			logger.Printf("Error: %+v\n", err)
			type code interface {
				ExitStatus() int
			}
			if c, ok := err.(code); ok {
				os.Exit(c.ExitStatus())
			}
			os.Exit(1)
		}
	}
	_ = handleError

	// Set MAGEFILE_VERBOSE so mg.Verbose() reflects the flag value.
	if args.Verbose {
		os.Setenv("MAGEFILE_VERBOSE", "1")
	} else {
		os.Setenv("MAGEFILE_VERBOSE", "0")
	}

	_log.SetFlags(0)
	if !args.Verbose {
		_log.SetOutput(_ioutil.Discard)
	}
	logger := _log.New(os.Stderr, "", 0)
	if args.List {
		if err := list(); err != nil {
			_log.Println(err)
			os.Exit(1)
		}
		return
	}

	if args.Help {
		if len(args.Args) < 1 {
			logger.Println("no target specified")
			os.Exit(2)
		}
		switch _strings.ToLower(args.Args[0]) {
			case "cataloginfo:changes":
				_fmt.Println("Changes returns the string true or false depending on the fact that the current branch contains changes compared to the main branch.")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage cataloginfo:changes\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "cataloginfo:validate":
				_fmt.Println("Validate validates all terraform projects")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage cataloginfo:validate\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "go:changes":
				_fmt.Println("Changes returns the string true or false depending on the fact that the current branch contains changes compared to the main branch.")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage go:changes\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "go:downloadmodules":
				_fmt.Println("DownloadModules download the go modules")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage go:downloadmodules\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "go:fetchconfigs":
				_fmt.Println("FetchConfigs syncs all configuration files into the repository. Currently syncs GolangCiConfig to the specified path relative to the repository root.")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage go:fetchconfigs <where>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "go:fetchgolangcilintconfig":
				_fmt.Println("FetchGolangCILintConfig writes the golangci-lint configuration file provided path relative to root if it doesn't already exist.")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage go:fetchgolangcilintconfig <where>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "go:fix":
				_fmt.Println("Fix runs auto fixes on the Go source code in the repository.  For details see [Go.LintFix].")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage go:fix\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "go:generate":
				_fmt.Println("Generate runs commands described by directives within existing files with the intent to generate Go code. Those commands can run any process but the intent is to create or update Go source files  For details see [golang.Generate].")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage go:generate\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "go:lint":
				_fmt.Println("Lint checks all Go source code for issues.  See [golang.Lint] for details.")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage go:lint\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "go:lintfix":
				_fmt.Println("LintFix fixes found issues (if it's supported by the linters)  For details see [golang.LintFix].")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage go:lintfix\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "go:test":
				_fmt.Println("Test automates testing the packages named by the import paths, see also: go test.  For details see [golang.Test].")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage go:test\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "go:validate":
				_fmt.Println("Validate runs validation check on the Go source code in the repository.  See [Go.Test] and [Go.Lint] for details.")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage go:validate\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "pallets:changes":
				_fmt.Println("Changes returns the string true or false depending on the fact that the current branch contains changes compared to the main branch.")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage pallets:changes\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "pallets:validate":
				_fmt.Println("Validate validates all terraform projects")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage pallets:validate\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "policybotconfig:changes":
				_fmt.Println("Changes returns the string true or false depending on the fact that the current branch contains changes compared to the main branch.")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage policybotconfig:changes\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "policybotconfig:validate":
				_fmt.Println("Validate validates all terraform projects")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage policybotconfig:validate\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "build":
				_fmt.Println("Build runs validate  For details see [Validate].")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage build\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "clean":
				_fmt.Println("Clean removes validate and build output.  Deletes the [core.OutputDir].")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage clean\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "fix":
				_fmt.Println("Fix fixes found issues (if it's supported by the linters)  For details see [Go.Fix].")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage fix\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "generate":
				_fmt.Println("Generate runs commands described by directives within existing files with the intent to generate Go code. Those commands can run any process but the intent is to create or update Go source files  For details see [Go.Generate].")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage generate\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				case "validate":
				_fmt.Println("Validate runs validation check on the source code in the repository.  For details see [Go.Validate].")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage validate\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
				default:
				logger.Printf("Unknown target: %q\n", args.Args[0])
				os.Exit(2)
		}
	}
	if len(args.Args) < 1 {
		if err := list(); err != nil {
			logger.Println("Error:", err)
			os.Exit(1)
		}
		return
	}
	for x := 0; x < len(args.Args); {
		target := args.Args[x]
		x++

		// resolve aliases
		switch _strings.ToLower(target) {
		
		}

		switch _strings.ToLower(target) {
		
		
		
			
				case "cataloginfo:changes":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"CatalogInfo:Changes\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "CatalogInfo:Changes")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.CatalogInfo{}.Changes(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "cataloginfo:validate":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"CatalogInfo:Validate\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "CatalogInfo:Validate")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.CatalogInfo{}.Validate(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "go:changes":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Go:Changes\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Go:Changes")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Go{}.Changes(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "go:downloadmodules":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Go:DownloadModules\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Go:DownloadModules")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Go{}.DownloadModules(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "go:fetchconfigs":
					expected := x + 1
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Go:FetchConfigs\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Go:FetchConfigs")
					}
					
			arg0 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Go{}.FetchConfigs(ctx, arg0)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "go:fetchgolangcilintconfig":
					expected := x + 1
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Go:FetchGolangCILintConfig\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Go:FetchGolangCILintConfig")
					}
					
			arg0 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Go{}.FetchGolangCILintConfig(ctx, arg0)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "go:fix":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Go:Fix\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Go:Fix")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Go{}.Fix(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "go:generate":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Go:Generate\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Go:Generate")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Go{}.Generate(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "go:lint":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Go:Lint\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Go:Lint")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Go{}.Lint(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "go:lintfix":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Go:LintFix\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Go:LintFix")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Go{}.LintFix(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "go:test":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Go:Test\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Go:Test")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Go{}.Test(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "go:validate":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Go:Validate\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Go:Validate")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Go{}.Validate(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "pallets:changes":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Pallets:Changes\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Pallets:Changes")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Pallets{}.Changes(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "pallets:validate":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Pallets:Validate\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Pallets:Validate")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Pallets{}.Validate(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "policybotconfig:changes":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"PolicyBotConfig:Changes\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "PolicyBotConfig:Changes")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.PolicyBotConfig{}.Changes(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "policybotconfig:validate":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"PolicyBotConfig:Validate\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "PolicyBotConfig:Validate")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.PolicyBotConfig{}.Validate(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "build":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Build\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Build")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Build(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "clean":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Clean\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Clean")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Clean(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "fix":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Fix\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Fix")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Fix(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "generate":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Generate\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Generate")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Generate(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
				case "validate":
					expected := x + 0
					if expected > len(args.Args) {
						// note that expected and args at this point include the arg for the target itself
						// so we subtract 1 here to show the number of args without the target.
						logger.Printf("not enough arguments for target \"Validate\", expected %v, got %v\n", expected-1, len(args.Args)-1)
						os.Exit(2)
					}
					if args.Verbose {
						logger.Println("Running target:", "Validate")
					}
					
				wrapFn := func(ctx context.Context) error {
					return golib_mageimport.Validate(ctx)
				}
				ret := runTarget(logger, wrapFn)
					handleError(logger, ret)
		default:
			logger.Printf("Unknown target specified: %q\n", target)
			os.Exit(2)
		}
	}
}




