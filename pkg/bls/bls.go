package bls

import (
	"crypto"
	"crypto/sha256"
	"encoding"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"strings"
)

// ErrNotSupport
var ErrNotSupport = errors.New("not support")

// ErrLimitSeedSize
var ErrLimitSeedSize = errors.New("seed size must be at least 32 bytes")

// ErrInvalidKey
var ErrInvalidKey = errors.New("invalid key")

// ErrInvalidSign
var ErrInvalidSign = errors.New("invalid signature")

// ErrDecryption represents a failure to decrypt a message.
var ErrDecryption = errors.New("crypto/bls: decryption error")

// ErrVerification represents a failure to verify a signature.
var ErrVerification = errors.New("crypto/bls: verification error")

// PrivateKey bls signature private key.
type PrivateKey struct {
	data []byte
}

// PrivateKey implements crypto.PrivateKey.
var _ crypto.PrivateKey = (*PrivateKey)(nil)

// PrivateKey implements crypto.Signer and crypto.Decrypter.
var _ crypto.Signer = (*PrivateKey)(nil)
var _ crypto.Decrypter = (*PrivateKey)(nil)

// PrivateKey implements encoding.BinaryMarshaler and encoding.BinaryUnmarshaler.
var _ encoding.BinaryMarshaler = (*PrivateKey)(nil)
var _ encoding.BinaryUnmarshaler = (*PrivateKey)(nil)

// PrivateKeyFromBytes generate from bytes
func PrivateKeyFromBytes(data []byte) (*PrivateKey, error) {
	var priv PrivateKey
	if err := priv.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	return &priv, nil
}

// PrivateKeyFromHex generate from hex string
func PrivateKeyFromHex(data string) (*PrivateKey, error) {
	if strings.HasPrefix(data, "0x") {
		data = data[2:]
	}
	pkBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return PrivateKeyFromBytes(pkBytes)
}

// MarshalBinary from private key.
// This method implements encoding.BinaryMarshaler.
func (priv *PrivateKey) MarshalBinary() ([]byte, error) {
	return priv.data, nil
}

// GenerateKey creates a new BLS signing key.
func GenerateKey(random io.Reader, bits int) (*PrivateKey, error) {
	if bits < 32 {
		return nil, ErrLimitSeedSize
	}

	seed := make([]byte, bits)
	for i := 0; i < bits; {
		n, err := random.Read(seed[i:])
		if err != nil {
			return nil, err
		}
		i += n
	}

	return GenerateKeyFromSeed(seed)
}

// PublicKey bls signature public key.
type PublicKey struct {
	data []byte
}

// PublicKey implements crypto.PublicKey.
var _ crypto.PublicKey = (*PublicKey)(nil)

// PublicKey implements encoding.BinaryMarshaler and encoding.BinaryUnmarshaler.
var _ encoding.BinaryMarshaler = (*PublicKey)(nil)
var _ encoding.BinaryUnmarshaler = (*PublicKey)(nil)

// PublicKeyFromBytes generate from bytes
func PublicKeyFromBytes(data []byte) (*PublicKey, error) {
	var pk PublicKey
	if err := pk.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	return &pk, nil
}

// PublicKeyFromHex generate from hex string
func PublicKeyFromHex(data string) (*PublicKey, error) {
	if strings.HasPrefix(data, "0x") {
		data = data[2:]
	}
	pkBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return PublicKeyFromBytes(pkBytes)
}

// GetFingerprint
func (pk *PublicKey) GetFingerprint() uint32 {
	if pk.data == nil {
		return 0
	}

	hashed := sha256.Sum256(pk.data)
	return binary.BigEndian.Uint32(hashed[:])
}

// MarshalBinary from public key.
// This method implements encoding.BinaryMarshaler.
func (pk *PublicKey) MarshalBinary() ([]byte, error) {
	return pk.data, nil
}

var errPrivateKeyInvalidSize = errors.New("private keys must be 32 bytes")

// checkPrivateKeyData
func checkPrivateKeyData(data []byte) error {
	if len(data) != 32 {
		return errPrivateKeyInvalidSize
	}
	return nil
}

var errG1InvalidSize = errors.New("g1 data must be 48 bytes")
var errG1InfinityMustBeCanonical = errors.New("g1 infinity element must be canonical")
var errG1NonInfinityInvalidStartWith = errors.New("g1 non-infinity element must start with 0b10")

// checkPublicKeyData verify G1 data
func checkPublicKeyData(data []byte) error {
	if len(data) != 48 {
		return errG1InvalidSize
	}

	if (data[0] & 0xc0) == 0xc0 {
		// representing infinity
		// enforce that infinity must be 0xc0000..00
		if data[0] != 0xc0 {
			return errG1InfinityMustBeCanonical
		}
		for i, b := range data {
			if i > 0 && b != 0x00 {
				return errG1InfinityMustBeCanonical
			}
		}
	} else {
		if (data[0] & 0xc0) != 0x80 {
			return errG1NonInfinityInvalidStartWith
		}
	}

	return nil
}

var errG2InvalidSize = errors.New("g2 data must be 96 bytes")
var errG2InvalidStartWith = errors.New("g2 element must always have 48th byte start with 0b000")
var errG2InfinityMustBeCanonical = errors.New("g2 infinity element must be canonical")
var errG2NonInvalidStartWith = errors.New("g2 non-inf element must have 0th byte start with 0b10")

// checkSignatureData verify G2 data
func checkSignatureData(data []byte) error {
	if len(data) != 96 {
		return errG2InvalidSize
	}
	if (data[48] & 0xe0) != 0x00 {
		return errG2InvalidStartWith
	}

	if (data[0] & 0xc0) == 0xc0 {
		// infinity
		// enforce that infinity must be 0xc0000..00
		if data[0] != 0xc0 {
			return errG2InfinityMustBeCanonical
		}
		for i, b := range data {
			if i > 0 && b != 0x00 {
				return errG2InfinityMustBeCanonical
			}
		}
	} else {
		if (data[0] & 0xc0) != 0x80 {
			return errG2NonInvalidStartWith
		}
	}

	return nil
}
