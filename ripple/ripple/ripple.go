package ripple

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"ripple_lib/ripple-sdk/crypto"
	"ripple_lib/ripple-sdk/data"
	"ripple_lib/ripple-sdk/websockets"

	rippleaddr "ripple_lib/ripple/rippleaddr"

	"github.com/sirupsen/logrus"
)

// Network define which the ripple network to use
type Network int

const (
	_           Network = iota
	MainNetwork         // MainNetwork ripple mainnet
	TestNetwork         // TestNetwork ripple testnet
)

type Ripple struct {
	Client *websockets.Remote
}

type SignerInfo struct {
	Account      string
	SingerWeight uint16
}

/*
NewRipple ...
Create a ripple websockets client
It's excpted to be used when need to connect to ripple network, e.g. chainnode
*/
func NewRipple(host string) (*Ripple, error) {
	client, err := websockets.NewRemote(host)
	if err != nil {
		logrus.Errorf("Fail to connet network %v", err)
		return nil, err
	}

	return &Ripple{
		Client: client,
	}, nil
}

/*
NewOfflineRipple ...
Create a ripple websockets offline client
It's expected to be used in wallet server, RC server, Sign server,
*/
func NewOfflineRipple() (*Ripple, error) {
	return &Ripple{
		Client: nil,
	}, nil
}

// Close close the client connection
func (r *Ripple) Close() {
	if r.Client != nil {
		r.Client.Close()
	}
}

/*
GetBalance ...
Get an account's balance
Currently, only support XRP
*/
func (r *Ripple) GetBalance(addr string) (string, error) {

	a, err := data.NewAccountFromAddress(addr)
	if err != nil {
		logrus.Errorf("Fail to covert address to account, err is %v", err)
		return "", err
	}

	result, err := r.Client.AccountInfo(*a)
	if err != nil {
		logrus.Errorf("Fail to get account %v's info, err is  %v", a, err)
		return "", err
	}

	if !result.AccountData.Balance.IsNative() {
		logrus.Errorf("Account %v's asset is not XRP", a, err)
		return "", err
	}

	v, err := result.AccountData.Balance.Native()
	return v.String(), nil
}

/*
GetBlockHeight ...
Get Ripple block height
*/
func (r *Ripple) GetBlockHeight() (uint32, error) {
	state, err := r.Client.ServerState()
	if err != nil {
		logrus.Errorf("Fail to get server  stat , err is %v", err)
		return 0, err
	}

	return state.State.ValidatedLedger.Sequence, nil
}

/*
GetSequence ...
Get an account's sequence
*/
func (r *Ripple) GetSequence(addr string) (uint32, error) {

	a, err := data.NewAccountFromAddress(addr)
	if err != nil {
		logrus.Errorf("Fail to covert address to account, err is %v", err)
		return 0, err
	}

	result, err := r.Client.AccountInfo(*a)
	if err != nil {
		logrus.Errorf("Fail to get account %v's info, err is  %v", a, err)
		return 0, err
	}

	return *result.AccountData.Sequence, nil

}

/*
GetAccountTx ...
Get an account's historical transactions
*/
func (r *Ripple) GetAccountTx(addr string) (*[]data.TransactionWithMetaData, int, error) {

	a, err := data.NewAccountFromAddress(addr)
	if err != nil {
		logrus.Errorf("Fail to covert address to account, err is %v", err)
		return nil, 0, err
	}

	var records []data.TransactionWithMetaData

	txchan := r.Client.AccountTx(*a, 20, -1, -1)

	for i := 0; ; i++ {
		record, isactive := <-txchan

		if isactive {
			records = append(records, *record)
		} else {
			fmt.Printf("After %v receive, channel is closed\n", i)
			return &records, i, nil
		}
	}
}

/*
GetTx ...
Get the tx by hash
*/
func (r *Ripple) GetTx(hash data.Hash256) (*websockets.TxResult, error) {
	return r.Client.Tx(hash)
}

/*
GetBlockByNumber ...
Get blockdata by number
*/
func (r *Ripple) GetBlockByNumber(height uint32) (*websockets.LedgerResultOnlyHash, error) {
	// return r.Client.Ledger(height, true)
	return r.Client.LedgerOnlyHash(height, true)
}

