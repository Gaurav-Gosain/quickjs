// Example: Concurrent Usage
//
// This example demonstrates using QuickJS-ng in concurrent Go programs:
// - Thread-safe runtime access
// - Multiple runtimes in parallel
// - Worker pool pattern
// - Safe concurrent evaluation
//
// Note: QuickJS itself is single-threaded, but the Go bindings provide
// mutex protection allowing safe access from multiple goroutines.
//
// Run with: go run ./examples/concurrent
package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Gaurav-Gosain/quickjs"
)

func main() {
	// === Single Runtime, Multiple Goroutines ===
	fmt.Println("=== Single Runtime, Multiple Goroutines ===")
	singleRuntimeDemo()

	// === Multiple Independent Runtimes ===
	fmt.Println("\n=== Multiple Independent Runtimes ===")
	multipleRuntimesDemo()

	// === Worker Pool Pattern ===
	fmt.Println("\n=== Worker Pool Pattern ===")
	workerPoolDemo()

	// === Concurrent Calculations ===
	fmt.Println("\n=== Concurrent Calculations ===")
	concurrentCalculationsDemo()
}

// singleRuntimeDemo shows that a single runtime can be safely accessed
// from multiple goroutines (operations are serialized internally).
func singleRuntimeDemo() {
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

	// Set up a shared counter
	ctx.Eval(`var counter = 0`)

	var wg sync.WaitGroup
	numGoroutines := 10
	incrementsPerGoroutine := 100

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				ctx.Eval(`counter++`)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	result, _ := ctx.Eval(`counter`)
	fmt.Printf("Final counter value: %s (expected: %d)\n",
		result.String(), numGoroutines*incrementsPerGoroutine)
	fmt.Printf("Time: %v\n", elapsed)
}

// multipleRuntimesDemo shows running multiple independent runtimes in parallel.
// Each runtime has its own memory space and can run truly in parallel.
func multipleRuntimesDemo() {
	numRuntimes := 4
	var wg sync.WaitGroup
	results := make([]string, numRuntimes)

	start := time.Now()

	for i := 0; i < numRuntimes; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Each goroutine gets its own runtime
			rt, err := quickjs.NewRuntime()
			if err != nil {
				results[id] = fmt.Sprintf("Error: %v", err)
				return
			}
			defer rt.Close()

			ctx, err := rt.NewContext()
			if err != nil {
				results[id] = fmt.Sprintf("Error: %v", err)
				return
			}
			defer ctx.Close()

			// Do some work
			code := fmt.Sprintf(`
				(() => {
					let sum = 0;
					for (let i = 0; i < 10000; i++) {
						sum += i * %d;
					}
					return sum;
				})()
			`, id+1)

			result, err := ctx.Eval(code)
			if err != nil {
				results[id] = fmt.Sprintf("Error: %v", err)
			} else {
				results[id] = result.String()
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("Results from %d parallel runtimes:\n", numRuntimes)
	for i, r := range results {
		fmt.Printf("  Runtime %d: %s\n", i, r)
	}
	fmt.Printf("Total time: %v\n", elapsed)
}

// workerPoolDemo demonstrates a worker pool pattern where multiple
// goroutines process jobs using a shared pool of JS runtimes.
func workerPoolDemo() {
	numWorkers := 3
	numJobs := 10

	// Create a pool of runtimes
	type jsWorker struct {
		rt  *quickjs.Runtime
		ctx *quickjs.Context
	}

	workers := make([]*jsWorker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		rt, _ := quickjs.NewRuntime()
		ctx, _ := rt.NewContext()
		workers[i] = &jsWorker{rt: rt, ctx: ctx}
	}
	defer func() {
		for _, w := range workers {
			w.ctx.Close()
			w.rt.Close()
		}
	}()

	// Job channel
	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int, w *jsWorker) {
			defer wg.Done()
			for job := range jobs {
				code := fmt.Sprintf(`"Worker %d processed job %d: " + (%d * %d)`,
					workerID, job, job, job)
				result, _ := w.ctx.Eval(code)
				results <- result.String()
			}
		}(i, workers[i])
	}

	// Send jobs
	start := time.Now()
	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	// Wait and collect results
	wg.Wait()
	close(results)
	elapsed := time.Since(start)

	fmt.Printf("Worker pool results (%d workers, %d jobs):\n", numWorkers, numJobs)
	for r := range results {
		fmt.Printf("  %s\n", r)
	}
	fmt.Printf("Total time: %v\n", elapsed)
}

// concurrentCalculationsDemo shows performing parallel calculations
// and combining results.
func concurrentCalculationsDemo() {
	type calculation struct {
		name string
		code string
	}

	calculations := []calculation{
		{"Fibonacci(25)", "(() => { const fib = n => n <= 1 ? n : fib(n-1) + fib(n-2); return fib(25); })()"},
		{"Sum 1-10000", "(() => { let s = 0; for (let i = 1; i <= 10000; i++) s += i; return s; })()"},
		{"Prime count to 1000", `(() => {
			const isPrime = n => {
				if (n < 2) return false;
				for (let i = 2; i * i <= n; i++) if (n % i === 0) return false;
				return true;
			};
			let count = 0;
			for (let i = 2; i <= 1000; i++) if (isPrime(i)) count++;
			return count;
		})()`},
		{"Factorial(15)", "(() => { let f = 1; for (let i = 2; i <= 15; i++) f *= i; return f; })()"},
	}

	results := make([]struct {
		name    string
		result  string
		elapsed time.Duration
	}, len(calculations))

	var wg sync.WaitGroup
	start := time.Now()

	for i, calc := range calculations {
		wg.Add(1)
		go func(idx int, c calculation) {
			defer wg.Done()

			calcStart := time.Now()

			rt, _ := quickjs.NewRuntime()
			defer rt.Close()
			ctx, _ := rt.NewContext()
			defer ctx.Close()

			result, err := ctx.Eval(c.code)
			calcElapsed := time.Since(calcStart)

			if err != nil {
				results[idx] = struct {
					name    string
					result  string
					elapsed time.Duration
				}{c.name, "Error: " + err.Error(), calcElapsed}
			} else {
				results[idx] = struct {
					name    string
					result  string
					elapsed time.Duration
				}{c.name, result.String(), calcElapsed}
			}
		}(i, calc)
	}

	wg.Wait()
	totalElapsed := time.Since(start)

	fmt.Println("Concurrent calculation results:")
	var sumIndividual time.Duration
	for _, r := range results {
		fmt.Printf("  %-25s = %-15s (took %v)\n", r.name, r.result, r.elapsed)
		sumIndividual += r.elapsed
	}
	fmt.Printf("\nSequential time would be: ~%v\n", sumIndividual)
	fmt.Printf("Actual parallel time:     %v\n", totalElapsed)
	fmt.Printf("Speedup:                  %.2fx\n", float64(sumIndividual)/float64(totalElapsed))
}
