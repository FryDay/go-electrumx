# go-electrumx [![GoDoc](https://godoc.org/github.com/FryDay/go-electrumx?status.svg)](https://godoc.org/github.com/FryDay/go-electrumx)

This repository is a fork of [ccampbell/go-electrum](https://github.com/ccampbell/go-electrum) which is a fork of [d4l3k/go-electrum](https://github.com/d4l3k/go-electrum), both of which are unmaintained now.

A pure Go [ElectrumX](https://electrumx.readthedocs.io/) bitcoin library. This makes it easy to write bitcoin based services using Go without having to run a full bitcoin node.

![go-electrumx](/media/logo.png)

Packages provided

- [electrumx](https://godoc.org/github.com/FryDay/go-electrumx/electrumx) - Library for using JSON-RPC to talk directly to ElectrumX servers.
- [wallet](https://godoc.org/github.com/FryDay/go-electrumx/wallet) - A bitcoin wallet built on [btcwallet](https://github.com/btcsuite/btcwallet) with ElectrumX as the backend.

## Usage

See [example/](https://github.com/FryDay/go-electrumx/tree/master/example) for more.

### electrumx [![GoDoc](https://godoc.org/github.com/FryDay/go-electrumx/electrumx?status.svg)](https://godoc.org/github.com/FryDay/go-electrumx/electrumx)

```bash
$ go get -u github.com/FryDay/go-electrumx/electrumx
```

```go
package main

import (
  "log"

  "github.com/FryDay/go-electrumx/electrumx"
)

func main() {
    // turn on debug mode
    electrumx.DebugMode = true

	node := electrumx.NewNode()
    // the specified ip is testnet server
	if err := node.ConnectTCP("39.104.125.149:9629"); err != nil {
		log.Fatal(err)
	}

    // please use bitcoin address accordant to the server environment
	balance, err := node.BlockchainAddressGetBalance("n4FyJMDYXJmPEm7cffFLrwLXvGWn8cW9q2")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Address balance: %+v", balance)
}
```

# License

go-electrumx is licensed under the MIT license.