/*
CreateSignleSignSignerList ...
Create a signle singer list
fee: amount in drops.
*/
func (r *Ripple) CreateSignleSignSignerList(addr string, signers []SignerInfo, seq uint32, fee string) (*data.SignerListSet, error) {
	var count uint32

	dfee, err := strconv.ParseInt(fee, 10, 64)
	if err != nil {
		logrus.Errorf("Fail to covert fee string to int64  err is %v", err)
		return nil, err
	}

	if dfee <= 0 {
		return nil, fmt.Errorf("fee is negative or zero %v", fee)
	}

	if len(signers) > 8 {
		return nil, fmt.Errorf("Fail to create singer list, singer's count: %v, > 8 ", len(signers))
	}

	account, err := data.NewAccountFromAddress(addr)
	if err != nil {
		logrus.Errorf("Fail to covert address to account, err is %v", err)
		return nil, err
	}

	var signer_list data.SignerListSet
	count = uint32(len(signers))

	if count != 0 {
		var signer_ex data.SignerEntryEx
		signer_list.SignerQuorum = count

		for _, signer := range signers {
			a, err := data.NewAccountFromAddress(signer.Account)
			if err != nil {
				logrus.Errorf("Fail to covert get account, err is %v", err)
				return nil, err
			}
			signer_ex.SignerEntry.Account = a
			signer_ex.SignerEntry.SignerWeight = &signer.SingerWeight
			signer_list.SignerEntries = append(signer_list.SignerEntries, signer_ex)
		}

	} else {
		logrus.Infof("**Intend to disable multiSign\n")
		signer_list.SignerEntries = nil
		signer_list.SignerQuorum = count
	}

	base := signer_list.GetBase()
	base.TransactionType = data.SIGNER_LIST_SET
	base.Account = (*account)
	b, _ := data.NewNativeValue(dfee)
	base.Fee = *b
	base.Sequence = seq

	return &signer_list, nil
}

func getHashFromSubmitResult(result *websockets.SubmitResult) (*data.Hash256, error) {
	var txid data.Hash256
	txid_string := (result.Tx.(map[string]interface{})["hash"])
	txid_binary, err := hex.DecodeString(txid_string.(string))
	if err != nil {
		return nil, err
	}
	copy(txid.Bytes(), txid_binary)
	return &txid, nil
}

/*
CreateSingleSignPayment ...
Create a signle signed payment
amount: in drops
fee: in drops.
tag: destination tag
*/
func (r *Ripple) CreateSingleSignPayment(from, to, amount, fee, memo string, tag, seq uint32) (*data.Payment, error) {

	var p data.Payment

	damount, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		logrus.Errorf("Fail to covert amount string to int64  err is %v", err)
		return nil, err
	}

	dfee, err := strconv.ParseInt(fee, 10, 64)
	if err != nil {
		logrus.Errorf("Fail to covert fee string to int64  err is %v", err)
		return nil, err
	}

	if damount <= 0 || dfee <= 0 {
		return nil, fmt.Errorf("Amount %v or fee %v is illegal", amount, fee)
	}

	account_from, err := data.NewAccountFromAddress(from)
	if err != nil {
		logrus.Errorf("Fail to covert address %v to account, err is %v", from, err)
		return nil, err
	}

	account_to, err := data.NewAccountFromAddress(to)
	if err != nil {
		logrus.Errorf("Fail to covert address %v to account, err is %v", to, err)
		return nil, err
	}

	//var tag uint32 = 102285
	p.Sequence = seq
	p.Destination = (*account_to)
	p.DestinationTag = &tag
	a, err := data.NewAmount(int64(damount))
	if err != nil {
		logrus.Errorf("Amount %v is illegal, err is %v", amount, err)
		return nil, err
	}
	p.Amount = *a
	base := p.GetBase()
	base.TransactionType = data.PAYMENT
	base.Account = (*account_from)
	b, err := data.NewNativeValue(dfee)
	if err != nil {
		logrus.Errorf("Fee %v is illegal, err is %v", fee, err)
		return nil, err
	}
	base.Fee = *b

	// if len(memo) != 0 {
	// 	var m data.Memo
	// 	m.SetTypeFromString("BHEX")
	// 	m.SetDataFromString(memo)
	// 	m.SetFormatFromString("12344321")
	// 	base.Memos = append(base.Memos, m)
	// }

	return &p, nil
}

