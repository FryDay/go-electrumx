package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/FryDay/go-electrumx/electrumx"
)

var (
	// Should specify a available server(IP:PORT) if connecting to the
	// following server failed.
	serverAddr = "blockstream.info:700"
)

// Wrap the 'real' main method to get defers to work properly.
func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	electrumx.DebugMode = true

	node := electrumx.NewNode()
	defer node.Disconnect()

	connectCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	if err := node.Connect(connectCtx, serverAddr, &tls.Config{}); err != nil {
		return fmt.Errorf("connect node: %w", err)
	}

	go func() {
		// Remember to drain errors! Since communication is async not all errors
		// will happen as a direct response to requests.
		err := <-node.Errors()
		log.Printf("ran into error: %s", err)
		node.Disconnect()
	}()

	//Send server.ping request in order to keep alive connection to
	//electrum server
	go func() {
		for {
			if err := node.Ping(ctx); err != nil {
				log.Fatal(err)
			}

			select {
			case <-time.After(5 * time.Second):
			case <-ctx.Done():
				return
			}

		}
	}()

	version, err := node.ServerVersion(ctx)
	if err != nil {
		return fmt.Errorf("server version: %w", err)
	}

	fmt.Printf("Version: %v\n\n", version)

	banner, err := node.ServerBanner(ctx)
	if err != nil {
		return fmt.Errorf("server banner: %w", err)
	}
	fmt.Printf("Banner: %s\n\n", banner)

	address, err := node.ServerDonationAddress(ctx)
	if err != nil {
		return fmt.Errorf("server donation address: %w", err)
	}
	fmt.Printf("Address: %s\n\n", address)

	peers, err := node.ServerPeersSubscribe(ctx)
	if err != nil {
		return fmt.Errorf("server peers subscribe: %w", err)
	}
	fmt.Printf("Peers: %+v\n\n", peers)

	headerChan, err := node.BlockchainHeadersSubscribe(ctx)
	if err != nil {
		return fmt.Errorf("blockchain headers: %w", err)
	}
	go func() {
		for header := range headerChan {
			fmt.Printf("Headers: %+v\n\n", header)
		}
	}()

	transaction, err := node.BlockchainTransactionGet(
		ctx, "f53a8b83f85dd1ce2a6ef4593e67169b90aaeb402b3cf806b37afc634ef71fbc", false)
	if err != nil {
		fmt.Printf("blockchain transaction get: %s\n", err)
	} else {
		fmt.Printf("Transaction: %+v\n\n", transaction)
	}

	// If you're connecting to a node that support address queries, change this
	// to true.
	nodeSupportsAddressQueries := false
	if nodeSupportsAddressQueries {
		balance, err := node.BlockchainAddressGetBalance(ctx, "bc1qv5wppm0xwarzwea9xxascc5rne7c0c296h7y5p")
		if err != nil {
			return err
		}

		fmt.Printf("address balance: %+v\n\n", balance)
	}

	const until = time.Minute * 15
	fmt.Printf("leaving connection open for %s\n", until)
	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-time.After(until):
	}

	return nil

}
