package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	runtime_pprof "runtime/pprof"
	"time"

	"github.com/Effi-S/go-blockchain/blockchain/block"
	"github.com/Effi-S/go-blockchain/blockchain/config"
	"github.com/Effi-S/go-blockchain/blockchain/proof"
	"rsc.io/quote"
)

func run() {
	start := time.Now()

	chain := block.GetBlockChain()

	chain.AddBlock(quote.Hello())
	chain.AddBlock(quote.Glass())
	chain.AddBlock(quote.Opt())
	chain.AddBlock(quote.Go())

	fmt.Println()
	for i, b := range chain.Blocks {
		fmt.Printf("Block %d: %s\n\t", i, b.Data)
		fmt.Printf("PrevHash: %x\n\t", b.PrevHash)
		fmt.Printf("Hash: %x\n\t", b.Hash)
		pow := proof.NewProofOfWork(b.PrevHash, b.Data)
		fmt.Printf("PoW: %v\n", pow.Validate(b.Nonce))
	}

	fmt.Printf("took %s\n", time.Since(start))
}

// Start pprof HTTP server in a goroutine
func StartPProfServerWithHandler() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", mux))
	}()

	// (Set mux to nil for default handler)
}

func MemoryProfileRun() {
	f, err := os.Create("heap.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	config.SetNumWorkers(20)
	config.SetDifficulty(10)
	fmt.Printf("Running with %d workers\n", config.NumWorkers())
	run()

	runtime_pprof.WriteHeapProfile(f)

	// Now You can run go tool pprof -http=:8080 heap.prof
}

func CPUProfileRun() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}

	config.SetNumWorkers(20)
	config.SetDifficulty(20)
	fmt.Printf("Running with %d workers\n", config.NumWorkers())

	runtime_pprof.StartCPUProfile(f)

	run()

	runtime_pprof.StopCPUProfile()
	f.Close()

	// Now You can run go tool pprof -http=:8080 cpu.prof
}

func main() {
	// 1.  PProf Server
	StartPProfServerWithHandler()
	config.SetNumWorkers(20)
	config.SetDifficulty(30)
	fmt.Printf("Running with %d workers\n", config.NumWorkers())
	run()

	// // 2. CPU Run
	// CPUProfileRun()

	// // 3. RAM Profile
	// MemoryProfileRun()
}
