package block

import (
	"sync"

	"github.com/Effi-S/go-blockchain/blockchain/config"
	"github.com/Effi-S/go-blockchain/blockchain/proof"
)

var (
	instance *BlockChain
	once     sync.Once
)

type BlockChain struct {
	Blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// NewBlock creates a new block with proof of work.
// The block can only be created with a valid proof of work result.
func NewBlock(data string, prevHash []byte, nonce int, hash []byte) *Block {
	return &Block{
		Data:     []byte(data),
		PrevHash: prevHash,
		Nonce:    nonce,
		Hash:     hash,
	}
}

// createBlock creates a block by performing proof of work.
// This is a convenience function that ensures blocks are only created with proof of work.
func createBlock(data string, prevHash []byte) *Block {
	numWorkers := config.NumWorkers()
	pow := proof.NewProofOfWork(prevHash, []byte(data))
	var powRes proof.PowResult
	if numWorkers == 1 {
		powRes = pow.Run()
	} else {
		powRes = pow.RunDistributed(numWorkers)
	}
	return NewBlock(data, prevHash, powRes.Nonce, powRes.Hash)
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]

	newBlock := createBlock(data, prevBlock.Hash)

	chain.Blocks = append(chain.Blocks, newBlock)
}

func GetBlockChain() *BlockChain {
	once.Do(func() {
		genesis := createBlock("Genesis", []byte{})
		instance = &BlockChain{Blocks: []*Block{genesis}}
	})
	return instance
}
