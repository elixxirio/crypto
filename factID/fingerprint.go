////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package factID

import (
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/xx_network/crypto/hasher"
)

// Salt used for fact hashing to prevent rainbow table attacks
var factSalt = []byte{212,50,182,4,23,108,57,191,246,154,93,174,127,144,17,96,174,26,209,40,244,231,237,115,117,117,163,115,157,42,65,223}

// Creates a fingerprint of a fact
func Fingerprint(f fact.Fact) []byte {
	h := hasher.BLAKE2.New()
	h.Write([]byte(f.Fact))
	h.Write([]byte(f.T.Stringify()))
	h.Write(factSalt)
	return h.Sum(nil)
}
