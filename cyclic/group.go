////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package cyclic

import (
	"crypto/sha256"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/crypto/large"
)

// Groups provide cyclic int operations that keep the return values confined to
// a finite field under modulo p
//TODO: EVENTUALLY WE NEED TO UPDATE THIS STRUCT AND REMOVE RAND, SEED, RNG, ETC... this is way too complex
type Group struct {
	prime       large.Int
	psub1       large.Int
	psub2       large.Int
	psub3       large.Int
	psub1factor large.Int
	seed        large.Int
	random      large.Int
	zero        large.Int
	one         large.Int
	two         large.Int
	G           large.Int
	rng         Random
	fingerprint uint64
}

const GroupFingerprintSize = 8

// NewGroup returns a group with the given prime, seed and generator
func NewGroup(p large.Int, s large.Int, g large.Int, rng Random) Group {
	h := sha256.New()
	h.Write(p.Bytes())
	h.Write(g.Bytes())

	hashVal := h.Sum(nil)[:GroupFingerprintSize]
	value := large.NewIntFromBytes(hashVal)
	return Group{
		prime:       p,
		psub1:       large.NewInt(0).Sub(p, large.NewInt(1)),
		psub2:       large.NewInt(0).Sub(p, large.NewInt(2)),
		psub3:       large.NewInt(0).Sub(p, large.NewInt(3)),
		psub1factor: large.NewInt(0).RightShift(large.NewInt(0).Sub(p, large.NewInt(1)), 1),

		seed:        s,
		random:      large.NewInt(0),
		zero:        large.NewInt(0),
		one:         large.NewInt(1),
		two:         large.NewInt(2),
		G:           g,
		rng:         rng,
		fingerprint: value.Uint64(),
	}
}

// Constructors for cyclicInt

// Create a new cyclicInt in the group from an int64 value
func (g *Group) NewInt(x int64) *Int {
	val := large.NewInt(x)
	n := &Int{value: val, fingerprint: g.fingerprint}
	return n
}

// Create a new cyclicInt in the group from a large.Int value
func (g *Group) NewIntFromLargeInt(x large.Int) *Int {
	n := &Int{value: x, fingerprint: g.fingerprint}
	return n
}

// Create a new cyclicInt in the group from a byte buffer
func (g *Group) NewIntFromBytes(buf []byte) *Int {
	val := large.NewIntFromBytes(buf)
	n := &Int{value: val, fingerprint: g.fingerprint}
	return n
}

// Create a new cyclicInt in the group from a string using the passed base
// returns nil if string cannot be parsed
func (g *Group) NewIntFromString(str string, base int) *Int {
	val := large.NewIntFromString(str, base)
	if val == nil {
		return nil
	}
	n := &Int{value: val, fingerprint: g.fingerprint}
	return n
}

// Create a new cyclicInt in the group with the Max4KBit value
func (g *Group) NewMaxInt() *Int {
	val := large.NewMaxInt()
	n := &Int{value: val, fingerprint: g.fingerprint}
	return n
}

// Create a new cyclicInt in the group from an uint64 value
func (g *Group) NewIntFromUInt(i uint64) *Int {
	val := large.NewIntFromUInt(i)
	n := &Int{value: val, fingerprint: g.fingerprint}
	return n
}

// Check if all cyclic Ints belong to the group and panic otherwise
func (g *Group) checkInts(ints ...*Int) {
	for _, i := range ints {
		if i.GetGroupFingerprint() != g.fingerprint {
			jww.FATAL.Panicf("cyclicInt being used in wrong group! "+
				"Group fingerprint is %d and cyclicInt has %d",
				g.fingerprint, i.GetGroupFingerprint())
		}
	}
}

// Get group fingerprint
func (g *Group) GetFingerprint() uint64 {
	return g.fingerprint
}

// Mul multiplies a and b within the group, putting the result in c
// and returning c
func (g *Group) Mul(a, b, c *Int) *Int {
	g.checkInts(a, b, c)
	c.value.Mod(c.value.Mul(a.value, b.value), g.prime)
	return c
}

// Inside returns true of the Int is within the group, false if it isn't
func (g *Group) Inside(a *Int) bool {
	g.checkInts(a)
	return a.value.Cmp(g.zero) == 1 && a.value.Cmp(g.prime) == -1
}

// ModP sets z ≡ x mod prime within the group and returns z.
func (g Group) ModP(x, z *Int) *Int {
	g.checkInts(x, z)
	z.value.Mod(x.value, g.prime)
	return z
}

// Inverse sets b equal to the inverse of a within the group and returns b
func (g *Group) Inverse(a, b *Int) *Int {
	g.checkInts(a, b)
	b.value.ModInverse(a.value, g.prime)
	return b
}

