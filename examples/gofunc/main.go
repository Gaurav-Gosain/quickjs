// Example: Go Function Bindings
//
// This example demonstrates how to expose Go functions to JavaScript:
// - Creating Go functions callable from JS
// - Passing arguments between Go and JS
// - Returning values from Go to JS
// - Building a mini calculator with Go-backed operations
//
// Run with: go run ./examples/gofunc
package main

import (
	"fmt"
	"log"
	"math"
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

	// === Basic Go Function ===
	fmt.Println("=== Basic Go Function ===")

	// Create a simple add function
	addFn := ctx.Function("add", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		if len(args) < 2 {
			return ctx.Int32(0)
		}
		a, _ := args[0].Int32()
		b, _ := args[1].Int32()
		return ctx.Int32(a + b)
	})

	// Set it as a global
	ctx.SetGlobal("add", addFn)

	// Call from JavaScript
	result, _ := ctx.Eval("add(10, 20)")
	fmt.Printf("add(10, 20) = %s\n", result.String())

	// === String Processing ===
	fmt.Println("\n=== String Processing ===")

	// Create a string manipulation function
	upperFn := ctx.Function("toUpperCase", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		if len(args) < 1 {
			return ctx.String("")
		}
		s := args[0].String()
		return ctx.String(strings.ToUpper(s))
	})
	ctx.SetGlobal("goToUpperCase", upperFn)

	result, _ = ctx.Eval(`goToUpperCase("hello world")`)
	fmt.Printf(`goToUpperCase("hello world") = %s\n`, result.String())

	// === Math Functions ===
	fmt.Println("\n=== Math Functions (Go-backed) ===")

	// Expose Go's math functions
	sqrtFn := ctx.Function("sqrt", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		if len(args) < 1 {
			return ctx.Float64(math.NaN())
		}
		x, _ := args[0].Float64()
		return ctx.Float64(math.Sqrt(x))
	})
	ctx.SetGlobal("goSqrt", sqrtFn)

	powFn := ctx.Function("pow", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		if len(args) < 2 {
			return ctx.Float64(math.NaN())
		}
		base, _ := args[0].Float64()
		exp, _ := args[1].Float64()
		return ctx.Float64(math.Pow(base, exp))
	})
	ctx.SetGlobal("goPow", powFn)

	result, _ = ctx.Eval("goSqrt(16)")
	fmt.Printf("goSqrt(16) = %s\n", result.String())

	result, _ = ctx.Eval("goPow(2, 10)")
	fmt.Printf("goPow(2, 10) = %s\n", result.String())

	// === Callback with Multiple Args ===
	fmt.Println("\n=== Multiple Arguments ===")

	sumFn := ctx.Function("sum", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		var total float64
		for _, arg := range args {
			v, _ := arg.Float64()
			total += v
		}
		return ctx.Float64(total)
	})
	ctx.SetGlobal("sum", sumFn)

	result, _ = ctx.Eval("sum(1, 2, 3, 4, 5)")
	fmt.Printf("sum(1, 2, 3, 4, 5) = %s\n", result.String())

	// === Using Go Functions in JS Expressions ===
	fmt.Println("\n=== Go Functions in JS Expressions ===")

	result, _ = ctx.Eval(`
		const numbers = [1, 4, 9, 16, 25];
		const roots = numbers.map(n => goSqrt(n));
		roots.join(", ")
	`)
	fmt.Printf("Square roots of [1,4,9,16,25] = [%s]\n", result.String())

	// === Building a Mini Calculator ===
	fmt.Println("\n=== Mini Calculator ===")

	// Create a calculator object with Go-backed methods
	ctx.Eval(`var calc = {}`)

	calcAdd := ctx.Function("calcAdd", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		a, _ := args[0].Float64()
		b, _ := args[1].Float64()
		return ctx.Float64(a + b)
	})

	calcMul := ctx.Function("calcMul", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		a, _ := args[0].Float64()
		b, _ := args[1].Float64()
		return ctx.Float64(a * b)
	})

	calcObj, _ := ctx.GetGlobal("calc")
	calcObj.Set("add", calcAdd)
	calcObj.Set("multiply", calcMul)

	result, _ = ctx.Eval("calc.add(5, 3)")
	fmt.Printf("calc.add(5, 3) = %s\n", result.String())

	result, _ = ctx.Eval("calc.multiply(4, 7)")
	fmt.Printf("calc.multiply(4, 7) = %s\n", result.String())

	// === Logging from JS to Go ===
	fmt.Println("\n=== Logging from JS to Go ===")

	var logs []string
	logFn := ctx.Function("goLog", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		for _, arg := range args {
			logs = append(logs, arg.String())
		}
		return ctx.Undefined()
	})
	ctx.SetGlobal("goLog", logFn)

	ctx.Eval(`
		goLog("Starting calculation...");
		const result = add(100, 200);
		goLog("Result:", result);
		goLog("Done!");
	`)

	fmt.Println("Captured logs from JS:")
	for i, msg := range logs {
		fmt.Printf("  [%d] %s\n", i, msg)
	}

	fmt.Println("\n=== Done ===")
}
