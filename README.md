# go-electrumx [![GoDoc](https://godoc.org/github.com/FryDay/go-electrumx/electrumx?status.svg)](https://godoc.org/github.com/FryDay/go-electrumx/electrumx)

This repository is a fork of [ccampbell/go-electrum](https://github.com/ccampbell/go-electrum) which is a fork of [d4l3k/go-electrum](https://github.com/d4l3k/go-electrum), both of which are unmaintained now.

A pure Go [ElectrumX](https://electrumx.readthedocs.io/) bitcoin client library. This makes it easy to write bitcoin based services using Go without having to run a full bitcoin node, or querying for bitcoin data that is not available through a regular full node. 

![go-electrumx](/media/logo.png)

Packages provided

- [electrumx](https://godoc.org/github.com/FryDay/go-electrumx/electrumx) - Dependency-free (except for the standard library) module for using JSON-RPC to talk directly to ElectrumX servers.
- [wallet](https://godoc.org/github.com/FryDay/go-electrumx/wallet) - A bitcoin wallet built on [btcwallet](https://github.com/btcsuite/btcwallet) with ElectrumX as the backend.

## Usage

See [examples](https://github.com/FryDay/go-electrumx/tree/master/examples) for more.

# Donate

If you find these packages useful, please consider flipping me a few satoshis :hugs:

<img src="https://github.com/FryDay/go-electrumx/raw/master/media/qr.png" alt="bc1qtsh9vfy2xc5n972lhh0rezkj0ny7kjktwsrl6z" title="bc1qtsh9vfy2xc5n972lhh0rezkj0ny7kjktwsrl6z" width="300px" />

# License

go-electrumx is licensed under the MIT license.
