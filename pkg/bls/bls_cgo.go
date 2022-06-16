//go:build cgo
// +build cgo

package bls

// #cgo CXXFLAGS: -fexceptions
// #cgo LDFLAGS: -L./c-bindings/libs
//
// #include "c-bindings/include/cbindings.h"
//
import "C"
import (
	"crypto"
	"io"
	"unsafe"

	// keep directories
	_ "chia-miner/pkg/bls/c-bindings/include"
	_ "chia-miner/pkg/bls/c-bindings/libs"
)

// GenerateKeyFromSeed creates a new BLS signing key.
func GenerateKeyFromSeed(seed []byte) (*PrivateKey, error) {
	if len(seed) < 32 {
		return nil, ErrLimitSeedSize
	}

	privT := C.key_gen((*C.uint8_t)(unsafe.Pointer(&seed[0])), C.size_t(len(seed)))
	defer C.key_destroy(privT)

	buf := make([]byte, 32)
	bufLen := C.key_bytes(privT, (*C.uint8_t)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	return &PrivateKey{data: buf[0:int(bufLen)]}, nil
}

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() crypto.PublicKey {
	privT := C.key_from_bytes((*C.uint8_t)(unsafe.Pointer(&priv.data[0])), C.size_t(len(priv.data)))
	defer C.key_destroy(privT)

	g1T := C.key_g1(privT)
	defer C.g1_destroy(g1T)

	buf := make([]byte, 48)
	bufLen := C.g1_bytes(g1T, (*C.uint8_t)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	return &PublicKey{data: buf[0:int(bufLen):int(bufLen)]}
}

// Sign message by private key.
// This method implements crypto.Signer.
// Ignore io.Reader and crypto.SignerOpts
func (priv *PrivateKey) Sign(_ io.Reader, message []byte, _ crypto.SignerOpts) ([]byte, error) {
	return priv.SignMessage(message)
}

// SignMessage message by private key.
func (priv *PrivateKey) SignMessage(message []byte) ([]byte, error) {
	privT := C.key_from_bytes((*C.uint8_t)(unsafe.Pointer(&priv.data[0])), C.size_t(len(priv.data)))
	if privT == nil {
		return nil, ErrInvalidKey
	}
	defer C.key_destroy(privT)

	g2T := C.key_sign(privT, (*C.uint8_t)(unsafe.Pointer(&message[0])), C.size_t(len(message)), (*C.G1ElementT)(unsafe.Pointer(nil)))
	if g2T == nil {
		return nil, ErrInvalidKey
	}
	defer C.g2_destroy(g2T)

	buf := make([]byte, 96)
	bufLen := C.g2_bytes(g2T, (*C.uint8_t)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	return buf[0:int(bufLen):int(bufLen)], nil
}

// SignMessageWithPrependPublicKey message by private key and prepend public key.
func (priv *PrivateKey) SignMessageWithPrependPublicKey(message []byte, prependPublicKey *PublicKey) ([]byte, error) {
	privT := C.key_from_bytes((*C.uint8_t)(unsafe.Pointer(&priv.data[0])), C.size_t(len(priv.data)))
	if privT == nil {
		return nil, ErrInvalidKey
	}
	defer C.key_destroy(privT)

	prependPkT := C.g1_from_bytes((*C.uint8_t)(unsafe.Pointer(&prependPublicKey.data[0])), C.size_t(len(prependPublicKey.data)))
	if prependPkT == nil {
		return nil, ErrInvalidKey
	}
	defer C.g1_destroy(prependPkT)

	g2T := C.key_sign(privT, (*C.uint8_t)(unsafe.Pointer(&message[0])), C.size_t(len(message)), prependPkT)
	if g2T == nil {
		return nil, ErrInvalidKey
	}
	defer C.g2_destroy(g2T)

	buf := make([]byte, 96)
	bufLen := C.g2_bytes(g2T, (*C.uint8_t)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	return buf[0:int(bufLen):int(bufLen)], nil
}

// Decrypt decrypts ciphertext with priv.
func (priv *PrivateKey) Decrypt(_ io.Reader, ciphertext []byte, _ crypto.DecrypterOpts) (plaintext []byte, err error) {
	return nil, ErrNotSupport
}

// DeriveChild derive child private key
func (priv *PrivateKey) DeriveChild(paths []int) *PrivateKey {
	if len(paths) == 0 {
		return &PrivateKey{data: priv.data[:len(priv.data):len(priv.data)]}
	}

	parentPkT := C.key_from_bytes((*C.uint8_t)(unsafe.Pointer(&priv.data[0])), C.size_t(len(priv.data)))
	for _, path := range paths {
		derivePkT := C.key_derive_child(parentPkT, C.uint32_t(path))
		C.key_destroy(parentPkT)
		parentPkT = derivePkT
	}
	defer C.key_destroy(parentPkT)

	buf := make([]byte, 32)
	bufLen := C.key_bytes(parentPkT, (*C.uint8_t)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	return &PrivateKey{data: buf[0:int(bufLen):int(bufLen)]}
}

// UnmarshalBinary to private key
// This method implements encoding.BinaryUnmarshaler.
func (priv *PrivateKey) UnmarshalBinary(data []byte) error {
	// verify
	if err := checkPrivateKeyData(data); err != nil {
		return err
	}
	privT := C.key_from_bytes((*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)))
	if privT == nil {
		return ErrInvalidKey
	}
	C.key_destroy(privT)

	priv.data = data[0:len(data):len(data)]
	return nil
}

// UnmarshalBinary to public key
// This method implements encoding.BinaryUnmarshaler.
func (pk *PublicKey) UnmarshalBinary(data []byte) error {
	// verify
	if err := checkPublicKeyData(data); err != nil {
		return err
	}
	g1T := C.g1_from_bytes((*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)))
	if g1T == nil {
		return ErrInvalidKey
	}
	C.g1_destroy(g1T)

	pk.data = data[0:len(data):len(data)]
	return nil
}

// Verify verifies a BLS signature.
func (pk *PublicKey) Verify(digest []byte, sig []byte) error {
	pkT := C.g1_from_bytes((*C.uint8_t)(unsafe.Pointer(&pk.data[0])), C.size_t(len(pk.data)))
	if pkT == nil {
		return ErrInvalidKey
	}
	defer C.g1_destroy(pkT)

	// verify
	if err := checkSignatureData(sig); err != nil {
		return err
	}
	sigT := C.g2_from_bytes((*C.uint8_t)(unsafe.Pointer(&sig[0])), C.size_t(len(sig)))
	if sigT == nil {
		return ErrInvalidSign
	}
	defer C.g2_destroy(sigT)

	ret := C.g1_verify(pkT, (*C.uint8_t)(unsafe.Pointer(&digest[0])), C.size_t(len(digest)), sigT)
	if ret != 0 {
		return ErrVerification
	}
	return nil
}

// Add combine 2 public keys
func (pk *PublicKey) Add(other *PublicKey) *PublicKey {
	pkT1 := C.g1_from_bytes((*C.uint8_t)(unsafe.Pointer(&pk.data[0])), C.size_t(len(pk.data)))
	defer C.g1_destroy(pkT1)

	pkT2 := C.g1_from_bytes((*C.uint8_t)(unsafe.Pointer(&other.data[0])), C.size_t(len(other.data)))
	defer C.g1_destroy(pkT2)

	pkT := C.g1_add2(pkT1, pkT2)

	buf := make([]byte, 48)
	bufLen := C.g1_bytes(pkT, (*C.uint8_t)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	return &PublicKey{data: buf[0:int(bufLen):int(bufLen)]}
}
