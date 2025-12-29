// Command qjs provides an interactive JavaScript REPL using QuickJS-ng.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Gaurav-Gosain/quickjs"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/chzyer/readline"
)

const version = "0.1.0"

// Styles
var (
	primaryColor   = lipgloss.Color("#7C3AED")
	secondaryColor = lipgloss.Color("#10B981")
	errorColor     = lipgloss.Color("#EF4444")
	warningColor   = lipgloss.Color("#F59E0B")
	infoColor      = lipgloss.Color("#3B82F6")
	dimColor       = lipgloss.Color("#6B7280")
	stringColor    = lipgloss.Color("#10B981")
	numberColor    = lipgloss.Color("#3B82F6")
	boolColor      = lipgloss.Color("#F59E0B")
	nullColor      = lipgloss.Color("#6B7280")

	logoStyle         = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	promptStyle       = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	continuationStyle = lipgloss.NewStyle().Foreground(dimColor)
	errorStyle        = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	errorMsgStyle     = lipgloss.NewStyle().Foreground(errorColor)
	successStyle      = lipgloss.NewStyle().Foreground(secondaryColor)
	infoStyle         = lipgloss.NewStyle().Foreground(infoColor)
	dimStyle          = lipgloss.NewStyle().Foreground(dimColor)
	stringStyle       = lipgloss.NewStyle().Foreground(stringColor)
	numberStyle       = lipgloss.NewStyle().Foreground(numberColor)
	boolStyle         = lipgloss.NewStyle().Foreground(boolColor)
	nullStyle         = lipgloss.NewStyle().Foreground(nullColor)
	cmdStyle          = lipgloss.NewStyle().Foreground(warningColor)
	titleStyle        = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Underline(true)
	resultStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#A78BFA"))
)

// Syntax highlighter
var (
	jsLexer     chroma.Lexer
	chromaStyle *chroma.Style
	formatter   chroma.Formatter
)

func initSyntaxHighlighter() {
	jsLexer = lexers.Get("javascript")
	if jsLexer == nil {
		jsLexer = lexers.Fallback
	}
	jsLexer = chroma.Coalesce(jsLexer)
	chromaStyle = styles.Get("dracula")
	if chromaStyle == nil {
		chromaStyle = styles.Fallback
	}
	formatter = formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}
}

