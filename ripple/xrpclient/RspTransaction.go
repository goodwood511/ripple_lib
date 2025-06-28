package xrpclient

type RspTransaction struct {
	Result Result `json:"result"`
}

type Result struct {
	Ledger      Ledger `json:"ledger"`
	Validated   bool   `json:"validated"`
	LedgerIndex int64  `json:"ledger_index"`
	LedgerHash  string `json:"ledger_hash"`
	Status      string `json:"status"`
}

type Ledger struct {
	CloseFlags          int64         `json:"close_flags"`
	LedgerIndex         string        `json:"ledger_index"`
	AccountHash         string        `json:"account_hash"`
	CloseTimeResolution int64         `json:"close_time_resolution"`
	CloseTime           int64         `json:"close_time"`
	Transactions        []Transaction `json:"transactions"`
	CloseTimeHuman      string        `json:"close_time_human"`
	LedgerHash          string        `json:"ledger_hash"`
	TotalCoins          string        `json:"total_coins"`
	Closed              bool          `json:"closed"`
	CloseTimeISO        string        `json:"close_time_iso"`
	ParentCloseTime     int64         `json:"parent_close_time"`
	ParentHash          string        `json:"parent_hash"`
	TransactionHash     string        `json:"transaction_hash"`
}

type Transaction struct {
	DeliverMax         string        `json:"DeliverMax"`
	Account            string        `json:"Account"`
	Destination        string        `json:"Destination"`
	MetaData           MetaData      `json:"metaData"`
	TransactionType    string        `json:"TransactionType"`
	TxnSignature       string        `json:"TxnSignature"`
	SigningPubKey      string        `json:"SigningPubKey"`
	Amount             string        `json:"Amount"`
	Fee                string        `json:"Fee"`
	Sequence           int64         `json:"Sequence"`
	DestinationTag     *int64        `json:"DestinationTag,omitempty"`
	Hash               string        `json:"hash"`
	LastLedgerSequence *int64        `json:"LastLedgerSequence,omitempty"`
	Flags              *int64        `json:"Flags,omitempty"`
	TicketSequence     *int64        `json:"TicketSequence,omitempty"`
	SourceTag          *int64        `json:"SourceTag,omitempty"`
	Memos              []MemoElement `json:"Memos,omitempty"`
}

type MemoElement struct {
	Memo MemoMemo `json:"Memo"`
}

type MemoMemo struct {
	MemoType   string `json:"MemoType"`
	MemoData   string `json:"MemoData"`
	MemoFormat string `json:"MemoFormat"`
}

type MetaData struct {
	AffectedNodes     []AffectedNode `json:"AffectedNodes"`
	TransactionResult string         `json:"TransactionResult"`
	TransactionIndex  int64          `json:"TransactionIndex"`
	DeliveredAmount   string         `json:"delivered_amount"`
}

type AffectedNode struct {
	ModifiedNode *ModifiedNode `json:"ModifiedNode,omitempty"`
	DeletedNode  *DeletedNode  `json:"DeletedNode,omitempty"`
	CreatedNode  *CreatedNode  `json:"CreatedNode,omitempty"`
}

type CreatedNode struct {
	LedgerIndex     string          `json:"LedgerIndex"`
	LedgerEntryType LedgerEntryType `json:"LedgerEntryType"`
	NewFields       NewFields       `json:"NewFields"`
}

type NewFields struct {
	Account  string `json:"Account"`
	Sequence int64  `json:"Sequence"`
	Balance  string `json:"Balance"`
}

type DeletedNode struct {
	LedgerIndex     string                 `json:"LedgerIndex"`
	FinalFields     DeletedNodeFinalFields `json:"FinalFields"`
	LedgerEntryType string                 `json:"LedgerEntryType"`
}

type DeletedNodeFinalFields struct {
	Account           string `json:"Account"`
	PreviousTxnLgrSeq int64  `json:"PreviousTxnLgrSeq"`
	OwnerNode         string `json:"OwnerNode"`
	TicketSequence    int64  `json:"TicketSequence"`
	PreviousTxnID     string `json:"PreviousTxnID"`
	Flags             int64  `json:"Flags"`
}

type ModifiedNode struct {
	LedgerIndex       string                  `json:"LedgerIndex"`
	FinalFields       ModifiedNodeFinalFields `json:"FinalFields"`
	PreviousFields    *PreviousFields         `json:"PreviousFields,omitempty"`
	PreviousTxnLgrSeq int64                   `json:"PreviousTxnLgrSeq"`
	LedgerEntryType   LedgerEntryType         `json:"LedgerEntryType"`
	PreviousTxnID     string                  `json:"PreviousTxnID"`
}

type ModifiedNodeFinalFields struct {
	Account       *string `json:"Account,omitempty"`
	OwnerCount    *int64  `json:"OwnerCount,omitempty"`
	Flags         int64   `json:"Flags"`
	Sequence      *int64  `json:"Sequence,omitempty"`
	Balance       *string `json:"Balance,omitempty"`
	TicketCount   *int64  `json:"TicketCount,omitempty"`
	Owner         *string `json:"Owner,omitempty"`
	IndexNext     *string `json:"IndexNext,omitempty"`
	IndexPrevious *string `json:"IndexPrevious,omitempty"`
	RootIndex     *string `json:"RootIndex,omitempty"`
}

type PreviousFields struct {
	Sequence    *int64 `json:"Sequence,omitempty"`
	Balance     string `json:"Balance"`
	OwnerCount  *int64 `json:"OwnerCount,omitempty"`
	TicketCount *int64 `json:"TicketCount,omitempty"`
}

type LedgerEntryType string

const (
	AccountRoot   LedgerEntryType = "AccountRoot"
	DirectoryNode LedgerEntryType = "DirectoryNode"
)
