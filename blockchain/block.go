////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"gitlab.com/elixxir/crypto/coin"
	"gitlab.com/elixxir/crypto/shuffle"
	"golang.org/x/crypto/blake2b"
	"sync"
)

const BlockHashLenBits = 256
const BlockHashLen = BlockHashLenBits / 8

// Array that holds hashes for the blockchain
type BlockHash [BlockHashLen]byte

// A single block in the blockchain
type Block struct {
	id           uint64
	hash         BlockHash
	previousHash BlockHash
	treeRoot     BlockHash
	created      []coin.Coin
	destroyed    []coin.Coin
	lifecycle    BlockLifecycle
	mutex        sync.Mutex
}

// A structure that holds a block's data, allows for easy serialization and deserialization
type serialBlock struct {
	ID           uint64
	Hash         []byte
	PreviousHash []byte
	TreeRoot     []byte
	Created      [][]byte
	Destroyed    [][]byte
}

// GenerateOriginBlock generates the origin block for the blockchain
func GenerateOriginBlock() *Block {
	b := Block{lifecycle: Raw}

	b.id = 0

	b.created = append(b.created, coin.Coin{})
	b.destroyed = append(b.destroyed, coin.Coin{})

	b.Bake([]coin.Seed{coin.Seed{}}, BlockHash{})

	return &b
}

// NextBlock creates the next block from the previous one
func (b *Block) NextBlock() (*Block, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.lifecycle != Baked {
		return &Block{}, ErrBaked
	}

	newBlock := Block{}

	copy(newBlock.previousHash[:], b.hash[:])

	newBlock.id = b.id + 1

	return &newBlock, nil
}

// GetCreated returns a copy of the created coins list
func (b *Block) GetCreated() []coin.Coin {
	b.mutex.Lock()
	cCopy := make([]coin.Coin, len(b.created))
	copy(cCopy, b.created)
	b.mutex.Unlock()
	return cCopy
}

// AddCreated adds an element to the created coins list
// Only works while the block is "Raw"
func (b *Block) AddCreated(c []coin.Coin) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.lifecycle != Raw {
		return ErrRaw
	}

	b.created = append(b.created, c...)

	return nil
}

// GetDestroyed returns a copy of the destroyed coins list
func (b *Block) GetDestroyed() []coin.Coin {
	b.mutex.Lock()
	cCopy := make([]coin.Coin, len(b.destroyed))
	copy(cCopy, b.destroyed)
	b.mutex.Unlock()
	return cCopy
}

// AddDestroyed adds a coin to the destroyed coins list
// Only works while the block is "Raw"
func (b *Block) AddDestroyed(c []coin.Coin) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.lifecycle != Raw {
		return ErrRaw
	}

	b.destroyed = append(b.destroyed, c...)

	return nil
}

// GetHash returns a copy of the block's hash
func (b *Block) GetHash() (BlockHash, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.lifecycle != Baked {
		return BlockHash{}, ErrBaked
	}

	var rtnBH BlockHash
	copy(rtnBH[:], b.hash[:])
	return rtnBH, nil
}

// GetPreviousHash returns a copy of the previous Block's hash
func (b *Block) GetPreviousHash() BlockHash {
	var rtnBH BlockHash
	b.mutex.Lock()
	copy(rtnBH[:], b.previousHash[:])
	b.mutex.Unlock()
	return rtnBH
}

// GetLifecycle returns the lifecycle state of the block
func (b *Block) GetLifecycle() BlockLifecycle {
	b.mutex.Lock()
	blc := b.lifecycle
	b.mutex.Unlock()
	return blc
}

// GetID returns the ID of the block
func (b *Block) GetID() uint64 {
	return b.id
}

// GetTreeRoot returns the treeRoot of the block
func (b *Block) GetTreeRoot() (BlockHash, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.lifecycle != Baked {
		return BlockHash{}, ErrBaked
	}

	var rtnTR BlockHash
	copy(rtnTR[:], b.treeRoot[:])
	return rtnTR, nil
}

// Bake permutes the coins and hashes the block
// Only runs if the lifecycle state is "Raw" and sets the state to "Baked"
func (b *Block) Bake(seedList []coin.Seed, treeRoot BlockHash) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.lifecycle != Raw {
		return ErrRaw
	}

	copy(b.treeRoot[:], treeRoot[:])

	//Shuffle the elements
	rawSeed := seedsToBytes(seedList)
	hb, err := blake2b.New256(nil)
	if err != nil {
		return err
	}

	shuffle.ShuffleSwap(rawSeed, len(b.created), func(i, j int) {
		b.created[i], b.created[j] = b.created[j], b.created[i]
	})

	//Hash the seed used for the destroy elements so the two lists aren't shuffled the same way
	hb.Write(rawSeed)
	destroySeed := hb.Sum(nil)
	shuffle.ShuffleSwap(destroySeed, len(b.destroyed), func(i, j int) {
		b.destroyed[i], b.destroyed[j] = b.destroyed[j], b.destroyed[i]
	})

	//Hash the elements
	h := sha256.New()
	h.Write(b.previousHash[:])
	h.Write(treeRoot[:])
	h.Write(coinsToBytes(b.created))
	h.Write(coinsToBytes(b.destroyed))
	hashBytes := h.Sum(nil)

	copy(b.hash[:BlockHashLen], hashBytes[:BlockHashLen])

	//Set the lifecycle to baked
	b.lifecycle = Baked

	return nil
}

// Serialize serializes the block and outputs a JSON string
// Only runs when the block is "Baked"
func (b *Block) Serialize() ([]byte, error) {
	if b.lifecycle != Baked {
		return []byte{}, ErrBaked
	}

	pb := serialBlock{
		ID:           b.id,
		Hash:         b.hash[:],
		PreviousHash: b.previousHash[:],
		TreeRoot:     b.treeRoot[:],
	}

	for indx := range b.created {
		pb.Created = append(pb.Created, b.created[indx][:])
	}

	for indx := range b.destroyed {
		pb.Destroyed = append(pb.Destroyed, b.destroyed[indx][:])
	}

	return json.Marshal(pb)
}

// Deserialize converts a serialized block to a block data structure
func Deserialize(sBlock []byte) (*Block, error) {
	sb := &serialBlock{}

	err := json.Unmarshal(sBlock, sb)

	if err != nil {
		return nil, err
	}

	b := Block{}

	b.mutex.Lock()

	copy(b.hash[:], sb.Hash)

	copy(b.previousHash[:], sb.PreviousHash)

	copy(b.treeRoot[:], sb.TreeRoot)

	for i := range sb.Created {
		newCoin := coin.Coin{}
		copy(newCoin[:], sb.Created[i])
		b.created = append(b.created, newCoin)
	}

	for i := range sb.Destroyed {
		newCoin := coin.Coin{}
		copy(newCoin[:], sb.Destroyed[i])
		b.destroyed = append(b.destroyed, newCoin)
	}

	b.id = sb.ID

	b.lifecycle = Baked
	b.mutex.Unlock()

	return &b, nil
}

// Private Helper Functions

// seedsToBytes serializes seeds into byte slices
func seedsToBytes(seedList []coin.Seed) []byte {
	var outBytes []byte

	for _, s := range seedList {
		outBytes = append(outBytes, s[:]...)
	}

	return outBytes
}

// coindsToBytes serializes coins into byte slices
func coinsToBytes(coinList []coin.Coin) []byte {
	var outBytes []byte

	for _, c := range coinList {
		outBytes = append(outBytes, c[:]...)
	}

	return outBytes
}
