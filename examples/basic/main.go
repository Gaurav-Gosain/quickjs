// Example: Basic QuickJS usage
//
// This example demonstrates the fundamental usage of the QuickJS-ng Go bindings:
// - Creating a runtime and context
// - Evaluating JavaScript code
// - Getting results back in Go
//
// Run with: go run ./examples/basic
package main

import (
	"fmt"
	"log"

	"github.com/Gaurav-Gosain/quickjs"
)

func main() {
	// Create a new JavaScript runtime
	// The runtime manages the WASM module and memory
	rt, err := quickjs.NewRuntime()
	if err != nil {
		log.Fatal(err)
	}
	defer rt.Close()

	// Create a context within the runtime
	// A context holds global variables and provides isolation
	ctx, err := rt.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Close()

	// Evaluate simple expressions
	fmt.Println("=== Basic Arithmetic ===")
	result, err := ctx.Eval("1 + 2 * 3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("1 + 2 * 3 = %s\n", result.String())

	// Get typed values
	val, _ := result.Int32()
	fmt.Printf("As int32: %d\n", val)

	// Evaluate more complex expressions
	fmt.Println("\n=== String Operations ===")
	result, err = ctx.Eval(`"Hello, " + "World!"`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("String concat: %s\n", result.String())

	// Check value types
	fmt.Println("\n=== Type Checking ===")
	result, _ = ctx.Eval("42")
	fmt.Printf("42 is number: %v\n", result.IsNumber())

	result, _ = ctx.Eval(`"hello"`)
	fmt.Printf(`"hello" is string: %v\n`, result.IsString())

	result, _ = ctx.Eval("true")
	fmt.Printf("true is bool: %v\n", result.IsBool())

	result, _ = ctx.Eval("[1, 2, 3]")
	fmt.Printf("[1,2,3] is array: %v\n", result.IsArray())

	result, _ = ctx.Eval("({x: 1})")
	fmt.Printf("{x:1} is object: %v\n", result.IsObject())

	// Define and call functions
	fmt.Println("\n=== Functions ===")
	_, err = ctx.Eval(`
		function greet(name) {
			return "Hello, " + name + "!";
		}
	`)
	if err != nil {
		log.Fatal(err)
	}

	result, err = ctx.Eval(`greet("QuickJS")`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("greet('QuickJS') = %s\n", result.String())

	// Multiple contexts are isolated
	fmt.Println("\n=== Context Isolation ===")
	ctx2, _ := rt.NewContext()
	defer ctx2.Close()

	ctx.Eval("var x = 100")
	ctx2.Eval("var x = 200")

	r1, _ := ctx.Eval("x")
	r2, _ := ctx2.Eval("x")
	fmt.Printf("ctx1.x = %s, ctx2.x = %s\n", r1.String(), r2.String())

	fmt.Println("\n=== Done ===")
}
