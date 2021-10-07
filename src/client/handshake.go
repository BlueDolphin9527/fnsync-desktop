package client

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/cxfksword/fnsync-desktop/config"
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

func (h *HandShake) Send(handshakeMsg entity.HandshakeMsg, targets []entity.Device, portIncrement int) {
	port := DEFAULT_PORT + portIncrement

	for _, target := range targets {
		local, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
		address := fmt.Sprintf("%s:%d", target.LastIp, port) // ipv4
		if strings.Contains(target.LastIp, ":") {
			address = fmt.Sprintf("[%s]:%d", target.LastIp, port) // ipv6
		}
		remote, err := net.ResolveUDPAddr("udp", address)
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

		handshakeJson := utils.ToJSON(handshakeMsg)
		log.Info().Msgf("Handshake. target: %s msg: %s", address, string(handshakeJson))
		_, err = conn.Write(handshakeJson)
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
	}

}

func (h *HandShake) BroadcastLocalNetwork(handshakeMsg entity.HandshakeMsg, portIncrement int) {
	port := DEFAULT_PORT + portIncrement
	broadcastIps := utils.GetLocalBroadcastInterface()
	for _, ip := range broadcastIps {
		local, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
		address := fmt.Sprintf("%s:%d", ip, port) // ipv4
		if strings.Contains(ip, ":") {
			address = fmt.Sprintf("[%s]:%d", ip, port) // ipv6
		}
		remote, err := net.ResolveUDPAddr("udp", address)
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

		handshakeJson := utils.ToJSON(handshakeMsg)
		log.Info().Msgf("Handshake broadcast. target: %s  msg: %s", address, string(handshakeJson))
		_, err = conn.Write(handshakeJson)
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
	}
}

func StartHandshake() {
	defer func() { log.Info().Msgf("Quit autoconnect handshake.") }()
	if !config.App.ConnectOnStartup {
		return
	}
	if len(config.App.Devices) <= 0 {
		return
	}

	targets := []entity.Device{}
	for _, v := range config.App.Devices {
		targets = append(targets, v)
	}
	targetIds := []string{}
	for _, v := range targets {
		targetIds = append(targetIds, v.Id)
	}
	handshakeMsg := msg.Builder.MakeHandshake(uuid.NewString(), targetIds)

	log.Info().Msgf("Start send autoconnect handshake...")

	timeout := AUTO_CONNECT_TIMEOUT_MILLS
	portIncrement := 0
	for timeout > 0 && len(targets) > 0 {

		handShake.Send(handshakeMsg, targets, portIncrement)
		handShake.BroadcastLocalNetwork(handshakeMsg, portIncrement)

		timeout -= 2500
		portIncrement++
		time.Sleep(2500 * time.Millisecond)
	}
}

func StartHandshakeDevice(v entity.Device) {
	targets := []entity.Device{v}
	targetIds := []string{v.Id}
	handshakeMsg := msg.Builder.MakeHandshake(uuid.NewString(), targetIds)

	log.Info().Msgf("Start send autoconnect handshake...")

	timeout := AUTO_CONNECT_TIMEOUT_MILLS
	portIncrement := 0
	for timeout > 0 && len(targets) > 0 {

		handShake.Send(handshakeMsg, targets, portIncrement)
		handShake.BroadcastLocalNetwork(handshakeMsg, portIncrement)

		timeout -= 2500
		portIncrement++
		time.Sleep(2500 * time.Millisecond)
	}
}
