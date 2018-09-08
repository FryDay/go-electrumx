package electrum

// ServerVersion returns the server's version.
// http://docs.electrum.org/en/latest/protocol.html#server-version
func (n *Node) ServerVersion() (string, error) {
	resp := &basicResp{}
	err := n.request("server.version", []string{}, resp)
	return resp.Result, err
}

// ServerBanner returns the server's banner.
// http://docs.electrum.org/en/latest/protocol.html#server-banner
func (n *Node) ServerBanner() (string, error) {
	resp := &basicResp{}
	err := n.request("server.banner", []string{}, resp)
	return resp.Result, err
}

// ServerDonationAddress returns the donation address of the server.
// http://docs.electrum.org/en/latest/protocol.html#server-donation-address
func (n *Node) ServerDonationAddress() (string, error) {
	resp := &basicResp{}
	err := n.request("server.donation_address", []string{}, resp)
	return resp.Result, err
}

// ServerPeersSubscribe requests peers from a server.
// http://docs.electrum.org/en/latest/protocol.html#server-peers-subscribe
func (n *Node) ServerPeersSubscribe() ([][]interface{}, error) {
	resp := &struct {
		Peers [][]interface{} `json:"result"`
	}{}
	err := n.request("server.peers.subscribe", []string{}, resp)
	return resp.Peers, err
}
