package xrpclient

import (
	"github.com/shopspring/decimal"
)

type accountInfoFull struct {
	Result struct {
		AccountData struct {
			Balance    string `json:"Balance"`
			OwnerCount uint32 `json:"OwnerCount"`
		} `json:"account_data"`
	} `json:"result"`
}

// GetSendableAmount returns the maximum XRP (in drops) that can be safely sent from the account,
// considering the base reserve and owner reserve.
func (c *Client) GetSendableAmount(address string) (decimal.Decimal, error) {
	req := map[string]interface{}{
		"method": "account_info",
		"params": []interface{}{
			map[string]interface{}{
				"account":      address,
				"ledger_index": "validated",
			},
		},
	}

	var res accountInfoFull
	_, err := c.client.R().SetBody(req).SetResult(&res).Post(c.rpcURL)
	if err != nil {
		return decimal.Zero, err
	}

	balance, err := decimal.NewFromString(res.Result.AccountData.Balance)
	if err != nil {
		return decimal.Zero, err
	}

	baseReserve := decimal.NewFromInt(1_000_000) // 1 XRP = 1,000,000 drops
	ownerReserve := decimal.NewFromInt(200_000).Mul(decimal.NewFromInt(int64(res.Result.AccountData.OwnerCount)))
	totalReserve := baseReserve.Add(ownerReserve)

	available := balance.Sub(totalReserve)
	if available.LessThan(decimal.Zero) {
		return decimal.Zero, nil
	}
	return available, nil
}
