package ripplegateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	NewAddressURL    = "%s/v1/account/new"
	SignAndPushTxURL = "%s/v1/account/signTx"
)

type Gateway struct {
	URL string
}

// BaseResp :
type BaseResp struct {
	ReturnCode int64  `json:"returnCode"`
	ReturnMsg  string `json:"returnMsg"`
}

// NewAddressResp :
type NewAddressResp struct {
	Base    BaseResp `json:"base"`
	Address string   `json:"address"`
	Pubkey  string   `json:"pubkey"`
}

type SignTxReq struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Fee    string `json:"fee"`
	Amount string `json:"amount"`
	Tag    uint32 `json:"tag"`
	Seq    uint32 `json:"seq"`
}

type SignTxResp struct {
	Base      BaseResp          `json:"base"`
	TxData    []byte            `json:"txdata"`
	Extension map[string]string `json:"extension"`
}

func NewRippleGateway(url string) *Gateway {
	return &Gateway{
		URL: url,
	}
}

func (receiver *Gateway) Close() {

}

// sendGetRequest :
func sendGetRequest(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request to %v, response code is %v, not 200", url, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}

// sendPostRequest :
func sendPostRequest(url string, header map[string]string, data interface{}, result interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Set("content-Type", "application/json")
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}

func (receiver *Gateway) NewAddress() (string, string, error) {
	var resp NewAddressResp
	url := fmt.Sprintf(NewAddressURL, receiver.URL)

	err := sendGetRequest(url, &resp)
	if err != nil {
		logrus.Errorf("error when create new address: %s", err.Error())
		return "", "", err
	}

	if resp.Base.ReturnCode != 0 {
		err := fmt.Errorf("error chain response(%v): %v", resp.Base.ReturnCode, resp.Base.ReturnMsg)
		logrus.Errorf("get create new address return error: %s", err.Error())
		return "", "", err
	}

	return resp.Address, resp.Pubkey, nil
}

func (receiver *Gateway) SignTx(from, to, amount, fee string, tag, seq uint32) ([]byte, error) {
	var resp SignTxResp

	reqInfo := SignTxReq{
		From:   from,
		To:     to,
		Amount: amount,
		Fee:    fee,
		Tag:    tag,
		Seq:    seq,
	}

	url := fmt.Sprintf(SignAndPushTxURL, receiver.URL)

	err := sendPostRequest(url, nil, reqInfo, &resp)
	if err != nil {
		logrus.Errorf("SignAndPushTx send request error: %v", err.Error())
		return nil, err
	}

	if resp.Base.ReturnCode != 0 {
		err := fmt.Errorf("error chain response(%v): %v", resp.Base.ReturnCode, resp.Base.ReturnMsg)
		logrus.Errorf("singTx return error: %v", err.Error())
		return nil, err
	}

	return resp.TxData, nil
}
