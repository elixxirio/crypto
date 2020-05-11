////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

// Contains a generic signing interface and implementations to sign the data
// as well as verify the signature

package signature

import (
	"crypto"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/crypto/csprng"
	"gitlab.com/elixxir/crypto/signature/rsa"
)

// Interface for signing generically
type GenericSignable interface {
	String() string // Designed to be identical to String() in grpc
	GetSig() []byte
	SetSig(newSignature []byte) error
	GetNonce() []byte
	SetNonce(newNonce []byte) error
	ClearSig()
}

// Sign takes a genericSignable object, marshals the data intended to be signed.
// It hashes that data and sets it as the signature of that object
func Sign(signable GenericSignable, privKey *rsa.PrivateKey) error {
	// Clear any value in the signature field
	signable.ClearSig()

	// Create rand for signing and nonce generation
	rand := csprng.NewSystemRNG()

	// Generate nonce
	ourNonce := make([]byte, 32)
	_, err := rand.Read(ourNonce)
	if err != nil {
		return errors.Errorf("Failed to generate nonce: %+v", err)
	}

	// Set nonce
	err = signable.SetNonce(ourNonce)
	if err != nil {
		return errors.Errorf("Unable to set nonce: %+v", err)
	}

	// Get the data that is to be signed (including nonce)
	data := signable.String()

	// Prepare to hash the data
	sha := crypto.SHA256
	h := sha.New()
	h.Write([]byte(data))

	// Sign the message
	signature, err := rsa.Sign(rand, privKey, sha, h.Sum(nil), nil)
	if err != nil {
		return errors.Errorf("Unable to sign message: %+v", err)
	}

	// And set the signature
	err = signable.SetSig(signature)
	if err != nil {
		return errors.Errorf("Unable to set signature: %+v", err)
	}

	return nil
}

// Verify takes the signature from the object and clears it out.
// It then re-creates the signature and compares it to the original signature.
// If the recreation matches the original signature it returns true,
// else it returns false
func Verify(verifiable GenericSignable, pubKey *rsa.PublicKey) error {
	// Take the signature from the object
	sig := verifiable.GetSig()

	// Clear the signature
	verifiable.ClearSig()

	// Get the data to replicate the signature
	data := verifiable.String()

	// Hash the data
	sha := crypto.SHA256
	h := sha.New()
	h.Write([]byte(data))
	ourHash := h.Sum(nil)

	// Reset signature so verify is not destructive
	err := verifiable.SetSig(sig)
	if err != nil {
		return errors.Errorf("Unable to reset signature: %+v", err)
	}

	// Verify the signature using our implementation
	err = rsa.Verify(pubKey, sha, ourHash, sig, nil)

	// And check for an error
	if err != nil {
		// If there is an error, then signature is invalid
		return err
	}
	// Otherwise it has been verified
	return nil

}
