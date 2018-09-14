package blockchain

import (
	"crypto/dsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gitlab.com/privategrity/crypto/coin"
	"gitlab.com/privategrity/crypto/shuffle"
	"golang.org/x/crypto/blake2b"
	"os"
	"sync"
)

const BlockHashLenBits = 256
const BlockHashLen = BlockHashLenBits / 8

// Array that holds hashes for the blockchain
type BlockHash [BlockHashLen]byte

var params *dsa.Parameters

func init() {
	params = new(dsa.Parameters)
	if err := dsa.GenerateParameters(params, rand.Reader, dsa.L1024N160); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

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
	signatures   []DsaSigTuple
}

// a structure that holds a blockchain's data for easy serialization and deserialization
type serialBlock struct {
	ID           uint64
	Hash         []byte
	PreviousHash []byte
	TreeRoot     []byte
	Created      [][]byte
	Destroyed    [][]byte
	Signatures   []DsaSigTuple
}

type DsaSigTuple struct {
	R []byte
	S []byte
}

// Generates the origin block for the blockchain
// TODO: get all the members to the team to add something to the origin block
func GenerateOriginBlock() *Block {
	b := Block{lifecycle: Raw}

	b.id = 0

	b.created = append(b.created, coin.Coin{})
	b.destroyed = append(b.destroyed, coin.Coin{})

	b.Bake([]coin.Seed{coin.Seed{}}, BlockHash{}, 8)

	return &b
}

// Creates the next block from the previous one
func (b *Block) NextBlock() (*Block, error) {
	b.mutex.Lock()
	if b.lifecycle != Baked {
		b.mutex.Unlock()
		return &Block{}, ErrBaked
	}

	newBlock := Block{}

	copy(newBlock.previousHash[:], b.hash[:])

	newBlock.id = b.id + 1

	b.mutex.Unlock()

	return &newBlock, nil
}

// Returns a copy of the created coins list
func (b *Block) GetCreated() []coin.Coin {
	b.mutex.Lock()
	cCopy := make([]coin.Coin, len(b.created))
	copy(cCopy, b.created)
	b.mutex.Unlock()
	return cCopy
}

// Adds an element to the created coins list
// Only works while the block is "Raw"
func (b *Block) AddCreated(c []coin.Coin) error {
	b.mutex.Lock()
	if b.lifecycle != Raw {
		b.mutex.Unlock()
		return ErrRaw
	}

	b.created = append(b.created, c...)

	b.mutex.Unlock()

	return nil
}

// Returns a copy of the destroyed coins list
func (b *Block) GetDestroyed() []coin.Coin {
	b.mutex.Lock()
	cCopy := make([]coin.Coin, len(b.destroyed))
	copy(cCopy, b.destroyed)
	b.mutex.Unlock()
	return cCopy
}

// Adds a coin to the destroyed coins list
// Only works while the block is "Raw"
func (b *Block) AddDestroyed(c []coin.Coin) error {
	b.mutex.Lock()
	if b.lifecycle != Raw {
		b.mutex.Unlock()
		return ErrRaw
	}

	b.destroyed = append(b.destroyed, c...)

	b.mutex.Unlock()

	return nil
}

// Returns a copy of the block's hash
func (b *Block) GetHash() (BlockHash, error) {
	b.mutex.Lock()
	if b.lifecycle != Baked {
		b.mutex.Unlock()
		return BlockHash{}, ErrBaked
	}

	var rtnBH BlockHash
	copy(rtnBH[:], b.hash[:])
	b.mutex.Unlock()
	return rtnBH, nil
}

// Returns a copy of the previous Block's hash
func (b *Block) GetPreviousHash() BlockHash {
	var rtnBH BlockHash
	b.mutex.Lock()
	copy(rtnBH[:], b.previousHash[:])
	b.mutex.Unlock()
	return rtnBH
}

// Returns the lifecycle state of the block
func (b *Block) GetLifecycle() BlockLifecycle {
	b.mutex.Lock()
	blc := b.lifecycle
	b.mutex.Unlock()
	return blc
}

// Returns the ID of the block
func (b *Block) GetID() uint64 {
	return b.id
}

// Returns the treeRoot of the block
func (b *Block) GetTreeRoot() (BlockHash, error) {
	b.mutex.Lock()
	if b.lifecycle != Baked {
		b.mutex.Unlock()
		return BlockHash{}, ErrBaked
	}

	var rtnTR BlockHash
	copy(rtnTR[:], b.treeRoot[:])
	b.mutex.Unlock()
	return rtnTR, nil
}

// Permutes the coins and hashes the block
// Only runs if the lifecycle state is "Raw" and sets the state to "Baked"
func (b *Block) Bake(seedList []coin.Seed, treeRoot BlockHash, numkeys int) error {
	b.mutex.Lock()

	if b.lifecycle != Raw {
		b.mutex.Unlock()
		return ErrRaw
	}

	copy(b.treeRoot[:], treeRoot[:])

	//Shuffle the elements
	rawSeed := seedsToBytes(seedList)
	hb, err := blake2b.New256(nil)
	if err != nil {
		return err
	}

	shuffle.ShufflePRNG(rawSeed, len(b.created), func(i, j int) {
		b.created[i], b.created[j] = b.created[j], b.created[i]
	})
	//Hash the seed used for the destroy elements so the two lists aren't shuffled the same way

	hb.Write(rawSeed)
	destroySeed := hb.Sum(nil)
	shuffle.ShufflePRNG(destroySeed, len(b.destroyed), func(i, j int) {
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

	for i := 0; i < numkeys; i++ {
		privatekey := new(dsa.PrivateKey)
		privatekey.PublicKey.Parameters = *params
		dsa.GenerateKey(privatekey, rand.Reader)
		r, s, err := dsa.Sign(rand.Reader, privatekey, hashBytes[:BlockHashLen])
		if err != nil {
			return err
		}
		b.signatures = append(b.signatures,
			DsaSigTuple{r.Bytes(), s.Bytes()})
	}

	//Set the lifecycle to baked
	b.lifecycle = Baked

	b.mutex.Unlock()
	return nil
}

// Serializes the block and outputs a JSON string
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
		Signatures:   b.signatures,
	}

	for indx := range b.created {
		pb.Created = append(pb.Created, b.created[indx][:])
	}

	for indx := range b.destroyed {
		pb.Destroyed = append(pb.Destroyed, b.destroyed[indx][:])
	}

	return json.Marshal(pb)
}

// Converts a serialized block to a block data structure
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

	b.signatures = sb.Signatures

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

//Private Helper Functions
func seedsToBytes(seedList []coin.Seed) []byte {
	var outBytes []byte

	for _, s := range seedList {
		outBytes = append(outBytes, s[:]...)
	}

	return outBytes
}

func coinsToBytes(coinList []coin.Coin) []byte {
	var outBytes []byte

	for _, c := range coinList {
		outBytes = append(outBytes, c[:]...)
	}

	return outBytes
}
