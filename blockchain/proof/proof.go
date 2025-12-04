package proof

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

// Take the data from the block
// create a counter (nonce) which starts at 0
// create a hash of the data plus the counter
// check the hash to see if it meets a set of requirements
// Requirements:
// The first few bytes must contain 0s

const Difficulty = 15

type ProofOfWork struct {
	PrevHash []byte
	Data     []byte
	Target   *big.Int
}

type PowResult struct {
	Nonce int
	Hash  []byte
}

func NewProofOfWork(prevHash []byte, data []byte) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{
		PrevHash: prevHash,
		Data:     data,
		Target:   target,
	}

	return pow
}

func ToHex(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.PrevHash,
		pow.Data,
		ToHex(int64(nonce)),
		ToHex(int64(Difficulty)),
	}, []byte{})

	return data
}

func (pow *ProofOfWork) Run() PowResult {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0
	for ; nonce < math.MaxInt64; nonce++ {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 {
			break
		}
	}
	return PowResult{Nonce: nonce, Hash: hash[:]}
}

func (pow *ProofOfWork) Validate(nonce int) bool {
	var hashInt big.Int
	data := pow.InitData(nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.Target) == -1
}

func (pow *ProofOfWork) RunDistributed(workers int) PowResult {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := make(chan PowResult)

	for i := 0; i < workers; i++ {
		go RunWorker(ctx, pow.InitData, pow.Target, i, workers, resultChan)
	}

	result := <-resultChan
	cancel() // stop all other workers

	return PowResult{Nonce: result.Nonce, Hash: result.Hash}
}