/*
CreateDisableMasterKey ...
Create a signle signed payment
Disable masterkey, this funciton should only be called after successfully set the multisigner
and is only Signle signed by master key
*/
func (r *Ripple) CreateDisableMasterKey(from, fee string, seq uint32) (*data.AccountSet, error) {

	var p data.AccountSet

	dfee, err := strconv.ParseInt(fee, 10, 64)
	if err != nil {
		logrus.Errorf("Fail to covert fee string to int64  err is %v", err)
		return nil, err
	}

	if dfee <= 0 {
		return nil, fmt.Errorf("fee %v is illegal", fee)
	}

	account_from, err := data.NewAccountFromAddress(from)
	if err != nil {
		logrus.Errorf("Fail to covert address %v to account, err is %v", from, err)
		return nil, err
	}

	var flag uint32 = 4
	p.Sequence = seq
	p.SetFlag = &flag

	base := p.GetBase()
	base.TransactionType = data.ACCOUNT_SET
	base.Account = (*account_from)
	b, err := data.NewNativeValue(dfee)
	if err != nil {
		logrus.Errorf("Fee %v is illegal, err is %v", fee, err)
		return nil, err
	}
	base.Fee = *b
	return &p, nil
}

/*
SignSingleSignTransaction ...
Sign a signle signed transaction with secrete
*/
func (r *Ripple) SignSingleSignTransaction(s data.Signer, screte string) error {

	var seed data.Seed
	var sequence uint32

	seed.UnmarshalText([]byte(screte))
	key := seed.Key(data.ECDSA)
	return data.Sign(s, key, &sequence)
}

/*
SignSingleSignTransactionWithPrivKey ...
Sign a signle transaction with privekey string
privkey: a huam readable privatekey "pxxxx"
*/
func (r *Ripple) SignSingleSignTransactionWithPrivKey(s data.Signer, privkey string) error {
	if !rippleaddr.CheckRipplePrivKey(privkey) {
		return fmt.Errorf("invalide privkey string %v", privkey)
	}

	b, err := crypto.Base58Decode(privkey, crypto.ALPHABET)
	if err != nil {
		return err
	}

	return data.SignWithPrivKey(s, b[1:len(b)-4])
}

/*
BroadcastSignleSignTransaction ...
push a single signed transaction to blockchain
*/
func (r *Ripple) BroadcastSignleSignTransaction(tx data.Transaction) (*data.Hash256, error) {
	submit_result, err := r.Client.Submit(tx)
	if err != nil {
		return nil, err
	}
	if !submit_result.EngineResult.Success() {
		return nil, fmt.Errorf("Fail to submit signle signed transaction, submit_result.EngineResult is %v", submit_result.EngineResultMessage)
	}

	txhash, err := getHashFromSubmitResult(submit_result)
	if err != nil {
		return nil, err
	}

	return txhash, nil
}

/*
CreateMultiSignPayment ...
Create a Multi signed transaction
amount: in drops
fee: in drops
*/
func (r *Ripple) CreateMultiSignPayment(from, to, amount, fee, memo string, tag, seq uint32) (*data.MultiSignPayment, error) {

	var p data.MultiSignPayment
	damount, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		logrus.Errorf("Fail to covert amount string to int64  err is %v", err)
		return nil, err
	}

	dfee, err := strconv.ParseInt(fee, 10, 64)
	if err != nil {
		logrus.Errorf("Fail to covert fee string to int64  err is %v", err)
		return nil, err
	}

	if damount <= 0 || dfee <= 0 {
		return nil, fmt.Errorf("Amount %v or fee %v is illegal", amount, fee)
	}

	account_from, err := data.NewAccountFromAddress(from)
	if err != nil {
		logrus.Errorf("Fail to covert address %v to account, err is %v", from, err)
		return nil, err
	}

	account_to, err := data.NewAccountFromAddress(to)
	if err != nil {
		logrus.Errorf("Fail to covert address %v to account, err is %v", to, err)
		return nil, err
	}

	p.Sequence = seq
	p.Destination = (*account_to)
	p.DestinationTag = &tag
	a, err := data.NewAmount(int64(damount))
	if err != nil {
		logrus.Errorf("Amount %v is illegal, err is %v", amount, err)
		return nil, err
	}
	p.Amount = *a

	base := p.GetBase()
	base.TransactionType = data.PAYMENT
	base.Account = (*account_from)
	b, err := data.NewNativeValue(dfee)
	if err != nil {
		logrus.Errorf("Fee %v is illegal, err is %v", fee, err)
		return nil, err
	}
	base.Fee = *b

	if len(memo) != 0 {
		var m data.Memo
		m.SetTypeFromString("BHEX")
		m.SetDataFromString(memo)
		m.SetFormatFromString("12344321")
		base.Memos = append(base.Memos, m)
	}

	//copy(p.GetPublicKey().Bytes(), "") //set the SigningPubKey to ""

	return &p, nil
}

