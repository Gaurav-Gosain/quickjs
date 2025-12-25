// Example: Error Handling
//
// This example demonstrates error handling in QuickJS-ng:
// - Catching JavaScript exceptions
// - Error messages and stack traces
// - Try/catch in JavaScript
// - Throwing errors from Go functions
// - Custom error types
//
// Run with: go run ./examples/errors
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Gaurav-Gosain/quickjs"
)

func main() {
	rt, err := quickjs.NewRuntime()
	if err != nil {
		log.Fatal(err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Close()

	// === Catching Syntax Errors ===
	fmt.Println("=== Catching Syntax Errors ===")

	_, err = ctx.Eval(`function broken( { }`)
	if err != nil {
		fmt.Printf("Syntax error caught: %v\n", err)
	}

	// === Catching Runtime Errors ===
	fmt.Println("\n=== Catching Runtime Errors ===")

	_, err = ctx.Eval(`undefinedVariable.foo`)
	if err != nil {
		fmt.Printf("Runtime error caught: %v\n", err)
	}

	// === Type Errors ===
	fmt.Println("\n=== Type Errors ===")

	_, err = ctx.Eval(`null.toString()`)
	if err != nil {
		fmt.Printf("Type error caught: %v\n", err)
	}

	// === Reference Errors ===
	fmt.Println("\n=== Reference Errors ===")

	_, err = ctx.Eval(`nonExistentFunction()`)
	if err != nil {
		fmt.Printf("Reference error caught: %v\n", err)
	}

	// === Try/Catch in JavaScript ===
	fmt.Println("\n=== Try/Catch in JavaScript ===")

	result, err := ctx.Eval(`
		(() => {
			try {
				throw new Error("Something went wrong!");
			} catch (e) {
				return "Caught: " + e.message;
			}
		})()
	`)
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", result.String())
	}

	// === Try/Catch/Finally ===
	fmt.Println("\n=== Try/Catch/Finally ===")

	result, _ = ctx.Eval(`
		(() => {
			let log = [];
			try {
				log.push("try block");
				throw new Error("oops");
			} catch (e) {
				log.push("catch: " + e.message);
			} finally {
				log.push("finally block");
			}
			return log.join(" -> ");
		})()
	`)
	fmt.Printf("Execution flow: %s\n", result.String())

	// === Custom Error Types ===
	fmt.Println("\n=== Custom Error Types ===")

	result, _ = ctx.Eval(`
		(() => {
			class ValidationError extends Error {
				constructor(message, field) {
					super(message);
					this.name = "ValidationError";
					this.field = field;
				}
			}
			
			try {
				throw new ValidationError("Invalid email format", "email");
			} catch (e) {
				if (e instanceof ValidationError) {
					return "Validation failed on field '" + e.field + "': " + e.message;
				}
				return "Unknown error: " + e.message;
			}
		})()
	`)
	fmt.Printf("Custom error: %s\n", result.String())

	// === Error with Stack Trace ===
	fmt.Println("\n=== Error with Stack Trace ===")

	result, _ = ctx.Eval(`
		(() => {
			function level1() { return level2(); }
			function level2() { return level3(); }
			function level3() { throw new Error("Deep error"); }
			
			try {
				level1();
			} catch (e) {
				return e.stack || e.message;
			}
		})()
	`)
	// Truncate stack trace for display
	stack := result.String()
	if len(stack) > 200 {
		stack = stack[:200] + "..."
	}
	fmt.Printf("Stack trace:\n%s\n", stack)

	// === Throwing Errors from Go Functions ===
	fmt.Println("\n=== Throwing Errors from Go Functions ===")

	// Create a Go function that validates input
	validateFn := ctx.Function("validateAge", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		if len(args) < 1 {
			// Return an error value
			errVal, _ := ctx.Eval(`new Error("Age is required")`)
			return errVal
		}

		age, err := args[0].Int32()
		if err != nil || age < 0 || age > 150 {
			errVal, _ := ctx.Eval(`new RangeError("Age must be between 0 and 150")`)
			return errVal
		}

		return ctx.Bool(true)
	})
	ctx.SetGlobal("validateAge", validateFn)

	// Test the validation
	testCases := []string{
		`validateAge(25)`,
		`validateAge(-5)`,
		`validateAge(200)`,
	}

	for _, tc := range testCases {
		result, _ := ctx.Eval(tc)
		if result.IsError() {
			fmt.Printf("%s -> Error: %s\n", tc, result.String())
		} else {
			fmt.Printf("%s -> %s\n", tc, result.String())
		}
	}

	// === Handling Errors in Callbacks ===
	fmt.Println("\n=== Handling Errors in Callbacks ===")

	var errors []string
	safeDivideFn := ctx.Function("safeDivide", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		if len(args) < 2 {
			errors = append(errors, "safeDivide: requires 2 arguments")
			return ctx.Undefined()
		}

		a, _ := args[0].Float64()
		b, _ := args[1].Float64()

		if b == 0 {
			errors = append(errors, fmt.Sprintf("safeDivide: division by zero (%.0f / %.0f)", a, b))
			return ctx.Undefined()
		}

		return ctx.Float64(a / b)
	})
	ctx.SetGlobal("safeDivide", safeDivideFn)

	ctx.Eval(`safeDivide(10, 2)`)
	ctx.Eval(`safeDivide(10, 0)`)
	ctx.Eval(`safeDivide(15)`)

	fmt.Println("Errors collected in Go:")
	for _, e := range errors {
		fmt.Printf("  - %s\n", e)
	}

	// === Promise Rejection ===
	fmt.Println("\n=== Promise Rejection ===")

	result, _ = ctx.Eval(`
		(() => {
			const p = new Promise((resolve, reject) => {
				reject(new Error("Promise was rejected"));
			});
			
			return p.catch(e => "Handled rejection: " + e.message);
		})()
	`)
	rt.ExecutePendingJobs()
	fmt.Printf("Promise rejection: %s\n", result.String())

	// === Error Boundary Pattern ===
	fmt.Println("\n=== Error Boundary Pattern ===")

	// Create a safe eval function
	safeEvalFn := ctx.Function("safeEval", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		if len(args) < 1 {
			return ctx.Object() // return empty result
		}

		code := args[0].String()
		result, err := ctx.Eval(code)

		resp := ctx.Object()
		if err != nil {
			resp.Set("success", ctx.Bool(false))
			resp.Set("error", ctx.String(err.Error()))
		} else {
			resp.Set("success", ctx.Bool(true))
			resp.Set("result", result)
		}
		return resp
	})
	ctx.SetGlobal("safeEval", safeEvalFn)

	// Test safe eval
	codes := []string{
		`"1 + 2"`,
		`"invalidSyntax((("`,
		`"Math.sqrt(16)"`,
	}

	for _, code := range codes {
		result, _ := ctx.Eval(`JSON.stringify(safeEval(` + code + `))`)
		// Clean up the output
		output := strings.ReplaceAll(result.String(), `\n`, " ")
		if len(output) > 80 {
			output = output[:80] + "..."
		}
		fmt.Printf("safeEval(%s) -> %s\n", code, output)
	}

	fmt.Println("\n=== Done ===")
}
