package main

import (
	"github.com/goodwood511/ripple_lib/ripple/xrpclient"
	"github.com/sirupsen/logrus"
)

func main() {
	client := xrpclient.NewClient("https://s.altnet.rippletest.net:51234/")

	latestBlock, err := client.GetLatestLedgerIndex()
	if err != nil {
		logrus.Warnln(err)
		return
	}
	logrus.Infoln(latestBlock)
	latestBlock = 8546532

	transactions, blockTime, status, err := client.GetLedgerTransactions(latestBlock)
	if err != nil {
		logrus.Warnln(err)
		return
	}
	logrus.Infoln("Block time:", blockTime, "status", status)
	payments := xrpclient.ParsePayments(transactions)

	for _, payment := range payments {
		logrus.Infoln("Hash:", payment.Hash,
			"From:", payment.Account,
			"To:", payment.Destination,
			"Amount (drops):", payment.Amount)
	}
}
