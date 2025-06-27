package xrpclient

type txResponse struct {
	Result map[string]interface{} `json:"result"`
}

func (c *Client) GetTransaction(hash string) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"method": "tx",
		"params": []interface{}{
			map[string]interface{}{
				"transaction": hash,
				"binary":      false,
			},
		},
	}
	var res txResponse
	_, err := c.client.R().SetBody(req).SetResult(&res).Post(c.rpcURL)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}
