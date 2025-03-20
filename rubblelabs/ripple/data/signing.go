package data

import (
	"ripple_lib/rubblelabs/ripple/crypto"

	"github.com/btcsuite/btcd/btcec"
)

// sign with secrete
func Sign(s Signer, key crypto.Key, sequence *uint32) error {
	s.InitialiseForSigning()
	copy(s.GetPublicKey().Bytes(), key.Public(sequence))
	hash, msg, err := SigningHash(s)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(key.Private(sequence), hash.Bytes(), append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return err
	}
	*s.GetSignature() = VariableLength(sig)
	hash, _, err = Raw(s)
	if err != nil {
		return err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())
	return nil
}

/* Sign with privatekey...
key: raw private key, in 33 bytes
*/

func SignWithPrivKey(s Signer, key []byte) error {

	s.InitialiseForSigning()
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), key)
	copy(s.GetPublicKey().Bytes(), pubKey.SerializeCompressed())

	hash, msg, err := SigningHash(s)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(privKey.Serialize(), hash.Bytes(), append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return err
	}
	*s.GetSignature() = VariableLength(sig)
	hash, _, err = Raw(s)
	if err != nil {
		return err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())
	return nil
}

// Todo
func MultiSignInSerial(s Signer, key crypto.Key, sequence *uint32) error {
	return nil

}

func MultiSignInParallel(s Signer, key crypto.Key, sequence *uint32) (MultiSignerEntryEx, error) {

	s.InitialiseForSigning()
	var signer MultiSignerEntryEx

	copy(s.GetPublicKey().Bytes(), "") //set the SigningPubKey to ""
	account := key.Id(sequence)
	hash, msg, err := MultiSignHash(s, account)

	//for ECDSA crypto,
	sig, err := crypto.Sign(key.Private(sequence), hash.Bytes(), append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return signer, err
	}

	signer.Signer.SigningPubKey = new(PublicKey)
	signer.Signer.TxnSignature = new(VariableLength)

	*signer.GetSignature() = VariableLength(sig)
	copy(signer.GetPublicKey().Bytes(), key.Public(sequence)) //利用
	copy(signer.GetAccount().Bytes(), key.Id(sequence))

	hash, _, err = Raw(s)

	if err != nil {
		return signer, err
	}
	copy(s.GetHash().Bytes(), hash.Bytes()) //copy hash
	return signer, nil
}

func MultiSignInParallelWithPrivKey(s Signer, key []byte) (MultiSignerEntryEx, error) {

	s.InitialiseForSigning()
	var signer MultiSignerEntryEx

	copy(s.GetPublicKey().Bytes(), "") //set the SigningPubKey to ""
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), key)

	account := crypto.Sha256RipeMD160(pubKey.SerializeCompressed())
	hash, msg, err := MultiSignHash(s, account)

	//for ECDSA crypto,
	sig, err := crypto.Sign(privKey.Serialize(), hash.Bytes(), append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return signer, err
	}

	signer.Signer.SigningPubKey = new(PublicKey)
	signer.Signer.TxnSignature = new(VariableLength)

	*signer.GetSignature() = VariableLength(sig)
	copy(signer.GetPublicKey().Bytes(), pubKey.SerializeCompressed()) //利用
	copy(signer.GetAccount().Bytes(), account)

	hash, _, err = Raw(s)

	if err != nil {
		return signer, err
	}
	copy(s.GetHash().Bytes(), hash.Bytes()) //copy hash
	return signer, nil
}

func CheckSignature(s Signer) (bool, error) {
	hash, msg, err := SigningHash(s)
	if err != nil {
		return false, err
	}
	return crypto.Verify(s.GetPublicKey().Bytes(), hash.Bytes(), msg, s.GetSignature().Bytes())
}
