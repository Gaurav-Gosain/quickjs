// Example: JSON Processing
//
// This example demonstrates JSON handling between Go and JavaScript:
// - Parsing JSON strings in JavaScript
// - Stringifying JavaScript objects
// - Exchanging complex data structures
// - Working with nested objects and arrays
//
// Run with: go run ./examples/json
package main

import (
	"encoding/json"
	"fmt"
	"log"

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

	// === Parse JSON in JavaScript ===
	fmt.Println("=== Parse JSON in JavaScript ===")

	result, err := ctx.Eval(`
		const data = JSON.parse('{"name": "Alice", "age": 30, "active": true}');
		data.name + " is " + data.age + " years old"
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed result: %s\n", result.String())

	// === Stringify JavaScript Objects ===
	fmt.Println("\n=== Stringify JavaScript Objects ===")

	result, err = ctx.Eval(`
		const obj = {
			name: "Bob",
			scores: [85, 92, 78],
			metadata: {
				created: "2024-01-15",
				version: 2
			}
		};
		JSON.stringify(obj)
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Stringified: %s\n", result.String())

	// Parse the JSON string in Go
	var data map[string]any
	if err := json.Unmarshal([]byte(result.String()), &data); err == nil {
		fmt.Printf("Parsed in Go: name=%v, scores=%v\n", data["name"], data["scores"])
	}

	// === Send JSON from Go to JavaScript ===
	fmt.Println("\n=== Send JSON from Go to JavaScript ===")

	goData := map[string]any{
		"users": []map[string]any{
			{"id": 1, "name": "Alice", "role": "admin"},
			{"id": 2, "name": "Bob", "role": "user"},
			{"id": 3, "name": "Charlie", "role": "user"},
		},
		"total": 3,
	}

	jsonBytes, _ := json.Marshal(goData)
	jsonStr := string(jsonBytes)

	// Set the JSON string as a global variable
	ctx.SetGlobal("rawJson", ctx.String(jsonStr))

	result, _ = ctx.Eval(`
		const users = JSON.parse(rawJson);
		users.users.map(u => u.name).join(", ")
	`)
	fmt.Printf("User names: %s\n", result.String())

	result, _ = ctx.Eval(`
		users.users.filter(u => u.role === "user").length
	`)
	fmt.Printf("Number of regular users: %s\n", result.String())

	// === Pretty Print JSON ===
	fmt.Println("\n=== Pretty Print JSON ===")

	result, _ = ctx.Eval(`
		const config = {
			server: {
				host: "localhost",
				port: 8080,
				ssl: false
			},
			database: {
				type: "postgres",
				connection: "postgres://localhost/mydb"
			},
			features: ["auth", "logging", "caching"]
		};
		JSON.stringify(config, null, 2)
	`)
	fmt.Printf("Pretty JSON:\n%s\n", result.String())

	// === Process Array of JSON Objects ===
	fmt.Println("\n=== Process Array of JSON Objects ===")

	result, _ = ctx.Eval(`
		const products = [
			{ name: "Widget", price: 25.99, qty: 100 },
			{ name: "Gadget", price: 49.99, qty: 50 },
			{ name: "Gizmo", price: 15.99, qty: 200 }
		];
		
		const totalValue = products.reduce((sum, p) => sum + (p.price * p.qty), 0);
		const summary = {
			products: products.length,
			totalValue: totalValue.toFixed(2),
			avgPrice: (products.reduce((s, p) => s + p.price, 0) / products.length).toFixed(2)
		};
		JSON.stringify(summary)
	`)
	fmt.Printf("Summary: %s\n", result.String())

	// === JSON Validation ===
	fmt.Println("\n=== JSON Validation ===")

	validateFn := ctx.Function("isValidJSON", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		if len(args) < 1 {
			return ctx.Bool(false)
		}
		str := args[0].String()
		var js json.RawMessage
		return ctx.Bool(json.Unmarshal([]byte(str), &js) == nil)
	})
	ctx.SetGlobal("isValidJSON", validateFn)

	result, _ = ctx.Eval(`isValidJSON('{"valid": true}')`)
	fmt.Printf(`isValidJSON('{"valid": true}') = %s\n`, result.String())

	result, _ = ctx.Eval(`isValidJSON('{invalid}')`)
	fmt.Printf(`isValidJSON('{invalid}') = %s\n`, result.String())

	// === Deep Clone with JSON ===
	fmt.Println("\n=== Deep Clone with JSON ===")

	result, _ = ctx.Eval(`
		const original = { a: 1, b: { c: 2, d: [3, 4] } };
		const clone = JSON.parse(JSON.stringify(original));
		clone.b.c = 999;
		"original.b.c=" + original.b.c + ", clone.b.c=" + clone.b.c
	`)
	fmt.Printf("Deep clone test: %s\n", result.String())

	// === Custom JSON Replacer ===
	fmt.Println("\n=== Custom JSON Replacer ===")

	result, _ = ctx.Eval(`
		const sensitiveData = {
			username: "alice",
			password: "secret123",
			email: "alice@example.com",
			apiKey: "sk-xxxxx"
		};
		
		// Redact sensitive fields
		JSON.stringify(sensitiveData, (key, value) => {
			if (key === "password" || key === "apiKey") {
				return "[REDACTED]";
			}
			return value;
		})
	`)
	fmt.Printf("Redacted JSON: %s\n", result.String())

	fmt.Println("\n=== Done ===")
}