func highlightCode(code string) string {
	if jsLexer == nil {
		return code
	}
	var buf bytes.Buffer
	iterator, err := jsLexer.Tokenise(nil, code)
	if err != nil {
		return code
	}
	if err := formatter.Format(&buf, chromaStyle, iterator); err != nil {
		return code
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

// REPL state
type replState struct {
	ctx         *quickjs.Context
	rt          *quickjs.Runtime
	rl          *readline.Instance
	showTiming  bool
	evalCount   int
	multiline   strings.Builder
	inMultiline bool
	startTime   time.Time
}

func main() {
	os.Exit(run())
}

func run() int {
	evalCode := flag.String("e", "", "evaluate code and exit")
	showVersion := flag.Bool("version", false, "show version")
	showHelp := flag.Bool("help", false, "show help")
	timing := flag.Bool("timing", false, "show execution time")
	flag.Parse()

	initSyntaxHighlighter()

	if *showVersion {
		printVersion()
		return 0
	}

	if *showHelp {
		printUsage()
		return 0
	}

	rt, err := quickjs.NewRuntime()
	if err != nil {
		fmt.Fprintln(os.Stderr, errorStyle.Render("Error:")+" failed to create runtime:", err)
		return 1
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		fmt.Fprintln(os.Stderr, errorStyle.Render("Error:")+" failed to create context:", err)
		return 1
	}
	defer ctx.Close()

	state := &replState{
		ctx:        ctx,
		rt:         rt,
		showTiming: *timing,
		startTime:  time.Now(),
	}

	if *evalCode != "" {
		result, duration, err := state.eval(*evalCode)
		if err != nil {
			printError(err)
			return 1
		}
		if !result.IsUndefined() {
			printValue(result)
		}
		if state.showTiming {
			printTiming(duration)
		}
		return 0
	}

	args := flag.Args()
	if len(args) > 0 {
		for _, filename := range args {
			if err := state.runFile(filename); err != nil {
				printError(err)
				return 1
			}
		}
		return 0
	}

	state.runREPL()
	return 0
}

func printVersion() {
	fmt.Println(logoStyle.Render("qjs") + dimStyle.Render(" v"+version))
	fmt.Println(dimStyle.Render("A modern JavaScript runtime powered by QuickJS-ng + WebAssembly"))
	fmt.Println(dimStyle.Render(fmt.Sprintf("Go %s, %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)))
}

func printUsage() {
	fmt.Println()
	fmt.Println(titleStyle.Render("qjs - QuickJS-ng JavaScript Runtime"))
	fmt.Println()

	fmt.Println(logoStyle.Render("USAGE"))
	fmt.Println("  qjs [options] [script.js] [arguments...]")
	fmt.Println()

	fmt.Println(logoStyle.Render("OPTIONS"))
	fmt.Println("  " + cmdStyle.Render("-e <code>") + "      Evaluate JavaScript code and exit")
	fmt.Println("  " + cmdStyle.Render("-timing") + "        Show execution time")
	fmt.Println("  " + cmdStyle.Render("-version") + "       Show version information")
	fmt.Println("  " + cmdStyle.Render("-help") + "          Show this help message")
	fmt.Println()

	fmt.Println(logoStyle.Render("REPL COMMANDS"))
	cmds := []struct{ cmd, desc string }{
		{".help", "Show help for REPL commands"},
		{".exit", "Exit the REPL"},
		{".clear", "Clear the screen"},
		{".examples", "Show JavaScript examples"},
		{".bench", "Run performance benchmarks"},
		{".timing", "Toggle timing display"},
		{".load <file>", "Load and execute a file"},
		{".info", "Show runtime information"},
		{".gc", "Trigger garbage collection"},
		{".reset", "Reset the context"},
	}
	for _, c := range cmds {
		fmt.Printf("  %s  %s\n", cmdStyle.Render(fmt.Sprintf("%-14s", c.cmd)), dimStyle.Render(c.desc))
	}
	fmt.Println()
}

func (s *replState) runFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filename, err)
	}

	start := time.Now()
	_, err = s.ctx.EvalFile(string(data), filename)
	duration := time.Since(start)

	if err != nil {
		return fmt.Errorf("%s: %w", filename, err)
	}

	if s.showTiming {
		printTiming(duration)
	}
	return nil
}

func (s *replState) eval(code string) (quickjs.Value, time.Duration, error) {
	start := time.Now()
	result, err := s.ctx.Eval(code)
	duration := time.Since(start)
	return result, duration, err
}

func (s *replState) runREPL() {
	historyFile := ""
	if home, err := os.UserHomeDir(); err == nil {
		historyFile = filepath.Join(home, ".qjs_history")
	}

	// Completions
	completions := []string{
		// Keywords
		"var", "let", "const", "function", "return", "if", "else", "for", "while",
		"do", "switch", "case", "default", "break", "continue", "try", "catch",
		"finally", "throw", "new", "delete", "typeof", "instanceof", "in", "this",
		"true", "false", "null", "undefined", "class", "extends", "super",
		"async", "await", "yield", "import", "export", "from", "as",
		// Globals
		"console", "Math", "JSON", "Object", "Array", "String", "Number", "Boolean",
		"Date", "RegExp", "Error", "Promise", "Map", "Set", "WeakMap", "WeakSet",
		"Symbol", "Proxy", "Reflect", "BigInt", "ArrayBuffer", "DataView",
		"parseInt", "parseFloat", "isNaN", "isFinite", "print",
		// Math
		"Math.abs", "Math.ceil", "Math.floor", "Math.round", "Math.sqrt", "Math.pow",
		"Math.min", "Math.max", "Math.random", "Math.sin", "Math.cos", "Math.PI",
		// JSON
		"JSON.parse", "JSON.stringify",
		// Object
		"Object.keys", "Object.values", "Object.entries", "Object.assign",
		"Object.freeze", "Object.seal", "Object.create",
		// Array
		"Array.isArray", "Array.from", "Array.of",
		// Promise
		"Promise.resolve", "Promise.reject", "Promise.all", "Promise.race",
		// Commands
		".help", ".exit", ".clear", ".examples", ".bench", ".timing", ".load",
		".info", ".gc", ".reset", ".history",
	}

	completer := readline.NewPrefixCompleter()
	for _, item := range completions {
		completer.Children = append(completer.Children, readline.PcItem(item))
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:            s.getPrompt(false),
		HistoryFile:       historyFile,
		HistoryLimit:      1000,
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, errorStyle.Render("Error:")+" failed to initialize readline:", err)
		os.Exit(1)
	}
	defer rl.Close()
	s.rl = rl

	printBanner()

	for {
		if s.inMultiline {
			rl.SetPrompt(s.getPrompt(true))
		} else {
			rl.SetPrompt(s.getPrompt(false))
		}

		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				if s.inMultiline {
					s.multiline.Reset()
					s.inMultiline = false
					fmt.Println()
					continue
				}
				continue
			}
			if err == io.EOF {
				fmt.Println()
				fmt.Println(dimStyle.Render("Goodbye!"))
				break
			}
			continue
		}

		if !s.inMultiline && strings.HasPrefix(line, ".") {
			s.handleCommand(line)
			continue
		}

		if s.inMultiline {
			if line == "" {
				code := s.multiline.String()
				s.multiline.Reset()
				s.inMultiline = false
				s.evalAndPrint(code)
			} else {
				s.multiline.WriteString(line)
				s.multiline.WriteString("\n")
			}
			continue
		}

		if line == "exit" || line == "quit" {
			fmt.Println(dimStyle.Render("Goodbye!"))
			break
		}

		if strings.HasSuffix(line, "\\") {
			s.multiline.WriteString(strings.TrimSuffix(line, "\\"))
			s.multiline.WriteString("\n")
			s.inMultiline = true
			continue
		}

		if needsContinuation(line) {
			s.multiline.WriteString(line)
			s.multiline.WriteString("\n")
			s.inMultiline = true
			continue
		}

		s.evalAndPrint(line)
	}
}

