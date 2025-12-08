# pprof Profiling Guide

This guide demonstrates various ways to use Go's `pprof` tool to profile the blockchain codebase.

## Table of Contents

1. [`pprof` Setup](#pprof-setup)
2. [Profiling Types](#profiling-types)
    - 2a. [CPU Profiling](#cpu-profiling)
    - 2b. [Memory Profiling](#memory-profiling)
    - 2c. [Goroutine Profiling](#goroutine-profiling)
    - 2d. [Block Profiling](#block-profiling)
    - 2e. [Mutex Profiling](#mutex-profiling)
    - 2f. [Trace Profiling](#trace-profiling)
3. [Web UI](#web-interface)
4. [CLI](#command-line-examples)
5. [Profiling in code](#profiling-in-code)

## `pprof` Setup

Add this to your `main.go`:

```go
import (
    "net/http"
    "net/http/pprof"
)

func setupPprof() {
    mux := http.NewServeMux()
    mux.HandleFunc("/debug/pprof/", pprof.Index)
    mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
    mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
    mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
    
    go func() {
		log.Println(http.ListenAndServe("localhost:6060", mux))
	}()
}
```
 [ Note: for default endpoints mux can just be nil]

## Profiling Types

## CPU Profiling

CPU profiling helps identify which functions consume the most CPU time
(In this example this might help optimize the proof-of-work mining.)

#### Using HTTP Endpoint

1. **Start your application** with pprof enabled

   
2. **Analyse via interactive mode**
 while program is still running:
   ```bash
   go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10
   ```

    Alternatively save the profile first:
   ```bash
   curl http://localhost:6060/debug/pprof/profile?seconds=10 > cpu.prof
   go tool pprof cpu.prof
   ```


This Opens a `pprof` terminal

Use `help` to view all commands


Some interactive mode commands:

```
(pprof) top10
(pprof) top10 -cum
(pprof) list RunWorker
(pprof) list Run
(pprof) web
(pprof) png > cpu_graph.png
```

### Example: Profile Proof-of-Work Mining

```bash
go run main.go
# Profile CPU for 60 seconds during mining
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=60

# ... Wait a minute ...

# In pprof:
(pprof) top20
# Shows functions like:
# - crypto/sha256.(*digest).checkSum
# - proof.RunWorker
# - proof.(*ProofOfWork).InitData

# Now If I want to see my packages only 
(pprof) focus=Effi-S/go-blockchain
(pprof) ignore=runtime|crypto|encoding|bytes|math|sync|context
(pprof) top20

```
Example Output:
```bash
      flat  flat%   sum%    cum   cum%
    328.39s  74.94% 74.94%  328.39s 74.94% github.com/Effi-S/go-blockchain/blockchain/proof.RunWorker
    11.23s   2.56%  77.50%  11.23s  2.56% github.com/Effi-S/go-blockchain/blockchain.(*Block).CalculateHash
     8.71s   1.99%  79.49%   8.71s  1.99% github.com/Effi-S/go-blockchain/blockchain.NewBlock
     6.12s   1.40%  80.89%   6.12s  1.40% github.com/Effi-S/go-blockchain/blockchain.(*Blockchain).MineBlock
     3.45s   0.79%  81.68%   3.45s  0.79% github.com/Effi-S/go-blockchain/utils.MarshalBlock
```

Brief explanation of the columns and what they mean:
| Name | Meaning | Example | Example meaning                                                                                                                                                                            
|--------|---------------|--------|------|
| flat   | CPU time spent in this function itself (excluding time in functions it called)                  | 328.39s for `RunWorker`         | 328 seconds executing code directly inside `RunWorker` (not inside `SHA256`, `malloc`, etc.)                                                                      |
| flat%  | `flat` time as a percentage of total profile time.                         | 74.94%                          | Almost ¾ of all CPU time on the machine was spent inside the body of `RunWorker`                                                                                                        |
| sum%   | sum of  `flat` time up to and including this row                     | 74.94% → 77.50% (top 2 lines)   | The top 2 functions alone account for 77.5% of all CPU time                                                                                                                             |
| cum    | Total CPU time spent in the function + everything it called (cumulative)                        | 328.39s for `RunWorker`         | Even though `RunWorker` calls `SHA256`, `bytes.Join`, `encoding/binary.Write`, etc., almost all of that time is spent inside `RunWorker` itself → your loop is extremely tight          |
| cum%   | `cum` time as a percentage of total profile time                                                 | 74.94%                          | Same as `flat%` in this case, meaning `RunWorker` is essentially a leaf function that does almost all the work itself    


## Memory Profiling

Memory profiling helps identify memory leaks and excessive allocations.

My understanding is that this is very common issue with long running code.

We do the same steps as in [CPU Profiling](#cpu-profiling) Just with our `heap` endpoint.


```bash
# Get current memory profile
go tool pprof http://localhost:6060/debug/pprof/heap?seconds=60
# Or 
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof

```
```
(pprof) top20
(pprof) top20 -cum
(pprof) list AddBlock
(pprof) list createBlock
(pprof) png > heap_graph.png
```

### Memory Comparison Example

```bash
# Get baseline
curl http://localhost:6060/debug/pprof/heap > heap_before.prof

# Run some operations (add blocks, mine, etc.)

# Get after profile
curl http://localhost:6060/debug/pprof/heap > heap_after.prof

# Compare
go tool pprof -base=heap_before.prof heap_after.prof
```

## Goroutine Profiling

Goroutine profiling is essential for understanding concurrent execution, especially with the distributed systems / workers.

```bash
# Get current goroutine stack traces
go tool pprof http://localhost:6060/debug/pprof/goroutine?seconds=60
# Or 
curl http://localhost:6060/debug/pprof/goroutine?seconds=60 > goroutine.prof
go tool pprof goroutine.prof
```
```
(pprof) top20
(pprof) list RunWorker
(pprof) traces
(pprof) web
```

### Example: Profile Distributed Workers

```bash
# While running with 20 workers
go tool pprof http://localhost:6060/debug/pprof/goroutine

(pprof) top
# Should show:
# - proof.RunWorker (20 goroutines)
# - runtime.gopark
# - etc.

(pprof) traces
# Shows full stack traces of all goroutines

(pprof) list RunWorker
# Shows where goroutines are spending time
```

### Count Goroutines

```bash
# Get count of goroutines
curl http://localhost:6060/debug/pprof/goroutine?debug=1 | grep "^goroutine" | wc -l

# Or in pprof:
go tool pprof http://localhost:6060/debug/pprof/goroutine
(pprof) top
```

## Block Profiling

Block profiling helps identify contention and blocking operations, useful for understanding synchronization overhead.

### Enable Block Profiling

Add to your code:

```go
import _ "runtime/pprof"

func main() {
    runtime.SetBlockProfileRate(1) // Log every blocking event
    // ... rest of code
}
```

### Using HTTP Endpoint

```bash
curl http://localhost:6060/debug/pprof/block > block.prof
go tool pprof block.prof
```

### Example: Find Blocking Operations

```bash
go tool pprof http://localhost:6060/debug/pprof/block

(pprof) top10
# Shows functions where goroutines are blocked
# Useful for finding channel contention, mutex contention, etc.

(pprof) list RunDistributed
# Check for blocking in distributed mining
```

## Mutex Profiling

Mutex profiling helps identify mutex contention, which can impact performance in concurrent code.

### Enable Mutex Profiling

```go
import _ "runtime/pprof"

func main() {
    runtime.SetMutexProfileFraction(1) // Sample every mutex event
    // ... rest of code
}
```

### Using HTTP Endpoint

```bash
curl http://localhost:6060/debug/pprof/mutex > mutex.prof
go tool pprof mutex.prof
```

### Example: Find Mutex Contention

```bash
go tool pprof http://localhost:6060/debug/pprof/mutex

(pprof) top10
# Shows mutexes with most contention

(pprof) list Get
# Check config.Get() for mutex contention
```

## Trace Profiling

Trace profiling provides a timeline view of program execution, showing goroutine creation, blocking, and scheduling.

1. **Generate a 5-second trace**

    ```bash
    curl http://localhost:6060/debug/pprof/trace?seconds=5 > trace.out
    go tool trace trace.out
    ```
2. Open The Trace viewer that is outputted to the terminal:

    For example: `localhost:39943`

Notice that this uses a different tool than `pprof` because it relies on existing infrastructure of Chrome

3. We can see *simple* **goroutine analysis** table here:
```localhost:39943/goroutines```

4. We can see a trace by clicking one of the traces in  `localhost:39943`

    Play around with `Flow Events` and `Processes` to view interesting info
    
**TODO: Better explain this**

## Web Interface


**View All Profiles** while program is running:

```bash
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=30
```

This opens a web interface at `http://localhost:8080` with:
- Top functions
- Graph view
- Flame graph
- Source view
- Disassembly view

### Flame Graph

```bash
# Generate flame graph
go tool pprof -http=:8080 -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=30
# Navigate to http://localhost:8080/ui/flamegraph
```


## Command-Line Examples

### Quick CPU Profile During Execution

```bash
# Profile for 30 seconds and open interactive mode
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

### Save Profile and Analyze Later

```bash
# Save CPU profile
go tool pprof -proto http://localhost:6060/debug/pprof/profile?seconds=30 > cpu.pb.gz

# Analyze later
go tool pprof cpu.pb.gz
```

### Generate Report

```bash
# Generate text report
go tool pprof -text http://localhost:6060/debug/pprof/profile?seconds=30

# Generate top 20 functions
go tool pprof -top http://localhost:6060/debug/pprof/profile?seconds=30

# Generate tree view
go tool pprof -tree http://localhost:6060/debug/pprof/profile?seconds=30
```

### Compare Profiles

```bash
# Compare two CPU profiles
go tool pprof -base=baseline.prof current.prof

# Compare memory profiles
go tool pprof -base=heap_before.prof heap_after.prof
```

## Common pprof Commands Reference

```
top[N]              - Show top N functions by metric
top[N] -cum         - Show top N functions by cumulative metric
list <function>     - Show annotated source code
web                 - Open graph in browser
png                 - Generate PNG graph
svg                 - Generate SVG graph
pdf                 - Generate PDF graph
traces              - Show all stack traces (goroutine profile)
tree                - Show call tree
disasm <function>   - Show disassembly
peek <regex>        - Show functions matching regex
```

## Profiling in code

### CPU Profile in Code

```go
import (
    "os"
    "runtime/pprof"
)

func profileCPU() {
    f, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Your code here
    run()
}
```

### Memory Profile in Code

```go
import (
    "os"
    "runtime/pprof"
)

func profileMemory() {
    f, err := os.Create("heap.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    // Your code here
    run()
    
    pprof.WriteHeapProfile(f)
}
```

### Custom Profile Labels

```go
import "runtime/pprof"

func profileWithLabels() {
    ctx := pprof.WithLabels(context.Background(), pprof.Labels("operation", "mining"))
    pprof.SetGoroutineLabels(ctx)
    
    // Your code here
    pow.RunDistributed(20)
}
```

## Misc. Examples 

### Example 1: Profile Block Creation

```bash
# Focus on AddBlock and createBlock functions
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

(pprof) list AddBlock
(pprof) list createBlock
(pprof) list NewBlock
```

### Example 2: Analyze Worker Efficiency

```bash
# 1. Enable block profiling (add runtime.SetBlockProfileRate(1))
# 2. Run with 20 workers
# 3. Profile blocks
go tool pprof http://localhost:6060/debug/pprof/block

(pprof) top10
# Check if workers are blocking on channels

(pprof) list RunDistributed
# See where blocking occurs
```


## Additional Resources

- [GopherCon 2019: Two Go Programs, Three Different Profiling Techniques - Dave Cheney](https://www.youtube.com/watch?v=nok0aYiGiYA)
- [Go pprof Documentation](https://golang.org/pkg/net/http/pprof/)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [Dave Cheney's pprof post](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
