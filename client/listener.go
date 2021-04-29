package client

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/cxfksword/fnsync-desktop/clipboard"
	"github.com/cxfksword/fnsync-desktop/config"
	"github.com/cxfksword/fnsync-desktop/entity"
	"github.com/cxfksword/fnsync-desktop/msg"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	AUTO_CONNECT_TIMEOUT_MILLS = 20000
)

var Listener *listener = &listener{
	code:         "*",
	connHandlers: sync.Map{},
}

type listener struct {
	code         string
	connHandlers sync.Map
}

func (l *listener) StartAccept() {
	log.Info().Msg("Start tcp connect listen...")
	serv, err := net.Listen("tcp", fmt.Sprintf(":%d", config.App.GetListenPort()))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer serv.Close()

	go StartHandshake()
	go l.startClipboardChangeWatch()
	go l.startConnectionWatch()
	for {
		// Listen for an incoming connection.
		conn, err := serv.Accept()
		log.Info().Msgf("Client accept: %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
		if err != nil {
			log.Error().Err(err).Msg("Error accepting")
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		c := entity.Device{
			Code:                l.code,
			TempotaryCodeHolder: l.code,
			Conn:                conn,
		}

		handler := msg.NewHandler(c)
		l.connHandlers.Store(uuid.NewString(), handler)
		go handler.Process()
	}
}

func (l *listener) startClipboardChangeWatch() {
	clipboard.OnChange(func(text string) {
		log.Trace().Msgf("Clipboard change: %s", text)

		if !config.App.ClipboardSync {
			return
		}

		l.connHandlers.Range(func(key, value interface{}) bool {
			v := value.(*msg.Handler)
			if v.IsAlive() {
				v.SendTextMsg(text, msg.MSG_TYPE_NEW_CLIPBOARD_DATA)
			}
			return true
		})
	})
}

// 删除已关闭的连接
func (l *listener) startConnectionWatch() {
	ticker := time.NewTicker(time.Minute)
	for {
		<-ticker.C

		l.connHandlers.Range(func(key, value interface{}) bool {
			v := value.(*msg.Handler)
			if !v.IsAlive() {
				l.connHandlers.Delete(key)
			}
			return true
		})
	}
}

func (l *listener) GetAliveDevices() []entity.Device {
	devices := []entity.Device{}
	l.connHandlers.Range(func(key, value interface{}) bool {
		v := value.(*msg.Handler)
		if v.IsAlive() {
			devices = append(devices, v.GetDevice())
		}
		return true
	})

	return devices
}

func (l *listener) RefreshCode() {
	l.code = uuid.NewString()
}

func (l *listener) GetCode() string {
	return l.code
}

func (l *listener) Terminate() {
	// TODO: close all client connection
	l.connHandlers.Range(func(key, value interface{}) bool {
		value.(*msg.Handler).NotifyStop()
		return true
	})
}