func (s *replState) getPrompt(continuation bool) string {
	if continuation {
		return continuationStyle.Render("... ")
	}
	return promptStyle.Render("qjs") + dimStyle.Render(" > ")
}

func printBanner() {
	logo := `
      ┌─┐  ┬┌─┐
      │─┼┐ │└─┐
      └─┘└└┘└─┘`

	fmt.Println(logoStyle.Render(logo))
	fmt.Println()
	fmt.Println(dimStyle.Render("  QuickJS-ng JavaScript Runtime v" + version))
	fmt.Println(dimStyle.Render("  Type ") + cmdStyle.Render(".help") + dimStyle.Render(" for commands, ") + cmdStyle.Render(".examples") + dimStyle.Render(" for examples"))
	fmt.Println()
}

func (s *replState) handleCommand(line string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return
	}

	cmd := strings.ToLower(parts[0])
	args := parts[1:]

	switch cmd {
	case ".help", ".h", ".?":
		s.cmdHelp()
	case ".exit", ".quit", ".q":
		fmt.Println(dimStyle.Render("Goodbye!"))
		os.Exit(0)
	case ".clear", ".cls":
		fmt.Print("\033[H\033[2J")
	case ".history":
		s.cmdHistory(args)
	case ".timing", ".time":
		s.showTiming = !s.showTiming
		if s.showTiming {
			fmt.Println(successStyle.Render("✓") + " Timing enabled")
		} else {
			fmt.Println(infoStyle.Render("○") + " Timing disabled")
		}
	case ".load", ".l":
		s.cmdLoad(args)
	case ".info", ".i":
		s.cmdInfo()
	case ".gc":
		s.cmdGC()
	case ".reset":
		s.cmdReset()
	case ".examples", ".ex":
		s.cmdExamples()
	case ".bench", ".benchmark":
		s.cmdBenchmark()
	default:
		fmt.Println(errorStyle.Render("Unknown command:") + " " + cmd)
		fmt.Println(dimStyle.Render("Type .help for available commands"))
	}
}

