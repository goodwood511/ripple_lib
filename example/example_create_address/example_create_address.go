package main

import (
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"ripple_lib/ripple/rippleaddr"
)

func main() {
	keyBin, address, pubkey, err := rippleaddr.NewAddress()
	if err != nil {
		logrus.Infoln("rippleaddr.NewAddress() err", err)
		return
	}

	key := hex.EncodeToString(keyBin)

	logrus.Infoln("key", key)
	logrus.Infoln("address", address)
	logrus.Infoln("pubkey", pubkey)
}
