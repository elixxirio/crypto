////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package e2e

import (
	"bytes"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/primitives/format"
	"gitlab.com/xx_network/primitives/id"
	"math/rand"
	"testing"
)

const messageSize = 150

// Tests if IsUnencrypted() correctly determines an encrypted message as
// encrypted.
func TestIsUnencrypted_EncryptedMessage(t *testing.T) {
	// Generate random byte slice
	randSlice := make([]byte, messageSize)
	rand.Read(randSlice)
	fpSlice := make([]byte, format.KeyFPLen)
	rand.Read(fpSlice)
	fpSlice[0] &= 0x7f
	macSlice := make([]byte, format.MacLen)
	rand.Read(macSlice)
	macSlice[0] &= 0x7f

	// Create message
	m := format.NewMessage(messageSize)
	// Set message payload
	m.SetPayloadA(randSlice)
	m.SetPayloadB(randSlice)

	//set the fingerprint
	fp := format.Fingerprint{}
	copy(fp[:], fpSlice)
	m.SetKeyFP(fp)

	// Set the MAC
	m.SetMac(macSlice)

	// Check the message
	unencrypted, uid := IsUnencrypted(m)

	if unencrypted == true {
		t.Errorf("IsUnencrypted() determined the message is "+
			"unencrypted when it is actually encrypted"+
			"\n\treceived: %v\n\texpected: %v",
			unencrypted, false)
	}

	if uid != nil {
		t.Errorf("IsUnencrypted() should not return a user id on an" +
			"encrypted message")
	}
}

// Tests if IsUnencrypted() correctly determines an unencrypted message as
// unencrypted.
func TestIsUnencrypted_UnencryptedMessage(t *testing.T) {
	// Generate random byte slice
	randSlice := make([]byte, messageSize)
	rand.Read(randSlice)
	fpSlice := make([]byte, format.KeyFPLen)
	rand.Read(fpSlice)
	fpSlice[0] &= 0x7f

	// Create message
	m := format.NewMessage(messageSize)

	// Set message payload
	m.SetPayloadA(randSlice)
	m.SetPayloadB(randSlice)

	// Create new hash
	h, _ := hash.NewCMixHash()
	h.Write(m.GetContents())
	payloadHash := h.Sum(nil)
	payloadHash[0] &= 0x3F

	// Set the MAC with the high bit from the fingerprint as the ID
	m.SetMac(payloadHash)

	//set the fingerprint
	fp := format.Fingerprint{}
	copy(fp[:], fpSlice)
	m.SetKeyFP(fp)

	// Check the message
	unencrypted, uid := IsUnencrypted(m)

	if unencrypted == false {
		t.Errorf("IsUnencrypted() determined the message is encrypted when it is actually unencrypted"+
			"\n\treceived: %v\n\texpected: %v",
			unencrypted, true)
	}

	expectedUID := id.ID{}
	copy(expectedUID[:], fpSlice[:])
	expectedUID[0] |= (payloadHash[0] & 0b01000000) << 1
	expectedUID.SetType(id.User)

	if !bytes.Equal(uid[:], expectedUID[:]) {
		t.Errorf("IsUnencrypted() returned the wrong userID"+
			"\n\treceived: %s\n\texpected: %s",
			uid, expectedUID)
	}
}

// Tests if SetUnencrypted() makes a message unencrypted according to
// IsUnencrypted().
func TestSetUnencrypted(t *testing.T) {
	// Generate random byte slice
	randSlice := make([]byte, messageSize)
	rand.Read(randSlice)
	macSlice := make([]byte, format.KeyFPLen)
	rand.Read(macSlice)
	macSlice[0] &= 0x7f

	// Create message
	m := format.NewMessage(messageSize)

	// Set message payload
	m.SetPayloadA(randSlice)
	m.SetPayloadB(randSlice)

	// Set the MAC
	m.SetMac(macSlice)

	uid := id.ID{}
	rand.Read(uid[:32])
	uid.SetType(id.User)

	SetUnencrypted(m, &uid)

	encrypted, rtnUid := IsUnencrypted(m)

	if encrypted == false {
		t.Errorf("SetUnencrypted() determined the message is encrypted"+
			" when it should be unencrypted\n\treceived: %v\n\texpected: %v",
			encrypted, true)
	}

	if !bytes.Equal(uid[:], rtnUid[:]) {
		t.Errorf("IsUnencrypted() returned the wrong userID"+
			"\n\treceived: %s\n\texpected: %s",
			rtnUid, uid)
	}
}
