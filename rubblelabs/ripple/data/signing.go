package data

import (
	"github.com/goodwood511/ripple_lib/rubblelabs/ripple/crypto"

	"github.com/btcsuite/btcd/btcec/v2"
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

	// ✅ 1. 解析私钥（新版）
	privKey, _ := btcec.PrivKeyFromBytes(key)
	pubKey := privKey.PubKey()

	// ✅ 2. 设置公钥
	copy(s.GetPublicKey().Bytes(), pubKey.SerializeCompressed())

	// ✅ 3. 获取待签名 hash + raw msg
	hash, msg, err := SigningHash(s)
	if err != nil {
		return err
	}

	// ✅ 4. 签名（调用你项目中的 crypto.Sign 函数）
	sig, err := crypto.Sign(privKey.Serialize(), hash.Bytes(), append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return err
	}

	// ✅ 5. 写入签名和 hash
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

	// ✅ 清空 SigningPubKey
	clear(s.GetPublicKey().Bytes())

	// ✅ 升级为 btcec/v2
	privKey, _ := btcec.PrivKeyFromBytes(key)
	pubKey := privKey.PubKey()

	// ✅ 获取 Account (公钥哈希)
	account := crypto.Sha256RipeMD160(pubKey.SerializeCompressed())

	// ✅ 构造签名哈希
	hash, msg, err := MultiSignHash(s, account)
	if err != nil {
		return signer, err
	}

	// ✅ 进行签名
	sig, err := crypto.Sign(privKey.Serialize(), hash.Bytes(), append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return signer, err
	}

	// ✅ 填充 signer 字段
	signer.Signer.SigningPubKey = new(PublicKey)
	signer.Signer.TxnSignature = new(VariableLength)

	*signer.GetSignature() = VariableLength(sig)
	copy(signer.GetPublicKey().Bytes(), pubKey.SerializeCompressed())
	copy(signer.GetAccount().Bytes(), account)

	// ✅ 计算最终交易 Hash
	hash, _, err = Raw(s)
	if err != nil {
		return signer, err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())

	return signer, nil
}

func CheckSignature(s Signer) (bool, error) {
	hash, msg, err := SigningHash(s)
	if err != nil {
		return false, err
	}
	return crypto.Verify(s.GetPublicKey().Bytes(), hash.Bytes(), msg, s.GetSignature().Bytes())
}
