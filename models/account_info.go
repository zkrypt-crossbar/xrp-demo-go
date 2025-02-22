package models

type AccountInfoResp struct {
	Result *AccountInfoResult
}

type AccountInfoResult struct {
	Validated          bool
	Status             string
	LedgerCurrentIndex int64            `json:"ledger_current_index"`
	AccountData        *AccountInfoData `json:"account_data"`

	ErrorMessage string `json:"error_message"`
}

type AccountInfoData struct {
	Account           string `json:"Account"`
	Balance           string `json:"Balance"`
	Flags             int    `json:"Flags"`
	LedgerEntryType   string `json:"LedgerEntryType"`
	OwnerCount        int    `json:"OwnerCount"`
	PreviousTxnID     string `json:"PreviousTxnID"`
	PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"`
	Sequence          uint32 `json:"Sequence"`

	Index string `json:"index"`
}

type AccountTransactionResp struct {
	Result       string
	Count        int
	Marker       string
	Transactions []*AccountTransactionData
}

type AccountTransactionData struct {
	Hash        string
	LedgerIndex int64 `json:"ledger_index"`
	Date        string
	Tx          *AccountTransactionTx
	Meta        *AccountTransactionMeta
}

type AccountTransactionTx struct {
	TransactionType    string
	Flags              int64
	Sequence           int64
	LastLedgerSequence int64
	Amount             string
	Fee                string
	SigningPubKey      string
	TxnSignature       string
	Account            string
	Destination        string
	DestinationTag     int64
}

type AccountTransactionMeta struct {
	TransactionIndex int64
	// AffectedNodes     []*AccountTransactionAffectedNodes
	TransactionResult string
	DeliveredAmount   string `json:"delivered_amount"`
}
