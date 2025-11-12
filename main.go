package main

import (
	"fmt"

	"github.com/Effi-S/go-blockchain/blockchain"
	"rsc.io/quote"
)

func main() {
	chain := blockchain.InitBlockChain()
	chain.AddBlock(quote.Hello())
	chain.AddBlock(quote.Glass())
	chain.AddBlock(quote.Opt())
	chain.AddBlock(quote.Go())

	for i, b := range chain.Blocks {
		fmt.Printf("Block %d: %s\n\t", i, b.Data)
		fmt.Printf("PrevHash: %x\n\t", b.PrevHash)
		fmt.Printf("Hash: %x\n", b.Hash)
	}
}
