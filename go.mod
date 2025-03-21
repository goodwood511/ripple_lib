module github.com/goodwood511/ripple_lib

go 1.23.0

toolchain go1.24.0

require (
	github.com/bits-and-blooms/bitset v1.22.0
	github.com/btcsuite/btcd v0.23.4
	github.com/golang/glog v1.2.4
	github.com/goodwood511/ripple_lib/ripple-sdk/crypto v0.0.0-20250321020749-2e530e8b6776
	github.com/gorilla/websocket v1.5.3
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/crypto v0.36.0
)

require (
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/sys v0.31.0 // indirect
)

replace github.com/btcsuite/btcd => github.com/btcsuite/btcd v0.22.1
