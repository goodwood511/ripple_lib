package xrpclient

// GetLatestLedgerIndex returns the latest validated ledger index
func (c *Client) GetLatestLedgerIndex() (uint64, error) {
	req := map[string]interface{}{
		"method": "ledger_current",
		"params": []interface{}{map[string]interface{}{}},
	}
	var res struct {
		Result struct {
			LedgerIndex uint64 `json:"ledger_current_index"`
		} `json:"result"`
	}
	_, err := c.client.R().SetBody(req).SetResult(&res).Post(c.rpcURL)
	if err != nil {
		return 0, err
	}
	return res.Result.LedgerIndex, nil
}

// GetLedgerTransactions returns the transactions in a specific ledger (only payments)
func (c *Client) GetLedgerTransactions(ledgerIndex uint64) ([]Transaction, error) {
	req := map[string]interface{}{
		"method": "ledger",
		"params": []interface{}{map[string]interface{}{
			"ledger_index": ledgerIndex,
			"transactions": true,
			"expand":       true,
		}},
	}
	res := RspTransaction{}
	_, err := c.client.R().SetBody(req).SetResult(&res).Post(c.rpcURL)
	if err != nil {
		return nil, err
	}
	return res.Result.Ledger.Transactions, nil
}

// ParsePayments filters payment transactions with XRP only (Amount is string)
func ParsePayments(txs []Transaction) []Transaction {
	var result []Transaction
	for _, tx := range txs {
		if tx.TransactionType != "Payment" {
			continue
		}
		// If Amount is string, it's XRP; if it's object, it's token
		if tx.Amount != "" && tx.Destination != "" && tx.Account != "" {
			result = append(result, tx)
		}
	}
	return result
}
