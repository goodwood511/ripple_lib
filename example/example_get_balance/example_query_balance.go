package main

import (
	"fmt"
	"github.com/goodwood511/ripple_lib/ripple/xrpclient"
)

func main() {
	client := xrpclient.NewClient("https://s.altnet.rippletest.net:51234/")
	addr := "rs55BYv4iNsErKuVrT3R5tnUktbnz68XrL"

	xrp, err := client.GetBalance(addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("XRP: %s\n", xrp.String())

	tokens, err := client.GetTokenBalances(addr)
	if err != nil {
		panic(err)
	}
	for token, balance := range tokens {
		fmt.Printf("%s: %.6f\n", token, balance)
	}
}
