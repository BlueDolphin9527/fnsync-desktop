package client

import (
	"fmt"
	"net"

	"github.com/cxfksword/fnsync-desktop/entity"
	"github.com/cxfksword/fnsync-desktop/msg"
	"github.com/cxfksword/fnsync-desktop/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	DEFAULT_PORT = 21365
)
var handShake *HandShake = &HandShake{}

type HandShake struct {
}

func (h *HandShake) Send(targets []entity.Device, portIncrement int) {
	targetIds := []string{}
	for _, v := range targets {
		targetIds = append(targetIds, v.Id)
	}
	handshake := msg.Builder.MakeHandshake(uuid.NewString(), targetIds)

	for _, target := range targets {
		port := DEFAULT_PORT + portIncrement
		local, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
		remote, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", target.LastIp, port))
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
		conn, err := net.DialUDP("udp", local, remote)
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
		defer conn.Close()

		handshakeJson := utils.ToJSON(handshake)
		_, err = conn.Write(handshakeJson)
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
	}

}

func (h *HandShake) BroadcastLocalNetwork(targets []entity.Device, portIncrement int) {
	targetIds := []string{}
	for _, v := range targets {
		targetIds = append(targetIds, v.Id)
	}
	handshake := msg.Builder.MakeHandshake(uuid.NewString(), targetIds)

	broadcastIps := utils.GetLocalBroadcastInterface()
	for _, ip := range broadcastIps {

		port := DEFAULT_PORT + portIncrement
		local, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
		remote, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
		conn, err := net.DialUDP("udp", local, remote)
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
		defer conn.Close()

		handshakeJson := utils.ToJSON(handshake)
		log.Info().Msgf("Handshake broadcast. target: %s  msg: %s", ip, string(handshakeJson))
		_, err = conn.Write(handshakeJson)
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
	}
}