// SetSeed sets a seed for use in random number generation
func (g *Group) SetSeed(k large.Int) {
	g.seed = k
}

// Random securely generates a random number within the group and sets r
// equal to it.
func (g *Group) Random(r *Int) *Int {
	g.checkInts(r)
	r.value.Add(g.seed, g.rng.Rand(g.random))
	r.value.Mod(r.value, g.psub2)
	r.value.Add(r.value, g.two)
	if !g.Inside(r) {
		jww.FATAL.Panicf("Random int is not in cyclic group: %s",
			r.value.TextVerbose(16, 0))
	}
	return r
}

// GetP returns a copy of the group's prime
func (g *Group) GetP() large.Int {
	n := large.NewInt(0)
	n.Set(g.prime)
	return n
}

// GetPSub1 returns a copy of the group's p-1
func (g *Group) GetPSub1() large.Int {
	n := large.NewInt(0)
	n.Set(g.psub1)
	return n
}

// GroupMul Multiplies all ints in the passed slice slc together and
// places the result in c
func (g Group) ArrayMul(slc []*Int, c *Int) *Int {
	g.checkInts(c)
	c.value.SetString("1", 10)

	for _, islc := range slc {
		g.checkInts(islc)
		g.Mul(c, islc, c)
	}

	return c
}

// Exp sets z = x**y mod p, and returns z.
func (g Group) Exp(x, y, z *Int) *Int {
	g.checkInts(x, y, z)
	z.value.Exp(x.value, y.value, g.prime)
	return z
}

// RandomCoprime randomly generates coprimes in the group (coprime
// against g.prime-1)
func (g *Group) RandomCoprime(r *Int) *Int {
	g.checkInts(r)
	for r.value.Set(g.psub1); !r.value.IsCoprime(g.psub1); {
		r.value.Add(g.seed, g.rng.Rand(g.random))
		r.value.Mod(r.value, g.psub3)
		r.value.Add(r.value, g.two)
	}
	return r
}

// RootCoprime sets z = y√x mod p, and returns z. Only works with y's
// coprime with g.prime-1 (g.psub1)
func (g Group) RootCoprime(x, y, z *Int) *Int {
	g.checkInts(x, y, z)
	z.value.ModInverse(y.value, g.psub1)
	g.Exp(x, z, z)
	return z
}

// Finds a number coprime with p-1 and who's modular exponential inverse is
// the number of prescribed bits value. Bits must be greater than 1.
// Only works when the prime is safe or strong
// Using a smaller bytes length is acceptable because modular logarithm algorithm's
// complexities derive primarily from the size of the prime defining the group
// not the size of the exponent.  More information can be found here:
// TODO: add link to doc
// The function will panic if bits >= log2(g.prime), so the caller MUST use
// a correct value of bits

func (g Group) FindSmallCoprimeInverse(z *Int, bits uint32) *Int {
	if bits >= uint32(g.prime.BitLen()) {
		jww.FATAL.Panicf("Requested bits: %d is greater than"+
			" or equal to group's prime: %d", bits, g.prime.BitLen())
	}

	g.checkInts(z)
	// RNG that ensures the output is an odd number between 2 and 2^(
	// bit*8) that is not equal to p-1/2.  This must occur because for a proper
	// modular inverse to exist within a group a number must have no common
	// factors with the number that defines the group.  Normally that would not
	// be a problem because the number that defines the group normally is a prime,
	// but we are inverting within a group defined by the even number p-1 to find the
	// modular exponential inverse, so the number must be chozen from a different set
	max := large.NewInt(0).Sub(
		large.NewInt(0).LeftShift(
			large.NewInt(1),
			uint(bits)-1),
		large.NewInt(1))
	rng := NewRandom(large.NewInt(2), max)

	for true {
		zinv := large.NewInt(0).Or(
			large.NewInt(0).LeftShift(
				rng.Rand(large.NewInt(0)),
				1),
			large.NewInt(1))

		// p-1 has one odd factor, (p-1)/2,  we must check that the generated number is not that
		if zinv.Cmp(g.psub1factor) == 0 {
			continue
		}

		//Modulo inverse zinv and check that the inverse exists
		if z.value.ModInverse(zinv, g.psub1) == nil {
			continue
		}

		zbytes := z.Bytes()

		// Checks if the lowest bit is 1, implying the value is odd.
		// Due to the fact that p is a safe prime, this means the value is
		// coprime with p minus 1 because its only has one odd factor, which is
		// also checked

		if zbytes[len(zbytes)-1]&0x01 == 1 {
			if zinv.Cmp(g.psub1factor) != 0 {
				break
			}
		}

	}

	return z
}