func (s *replState) cmdHelp() {
	fmt.Println()
	fmt.Println(titleStyle.Render("Commands"))
	fmt.Println()

	cmds := []struct{ cmd, desc string }{
		{".help", "Show this help message"},
		{".exit", "Exit the REPL"},
		{".clear", "Clear the screen"},
		{".history [n]", "Show last n commands (default: 20)"},
		{".examples", "Show JavaScript examples"},
		{".bench", "Run performance benchmarks"},
		{".timing", "Toggle execution timing"},
		{".load <file>", "Load and execute a JavaScript file"},
		{".info", "Show runtime information"},
		{".gc", "Trigger garbage collection"},
		{".reset", "Reset context (clear all variables)"},
	}

	for _, c := range cmds {
		fmt.Printf("  %s  %s\n", cmdStyle.Render(fmt.Sprintf("%-14s", c.cmd)), dimStyle.Render(c.desc))
	}

	fmt.Println()
	fmt.Println(titleStyle.Render("Keyboard"))
	fmt.Println()
	shortcuts := []struct{ key, desc string }{
		{"↑/↓", "Navigate history"},
		{"Ctrl+R", "Search history"},
		{"Ctrl+C", "Cancel input"},
		{"Ctrl+D", "Exit REPL"},
		{"Tab", "Autocomplete"},
	}
	for _, sc := range shortcuts {
		fmt.Printf("  %s  %s\n", cmdStyle.Render(fmt.Sprintf("%-12s", sc.key)), dimStyle.Render(sc.desc))
	}
	fmt.Println()
}

func (s *replState) cmdExamples() {
	fmt.Println()
	fmt.Println(titleStyle.Render("JavaScript Examples (ES2023+)"))
	fmt.Println()

	examples := []struct {
		title string
		code  string
		desc  string
	}{
		{
			"Arrow Functions",
			"const add = (a, b) => a + b; add(2, 3)",
			"Concise function syntax",
		},
		{
			"Template Literals",
			"`Hello, ${'World'}!`",
			"String interpolation",
		},
		{
			"Destructuring",
			"const {x, y} = {x: 1, y: 2}; x + y",
			"Extract values from objects",
		},
		{
			"Spread Operator",
			"[...[1, 2], ...[3, 4]]",
			"Spread arrays",
		},
		{
			"Classes",
			"class Animal { constructor(name) { this.name = name; } speak() { return this.name; } } new Animal('Dog').speak()",
			"ES6 classes",
		},
		{
			"Promises",
			"Promise.resolve(42).then(x => x * 2)",
			"Promise handling",
		},
		{
			"Map & Set",
			"new Set([1, 2, 2, 3]).size",
			"Built-in collections",
		},
		{
			"BigInt",
			"BigInt(Number.MAX_SAFE_INTEGER) + 1n",
			"Arbitrary precision integers",
		},
		{
			"Optional Chaining",
			"const obj = {a: {b: 1}}; obj?.a?.b ?? 'default'",
			"Safe property access",
		},
		{
			"Nullish Coalescing",
			"null ?? 'fallback'",
			"Default for null/undefined",
		},
		{
			"Array Methods",
			"[1, 2, 3].map(x => x * 2).filter(x => x > 2)",
			"Functional array processing",
		},
		{
			"Object Methods",
			"Object.entries({a: 1, b: 2}).map(([k, v]) => k + v)",
			"Object iteration",
		},
		{
			"Symbol",
			"const sym = Symbol('desc'); typeof sym",
			"Unique identifiers",
		},
		{
			"Proxy",
			"new Proxy({x: 1}, {get: (t, p) => t[p] * 2}).x",
			"Object interception",
		},
		{
			"Reflect",
			"Reflect.has({a: 1}, 'a')",
			"Meta-programming",
		},
		{
			"Recursion",
			"const fib = n => n <= 1 ? n : fib(n-1) + fib(n-2); fib(10)",
			"Recursive Fibonacci",
		},
		{
			"Closures",
			"const counter = () => { let n = 0; return () => ++n; }; const c = counter(); [c(), c(), c()]",
			"Function closures",
		},
		{
			"Higher-Order",
			"const twice = (f, x) => f(f(x)); twice(x => x * 2, 3)",
			"Functions as values",
		},
		{
			"JSON",
			`JSON.parse('{"x": 1}').x`,
			"Parse JSON",
		},
		{
			"Error Handling",
			"try { throw new Error('oops'); } catch(e) { e.message }",
			"Exception handling",
		},
	}

	for i, ex := range examples {
		fmt.Printf("  %s%d. %s%s\n", dimStyle.Render(""), i+1, titleStyle.Render(ex.title), "")
		fmt.Printf("     %s\n", dimStyle.Render(ex.desc))
		fmt.Printf("     %s\n", highlightCode(ex.code))

		result, _, err := s.eval(ex.code)
		if err != nil {
			fmt.Printf("     %s %s\n", errorStyle.Render("→"), errorMsgStyle.Render(err.Error()))
		} else {
			fmt.Printf("     %s %s\n", resultStyle.Render("→"), formatResult(result))
		}
		fmt.Println()
	}
}

