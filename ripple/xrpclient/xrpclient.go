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
		// 错误时存在这些字段
		Error        string `json:"error,omitempty"`
		ErrorCode    int    `json:"error_code,omitempty"`
		ErrorMessage string `json:"error_message,omitempty"`
		Status       string `json:"status,omitempty"` // "error" or "success"
	} `json:"result"`
}

/*
	{
	    "result": {
	        "account": "rKp7KgcYjdEQepQLc27ZHz76ukLwE4S1CN",
	        "error": "actNotFound",
	        "error_code": 19,
	        "error_message": "Account not found.",
	        "ledger_hash": "D59CDB76635F01A7BA054D8CFACAACA9D7ECF17C67ACB065CC55208DAEC18BF2",
	        "ledger_index": 8629213,
	        "request": {
	            "account": "rKp7KgcYjdEQepQLc27ZHz76ukLwE4S1CN",
	            "command": "account_info",
	            "ledger_index": "validated"
	        },
	        "status": "error",
	        "validated": true
	    }
	}
*/
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
	if res.Result.Error != "" {
		if res.Result.ErrorCode == 19 {
			return decimal.Zero, nil
		}
		return decimal.Zero, fmt.Errorf("%s: %s", res.Result.Error, res.Result.ErrorMessage)
	}
	return decimal.NewFromString(res.Result.AccountData.Balance)
}

type accountLinesResponse struct {
	Result struct {
		// 错误时存在这些字段
		Error        string `json:"error,omitempty"`
		ErrorCode    int    `json:"error_code,omitempty"`
		ErrorMessage string `json:"error_message,omitempty"`
		Status       string `json:"status,omitempty"` // "error" or "success"
		Lines        []struct {
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
