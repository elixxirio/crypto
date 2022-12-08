////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package dh

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNike(t *testing.T) {
	alicePrivateKey, alicePublicKey := DHNIKE.NewKeypair()
	bobPrivateKey, bobPublicKey := DHNIKE.NewKeypair()

	secret1 := alicePrivateKey.DeriveSecret(bobPublicKey)
	secret2 := bobPrivateKey.DeriveSecret(alicePublicKey)

	require.Equal(t, secret1, secret2)
}

func TestPrivateKeyMarshaling(t *testing.T) {
	alicePrivateKey, _ := DHNIKE.NewKeypair()

	alicePrivateKeyBytes := alicePrivateKey.Bytes()
	alice2PrivateKey, _ := DHNIKE.NewKeypair()

	err := alice2PrivateKey.FromBytes(alicePrivateKeyBytes)
	require.NoError(t, err)

	alice2PrivateKeyBytes := alice2PrivateKey.Bytes()

	require.Equal(t, alice2PrivateKeyBytes, alicePrivateKeyBytes)

	alice3PrivateKey, err := DHNIKE.UnmarshalBinaryPrivateKey(alice2PrivateKeyBytes)
	require.NoError(t, err)

	alice3PrivateKeyBytes := alice3PrivateKey.Bytes()

	require.Equal(t, alice3PrivateKeyBytes, alice2PrivateKeyBytes)
	require.Equal(t, len(alice3PrivateKeyBytes), DHNIKE.PrivateKeySize())
}

func TestPublicKeyMarshaling(t *testing.T) {
	_, alicePublicKey := DHNIKE.NewKeypair()

	alicePublicKeyBytes := alicePublicKey.Bytes()
	_, alice2PublicKey := DHNIKE.NewKeypair()

	err := alice2PublicKey.FromBytes(alicePublicKeyBytes)
	require.NoError(t, err)

	alice2PublicKeyBytes := alice2PublicKey.Bytes()

	require.Equal(t, alice2PublicKeyBytes, alicePublicKeyBytes)

	alice3PublicKey, err := DHNIKE.UnmarshalBinaryPublicKey(alice2PublicKeyBytes)
	require.NoError(t, err)

	alice3PublicKeyBytes := alice3PublicKey.Bytes()

	require.Equal(t, alice3PublicKeyBytes, alice2PublicKeyBytes)
	require.Equal(t, len(alice3PublicKeyBytes), DHNIKE.PublicKeySize())
}

func TestPublicKey_Reset(t *testing.T) {
	_, alicePublicKey := DHNIKE.NewKeypair()
	alicePublicKey.Reset()
	if alicePublicKey.Bytes() != nil {
		t.Fatalf("After reset, key should be nil!")
	}
}

func TestPrivateKey_Reset(t *testing.T) {
	alicePrivKey, _ := DHNIKE.NewKeypair()
	alicePrivKey.Reset()
	if alicePrivKey.Bytes() != nil {
		t.Fatalf("After reset, key should be nil!")
	}

}

func TestPrivateKey_Scheme(t *testing.T) {
	alicePrivKey, _ := DHNIKE.NewKeypair()

	if !reflect.DeepEqual(alicePrivKey.Scheme(), DHNIKE) {
		t.Fatalf("GetScheme failed to retrieve expected value")
	}
}

func TestPublicKey_Scheme(t *testing.T) {
	_, alicePubKey := DHNIKE.NewKeypair()

	if !reflect.DeepEqual(alicePubKey.Scheme(), DHNIKE) {
		t.Fatalf("GetScheme failed to retrieve expected value")
	}

}