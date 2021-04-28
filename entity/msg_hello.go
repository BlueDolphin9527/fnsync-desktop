package entity

type HelloMsg struct {
	PeerId        string `json:"peerid"`
	PeerName      string `json:"peername"`
	PeerPort      int    `json:"peerport"`
	OldConnection bool   `json:"oldconnection,omitempty"`
	Token         string `json:"token"`
	Ips           string `json:"ips,omitempty"`
}
