package electrumx

import "context"

// TODO implement
// ServerAddPeer add a peer (but only if the peer resolves to the source).
// method: "server.add_peer"
//
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#server-add-peer
func (n *Node) ServerAddPeer() error {
	return ErrNotImplemented
}

// ServerBanner returns the server banner text.
//
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#server-banner
func (n *Node) ServerBanner(ctx context.Context) (string, error) {
	resp := &basicResp{}
	err := n.request(ctx, "server.banner", []interface{}{}, resp)
	return resp.Result, err
}

// ServerDonationAddress returns the donation address of the server.
//
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#server-donation-address
func (n *Node) ServerDonationAddress(ctx context.Context) (string, error) {
	resp := &basicResp{}
	err := n.request(ctx, "server.donation_address", []interface{}{}, resp)
	return resp.Result, err
}

// TODO implement
// ServerFeatures returns the server features dictionary.
// method: "server.features"
//
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#server-features
func (n *Node) ServerFeatures() error {
	return ErrNotImplemented
}

// ServerPeersSubscribe requests peers from a server.
//
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#server-peers-subscribe
func (n *Node) ServerPeersSubscribe(ctx context.Context) ([][]interface{}, error) {
	resp := &struct {
		Peers [][]interface{} `json:"result"`
	}{}
	err := n.request(ctx, "server.peers.subscribe", []interface{}{}, resp)
	return resp.Peers, err
}

// TODO return concrete struct instead of unnamed string slice
// ServerVersion returns the server's version.
//
// http://docs.electrum.org/en/latest/protocol.html#server-version
func (n *Node) ServerVersion(ctx context.Context) ([]string, error) {
	resp := &struct {
		Result []string `json:"result"`
	}{}
	err := n.request(ctx, "server.version", []interface{}{}, resp)
	return resp.Result, err
}

func (n *Node) Ping(ctx context.Context) error {
	err := n.request(ctx, "server.ping", []interface{}{}, nil)

	return err
}
