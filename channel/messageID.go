package channel

import (
	"bytes"
	"encoding/base64"
	"golang.org/x/crypto/blake2b"

	jww "github.com/spf13/jwalterweatherman"

)

const (
	MessageIDLen  = 32
	messageIDSalt = "ChannelsMessageIdSalt"
)

type MessageID [MessageIDLen]byte

// MakeMessageID returns the ID for the given serialized message
// Due to the fact that messages contain the round they are sent in,
// they are replay resistant. This property, when combined with the collision
// resistance of the hash function, ensures that an adversary will not be able
// to cause multiple messages to have the same ID
func MakeMessageID(message []byte)MessageID{
	h, err := blake2b.New256(nil)
	if err!=nil{
		jww.FATAL.Panicf("Failed to get ")
	}
	h.Write(message)
	h.Write([]byte(messageIDSalt))
	midBytes := h.Sum(nil)

	mid := MessageID{}
	copy(mid[:],midBytes)
	return mid
}

// Equals checks if two message IDs which are the same
//Not constant time
func (mid MessageID)Equals(mid2 MessageID)bool{
	return bytes.Equal(mid[:],mid2[:])
}

// String prints a base64 encoded message ID for debugging
// Adheres to the go stringer interface
func (mid MessageID)String()string{
	return "ChMsgID-" + base64.StdEncoding.EncodeToString(mid[:])
}
