package rippleaddr

import (
	"fmt"
	"github.com/goodwood511/ripple_lib/ripple-sdk/crypto"

	"github.com/btcsuite/btcd/btcec"
)

func RippleGenerateKey(s string) (crypto.Key, error) {

	seed, err := crypto.GenerateFamilySeed(s)
	if err != nil {
		return nil, fmt.Errorf("Fail to generate new seed %v", err)
	}

	key, err := crypto.NewECDSAKey(seed.Payload())
	if err != nil {
		return nil, fmt.Errorf("Fail to generate new seed %v", err)
	}

	return key, nil
}

func CheckRippleSeed(seed string) bool {
	hash, err := crypto.NewRippleHashCheck(seed, crypto.RIPPLE_FAMILY_SEED)
	if err != nil {
		//fmt.Printf("CheckRippleAddr, err is %v\n", err)
		return false
	}

	if hash.String() == seed {
		return true
	}
	return false
}

/*
	Generate keys and addr from ripple seed

seed: human readable string, ex "sh7pek1W31vHCshtWo6hhksmCg7DG"
*/
func RippleSeedToKeysAndAddr(seed string, sequence *uint32) (privKey string, pubKey string, addr string, err error) {
	if !CheckRippleSeed(seed) {
		return "", "", "", err
	}

	b, err := crypto.Base58Decode(seed, crypto.ALPHABET)
	if err != nil {
		return "", "", "", err
	}

	key, err := crypto.NewECDSAKey(b[1 : len(b)-4])

	priv, err := crypto.AccountPrivateKey(key, sequence)
	if err != nil {
		return "", "", "", fmt.Errorf("Fail to generate new private key %v", err)
	}

	pub, err := crypto.AccountPublicKey(key, sequence)
	if err != nil {
		return "", "", "", fmt.Errorf("Fail to generate new public key %v", err)
	}

	a, err := crypto.AccountId(key, sequence)
	if err != nil {
		return "", "", "", fmt.Errorf("Fail to generate new addr %v", err)
	}

	return priv.String(), pub.String(), a.String(), nil
}

// Generate keys and addr from arbitrary string
func RippleGenerateKeysAndAddr(s string, sequence *uint32) (privKey string, pubKey string, addr string, err error) {

	seed, err := crypto.GenerateFamilySeed(s)
	if err != nil {
		return "", "", "", fmt.Errorf("Fail to generate new seed %v", err)
	}
	key, err := crypto.NewECDSAKey(seed.Payload())
	if err != nil {
		return "", "", "", fmt.Errorf("Fail to generate new seed %v", err)
	}

	priv, err := crypto.AccountPrivateKey(key, sequence)
	if err != nil {
		return "", "", "", fmt.Errorf("Fail to generate new private key %v", err)
	}

	pub, err := crypto.AccountPublicKey(key, sequence)
	if err != nil {
		return "", "", "", fmt.Errorf("Fail to generate new public key %v", err)
	}

	a, err := crypto.AccountId(key, sequence)
	if err != nil {
		return "", "", "", fmt.Errorf("Fail to generate new addr %v", err)
	}

	return priv.String(), pub.String(), a.String(), nil
}

func CheckRippleAddr(addr string) bool {
	hash, err := crypto.NewRippleHashCheck(addr, crypto.RIPPLE_ACCOUNT_ID)
	if err != nil {
		//fmt.Printf("CheckRippleAddr, err is %v\n", err)
		return false
	}

	if hash.String() == addr {
		return true
	}
	return false
}

func CheckRipplePubKey(pubKey string) bool {
	hash, err := crypto.NewRippleHashCheck(pubKey, crypto.RIPPLE_ACCOUNT_PUBLIC)
	if err != nil {
		//fmt.Printf("CheckRippleAddr, err is %v\n", err)
		return false
	}

	if hash.String() == pubKey {
		return true
	}
	return false
}

func CheckRipplePrivKey(privKey string) bool {
	hash, err := crypto.NewRippleHashCheck(privKey, crypto.RIPPLE_ACCOUNT_PRIVATE)
	if err != nil {
		//fmt.Printf("CheckRippleAddr, err is %v\n", err)
		return false
	}

	if hash.String() == privKey {
		return true
	}
	return false
}

func RipplePubKeyToAddr(pubKey string) (string, error) {

	if !CheckRipplePubKey(pubKey) {
		return "", fmt.Errorf("Not a Ripple public key")
	}

	b, err := crypto.Base58Decode(pubKey, crypto.ALPHABET)
	if err != nil {
		return " ", err
	}

	// fmt.Printf("public key binary is %v %X \n", b1, b1[1:len(b1)-4])
	// fmt.Printf("Public key is %v\n", pubKey)
	b1 := crypto.Sha256RipeMD160(b[1 : len(b)-4])
	a, err := crypto.NewAccountId(b1)
	if err != nil {
		return "", err
	}

	return a.String(), nil
}

func RipplePrivKeyToPub(privKey string) (string, error) {

	if !CheckRipplePrivKey(privKey) {
		return "", fmt.Errorf("Not a Ripple private key")
	}

	b, err := crypto.Base58Decode(privKey, crypto.ALPHABET)
	if err != nil {
		return " ", err
	}

	k, _ := btcec.PrivKeyFromBytes(btcec.S256(), b[1:len(b)-4])
	pubKey, err := crypto.NewAccountPublicKey(k.PubKey().SerializeCompressed())
	if err != nil {
		return " ", err
	}
	return pubKey.String(), nil
}

func CheckRipplePrivKeyPub(privKey string, pubKey string) (bool, error) {

	pub, err := RipplePrivKeyToPub(privKey)
	if err != nil {
		return false, err
	}

	if pub == pubKey {
		return true, nil
	} else {
		return false, nil
	}
}

// NewAddress new address return key address pubkey error
func NewAddress() ([]byte, string, string, error) {
	ecdsaKey, err := crypto.NewECDSAKey(nil)
	if err != nil {
		return nil, "", "", err
	}
	key := ecdsaKey.Serialize()

	k, _ := btcec.PrivKeyFromBytes(btcec.S256(), key)
	pubKey, err := crypto.NewAccountPublicKey(k.PubKey().SerializeCompressed())
	if err != nil {
		return nil, "", "", err
	}

	addr, err := RipplePubKeyToAddr(pubKey.String())
	if err != nil {
		return nil, "", "", err
	}

	return key, addr, pubKey.String(), nil
}
