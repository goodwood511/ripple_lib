package crypto

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ed25519"
)

func Sign(privateKey, hash, msg []byte) ([]byte, error) {
	switch len(privateKey) {
	case ed25519.PrivateKeySize:
		return signEd25519(privateKey, msg)
	case btcec.PrivKeyBytesLen:
		return signECDSA(privateKey, hash)
	default:
		return nil, fmt.Errorf("Unknown private key format")
	}
}

func Verify(publicKey, hash, msg, signature []byte) (bool, error) {
	switch publicKey[0] {
	case 0xED:
		return verifyEd25519(publicKey, signature, msg)
	case 0x02, 0x03:
		return verifyECDSA(publicKey, signature, hash)
	default:
		return false, fmt.Errorf("Unknown public key format")
	}
}

func signEd25519(privateKey, msg []byte) ([]byte, error) {
	// 验证私钥长度
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid private key size")
	}

	// 使用切片直接转换
	var p ed25519.PrivateKey = privateKey

	// 返回签名结果
	signature := ed25519.Sign(p, msg)
	return signature, nil
}

func verifyEd25519(pubKey, signature, msg []byte) (bool, error) {
	var (
		p [ed25519.PublicKeySize]byte
		s [ed25519.SignatureSize]byte
	)
	switch {
	case len(pubKey) != ed25519.PublicKeySize+1:
		return false, fmt.Errorf("Wrong public key length: %d", len(pubKey))
	case pubKey[0] != 0xED:
		return false, fmt.Errorf("Wrong public format:")
	case len(signature) != ed25519.SignatureSize:
		return false, fmt.Errorf("Wrong Signature length: %d", len(signature))
	default:
		copy(p[:], pubKey[1:])
		copy(s[:], signature)
		return ed25519.Verify(p[:], msg, s[:]), nil
	}
}

// Returns DER encoded signature from input hash
func signECDSA(privateKey, hash []byte) ([]byte, error) {
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKey)
	sig, err := priv.Sign(hash)
	if err != nil {
		return nil, err
	}
	return sig.Serialize(), nil
}

// Verifies a hash using DER encoded signature
func verifyECDSA(pubKey, signature, hash []byte) (bool, error) {
	sig, err := btcec.ParseDERSignature(signature, btcec.S256())
	if err != nil {
		return false, err
	}
	pk, err := btcec.ParsePubKey(pubKey, btcec.S256())
	if err != nil {
		return false, nil
	}
	return sig.Verify(hash, pk), nil
}