/*
SignMultiSignTransactionInSerial ...
Sign a Multi signed transaction in serial with secrete
*/
func (r *Ripple) SignMultiSignTransactionInSerial(s data.Signer, screte string) error {
	return nil
}

/*
SignMultiSignTransactionInParallel ...
Sign a Multi signed transaction in parallel with secrete
*/
func (r *Ripple) SignMultiSignTransactionInParallel(s data.Signer, screte string) (data.MultiSignerEntryEx, error) {

	var seed data.Seed
	var sequence uint32

	seed.UnmarshalText([]byte(screte))
	key := seed.Key(data.ECDSA)

	return data.MultiSignInParallel(s, key, &sequence)
}

/*
SignMultiSignTransactionInParallelWithPrivKey ...
Sign a Multi signed transaction in parallel with private key
private key:  a huam readable privatekey "pxxxx"
*/
func (r *Ripple) SignMultiSignTransactionInParallelWithPrivKey(s data.Signer, privkey string) (data.MultiSignerEntryEx, error) {

	var signer data.MultiSignerEntryEx
	if !rippleaddr.CheckRipplePrivKey(privkey) {
		return signer, fmt.Errorf("invalide privkey string %v", privkey)
	}

	b, err := crypto.Base58Decode(privkey, crypto.ALPHABET)
	if err != nil {
		return signer, err
	}
	return data.MultiSignInParallelWithPrivKey(s, b[1:len(b)-4])
}

func sortSignerEntryEx(s []data.MultiSignerEntryEx) {

	for i := 0; i < len(s)-1; i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i].Signer.Account.Compare(s[j].Signer.Account) > 0 {
				s[i], s[j] = s[j], s[i]
			}
		}

	}
}

/*
MergeMultiSignSignatures ...
Merge signatures togther for a MultiSign transaction
signers: signers
*/
func (r *Ripple) MergeMultiSignSignatures(tx data.MultiSignTransaction, signers []data.MultiSignerEntryEx) error {

	// if len(signers) != NUM_OF_MULTISIGNER {
	// 	return fmt.Errorf("signer num does not equalt to %v", NUM_OF_MULTISIGNER)
	// }

	sortSignerEntryEx(signers)
	for _, s := range signers {
		tx.GetBase().Signers = append(tx.GetBase().Signers, s)
	}
	return nil
}

/*
BroadcastMultiSignTransaction ...
Push a MultiSigned transaction to blockchain, return the txhash for future query
*/
func (r *Ripple) BroadcastMultiSignTransaction(tx data.MultiSignTransaction) (*data.Hash256, error) {

	submitResult, err := r.Client.SubmitMultiSign(tx)
	if err != nil || submitResult == nil {
		errInfo := fmt.Sprintf("submit multi sign error:%v ", err)
		return nil, errors.New(errInfo)
	}
	txhash, err := getHashFromSubmitResult(submitResult)
	if err != nil {
		logrus.Errorf("getHashFromSubmitResult err %v", err)
		return nil, err
	}

	if !submitResult.EngineResult.Success() && !submitResult.EngineResult.Queued() {
		return txhash, fmt.Errorf("Fail to submit multi signed transaction, submitResult.EngineResult is %v", submitResult.EngineResultMessage)
	}
	logrus.Infof("Submit result is %v", submitResult.EngineResultMessage)
	return txhash, nil
}