func (s *replState) cmdBenchmark() {
	fmt.Println()
	fmt.Println(titleStyle.Render("Performance Benchmarks"))
	fmt.Println()

	benchmarks := []struct {
		name string
		code string
	}{
		{
			"Empty loop (1M iterations)",
			"(() => { let i = 0; while (i < 1000000) { i++; } return i; })()",
		},
		{
			"Arithmetic (100K ops)",
			"(() => { let sum = 0; for (let i = 0; i < 100000; i++) { sum += i * 2 + 1; } return sum; })()",
		},
		{
			"String concatenation (10K)",
			`(() => { let s = ""; for (let i = 0; i < 10000; i++) { s += "x"; } return s.length; })()`,
		},
		{
			"Array push (10K)",
			"(() => { const arr = []; for (let i = 0; i < 10000; i++) { arr.push(i); } return arr.length; })()",
		},
		{
			"Object property access (100K)",
			"(() => { const obj = {x: 1, y: 2, z: 3}; let sum = 0; for (let i = 0; i < 100000; i++) { sum += obj.x + obj.y + obj.z; } return sum; })()",
		},
		{
			"Function calls (100K)",
			"(() => { const add = (a, b) => a + b; let sum = 0; for (let i = 0; i < 100000; i++) { sum = add(sum, 1); } return sum; })()",
		},
		{
			"Fibonacci(25)",
			"(() => { const fib = n => n <= 1 ? n : fib(n-1) + fib(n-2); return fib(25); })()",
		},
		{
			"Prime sieve (1000)",
			`(() => {
				const sieve = n => {
					const primes = [];
					const isPrime = Array(n + 1).fill(true);
					isPrime[0] = isPrime[1] = false;
					for (let i = 2; i <= n; i++) {
						if (isPrime[i]) {
							primes.push(i);
							for (let j = i * i; j <= n; j += i) isPrime[j] = false;
						}
					}
					return primes.length;
				};
				return sieve(1000);
			})()`,
		},
	}

	fmt.Println(dimStyle.Render("  Running benchmarks..."))
	fmt.Println()

	var totalTime time.Duration

	for _, bench := range benchmarks {
		start := time.Now()
		result, _, err := s.eval(bench.code)
		duration := time.Since(start)
		totalTime += duration

		if err != nil {
			fmt.Printf("  %s %s\n", errorStyle.Render("✗"), bench.name)
			fmt.Printf("    %s\n", errorMsgStyle.Render(err.Error()))
		} else {
			var timeStyle lipgloss.Style
			switch {
			case duration < 10*time.Millisecond:
				timeStyle = successStyle
			case duration < 100*time.Millisecond:
				timeStyle = lipgloss.NewStyle().Foreground(warningColor)
			default:
				timeStyle = errorStyle
			}

			fmt.Printf("  %s %s\n", successStyle.Render("✓"), bench.name)
			fmt.Printf("    Result: %s  Time: %s\n",
				dimStyle.Render(formatResultShort(result)),
				timeStyle.Render(duration.String()))
		}
		fmt.Println()
	}

	fmt.Println(dimStyle.Render("  ─────────────────────────────────────"))
	fmt.Printf("  Total time: %s\n", infoStyle.Render(totalTime.String()))
	fmt.Println()
}

