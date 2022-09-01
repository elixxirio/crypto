package broadcast

import (
	"encoding/json"
	"errors"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"

	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"

	"gitlab.com/elixxir/crypto/cmix"
)

const hkdfInfo = "XX_Network_Broadcast_Channel_HKDF_Blake2b"

// Channel is a multicast communication channel that retains the
// various privacy notions that this mix network provides.
type Channel struct {
	ReceptionID *id.ID
	Name        string
	Description string
	Salt        []byte
	RsaPubKey   *rsa.PublicKey
	Secret      []byte

	// Only appears in memory, is not contained in the marshalled version.
	// Lazily evaluated on first use.
	// key = H(ReceptionID)
	key []byte
}

func NewChannel(name, description string, rng csprng.Source) (*Channel, *rsa.PrivateKey, error) {
	pk, err := rsa.GenerateKey(rng, 4096)
	if err != nil {
		return nil, nil, err
	}
	salt := cmix.NewSalt(rng, 512)

	secret := make([]byte, 32)
	_, err = rng.Read(secret)
	if err != nil {
		panic(err)
	}

	channelID, key, err := NewChannelID(name, description, salt, rsa.CreatePublicKeyPem(pk.GetPublic()), secret)
	if err != nil {
		return nil, nil, err
	}

	return &Channel{
		ReceptionID: channelID,
		Name:        name,
		Description: description,
		Salt:        salt,
		RsaPubKey:   pk.GetPublic(),
		key:         key,
		Secret:      secret,
	}, pk, nil
}

func UnmarshalChannel(data []byte) (*Channel, error) {
	var c Channel
	return &c, json.Unmarshal(data, &c)
}

// Marshal serialises the Symmetric object into JSON.
func (c *Channel) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Channel) label() []byte {
	return append([]byte(c.Name), []byte(c.Description)...)
}

func deriveIntermediary(name, description string, salt, rsaPub, secret []byte) []byte {
	h, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}
	_, err = h.Write([]byte(name))
	if err != nil {
		panic(err)
	}
	_, err = h.Write([]byte(description))
	if err != nil {
		panic(err)
	}
	_, err = h.Write(rsaPub)
	if err != nil {
		panic(err)
	}
	// secret is hashed first so that
	// we can share all the inputs to the
	// hkdf without giving out the secret.
	secretHash := blake2b.Sum256(secret)
	_, err = h.Write(secretHash[:])
	if err != nil {
		panic(err)
	}
	_, err = h.Write(salt)
	if err != nil {
		panic(err)
	}
	return h.Sum(nil)
}

// NewChannelID creates a new channel ID.
func NewChannelID(name, description string, salt, rsaPub, secret []byte) (*id.ID, []byte, error) {

	if len(secret) != 32 {
		return nil, nil, errors.New("NewChannelID secret must be 32 bytes long.")
	}

	hkdfHash := func() hash.Hash {
		hash, err := blake2b.New256(nil)
		if err != nil {
			panic(err)
		}
		return hash
	}

	hkdf1 := hkdf.New(hkdfHash,
		deriveIntermediary(name, description, salt, rsaPub, secret),
		salt,
		[]byte(hkdfInfo))

	identityBytes := make([]byte, 32)
	_, err := io.ReadFull(hkdf1, identityBytes)
	if err != nil {
		panic(err)
	}

	sid := &id.ID{}
	copy(sid[:], identityBytes)
	sid.SetType(id.User)

	hkdf2 := hkdf.New(hkdfHash,
		secret,
		deriveIntermediary(name, description, salt, rsaPub, secret),
		[]byte(hkdfInfo))

	key := make([]byte, 32)
	_, err = io.ReadFull(hkdf2, key)
	if err != nil {
		panic(err)
	}

	return sid, key, nil
}

type channelDisk struct {
	ReceptionID *id.ID
	Name        string
	Description string
	Salt        []byte
	RsaPubKey   *rsa.PublicKey
	Secret      []byte
	key         []byte
}

func (c *Channel) MarshalJson() ([]byte, error) {
	return json.Marshal(channelDisk{
		ReceptionID: c.ReceptionID,
		Name:        c.Name,
		Description: c.Description,
		Salt:        c.Salt,
		RsaPubKey:   c.RsaPubKey,
		Secret:      c.Secret,
		key:         c.key,
	})

}

func (c *Channel) UnmarshalJson(b []byte) error {
	cDisk := &channelDisk{}
	err := json.Unmarshal(b, cDisk)
	if err != nil {
		return err
	}

	*c = Channel{
		ReceptionID: cDisk.ReceptionID,
		Name:        cDisk.Name,
		Description: cDisk.Description,
		Salt:        cDisk.Salt,
		RsaPubKey:   cDisk.RsaPubKey,
		Secret:      cDisk.Secret,
		key:         cDisk.key,
	}

	return nil

}
