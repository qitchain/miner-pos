// +build !cgo

package bls

// Public for crypto.Signer and crypto.Decrypter interfaces
func (priv *PrivateKey) Public() crypto.PublicKey {
	panic(ErrNotSupport)
}

// Sign for crypto.Signer interface
func (priv *PrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	return nil, ErrNotSupport
}

// Decrypt for crypto.Decrypter interface
func (priv *PrivateKey) Decrypt(rand io.Reader, ciphertext []byte, opts crypto.DecrypterOpts) (plaintext []byte, err error) {
	return nil, ErrNotSupport
}
