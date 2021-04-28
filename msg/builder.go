package msg

import (
	"github.com/cxfksword/fnsync-desktop/config"
	"github.com/cxfksword/fnsync-desktop/entity"
	"github.com/cxfksword/fnsync-desktop/utils"
)

var Builder *msgBuilder = &msgBuilder{}

type msgBuilder struct {
}

func (m *msgBuilder) MakeHello(token string) entity.HelloMsg {
	hello := entity.HelloMsg{
		PeerId:   config.App.MachineId,
		PeerName: utils.MachineName(),
		PeerPort: config.App.GetListenPort(),
		Token:    token,
		Ips:      utils.GetAllInterface(),
	}
	return hello
}

func (m *msgBuilder) MakeHandshake(token string, targets []string) entity.HandshakeMsg {
	handshake := entity.HandshakeMsg{
		PeerId:        config.App.MachineId,
		PeerName:      utils.MachineName(),
		PeerPort:      config.App.GetListenPort(),
		OldConnection: true,
		Token:         token,
		Target:        targets,
	}
	return handshake
}
