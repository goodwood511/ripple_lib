package crypto

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
)

var (
	order = btcec.S256().N
	zero  = big.NewInt(0)
	one   = big.NewInt(1)
)

type ecdsaKey struct {
	privKey *btcec.PrivateKey
}

func newKey(seed []byte) *btcec.PrivateKey {
	inc := big.NewInt(0).SetBytes(seed)
	inc.Lsh(inc, 32)
	for key := big.NewInt(0); ; inc.Add(inc, one) {
		key.SetBytes(Sha512Half(inc.Bytes()))
		if key.Cmp(zero) > 0 && key.Cmp(order) < 0 {
			privKey, _ := btcec.PrivKeyFromBytes(key.Bytes())
			return privKey
		}
	}
}

// If seed is nil, generate a random one
func NewECDSAKey(seed []byte) (*ecdsaKey, error) {
	if seed == nil {
		seed = make([]byte, 16)
		if _, err := rand.Read(seed); err != nil {
			return nil, err
		}
	}
	return &ecdsaKey{newKey(seed)}, nil
}

func (k *ecdsaKey) generateKey(sequence uint32) *btcec.PrivateKey {
	// 1. 构造种子
	seed := make([]byte, btcec.PubKeyBytesLenCompressed+4)
	copy(seed, k.PubKey().SerializeCompressed())
	binary.BigEndian.PutUint32(seed[btcec.PubKeyBytesLenCompressed:], sequence)

	// 2. 生成派生私钥
	child := newKey(seed)

	// 3. 私钥加法（mod N）
	var sum btcec.ModNScalar
	sum.Set(&k.PrivKey().Key) // 获取主私钥 scalar
	sum.Add(&child.Key)       // 加上派生 scalar，结果仍为 ModNScalar

	// 4. 生成最终私钥
	derived := btcec.PrivKeyFromScalar(&sum)
	return derived
}

func (k *ecdsaKey) Id(sequence *uint32) []byte {
	if sequence == nil {
		return Sha256RipeMD160(k.PubKey().SerializeCompressed())
	}
	return Sha256RipeMD160(k.Public(sequence))
}
func (k *ecdsaKey) PrivKey() *btcec.PrivateKey {
	return k.privKey
}

func (k *ecdsaKey) PubKey() *btcec.PublicKey {
	return k.privKey.PubKey()
}

func (k *ecdsaKey) Private(sequence *uint32) []byte {
	var priv *btcec.PrivateKey
	if sequence == nil {
		priv = k.privKey
	} else {
		priv = k.generateKey(*sequence)
	}
	return priv.Serialize()
}

func (k *ecdsaKey) Public(sequence *uint32) []byte {
	if sequence == nil {
		return k.PubKey().SerializeCompressed()
	}
	return k.generateKey(*sequence).PubKey().SerializeCompressed()
}
