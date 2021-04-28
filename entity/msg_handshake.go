package entity

type HandshakeMsg struct {
	PeerId        string   `json:"peerid"`
	PeerName      string   `json:"peername"`
	PeerPort      int      `json:"peerport"`
	OldConnection bool     `json:"oldconnection,omitempty"`
	Target        []string `json:"target,omitempty"`
	Token         string   `json:"token"`
}
