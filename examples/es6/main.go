// Example: ES6+ Features
//
// This example demonstrates modern JavaScript features supported by QuickJS-ng:
// - Arrow functions, Classes, Destructuring, Template literals
// - Spread operator, Promises, Map/Set/Symbol
// - Optional chaining, Nullish coalescing, BigInt
//
// Run with: go run ./examples/es6
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

	// === Arrow Functions ===
	fmt.Println("=== Arrow Functions ===")

	result, _ := ctx.Eval(`
		const add = (a, b) => a + b;
		const square = x => x * x;
		const greet = () => "Hello!";
		[add(2, 3), square(4), greet()].join(", ")
	`)
	fmt.Printf("Arrow functions: %s\n", result.String())

	// === Template Literals ===
	fmt.Println("\n=== Template Literals ===")

	result, _ = ctx.Eval("const tplName = 'World'; const tplCount = 42; `Hello, ${tplName}! The answer is ${tplCount}.`")
	fmt.Printf("Template literal: %s\n", result.String())

	// === Destructuring ===
	fmt.Println("\n=== Destructuring ===")

	result, _ = ctx.Eval(`
		(() => {
			const { x, y, z = 10 } = { x: 1, y: 2 };
			const [first, second, ...rest] = [1, 2, 3, 4, 5];
			return "x=" + x + ", y=" + y + ", z=" + z + ", first=" + first + ", rest=[" + rest + "]";
		})()
	`)
	fmt.Printf("Destructuring: %s\n", result.String())

	// === Spread Operator ===
	fmt.Println("\n=== Spread Operator ===")

	result, _ = ctx.Eval(`
		(() => {
			const arr1 = [1, 2, 3];
			const arr2 = [4, 5, 6];
			const combined = [...arr1, ...arr2];
			const obj1 = { a: 1, b: 2 };
			const obj2 = { ...obj1, c: 3 };
			return "array: [" + combined + "], object: " + JSON.stringify(obj2);
		})()
	`)
	fmt.Printf("Spread: %s\n", result.String())

	// === Classes ===
	fmt.Println("\n=== Classes ===")

	result, _ = ctx.Eval(`
		(() => {
			class Animal {
				constructor(name) { this.name = name; }
				speak() { return this.name + " makes a sound"; }
			}
			class Dog extends Animal {
				constructor(name, breed) { super(name); this.breed = breed; }
				speak() { return this.name + " barks!"; }
				static species() { return "Canis familiaris"; }
			}
			const dog = new Dog("Rex", "German Shepherd");
			return [dog.speak(), dog.breed, Dog.species()].join(" | ");
		})()
	`)
	fmt.Printf("Classes: %s\n", result.String())

	// === Map and Set ===
	fmt.Println("\n=== Map and Set ===")

	result, _ = ctx.Eval(`
		(() => {
			const map = new Map();
			map.set("a", 1);
			map.set("b", 2);
			map.set("c", 3);
			const set = new Set([1, 2, 2, 3, 3, 3]);
			return "Map size: " + map.size + ", Set size: " + set.size + ", Set values: [" + [...set] + "]";
		})()
	`)
	fmt.Printf("Map & Set: %s\n", result.String())

	// === Symbol ===
	fmt.Println("\n=== Symbol ===")

	result, _ = ctx.Eval(`
		(() => {
			const sym1 = Symbol("description");
			const sym2 = Symbol("description");
			return "sym1 === sym2: " + (sym1 === sym2) + ", typeof: " + (typeof sym1);
		})()
	`)
	fmt.Printf("Symbol: %s\n", result.String())

	// === Promises ===
	fmt.Println("\n=== Promises ===")

	result, _ = ctx.Eval(`Promise.resolve(42)`)
	fmt.Printf("Promise created: %v (isObject: %v)\n", result.String(), result.IsObject())

	// === Optional Chaining ===
	fmt.Println("\n=== Optional Chaining ===")

	result, _ = ctx.Eval(`
		(() => {
			const user = { name: "Alice", address: { city: "Wonderland" } };
			const city = user?.address?.city ?? "Unknown";
			const zip = user?.address?.zip ?? "N/A";
			return "city: " + city + ", zip: " + zip;
		})()
	`)
	fmt.Printf("Optional chaining: %s\n", result.String())

	// === Nullish Coalescing ===
	fmt.Println("\n=== Nullish Coalescing ===")

	result, _ = ctx.Eval(`
		(() => {
			const a = null ?? "default";
			const b = undefined ?? "default";
			const c = 0 ?? "default";
			const d = "" ?? "default";
			const e = false ?? "default";
			return [a, b, c, d, e].join(", ");
		})()
	`)
	fmt.Printf("Nullish coalescing: %s\n", result.String())

	// === BigInt ===
	fmt.Println("\n=== BigInt ===")

	result, _ = ctx.Eval(`
		(() => {
			const big1 = 9007199254740991n;
			const big2 = big1 + 1n;
			const big3 = big1 * 2n;
			return "big1: " + big1 + ", big1+1: " + big2 + ", big1*2: " + big3;
		})()
	`)
	fmt.Printf("BigInt: %s\n", result.String())

	// === Modern Array Methods ===
	fmt.Println("\n=== Modern Array Methods ===")

	result, _ = ctx.Eval(`
		(() => {
			const numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
			const doubled = numbers.map(n => n * 2);
			const evens = numbers.filter(n => n % 2 === 0);
			const sum = numbers.reduce((acc, n) => acc + n, 0);
			const found = numbers.find(n => n > 5);
			return JSON.stringify({
				doubled: doubled.slice(0, 3),
				evens: evens,
				sum: sum,
				found: found
			});
		})()
	`)
	fmt.Printf("Array methods: %s\n", result.String())

	// === Object Methods ===
	fmt.Println("\n=== Object Methods ===")

	result, _ = ctx.Eval(`
		(() => {
			const obj = { a: 1, b: 2, c: 3 };
			return JSON.stringify({
				keys: Object.keys(obj),
				values: Object.values(obj),
				entries: Object.entries(obj)
			});
		})()
	`)
	fmt.Printf("Object methods: %s\n", result.String())

	// === Proxy ===
	fmt.Println("\n=== Proxy ===")

	result, _ = ctx.Eval(`
		(() => {
			const target = { x: 10, y: 20 };
			const handler = {
				get(obj, prop) {
					return prop in obj ? obj[prop] * 2 : "not found";
				}
			};
			const proxy = new Proxy(target, handler);
			return "x: " + proxy.x + ", y: " + proxy.y + ", z: " + proxy.z;
		})()
	`)
	fmt.Printf("Proxy: %s\n", result.String())

	fmt.Println("\n=== Done ===")
}
