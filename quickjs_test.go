package quickjs

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestNewRuntime(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	rt.Close()
}

func TestNewContext(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	ctx.Close()
}

func TestMultipleContexts(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	// Create multiple contexts
	ctx1, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx1.Close()

	ctx2, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx2.Close()

	// Set different values in each context
	_, err = ctx1.Eval("var x = 42;")
	if err != nil {
		t.Fatalf("Eval error = %v", err)
	}

	_, err = ctx2.Eval("var x = 100;")
	if err != nil {
		t.Fatalf("Eval error = %v", err)
	}

	// Verify isolation
	result1, err := ctx1.Eval("x")
	if err != nil {
		t.Fatalf("Eval error = %v", err)
	}
	if result1.String() != "42" {
		t.Errorf("ctx1 x = %q, want %q", result1.String(), "42")
	}

	result2, err := ctx2.Eval("x")
	if err != nil {
		t.Fatalf("Eval error = %v", err)
	}
	if result2.String() != "100" {
		t.Errorf("ctx2 x = %q, want %q", result2.String(), "100")
	}
}

// ============================================================================
// Basic JavaScript Evaluation
// ============================================================================

func TestEvalBasicArithmetic(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"1 + 2", "3"},
		{"10 - 3", "7"},
		{"4 * 5", "20"},
		{"20 / 4", "5"},
		{"10 % 3", "1"},
		{"2 ** 10", "1024"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestEvalKeywords(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"true", "true"},
		{"false", "false"},
		{"null", "null"},
		{"typeof 42", "number"},
		{"typeof 'hello'", "string"},
		{"typeof true", "boolean"},
		{"typeof null", "object"},
		{"typeof undefined", "undefined"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestEvalStrings(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"hello" + " " + "world"`, "hello world"},
		{`"hello".toUpperCase()`, "HELLO"},
		{`"HELLO".toLowerCase()`, "hello"},
		{`"hello".length`, "5"},
		{`"hello world".indexOf("world")`, "6"},
		{`"hello".charAt(1)`, "e"},
		{`"hello".substring(1, 4)`, "ell"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestEvalArrays(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"[1, 2, 3].length", "3"},
		{"[1, 2, 3][0]", "1"},
		{"[1, 2, 3][2]", "3"},
		{"[1, 2, 3].join(',')", "1,2,3"},
		{"[3, 1, 2].sort().join(',')", "1,2,3"},
		{"[1, 2, 3].reverse().join(',')", "3,2,1"},
		{"[1, 2, 3].indexOf(2)", "1"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestEvalObjects(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{`var obj = {a: 1, b: 2}; obj.a`, "1"},
		{`var obj = {a: 1, b: 2}; obj.b`, "2"},
		{`var obj = {a: 1, b: 2}; obj.a + obj.b`, "3"},
		{`var obj = {}; obj.x = 42; obj.x`, "42"},
		{`JSON.stringify({x: 1, y: 2})`, `{"x":1,"y":2}`},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestEvalFunctions(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"function add(a, b) { return a + b; } add(2, 3)", "5"},
		{"var multiply = function(a, b) { return a * b; }; multiply(4, 5)", "20"},
		{"function factorial(n) { if (n <= 1) return 1; return n * factorial(n - 1); } factorial(5)", "120"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestEvalMath(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"Math.sqrt(16)", "4"},
		{"Math.abs(-5)", "5"},
		{"Math.floor(3.7)", "3"},
		{"Math.ceil(3.2)", "4"},
		{"Math.round(3.5)", "4"},
		{"Math.min(1, 2, 3)", "1"},
		{"Math.max(1, 2, 3)", "3"},
		{"Math.pow(2, 10)", "1024"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestEvalTryCatch(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	result, err := ctx.Eval(`
		try {
			throw new Error("test error");
		} catch (e) {
			e.message
		}
	`)
	if err != nil {
		t.Fatalf("Eval error = %v", err)
	}
	if result.String() != "test error" {
		t.Errorf("Eval = %q, want %q", result.String(), "test error")
	}
}

// ============================================================================
// ES6+ Features
// ============================================================================

func TestES6ArrowFunctions(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"const add = (a, b) => a + b; add(2, 3)", "5"},
		{"const square = x => x * x; square(4)", "16"},
		{"const greet = () => 'hello'; greet()", "hello"},
		{"[1, 2, 3].map(x => x * 2).join(',')", "2,4,6"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestES6LetConst(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"let x = 10; x", "10"},
		{"const y = 20; y", "20"},
		{"let a = 1; { let a = 2; } a", "1"}, // Block scoping
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestES6TemplateLiterals(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"`hello`", "hello"},
		{"`1 + 2 = ${1 + 2}`", "1 + 2 = 3"},
		{"const name = 'World'; `Hello, ${name}!`", "Hello, World!"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestES6Destructuring(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"const [a, b] = [1, 2]; a + b", "3"},
		{"const {x, y} = {x: 10, y: 20}; x + y", "30"},
		{"const [first, ...rest] = [1, 2, 3, 4]; rest.join(',')", "2,3,4"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestES6Spread(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"[...[1, 2], ...[3, 4]].join(',')", "1,2,3,4"},
		{"(() => { const arr = [1, 2, 3]; return [...arr, 4, 5].join(','); })()", "1,2,3,4,5"},
		{"({...{a: 1}, ...{b: 2}}).b", "2"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestES6Classes(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			"basic class",
			`(() => { class Animal { constructor(name) { this.name = name; } speak() { return this.name; } } return new Animal('Dog').speak(); })()`,
			"Dog",
		},
		{
			"class inheritance",
			`(() => {
				class Animal { constructor(name) { this.name = name; } }
				class Dog extends Animal { bark() { return this.name + ' barks'; } }
				return new Dog('Rex').bark();
			})()`,
			"Rex barks",
		},
		{
			"static method",
			`(() => { class Calculator { static add(a, b) { return a + b; } } return Calculator.add(5, 3); })()`,
			"8",
		},
		{
			"getter/setter",
			`(() => {
				class Circle {
					constructor(r) { this._r = r; }
					get area() { return 3.14159 * this._r * this._r; }
				}
				return Math.floor(new Circle(10).area);
			})()`,
			"314",
		},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("%s: Eval error = %v", tt.name, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("%s: got %q, want %q", tt.name, result.String(), tt.expected)
		}
	}
}

func TestES6Promises(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			"Promise.resolve",
			`let result = 0; Promise.resolve(42).then(x => { result = x; }); result`,
			"0", // Promise hasn't resolved yet
		},
		{
			"Promise creation",
			`typeof new Promise((resolve, reject) => {})`,
			"object",
		},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("%s: Eval error = %v", tt.name, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("%s: got %q, want %q", tt.name, result.String(), tt.expected)
		}
	}
}

func TestES6MapSet(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"new Set([1, 2, 2, 3]).size", "3"},
		{"new Map([['a', 1], ['b', 2]]).get('b')", "2"},
		{"new Map([['a', 1]]).has('a')", "true"},
		{"new Set([1, 2, 3]).has(2)", "true"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestES6Symbol(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"typeof Symbol('test')", "symbol"},
		{"Symbol('a') === Symbol('a')", "false"},
		{"Symbol.for('global') === Symbol.for('global')", "true"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestES6Proxy(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	result, err := ctx.Eval(`
		const handler = {
			get: (target, prop) => target[prop] * 2
		};
		const target = { x: 21 };
		const proxy = new Proxy(target, handler);
		proxy.x
	`)
	if err != nil {
		t.Fatalf("Eval error = %v", err)
	}
	if result.String() != "42" {
		t.Errorf("Proxy get trap: got %q, want %q", result.String(), "42")
	}
}

func TestES2020BigInt(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"typeof 1n", "bigint"},
		{"1n + 2n", "3"},
		{"BigInt(100)", "100"},
		{"(2n ** 64n).toString()", "18446744073709551616"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestES2020OptionalChaining(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"(() => { const obj = {a: {b: 1}}; return obj?.a?.b; })()", "1"},
		{"(() => { const obj = {a: {b: 1}}; return obj?.x?.y; })()", "undefined"},
		{"(() => { const arr = [1, 2, 3]; return arr?.[1]; })()", "2"},
		{"null?.foo", "undefined"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

func TestES2020NullishCoalescing(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		code     string
		expected string
	}{
		{"null ?? 'default'", "default"},
		{"undefined ?? 'default'", "default"},
		{"0 ?? 'default'", "0"},
		{"'' ?? 'default'", ""},
		{"false ?? 'default'", "false"},
	}

	for _, tt := range tests {
		result, err := ctx.Eval(tt.code)
		if err != nil {
			t.Errorf("Eval(%q) error = %v", tt.code, err)
			continue
		}
		if result.String() != tt.expected {
			t.Errorf("Eval(%q) = %q, want %q", tt.code, result.String(), tt.expected)
		}
	}
}

// ============================================================================
// Value Types
// ============================================================================

func TestValueTypes(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Test integer
	intVal, _ := ctx.Eval("42")
	if !intVal.IsNumber() {
		t.Errorf("42 should be number")
	}
	i, _ := intVal.Int32()
	if i != 42 {
		t.Errorf("Int32() = %d, want 42", i)
	}

	// Test float
	floatVal, _ := ctx.Eval("3.14")
	if !floatVal.IsNumber() {
		t.Errorf("3.14 should be number")
	}
	f, _ := floatVal.Float64()
	if f != 3.14 {
		t.Errorf("Float64() = %f, want 3.14", f)
	}

	// Test string
	strVal, _ := ctx.Eval(`"hello"`)
	if !strVal.IsString() {
		t.Errorf(`"hello" should be string`)
	}
	if strVal.String() != "hello" {
		t.Errorf("String() = %q, want %q", strVal.String(), "hello")
	}

	// Test boolean
	boolVal, _ := ctx.Eval("true")
	if !boolVal.IsBool() {
		t.Errorf("true should be bool")
	}
	if !boolVal.Bool() {
		t.Errorf("Bool() = false, want true")
	}

	// Test null
	nullVal, _ := ctx.Eval("null")
	if !nullVal.IsNull() {
		t.Errorf("null should be null")
	}

	// Test undefined
	undefVal, _ := ctx.Eval("undefined")
	if !undefVal.IsUndefined() {
		t.Errorf("undefined should be undefined")
	}

	// Test function
	funcVal, _ := ctx.Eval("(function() {})")
	if !funcVal.IsFunction() {
		t.Errorf("function should be function")
	}

	// Test array
	arrVal, _ := ctx.Eval("[1, 2, 3]")
	if !arrVal.IsArray() {
		t.Errorf("[] should be array")
	}

	// Test object
	objVal, _ := ctx.Eval("({a: 1})")
	if !objVal.IsObject() {
		t.Errorf("{} should be object")
	}
}

// ============================================================================
// Value Creation
// ============================================================================

func TestValueCreation(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Test Int32
	intVal := ctx.Int32(42)
	if intVal.String() != "42" {
		t.Errorf("Int32(42).String() = %q, want %q", intVal.String(), "42")
	}

	// Test Float64
	floatVal := ctx.Float64(3.14)
	if floatVal.String() != "3.14" {
		t.Errorf("Float64(3.14).String() = %q, want %q", floatVal.String(), "3.14")
	}

	// Test String
	strVal := ctx.String("hello")
	if strVal.String() != "hello" {
		t.Errorf("String(\"hello\").String() = %q, want %q", strVal.String(), "hello")
	}

	// Test Bool
	boolVal := ctx.Bool(true)
	if !boolVal.Bool() {
		t.Errorf("Bool(true).Bool() = false, want true")
	}

	// Test Null
	nullVal := ctx.Null()
	if !nullVal.IsNull() {
		t.Errorf("Null().IsNull() = false, want true")
	}

	// Test Undefined
	undefVal := ctx.Undefined()
	if !undefVal.IsUndefined() {
		t.Errorf("Undefined().IsUndefined() = false, want true")
	}
}

// ============================================================================
// Object Operations
// ============================================================================

func TestObjectOperations(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	obj := ctx.Object()

	// Set properties
	if err := obj.Set("x", ctx.Int32(42)); err != nil {
		t.Fatalf("Set error = %v", err)
	}
	if err := obj.Set("y", ctx.String("hello")); err != nil {
		t.Fatalf("Set error = %v", err)
	}

	// Get properties
	x, err := obj.Get("x")
	if err != nil {
		t.Fatalf("Get error = %v", err)
	}
	if x.String() != "42" {
		t.Errorf("Get(\"x\") = %q, want %q", x.String(), "42")
	}

	y, err := obj.Get("y")
	if err != nil {
		t.Fatalf("Get error = %v", err)
	}
	if y.String() != "hello" {
		t.Errorf("Get(\"y\") = %q, want %q", y.String(), "hello")
	}

	// Has property
	if !obj.Has("x") {
		t.Errorf("Has(\"x\") = false, want true")
	}
	if obj.Has("z") {
		t.Errorf("Has(\"z\") = true, want false")
	}
}

func TestArrayOperations(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	arr := ctx.Array()

	// Set elements
	if err := arr.SetIdx(0, ctx.Int32(10)); err != nil {
		t.Fatalf("SetIdx error = %v", err)
	}
	if err := arr.SetIdx(1, ctx.Int32(20)); err != nil {
		t.Fatalf("SetIdx error = %v", err)
	}
	if err := arr.SetIdx(2, ctx.Int32(30)); err != nil {
		t.Fatalf("SetIdx error = %v", err)
	}

	// Get length
	if arr.Len() != 3 {
		t.Errorf("Len() = %d, want 3", arr.Len())
	}

	// Get elements
	elem, err := arr.GetIdx(1)
	if err != nil {
		t.Fatalf("GetIdx error = %v", err)
	}
	if elem.String() != "20" {
		t.Errorf("GetIdx(1) = %q, want %q", elem.String(), "20")
	}
}

// ============================================================================
// Function Calling
// ============================================================================

func TestCallFunction(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Define a function
	_, err = ctx.Eval("function add(a, b) { return a + b; }")
	if err != nil {
		t.Fatalf("Eval error = %v", err)
	}

	// Get the function from global
	addFunc, err := ctx.GetGlobal("add")
	if err != nil {
		t.Fatalf("GetGlobal error = %v", err)
	}

	if !addFunc.IsFunction() {
		t.Fatalf("add should be a function")
	}

	// Call the function
	result, err := addFunc.Call(ctx.Undefined(), ctx.Int32(5), ctx.Int32(3))
	if err != nil {
		t.Fatalf("Call error = %v", err)
	}

	if result.String() != "8" {
		t.Errorf("add(5, 3) = %q, want %q", result.String(), "8")
	}
}

// ============================================================================
// Go Function Binding
// ============================================================================

func TestGoFunction(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Create a Go function
	addFn := ctx.Function("add", func(c *Context, this Value, args []Value) Value {
		if len(args) < 2 {
			return c.Int32(0)
		}
		a, _ := args[0].Int32()
		b, _ := args[1].Int32()
		return c.Int32(a + b)
	})

	// Set it as a global
	if err := ctx.SetGlobal("goAdd", addFn); err != nil {
		t.Fatalf("SetGlobal error = %v", err)
	}

	// Call it from JavaScript
	result, err := ctx.Eval("goAdd(10, 20)")
	if err != nil {
		t.Fatalf("Eval error = %v", err)
	}

	if result.String() != "30" {
		t.Errorf("goAdd(10, 20) = %q, want %q", result.String(), "30")
	}
}

func TestGoFunctionWithStrings(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Create a Go function that concatenates strings
	concatFn := ctx.Function("concat", func(c *Context, this Value, args []Value) Value {
		result := ""
		for _, arg := range args {
			result += arg.String()
		}
		return c.String(result)
	})

	if err := ctx.SetGlobal("goConcat", concatFn); err != nil {
		t.Fatalf("SetGlobal error = %v", err)
	}

	result, err := ctx.Eval(`goConcat("Hello, ", "World!")`)
	if err != nil {
		t.Fatalf("Eval error = %v", err)
	}

	if result.String() != "Hello, World!" {
		t.Errorf("goConcat = %q, want %q", result.String(), "Hello, World!")
	}
}

// ============================================================================
// JSON
// ============================================================================

func TestJSON(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Parse JSON
	obj, err := ctx.ParseJSON(`{"name": "John", "age": 30}`)
	if err != nil {
		t.Fatalf("ParseJSON error = %v", err)
	}

	name, err := obj.Get("name")
	if err != nil {
		t.Fatalf("Get error = %v", err)
	}
	if name.String() != "John" {
		t.Errorf("name = %q, want %q", name.String(), "John")
	}

	// Stringify
	jsonStr, err := obj.JSONStringify()
	if err != nil {
		t.Fatalf("JSONStringify error = %v", err)
	}
	if !strings.Contains(jsonStr, "John") {
		t.Errorf("JSONStringify should contain 'John', got %q", jsonStr)
	}
}

// ============================================================================
// Print/Console
// ============================================================================

func TestPrint(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	var logs []string
	rt.SetLogFunc(func(msg string) {
		logs = append(logs, msg)
	})

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	_, err = ctx.Eval(`print("hello"); print("world");`)
	if err != nil {
		t.Fatalf("Eval error = %v", err)
	}

	allLogs := strings.Join(logs, "")
	if !strings.Contains(allLogs, "hello") {
		t.Errorf("logs should contain %q, got %v", "hello", logs)
	}
	if !strings.Contains(allLogs, "world") {
		t.Errorf("logs should contain %q, got %v", "world", logs)
	}
}

// ============================================================================
// Concurrency
// ============================================================================

func TestParallelRuntimes(t *testing.T) {
	const numGoroutines = 10
	const iterationsPerGoroutine = 5

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*iterationsPerGoroutine)

	for g := range numGoroutines {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for i := range iterationsPerGoroutine {
				rt, err := NewRuntime()
				if err != nil {
					errors <- fmt.Errorf("goroutine %d, iter %d: NewRuntime error: %w", goroutineID, i, err)
					continue
				}

				ctx, err := rt.NewContext()
				if err != nil {
					rt.Close()
					errors <- fmt.Errorf("goroutine %d, iter %d: NewContext error: %w", goroutineID, i, err)
					continue
				}

				code := fmt.Sprintf("var x = %d * %d; x + 1", goroutineID, i)
				expected := goroutineID*i + 1

				result, err := ctx.Eval(code)
				if err != nil {
					ctx.Close()
					rt.Close()
					errors <- fmt.Errorf("goroutine %d, iter %d: Eval error: %w", goroutineID, i, err)
					continue
				}

				val, err := result.Int32()
				if err != nil {
					ctx.Close()
					rt.Close()
					errors <- fmt.Errorf("goroutine %d, iter %d: Int32 error: %w", goroutineID, i, err)
					continue
				}

				if int(val) != expected {
					errors <- fmt.Errorf("goroutine %d, iter %d: got %d, want %d", goroutineID, i, val, expected)
				}

				ctx.Close()
				rt.Close()
			}
		}(g)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

func TestConcurrentEvalSameContext(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Initialize counter
	_, err = ctx.Eval("var counter = 0")
	if err != nil {
		t.Fatalf("Eval init error: %v", err)
	}

	const numGoroutines = 10
	const incrementsPerGoroutine = 10

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*incrementsPerGoroutine)

	for range numGoroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range incrementsPerGoroutine {
				_, err := ctx.Eval("counter++")
				if err != nil {
					errors <- err
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Concurrent eval error: %v", err)
	}

	// Verify final counter value
	result, err := ctx.Eval("counter")
	if err != nil {
		t.Fatalf("Final eval error: %v", err)
	}

	val, _ := result.Int32()
	expected := numGoroutines * incrementsPerGoroutine
	if int(val) != expected {
		t.Errorf("Final counter = %d, want %d", val, expected)
	}
}

// ============================================================================
// Edge Cases and Error Handling Tests
// ============================================================================

func TestEvalEmptyString(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	result, err := ctx.Eval("")
	if err != nil {
		t.Fatalf("Eval('') error = %v", err)
	}
	if !result.IsUndefined() {
		t.Errorf("Eval('') = %v, want undefined", result.String())
	}
}

func TestEvalSyntaxError(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	_, err = ctx.Eval("function broken( { }")
	if err == nil {
		t.Error("Expected syntax error, got nil")
	}
}

func TestEvalReferenceError(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	_, err = ctx.Eval("undefinedVariable")
	if err == nil {
		t.Error("Expected reference error, got nil")
	}
}

func TestEvalTypeError(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	_, err = ctx.Eval("null.foo()")
	if err == nil {
		t.Error("Expected type error, got nil")
	}
}

func TestValueConversionErrors(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Test Int32 conversion of non-number
	strVal, _ := ctx.Eval(`"hello"`)
	_, _ = strVal.Int32() // May or may not error; just verify no panic

	// Test Float64 conversion
	_, _ = strVal.Float64() // May or may not error; just verify no panic

	// Test on object
	objVal, _ := ctx.Eval(`({x: 1})`)
	_ = objVal.String() // Should not panic
}

func TestNullAndUndefined(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Test null
	nullVal, _ := ctx.Eval("null")
	if !nullVal.IsNull() {
		t.Error("Expected IsNull() = true")
	}
	if nullVal.IsUndefined() {
		t.Error("null should not be undefined")
	}

	// Test undefined
	undefVal, _ := ctx.Eval("undefined")
	if !undefVal.IsUndefined() {
		t.Error("Expected IsUndefined() = true")
	}
	if undefVal.IsNull() {
		t.Error("undefined should not be null")
	}

	// Test created values
	ctxNull := ctx.Null()
	if !ctxNull.IsNull() {
		t.Error("ctx.Null() should be null")
	}

	ctxUndef := ctx.Undefined()
	if !ctxUndef.IsUndefined() {
		t.Error("ctx.Undefined() should be undefined")
	}
}

func TestLargeNumbers(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Test large integer
	result, _ := ctx.Eval("Number.MAX_SAFE_INTEGER")
	val, _ := result.Float64()
	if val != 9007199254740991 {
		t.Errorf("MAX_SAFE_INTEGER = %v, want 9007199254740991", val)
	}

	// Test negative numbers
	result, _ = ctx.Eval("-2147483648")
	intVal, _ := result.Int32()
	if intVal != -2147483648 {
		t.Errorf("Min int32 = %v, want -2147483648", intVal)
	}

	// Test infinity
	result, _ = ctx.Eval("Infinity")
	str := result.String()
	if str != "Infinity" {
		t.Errorf("Infinity = %v, want 'Infinity'", str)
	}

	// Test NaN
	result, _ = ctx.Eval("NaN")
	str = result.String()
	if str != "NaN" {
		t.Errorf("NaN = %v, want 'NaN'", str)
	}
}

func TestSpecialStrings(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", `""`, ""},
		{"unicode", `"你好世界"`, "你好世界"},
		{"emoji", `"Hello 👋 World 🌍"`, "Hello 👋 World 🌍"},
		{"newlines", `"line1\nline2"`, "line1\nline2"},
		{"tabs", `"col1\tcol2"`, "col1\tcol2"},
		{"quotes", `"say \"hello\""`, `say "hello"`},
		{"backslash", `"path\\to\\file"`, `path\to\file`},
		// Note: null characters truncate C strings, so "a\x00b" becomes "a"
		// This is expected behavior with the C bridge
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ctx.Eval(tt.input)
			if err != nil {
				t.Fatalf("Eval error: %v", err)
			}
			if result.String() != tt.expected {
				t.Errorf("got %q, want %q", result.String(), tt.expected)
			}
		})
	}
}

func TestDeepNesting(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Create deeply nested object
	result, err := ctx.Eval(`
		(() => {
			let obj = { value: 42 };
			for (let i = 0; i < 100; i++) {
				obj = { nested: obj };
			}
			// Access the deep value
			let current = obj;
			for (let i = 0; i < 100; i++) {
				current = current.nested;
			}
			return current.value;
		})()
	`)
	if err != nil {
		t.Fatalf("Deep nesting eval error: %v", err)
	}
	if result.String() != "42" {
		t.Errorf("Deep nested value = %v, want 42", result.String())
	}
}

func TestLargeArray(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Create and sum a large array
	result, err := ctx.Eval(`
		(() => {
			const arr = [];
			for (let i = 0; i < 10000; i++) {
				arr.push(i);
			}
			return arr.reduce((a, b) => a + b, 0);
		})()
	`)
	if err != nil {
		t.Fatalf("Large array eval error: %v", err)
	}

	val, _ := result.Float64()
	expected := float64(10000 * 9999 / 2) // Sum of 0 to 9999
	if val != expected {
		t.Errorf("Large array sum = %v, want %v", val, expected)
	}
}

func TestGoFunctionWithManyArgs(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Create a function that sums all arguments
	sumFn := ctx.Function("sumAll", func(ctx *Context, this Value, args []Value) Value {
		var sum float64
		for _, arg := range args {
			v, _ := arg.Float64()
			sum += v
		}
		return ctx.Float64(sum)
	})
	ctx.SetGlobal("sumAll", sumFn)

	// Test with many arguments
	result, err := ctx.Eval("sumAll(1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20)")
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}

	val, _ := result.Float64()
	if val != 210 { // Sum of 1 to 20
		t.Errorf("sumAll(1..20) = %v, want 210", val)
	}
}

func TestGoFunctionReturnsError(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Create a function that always returns undefined
	fn := ctx.Function("alwaysUndefined", func(ctx *Context, this Value, args []Value) Value {
		return ctx.Undefined()
	})
	ctx.SetGlobal("alwaysUndefined", fn)

	result, _ := ctx.Eval("alwaysUndefined()")
	if !result.IsUndefined() {
		t.Errorf("Expected undefined, got %v", result.String())
	}
}

func TestObjectPropertyChain(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Create nested object from Go
	root := ctx.Object()
	level1 := ctx.Object()
	level2 := ctx.Object()

	level2.Set("value", ctx.Int32(42))
	level1.Set("child", level2)
	root.Set("child", level1)

	ctx.SetGlobal("root", root)

	result, _ := ctx.Eval("root.child.child.value")
	val, _ := result.Int32()
	if val != 42 {
		t.Errorf("Nested value = %v, want 42", val)
	}
}

func TestArrayOperationsFromGo(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Create array from Go
	arr := ctx.Array()
	for i := range 5 {
		arr.SetIdx(i, ctx.Int32(int32(i*10)))
	}
	ctx.SetGlobal("arr", arr)

	// Verify length
	if arr.Len() != 5 {
		t.Errorf("Array length = %d, want 5", arr.Len())
	}

	// Verify elements
	for i := range 5 {
		elem, _ := arr.GetIdx(i)
		val, _ := elem.Int32()
		if val != int32(i*10) {
			t.Errorf("arr[%d] = %d, want %d", i, val, i*10)
		}
	}

	// Test JS operations on the array
	result, _ := ctx.Eval("arr.reduce((a, b) => a + b, 0)")
	sum, _ := result.Int32()
	if sum != 100 { // 0+10+20+30+40
		t.Errorf("Array sum = %d, want 100", sum)
	}
}

func TestClosurePreservation(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Create a closure
	_, err = ctx.Eval(`
		var createCounter = function() {
			var count = 0;
			return function() {
				return ++count;
			};
		};
		var counter = createCounter();
	`)
	if err != nil {
		t.Fatalf("Closure creation error: %v", err)
	}

	// Call multiple times and verify closure preserves state
	for i := 1; i <= 5; i++ {
		result, err := ctx.Eval("counter()")
		if err != nil {
			t.Fatalf("Counter call error: %v", err)
		}
		val, _ := result.Int32()
		if val != int32(i) {
			t.Errorf("counter() call %d = %d, want %d", i, val, i)
		}
	}
}

func TestMultipleGoFunctions(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Register multiple Go functions
	ctx.SetGlobal("goAdd", ctx.Function("add", func(ctx *Context, this Value, args []Value) Value {
		a, _ := args[0].Int32()
		b, _ := args[1].Int32()
		return ctx.Int32(a + b)
	}))

	ctx.SetGlobal("goMul", ctx.Function("mul", func(ctx *Context, this Value, args []Value) Value {
		a, _ := args[0].Int32()
		b, _ := args[1].Int32()
		return ctx.Int32(a * b)
	}))

	ctx.SetGlobal("goNeg", ctx.Function("neg", func(ctx *Context, this Value, args []Value) Value {
		a, _ := args[0].Int32()
		return ctx.Int32(-a)
	}))

	// Use them together
	result, err := ctx.Eval("goNeg(goAdd(goMul(3, 4), 5))")
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}

	val, _ := result.Int32()
	if val != -17 {
		t.Errorf("Result = %d, want -17 (expected -(3*4 + 5))", val)
	}
}

// ============================================================================
// Stress Tests
// ============================================================================

func TestStressManyEvals(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Run many evaluations
	for i := range 1000 {
		code := fmt.Sprintf("%d + %d", i, i*2)
		result, err := ctx.Eval(code)
		if err != nil {
			t.Fatalf("Eval error at iteration %d: %v", i, err)
		}
		val, _ := result.Int32()
		if val != int32(i*3) {
			t.Fatalf("Result at iteration %d = %d, want %d", i, val, i*3)
		}
	}
}

func TestStressManyObjects(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Create many objects
	for i := range 500 {
		obj := ctx.Object()
		obj.Set("id", ctx.Int32(int32(i)))
		obj.Set("name", ctx.String(fmt.Sprintf("object_%d", i)))

		// Verify
		idVal, _ := obj.Get("id")
		id, _ := idVal.Int32()
		if id != int32(i) {
			t.Fatalf("Object %d has wrong id: %d", i, id)
		}
	}
}

func TestStressManyGoCallbacks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	callCount := 0
	fn := ctx.Function("increment", func(ctx *Context, this Value, args []Value) Value {
		callCount++
		return ctx.Int32(int32(callCount))
	})
	ctx.SetGlobal("increment", fn)

	// Call the Go function many times from JS
	_, err = ctx.Eval(`
		for (let i = 0; i < 500; i++) {
			increment();
		}
	`)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}

	if callCount != 500 {
		t.Errorf("Call count = %d, want 500", callCount)
	}
}

func TestStressRapidContextCreation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	// Create and destroy many contexts
	for i := range 100 {
		ctx, err := rt.NewContext()
		if err != nil {
			t.Fatalf("NewContext error at iteration %d: %v", i, err)
		}

		result, err := ctx.Eval("42")
		if err != nil {
			ctx.Close()
			t.Fatalf("Eval error at iteration %d: %v", i, err)
		}

		val, _ := result.Int32()
		if val != 42 {
			ctx.Close()
			t.Fatalf("Result at iteration %d = %d, want 42", i, val)
		}

		ctx.Close()
	}
}

// ============================================================================
// Race Condition Tests (run with -race)
// ============================================================================

func TestRaceMultipleRuntimes(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			rt, err := NewRuntime()
			if err != nil {
				t.Errorf("Goroutine %d: NewRuntime error: %v", id, err)
				return
			}
			defer rt.Close()

			ctx, err := rt.NewContext()
			if err != nil {
				t.Errorf("Goroutine %d: NewContext error: %v", id, err)
				return
			}
			defer ctx.Close()

			for j := range 10 {
				code := fmt.Sprintf("%d * %d", id, j)
				_, err := ctx.Eval(code)
				if err != nil {
					t.Errorf("Goroutine %d: Eval error: %v", id, err)
					return
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestRaceConcurrentReads(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Create an object
	ctx.Eval(`var data = {a: 1, b: 2, c: 3}`)

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10 {
				ctx.Eval("data.a + data.b + data.c")
			}
		}()
	}

	wg.Wait()
}

func TestRaceGoCallback(t *testing.T) {
	rt, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	var mu sync.Mutex
	counter := 0

	fn := ctx.Function("safeIncrement", func(ctx *Context, this Value, args []Value) Value {
		mu.Lock()
		counter++
		mu.Unlock()
		return ctx.Int32(int32(counter))
	})
	ctx.SetGlobal("safeIncrement", fn)

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10 {
				ctx.Eval("safeIncrement()")
			}
		}()
	}

	wg.Wait()

	if counter != 100 {
		t.Errorf("Counter = %d, want 100", counter)
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkEval(b *testing.B) {
	for b.Loop() {
		rt, err := NewRuntime()
		if err != nil {
			b.Fatalf("NewRuntime() error = %v", err)
		}

		ctx, err := rt.NewContext()
		if err != nil {
			rt.Close()
			b.Fatalf("NewContext() error = %v", err)
		}

		result, err := ctx.Eval("1 + 2")
		if err != nil {
			ctx.Close()
			rt.Close()
			b.Fatalf("Eval error = %v", err)
		}
		_ = result.String()

		ctx.Close()
		rt.Close()
	}
}

func BenchmarkEvalComplex(b *testing.B) {
	code := `
		function fib(n) {
			if (n <= 1) return n;
			return fib(n - 1) + fib(n - 2);
		}
		fib(10)
	`

	for b.Loop() {
		rt, err := NewRuntime()
		if err != nil {
			b.Fatalf("NewRuntime() error = %v", err)
		}

		ctx, err := rt.NewContext()
		if err != nil {
			rt.Close()
			b.Fatalf("NewContext() error = %v", err)
		}

		_, err = ctx.Eval(code)
		if err != nil {
			ctx.Close()
			rt.Close()
			b.Fatalf("Eval error = %v", err)
		}

		ctx.Close()
		rt.Close()
	}
}

// BenchmarkEvalReuse benchmarks evaluation with runtime reuse
func BenchmarkEvalReuse(b *testing.B) {
	rt, err := NewRuntime()
	if err != nil {
		b.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		b.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	b.ResetTimer()
	for b.Loop() {
		result, err := ctx.Eval("1 + 2")
		if err != nil {
			b.Fatalf("Eval error = %v", err)
		}
		_ = result.String()
	}
}

// BenchmarkEvalFibonacci benchmarks Fibonacci calculation with reuse
func BenchmarkEvalFibonacci(b *testing.B) {
	rt, err := NewRuntime()
	if err != nil {
		b.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		b.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	// Define the function once
	_, err = ctx.Eval(`function fib(n) { return n <= 1 ? n : fib(n-1) + fib(n-2); }`)
	if err != nil {
		b.Fatalf("Function definition error = %v", err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := ctx.Eval("fib(20)")
		if err != nil {
			b.Fatalf("Eval error = %v", err)
		}
	}
}

// BenchmarkGoCallback benchmarks Go function callbacks
func BenchmarkGoCallback(b *testing.B) {
	rt, err := NewRuntime()
	if err != nil {
		b.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		b.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	fn := ctx.Function("add", func(ctx *Context, this Value, args []Value) Value {
		a, _ := args[0].Int32()
		b, _ := args[1].Int32()
		return ctx.Int32(a + b)
	})
	ctx.SetGlobal("add", fn)

	b.ResetTimer()
	for b.Loop() {
		_, err := ctx.Eval("add(1, 2)")
		if err != nil {
			b.Fatalf("Eval error = %v", err)
		}
	}
}

// BenchmarkObjectCreation benchmarks creating JS objects from Go
func BenchmarkObjectCreation(b *testing.B) {
	rt, err := NewRuntime()
	if err != nil {
		b.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		b.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	b.ResetTimer()
	for b.Loop() {
		obj := ctx.Object()
		obj.Set("x", ctx.Int32(1))
		obj.Set("y", ctx.String("test"))
	}
}

// BenchmarkJSONParse benchmarks JSON parsing
func BenchmarkJSONParse(b *testing.B) {
	rt, err := NewRuntime()
	if err != nil {
		b.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		b.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	b.ResetTimer()
	for b.Loop() {
		_, err := ctx.Eval(`JSON.parse('{"name":"test","value":123,"nested":{"a":1,"b":2}}')`)
		if err != nil {
			b.Fatalf("Eval error = %v", err)
		}
	}
}

// BenchmarkArrayOperations benchmarks array operations
func BenchmarkArrayOperations(b *testing.B) {
	rt, err := NewRuntime()
	if err != nil {
		b.Fatalf("NewRuntime() error = %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		b.Fatalf("NewContext() error = %v", err)
	}
	defer ctx.Close()

	b.ResetTimer()
	for b.Loop() {
		_, err := ctx.Eval(`[1,2,3,4,5].map(x => x * 2).filter(x => x > 4).reduce((a,b) => a+b, 0)`)
		if err != nil {
			b.Fatalf("Eval error = %v", err)
		}
	}
}
