package xrpclient

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
)

type Client struct {
	rpcURL string
	client *resty.Client
}

func NewClient(rpcURL string) *Client {
	return &Client{
		rpcURL: rpcURL,
		client: resty.New(),
	}
}

type accountInfoResponse struct {
	Result struct {
		AccountData struct {
			Balance string `json:"Balance"`
		} `json:"account_data"`
	} `json:"result"`
}

func (c *Client) GetBalance(address string) (decimal.Decimal, error) {
	req := map[string]interface{}{
		"method": "account_info",
		"params": []interface{}{
			map[string]interface{}{
				"account":      address,
				"ledger_index": "validated",
			},
		},
	}
	var res accountInfoResponse
	_, err := c.client.R().SetBody(req).SetResult(&res).Post(c.rpcURL)
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromString(res.Result.AccountData.Balance)
}

type accountLinesResponse struct {
	Result struct {
		Lines []struct {
			Currency string `json:"currency"`
			Issuer   string `json:"account"`
			Balance  string `json:"balance"`
		} `json:"lines"`
	} `json:"result"`
}

func (c *Client) GetTokenBalances(address string) (map[string]float64, error) {
	req := map[string]interface{}{
		"method": "account_lines",
		"params": []interface{}{
			map[string]interface{}{
				"account": address,
			},
		},
	}
	var res accountLinesResponse
	_, err := c.client.R().SetBody(req).SetResult(&res).Post(c.rpcURL)
	if err != nil {
		return nil, err
	}
	result := make(map[string]float64)
	for _, line := range res.Result.Lines {
		var value float64
		fmt.Sscanf(line.Balance, "%f", &value)
		key := fmt.Sprintf("%s.%s", line.Currency, line.Issuer)
		result[key] = value
	}
	return result, nil
}
