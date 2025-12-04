package main

import (
	"fmt"
	"time"

	"github.com/Effi-S/go-blockchain/blockchain/block"
	"github.com/Effi-S/go-blockchain/blockchain/config"
	"github.com/Effi-S/go-blockchain/blockchain/proof"
	"rsc.io/quote"
)

func main() {
	start := time.Now()

	config.Init(config.Default())

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
