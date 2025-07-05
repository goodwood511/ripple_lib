package main

import (
	"github.com/goodwood511/ripple_lib/ripple-sdk/data"
	"github.com/goodwood511/ripple_lib/ripple/ripple"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func main() {

	host := os.Getenv("wshost")
	fromAddress := os.Getenv("fromAddress")

	r, err := ripple.NewRipple(host)
	if err != nil {
		logrus.Warnln("NewRipple err", err.Error(), "host", host)
		return
	}

	address, err := data.NewAccountFromAddress(fromAddress)
	if err != nil {
		logrus.Warnln("fetchSignatureFormRemote NewAccountFromAddress err", err.Error(), "fromAddress", fromAddress)
		return
	}

	accountInfo, err := r.Client.AccountInfo(*address)
	if err != nil {
		logrus.Warnln("fetchSignatureFormRemote AccountInfo err", err.Error(), "fromAddress", fromAddress)
		return
	}
	logrus.Infoln("fetchSignatureFormRemote AccountInfo accountInfo", accountInfo)
	seq := strconv.Itoa(int(*accountInfo.AccountData.Sequence))

	logrus.Infoln("seq", seq)

	payment := &data.Payment{}
	hashObj, err := r.BroadcastSignleSignTransaction(payment)
	if err != nil {
		logrus.Warnln("fetchSignatureFormRemote BroadcastSignleSignTransaction err", err)
		return
	}

	logrus.Infoln("hashObj.String()", hashObj.String())

}
