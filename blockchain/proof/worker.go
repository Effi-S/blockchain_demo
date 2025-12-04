package proof

import (
	"context"
	"crypto/sha256"
	"math"
	"math/big"
)

// InitDataFunc is a function type that generates data for a given nonce
type InitDataFunc func(nonce int) []byte

// RunWorker runs a worker that performs proof of work
func RunWorker(
	ctx context.Context,
	initData InitDataFunc,
	target *big.Int,
	startNonce int,
	step int,
	result chan<- PowResult,
) {
	var hashInt big.Int
	var hash [32]byte
	for nonce := startNonce; nonce < math.MaxInt64; nonce += step {
		select {
		case <-ctx.Done():
			return // stop immediately
		default:
		}

		data := initData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			result <- PowResult{Nonce: nonce, Hash: hash[:]}
			return
		}
	}
}
