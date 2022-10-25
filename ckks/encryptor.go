package ckks

import (
	"github.com/tuneinsight/lattigo/v4/rlwe"
	"github.com/tuneinsight/lattigo/v4/utils"
)

// Encryptor an encryption interface for the CKKS scheme.
type Encryptor interface {
	Encrypt(pt *rlwe.Plaintext, ct *rlwe.Ciphertext)
	EncryptNew(pt *rlwe.Plaintext) (ct *rlwe.Ciphertext)
	EncryptZero(ct *rlwe.Ciphertext)
	EncryptZeroNew(level int) (ct *rlwe.Ciphertext)
	ShallowCopy() Encryptor
	WithKey(key interface{}) Encryptor
}

// PRNGEncryptor is an interface for encrypting BFV ciphertexts from a secret-key and
// a pre-determined PRNG. An Encryptor constructed from a secret-key complies to this
// interface.
type PRNGEncryptor interface {
	Encryptor
	WithPRNG(prng utils.PRNG) PRNGEncryptor
}

type encryptor struct {
	rlwe.Encryptor
	params Parameters
}

// NewEncryptor instantiates a new Encryptor for the CKKS scheme. The key argument can
// be *rlwe.PublicKey, *rlwe.SecretKey or nil.
func NewEncryptor(params Parameters, key interface{}) Encryptor {
	return &encryptor{rlwe.NewEncryptor(params.Parameters, key), params}
}

// NewPRNGEncryptor creates a new PRNGEncryptor instance that encrypts BFV ciphertexts from a secret-key and
// a PRNG.
func NewPRNGEncryptor(params Parameters, key *rlwe.SecretKey) PRNGEncryptor {
	enc := rlwe.NewPRNGEncryptor(params.Parameters, key)
	return &encryptor{enc, params}
}

// Encrypt encrypts the input plaintext and write the result on ciphertext.
// The level of the output ciphertext is min(plaintext.Level(), ct.Level()).
func (enc *encryptor) Encrypt(pt *rlwe.Plaintext, ct *rlwe.Ciphertext) {
	enc.Encryptor.Encrypt(pt, ct)
}

// EncryptNew encrypts the input plaintext returns the result as a newly allocated ciphertext.
func (enc *encryptor) EncryptNew(pt *rlwe.Plaintext) (ct *rlwe.Ciphertext) {
	ct = rlwe.NewCiphertextNTT(enc.params.Parameters, 1, pt.Level())
	enc.Encryptor.Encrypt(pt, ct)
	return
}

// EncryptZero generates an encryption of zero at the level and scale of ct, and writes the result on ctOut.
// Note that the Scale field of an encryption of zero can be changed arbitrarily, without requiring a Rescale.
func (enc *encryptor) EncryptZero(ct *rlwe.Ciphertext) {
	enc.Encryptor.EncryptZero(ct)
}

// EncryptZero generates an encryption of zero at the given level and scale and returns the
// result as a newly allocated ciphertext.
// Note that the Scale field of an encryption of zero can be changed arbitrarily, without requiring a Rescale.
func (enc *encryptor) EncryptZeroNew(level int) (ct *rlwe.Ciphertext) {
	ct = rlwe.NewCiphertextNTT(enc.params.Parameters, 1, level)
	enc.Encryptor.EncryptZero(ct)
	return
}

// ShallowCopy creates a shallow copy of this encryptor in which all the read-only data-structures are
// shared with the receiver and the temporary buffers are reallocated. The receiver and the returned
// Encryptors can be used concurrently.
func (enc *encryptor) ShallowCopy() Encryptor {
	return &encryptor{enc.Encryptor.ShallowCopy(), enc.params}
}

// WithKey creates a shallow copy of this encryptor with a new key in which all the read-only data-structures are
// shared with the receiver and the temporary buffers are reallocated. The receiver and the returned
// Encryptors can be used concurrently.
// Key can be *rlwe.PublicKey or *rlwe.SecretKey.
func (enc *encryptor) WithKey(key interface{}) Encryptor {
	return &encryptor{enc.Encryptor.WithKey(key), enc.params}
}

func (enc *encryptor) WithPRNG(prng utils.PRNG) PRNGEncryptor {
	if prngEnc, ok := enc.Encryptor.(rlwe.PRNGEncryptor); ok {
		return &encryptor{prngEnc.WithPRNG(prng), enc.params}
	}
	return nil
}
