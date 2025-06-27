package main

import (
	"github.com/goodwood511/ripple_lib/ripple/xrpclient"
	"github.com/sirupsen/logrus"
)

func main() {
	client := xrpclient.NewClient("https://s.altnet.rippletest.net:51234/")
	hash := "C30A073A9B3068232D769A9B10340B7E8445EBA9D6A208B60A6B2A57154114EB"

	rsp, err := client.GetTransaction(hash)
	if err != nil {
		panic(err)
	}

	logrus.Infoln(rsp)
}
