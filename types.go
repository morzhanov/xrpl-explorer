package main

type LedgerSubscriptionStartMessage struct {
	Id     string                `json:"id"`
	Result LedgerSubscriptionRes `json:"result"`
	Status string                `json:"status"`
	Type   string                `json:"type"`
}

type LedgerSubscriptionRes struct {
	FeeBase          int    `json:"fee_base"`
	FeeRef           int    `json:"fee_ref"`
	LedgerHash       string `json:"ledger_hash"`
	LedgerIndex      int    `json:"ledger_index"`
	LedgerTime       int    `json:"ledger_time"`
	ReserveBase      int    `json:"reserve_base"`
	ReserveInc       int    `json:"reserve_inc"`
	ValidatedLedgers string `json:"validated_ledgers"`
}

type TransactionSubscriptionStartMessage struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

type TransactionSubscriptionRes struct {
	EngineResult        string `json:"engine_result"`
	EngineResultCode    int    `json:"engine_result_code"`
	EngineResultMessage string `json:"engine_result_message"`
	LedgerHash          string `json:"ledger_hash"`
	LedgerIndex         int    `json:"ledger_index"`
	Meta                struct {
		AffectedNodes []struct {
			ModifiedNode struct {
				FinalFields struct {
					Account    string `json:"Account"`
					Balance    string `json:"Balance"`
					Flags      int    `json:"Flags"`
					OwnerCount int    `json:"OwnerCount"`
					Sequence   int    `json:"Sequence"`
				} `json:"FinalFields"`
				LedgerEntryType string `json:"LedgerEntryType"`
				LedgerIndex     string `json:"LedgerIndex"`
				PreviousFields  struct {
					Balance  string `json:"Balance"`
					Sequence int    `json:"Sequence"`
				} `json:"PreviousFields"`
				PreviousTxnID     string `json:"PreviousTxnID"`
				PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"`
			} `json:"ModifiedNode"`
		} `json:"AffectedNodes"`
		TransactionIndex  int    `json:"TransactionIndex"`
		TransactionResult string `json:"TransactionResult"`
	} `json:"meta"`
	Status      string `json:"status"`
	Transaction struct {
		Account            string `json:"Account"`
		Fee                string `json:"Fee"`
		Flags              int64  `json:"Flags"`
		LastLedgerSequence int    `json:"LastLedgerSequence"`
		Sequence           int    `json:"Sequence"`
		SigningPubKey      string `json:"SigningPubKey"`
		TakerGets          string `json:"TakerGets"`
		TakerPays          struct {
			Currency string `json:"currency"`
			Issuer   string `json:"issuer"`
			Value    string `json:"value"`
		} `json:"TakerPays"`
		TransactionType string `json:"TransactionType"`
		TxnSignature    string `json:"TxnSignature"`
		Date            int    `json:"date"`
		Hash            string `json:"hash"`
		OwnerFunds      string `json:"owner_funds"`
	} `json:"transaction"`
	Type      string `json:"type"`
	Validated bool   `json:"validated"`
}
