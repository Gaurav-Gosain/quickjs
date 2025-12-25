//go:build go1.18 && !race
// +build go1.18,!race

package quickjs

import (
	"testing"
	"unicode/utf8"
)

// Note: Fuzz tests are disabled when running with -race because
// creating many runtimes is slow with race detection enabled.

// FuzzEval tests that arbitrary input doesn't cause panics
func FuzzEval(f *testing.F) {
	// Seed corpus with various JavaScript snippets
	seeds := []string{
		"",
		"1",
		"1 + 2",
		"null",
		"undefined",
		"true",
		"false",
		`"hello"`,
		"[]",
		"{}",
		"function(){}",
		"() => {}",
		"class A {}",
		"var x = 1",
		"let x = 1",
		"const x = 1",
		"if (true) {}",
		"for (;;) break",
		"while (false) {}",
		"try {} catch(e) {}",
		"throw new Error()",
		"new Date()",
		"Math.PI",
		"JSON.stringify({})",
		"Object.keys({})",
		"Array.isArray([])",
		"Promise.resolve(1)",
		"async function f() {}",
		"function* g() {}",
		"Symbol('x')",
		"BigInt(1)",
		"/regex/",
		"0x1234",
		"0b1010",
		"0o777",
		"1e10",
		"1.5e-10",
		"'\\n\\t\\r'",
		"`template`",
		"a?.b?.c",
		"a ?? b",
		"...[]",
		"({...{}})",
		"[...[]]",
		"{ get x() {}, set x(v) {} }",
		"class A extends B {}",
		"import.meta",
		"new.target",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, code string) {
		// Skip invalid UTF-8 as JavaScript expects valid strings
		if !utf8.ValidString(code) {
			return
		}

		rt, err := NewRuntime()
		if err != nil {
			return // Runtime creation issues are acceptable
		}
		defer rt.Close()

		ctx, err := rt.NewContext()
		if err != nil {
			return // Context creation issues are acceptable
		}
		defer ctx.Close()

		// The main test: evaluation shouldn't panic
		// Errors are expected for invalid JS, that's fine
		result, _ := ctx.Eval(code)

		// If we got a result, try to access it (shouldn't panic)
		if result.ptr != 0 {
			_ = result.String()
			_ = result.IsUndefined()
			_ = result.IsNull()
			_ = result.IsBool()
			_ = result.IsNumber()
			_ = result.IsString()
			_ = result.IsObject()
			_ = result.IsArray()
			_ = result.IsFunction()
			_ = result.IsError()
		}
	})
}

// FuzzJSONParse tests JSON parsing with arbitrary input
func FuzzJSONParse(f *testing.F) {
	seeds := []string{
		"{}",
		"[]",
		"null",
		"true",
		"false",
		"0",
		"1",
		"-1",
		"1.5",
		`""`,
		`"hello"`,
		`{"a":1}`,
		`[1,2,3]`,
		`{"nested":{"deep":true}}`,
		`[{"a":1},{"b":2}]`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, json string) {
		if !utf8.ValidString(json) {
			return
		}

		rt, err := NewRuntime()
		if err != nil {
			return
		}
		defer rt.Close()

		ctx, err := rt.NewContext()
		if err != nil {
			return
		}
		defer ctx.Close()

		// Try to parse - errors are expected for invalid JSON
		_, _ = ctx.Eval("JSON.parse(" + escapeJSString(json) + ")")
	})
}

// FuzzGoFunction tests Go function callbacks with arbitrary input
func FuzzGoFunction(f *testing.F) {
	seeds := []string{
		"",
		"a",
		"hello world",
		"123",
		"true",
		"null",
		`{"key":"value"}`,
		"[1,2,3]",
		"function(){}",
		"'single quotes'",
		`"double quotes"`,
		"mixed\ttabs\nand\nnewlines",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		if !utf8.ValidString(input) {
			return
		}

		rt, err := NewRuntime()
		if err != nil {
			return
		}
		defer rt.Close()

		ctx, err := rt.NewContext()
		if err != nil {
			return
		}
		defer ctx.Close()

		// Create a Go function that processes the input
		fn := ctx.Function("process", func(ctx *Context, this Value, args []Value) Value {
			if len(args) == 0 {
				return ctx.Undefined()
			}
			// Just return the string back
			return ctx.String(args[0].String())
		})
		ctx.SetGlobal("process", fn)

		// Call with the fuzzed input
		code := "process(" + escapeJSString(input) + ")"
		_, _ = ctx.Eval(code)
	})
}

// escapeJSString escapes a string for use in JavaScript
func escapeJSString(s string) string {
	result := `"`
	for _, r := range s {
		switch r {
		case '\\':
			result += `\\`
		case '"':
			result += `\"`
		case '\n':
			result += `\n`
		case '\r':
			result += `\r`
		case '\t':
			result += `\t`
		default:
			if r < 32 || r == 127 {
				result += `\x` + string("0123456789abcdef"[r>>4]) + string("0123456789abcdef"[r&0xf])
			} else {
				result += string(r)
			}
		}
	}
	result += `"`
	return result
}
