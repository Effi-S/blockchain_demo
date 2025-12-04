package main

import (
	"fmt"
	"time"

	"github.com/Effi-S/go-blockchain/blockchain/block"
	"github.com/Effi-S/go-blockchain/blockchain/proof"
	"rsc.io/quote"
)

func main() {
	start := time.Now()

	numWorkers := 12
	chain := block.GetBlockChain()

	chain.AddBlock(quote.Hello(), numWorkers)
	chain.AddBlock(quote.Glass(), numWorkers)
	chain.AddBlock(quote.Opt(), numWorkers)
	chain.AddBlock(quote.Go(), numWorkers)

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
