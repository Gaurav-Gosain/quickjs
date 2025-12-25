// Example: Object and Array Manipulation
//
// This example demonstrates working with JavaScript objects and arrays:
// - Creating objects and arrays from Go
// - Getting and setting properties
// - Iterating over object properties
// - Array manipulation
// - Nested structures
//
// Run with: go run ./examples/objects
package main

import (
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

	// === Create Objects from Go ===
	fmt.Println("=== Create Objects from Go ===")

	// Create an empty object
	obj := ctx.Object()
	obj.Set("name", ctx.String("Alice"))
	obj.Set("age", ctx.Int32(30))
	obj.Set("active", ctx.Bool(true))

	// Set it as a global
	ctx.SetGlobal("person", obj)

	// Access from JavaScript
	result, _ := ctx.Eval(`person.name + " is " + person.age + " years old"`)
	fmt.Printf("Person: %s\n", result.String())

	// === Create Arrays from Go ===
	fmt.Println("\n=== Create Arrays from Go ===")

	arr := ctx.Array()
	arr.SetIdx(0, ctx.Int32(10))
	arr.SetIdx(1, ctx.Int32(20))
	arr.SetIdx(2, ctx.Int32(30))

	ctx.SetGlobal("numbers", arr)

	result, _ = ctx.Eval(`numbers.map(n => n * 2).join(", ")`)
	fmt.Printf("Doubled: %s\n", result.String())

	// === Get Properties ===
	fmt.Println("\n=== Get Properties ===")

	nameVal, _ := obj.Get("name")
	ageVal, _ := obj.Get("age")
	activeVal, _ := obj.Get("active")

	fmt.Printf("name: %s (isString: %v)\n", nameVal.String(), nameVal.IsString())
	fmt.Printf("age: %s (isNumber: %v)\n", ageVal.String(), ageVal.IsNumber())
	fmt.Printf("active: %s (isBool: %v)\n", activeVal.String(), activeVal.IsBool())

	// === Check Property Existence ===
	fmt.Println("\n=== Check Property Existence ===")

	hasName := obj.Has("name")
	hasEmail := obj.Has("email")
	fmt.Printf("has 'name': %v, has 'email': %v\n", hasName, hasEmail)

	// === Array Index Access ===
	fmt.Println("\n=== Array Index Access ===")

	ctx.Eval(`var fruits = ["apple", "banana", "cherry", "date"]`)
	fruits, _ := ctx.GetGlobal("fruits")

	for i := 0; i < 4; i++ {
		item, _ := fruits.GetIdx(i)
		fmt.Printf("fruits[%d] = %s\n", i, item.String())
	}

	// === Array Length ===
	fmt.Println("\n=== Array Length ===")

	length := fruits.Len()
	fmt.Printf("fruits.length = %d\n", length)

	// === Nested Objects ===
	fmt.Println("\n=== Nested Objects ===")

	ctx.Eval(`
		var company = {
			name: "TechCorp",
			address: {
				street: "123 Main St",
				city: "San Francisco",
				country: "USA"
			},
			employees: [
				{ name: "Alice", role: "Engineer" },
				{ name: "Bob", role: "Designer" },
				{ name: "Charlie", role: "Manager" }
			]
		}
	`)

	company, _ := ctx.GetGlobal("company")

	// Access nested properties
	address, _ := company.Get("address")
	city, _ := address.Get("city")
	fmt.Printf("Company city: %s\n", city.String())

	// Access array of objects
	employees, _ := company.Get("employees")
	emp0, _ := employees.GetIdx(0)
	emp0Name, _ := emp0.Get("name")
	emp0Role, _ := emp0.Get("role")
	fmt.Printf("First employee: %s (%s)\n", emp0Name.String(), emp0Role.String())

	// === Modify Object Properties ===
	fmt.Println("\n=== Modify Object Properties ===")

	result, _ = ctx.Eval(`person.age`)
	fmt.Printf("Before: person.age = %s\n", result.String())

	obj.Set("age", ctx.Int32(31))

	result, _ = ctx.Eval(`person.age`)
	fmt.Printf("After: person.age = %s\n", result.String())

	// === Delete Properties ===
	fmt.Println("\n=== Delete Properties ===")

	obj.Set("temporary", ctx.String("will be deleted"))
	result, _ = ctx.Eval(`person.temporary`)
	fmt.Printf("Before delete: person.temporary = %s\n", result.String())

	obj.Delete("temporary")
	result, _ = ctx.Eval(`person.temporary`)
	fmt.Printf("After delete: person.temporary = %s\n", result.String())

	// === Object Keys ===
	fmt.Println("\n=== Object Keys ===")

	result, _ = ctx.Eval(`Object.keys(person).join(", ")`)
	fmt.Printf("person keys: %s\n", result.String())

	// === Build Complex Object from Go ===
	fmt.Println("\n=== Build Complex Object from Go ===")

	config := ctx.Object()

	server := ctx.Object()
	server.Set("host", ctx.String("localhost"))
	server.Set("port", ctx.Int32(8080))
	config.Set("server", server)

	features := ctx.Array()
	features.SetIdx(0, ctx.String("auth"))
	features.SetIdx(1, ctx.String("logging"))
	features.SetIdx(2, ctx.String("caching"))
	config.Set("features", features)

	config.Set("debug", ctx.Bool(true))

	ctx.SetGlobal("config", config)

	result, _ = ctx.Eval(`JSON.stringify(config, null, 2)`)
	fmt.Printf("Built config:\n%s\n", result.String())

	// === Call Object Methods ===
	fmt.Println("\n=== Call Object Methods ===")

	ctx.Eval(`
		var calculator = {
			value: 0,
			add: function(n) { this.value += n; return this; },
			subtract: function(n) { this.value -= n; return this; },
			multiply: function(n) { this.value *= n; return this; },
			getResult: function() { return this.value; }
		}
	`)

	result, _ = ctx.Eval(`calculator.add(10).multiply(3).subtract(5).getResult()`)
	fmt.Printf("Calculator result: %s\n", result.String())

	// === Array Methods from Go ===
	fmt.Println("\n=== Array Methods from Go ===")

	ctx.Eval(`var scores = [85, 92, 78, 96, 88]`)

	// Call push via Eval (simpler)
	ctx.Eval(`scores.push(91)`)

	result, _ = ctx.Eval(`scores.join(", ")`)
	fmt.Printf("Scores after push: %s\n", result.String())

	result, _ = ctx.Eval(`Math.max(...scores)`)
	fmt.Printf("Max score: %s\n", result.String())

	result, _ = ctx.Eval(`scores.reduce((a, b) => a + b, 0) / scores.length`)
	fmt.Printf("Average score: %s\n", result.String())

	fmt.Println("\n=== Done ===")
}