func (s *replState) cmdHistory(args []string) {
	n := 20
	if len(args) > 0 {
		if parsed, err := strconv.Atoi(args[0]); err == nil && parsed > 0 {
			n = parsed
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		histFile := filepath.Join(home, ".qjs_history")
		data, err := os.ReadFile(histFile)
		if err != nil {
			fmt.Println(dimStyle.Render("No history"))
			return
		}

		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		start := 0
		if len(lines) > n {
			start = len(lines) - n
		}

		fmt.Println()
		fmt.Println(titleStyle.Render("History"))
		fmt.Println()
		for i, line := range lines[start:] {
			num := start + i + 1
			fmt.Printf("  %s  %s\n", dimStyle.Render(fmt.Sprintf("%4d", num)), highlightCode(line))
		}
		fmt.Println()
	}
}

func (s *replState) cmdLoad(args []string) {
	if len(args) == 0 {
		fmt.Println(errorStyle.Render("Usage:") + " .load <filename>")
		return
	}

	filename := args[0]
	fmt.Println(dimStyle.Render("Loading " + filename + "..."))

	if err := s.runFile(filename); err != nil {
		printError(err)
	} else {
		fmt.Println(successStyle.Render("✓") + " Loaded successfully")
	}
}

func (s *replState) cmdInfo() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	uptime := time.Since(s.startTime)

	fmt.Println()
	fmt.Println(titleStyle.Render("Runtime Information"))
	fmt.Println()

	info := []struct{ label, value string }{
		{"Version", version},
		{"Engine", "QuickJS-ng (ES2023+)"},
		{"Go Version", runtime.Version()},
		{"OS/Arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)},
		{"Go Heap", fmt.Sprintf("%.2f MB", float64(memStats.HeapAlloc)/1024/1024)},
		{"Go Sys", fmt.Sprintf("%.2f MB", float64(memStats.Sys)/1024/1024)},
		{"GC Runs", fmt.Sprintf("%d", memStats.NumGC)},
		{"Evaluations", fmt.Sprintf("%d", s.evalCount)},
		{"Uptime", uptime.Round(time.Second).String()},
	}

	for _, i := range info {
		fmt.Printf("  %s  %s\n", dimStyle.Render(fmt.Sprintf("%-14s", i.label)), i.value)
	}
	fmt.Println()
}

func (s *replState) cmdGC() {
	fmt.Println(dimStyle.Render("Running garbage collection..."))

	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	start := time.Now()
	if err := s.rt.RunGC(); err != nil {
		printError(err)
		return
	}
	runtime.GC()
	duration := time.Since(start)

	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	freed := max(int64(before.HeapAlloc)-int64(after.HeapAlloc), 0)

	fmt.Println(successStyle.Render("✓") + fmt.Sprintf(" GC completed in %v (freed ~%.2f KB)", duration, float64(freed)/1024))
}

func (s *replState) cmdReset() {
	fmt.Println(dimStyle.Render("Resetting context..."))

	s.ctx.Close()

	ctx, err := s.rt.NewContext()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	s.ctx = ctx
	s.evalCount = 0

	runtime.GC()

	fmt.Println(successStyle.Render("✓") + " Context reset")
}

func (s *replState) evalAndPrint(code string) {
	code = strings.TrimSpace(code)
	if code == "" {
		return
	}

	s.evalCount++

	result, duration, err := s.eval(code)
	if err != nil {
		printError(err)
		return
	}

	if !result.IsUndefined() {
		printValue(result)
	}

	if s.showTiming {
		printTiming(duration)
	}
}

func formatResult(v quickjs.Value) string {
	str := v.String()
	switch {
	case v.IsNull():
		return nullStyle.Render(str)
	case v.IsBool():
		return boolStyle.Render(str)
	case v.IsNumber():
		return numberStyle.Render(str)
	case v.IsString():
		return stringStyle.Render("\"" + str + "\"")
	case v.IsFunction():
		return dimStyle.Render("[Function]")
	case v.IsError():
		return errorStyle.Render(str)
	case v.IsBigInt():
		return numberStyle.Render(str + "n")
	case v.IsSymbol():
		return dimStyle.Render(str)
	default:
		return str
	}
}

func formatResultShort(v quickjs.Value) string {
	str := v.String()
	if len(str) > 50 {
		str = str[:47] + "..."
	}
	return str
}

func printValue(v quickjs.Value) {
	fmt.Println(formatResult(v))
}

func printError(err error) {
	fmt.Println()
	fmt.Println(errorStyle.Render("Error"))
	fmt.Println(errorMsgStyle.Render(err.Error()))
	fmt.Println()
}

func printTiming(duration time.Duration) {
	var style lipgloss.Style
	switch {
	case duration < 10*time.Millisecond:
		style = successStyle
	case duration < 100*time.Millisecond:
		style = lipgloss.NewStyle().Foreground(warningColor)
	default:
		style = errorStyle
	}
	fmt.Println(style.Render(fmt.Sprintf("⏱  %v", duration)))
}

func needsContinuation(line string) bool {
	opens := 0
	inString := false
	var stringChar byte

	for i := 0; i < len(line); i++ {
		ch := line[i]
		if inString {
			if ch == stringChar && (i == 0 || line[i-1] != '\\') {
				inString = false
			}
			continue
		}
		switch ch {
		case '"', '\'', '`':
			inString = true
			stringChar = ch
		case '{', '(', '[':
			opens++
		case '}', ')', ']':
			opens--
		}
	}
	return opens > 0 || inString
}
