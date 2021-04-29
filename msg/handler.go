package msg

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/cxfksword/fnsync-desktop/app"
	"github.com/cxfksword/fnsync-desktop/clipboard"
	"github.com/cxfksword/fnsync-desktop/config"
	"github.com/cxfksword/fnsync-desktop/entity"
	"github.com/cxfksword/fnsync-desktop/utils"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	device            entity.Device
	encryptionManager *EncryptionManager
	streamReader      *bufio.Reader
	sendMsgCh         chan entity.Msg
	stopCh            chan bool
}

func NewHandler(client entity.Device) *Handler {
	return &Handler{
		device:            client,
		streamReader:      bufio.NewReader(client.Conn),
		encryptionManager: NewEncryptionManager([]byte(client.Code)),
		sendMsgCh:         make(chan entity.Msg, 10),
		stopCh:            make(chan bool, 1),
	}
}

func (h *Handler) Process() {
	defer h.close()

	h.setReadDeadline(time.Now().Add(5 * time.Second))
	data, err := h.readPackage()
	if err != nil {
		return
	}

	log.Debug().Msgf("Accept client device: %s", string(data))
	acceptMsg, err := h.readPhoneId(data)
	if err != nil {
		return
	}

	h.device.Id = acceptMsg.PhoneId
	h.device.Oldconnection = acceptMsg.Oldconnection
	if acceptMsg.Oldconnection {
		h.device.Code = config.App.Devices[h.device.Id].Code
		h.encryptionManager = NewEncryptionManager([]byte(h.device.Code))
	} else {
		if acceptMsg.Key != "" {
			h.device.TempotaryCodeHolder += acceptMsg.Key
			h.device.Code = h.device.TempotaryCodeHolder
			h.encryptionManager = NewEncryptionManager([]byte(h.device.Code))
		}
	}

	data, err = h.readPackage()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	authMsg, err := h.readAuthenticate(data)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	if h.device.Id != authMsg.PhoneId {
		log.Info().Msgf("Phoneid not match. device_id:%s phoneid:%s", h.device.Id, authMsg.PhoneId)
		return
	}

	h.replyBack()
	h.device.Name = authMsg.PhoneName

	h.setReadDeadline(time.Now().Add(30 * time.Second))
	data, err = h.readPackage()
	if err != nil {
		return
	}
	acceptedMsg, err := h.readMsg(data)
	if err != nil {
		return
	}
	if acceptedMsg.MsgType != MSG_TYPE_CONNECTION_ACCEPTED {
		log.Warn().Msgf("Illegal Phone. [%s](%s)", h.device.Name, h.device.Id)
		return
	}

	go h.UpdateDevice()
	go h.StartSendEventLoop()
	go h.StartRecieveEventLoop()
	h.waitToStop()
}

func (h *Handler) UpdateDevice() {
	// save connected device
	h.device.IsAlive = true
	h.device.LastIp = h.device.Conn.RemoteAddr().(*net.TCPAddr).IP.String()
	config.App.SaveDevice(h.device)

	UIMsgHandler.Send(entity.UINotifyMsg{
		Title:   app.Name,
		Message: fmt.Sprintf("%s已连接", h.device.Name),
	})
	UIMsgHandler.Send(entity.UIUpdateStatusMsg{})
}

func (h *Handler) StartRecieveEventLoop() {
	defer func() { log.Info().Msgf("Quit RecieveEventLoop... [%s](%s)", h.device.Name, h.device.Id) }()
	log.Info().Msgf("Start RecieveEventLoop... [%s](%s)", h.device.Name, h.device.Id)

	// set no timeout
	h.setReadDeadline(time.Time{})
	for {
		data, err := h.readPackage()
		if err != nil {
			return
		}

		json, err := h.encryptionManager.Decrypt(data)
		log.Info().Msgf("Recieve msg: %s [%s](%s)", string(json), h.device.Name, h.device.Id)
		if err != nil {
			log.Error().Err(err).Msg("")
			continue
		}

		encryptMsg := entity.EncryptMsg{}
		if err := utils.FromJSON(json, &encryptMsg); err != nil {
			log.Error().Err(err).Msg("")
			continue
		}
		var msg entity.Msg
		if err := utils.FromJSON([]byte(encryptMsg.Data), &msg); err != nil {
			log.Error().Err(err).Msg("")
			continue
		}

		switch msg.MsgType {
		case MSG_TYPE_TEXT_CAST:
			if !config.App.ClipboardSync {
				continue
			}
			if err := clipboard.Set(msg.Text); err != nil {
				log.Error().Err(err).Msg("")
			}
		case MSG_TYPE_DISCONNECT_BY_PEER:
			log.Info().Msgf("Quit by client. [%s](%s)", h.device.Name, h.device.Id)
			h.NotifyStop()
			return
		case MSG_TYPE_LOCK_SCREEN:
		case MSG_TYPE_HELLO:
			// 用于终端判断连接是否还生效，一段时间没响应，终端会关闭连接重连
			h.helloBack()
		}
	}
}

func (h *Handler) StartSendEventLoop() {
	defer func() { log.Info().Msgf("Quit SendEventLoop... [%s](%s)", h.device.Name, h.device.Id) }()
	log.Info().Msgf("Start SendEventLoop... [%s](%s)", h.device.Name, h.device.Id)

	for {
		msg := <-h.sendMsgCh

		msgEncrypt, _ := h.encryptionManager.Encrypt(utils.ToJSON(msg))
		h.writePackage(msgEncrypt)
	}
}

func (h *Handler) readPhoneId(data []byte) (*entity.AcceptMsg, error) {
	accept := entity.AcceptMsg{}

	err := utils.FromJSON(data, &accept)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &accept, nil
}

func (h *Handler) readAuthenticate(data []byte) (*entity.AuthMsg, error) {
	json, err := h.encryptionManager.Decrypt(data)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Accept authenticate data: %s", string(json))
	encryptMsg := entity.EncryptMsg{}
	err = utils.FromJSON(json, &encryptMsg)
	if err != nil {
		return nil, err
	}

	authMsg := entity.AuthMsg{}
	err = utils.FromJSON([]byte(encryptMsg.Data), &authMsg)
	if err != nil {
		return nil, err
	}

	return &authMsg, nil
}

func (h *Handler) readMsg(data []byte) (*entity.Msg, error) {
	json, err := h.encryptionManager.Decrypt(data)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Read encrypt msg: %s [%s](%s)", string(json), h.device.Name, h.device.Id)
	encryptMsg := entity.EncryptMsg{}
	err = utils.FromJSON(json, &encryptMsg)
	if err != nil {
		return nil, err
	}

	msg := entity.Msg{}
	err = utils.FromJSON([]byte(encryptMsg.Data), &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (h *Handler) replyBack() {
	reply := entity.PearReplyMsg{
		PeerId: config.App.MachineId,
	}

	replyJson, _ := h.encryptionManager.Encrypt(utils.ToJSON(reply))
	h.writePackage(replyJson)
}

func (h *Handler) helloBack() {
	reply := entity.Msg{
		MsgType: MSG_TYPE_NONCE,
	}

	replyJson, _ := h.encryptionManager.Encrypt(utils.ToJSON(reply))
	h.writePackage(replyJson)
}

func (h *Handler) readPackageLength() ([]byte, error) {
	return h.readStream(4)
}

func (h *Handler) readPackage() ([]byte, error) {
	lenBytes, err := h.readPackageLength()
	if err != nil {
		return nil, err
	}
	packageLength := utils.BytesToInt(lenBytes)
	return h.readStream(packageLength)
}

func (h *Handler) writePackage(data []byte) {
	// append raw data length
	length := len(data)
	lengthBytes := utils.IntToBytes(length)
	_, err := h.device.Conn.Write(lengthBytes)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	_, err = h.device.Conn.Write(data)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
}

func (h *Handler) setReadDeadline(deadline time.Time) {
	err := h.device.Conn.SetReadDeadline(deadline)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
}

func (h *Handler) readStream(size int) ([]byte, error) {
	buf := make([]byte, size) // big buffer
	n, err := io.ReadFull(h.streamReader, buf)
	if err != nil {
		if err == io.EOF {
			log.Trace().Msg("READ io.EOF")
		} else {
			log.Error().Err(err).Msg("")
		}
		return nil, err
	}
	if n <= 0 {
		return nil, fmt.Errorf("disconnect_by_peer")
	}
	return buf, nil
}

func (h *Handler) sendTerminateMsg() {
	terminate := entity.Msg{
		MsgType: MSG_TYPE_DISCONNECT_BY_PEER,
	}

	terminateJson, _ := h.encryptionManager.Encrypt(utils.ToJSON(terminate))
	h.writePackage(terminateJson)
}

func (h *Handler) SendTextMsg(text string, msgType string) {
	textMsg := entity.Msg{
		Text:    text,
		MsgType: msgType,
	}

	h.sendMsgCh <- textMsg
}

func (h *Handler) IsAlive() bool {
	return h.device.IsAlive
}

func (h *Handler) GetDevice() entity.Device {
	return h.device
}

func (h *Handler) GetCode() string {
	return h.device.Code
}

func (h *Handler) close() {
	log.Info().Msgf("Disconnect_by_peer: %s -> %s", h.device.Conn.RemoteAddr(), h.device.Conn.LocalAddr())

	h.sendTerminateMsg()
	h.device.Conn.Close()
	h.device.IsAlive = false
	UIMsgHandler.Send(entity.UIUpdateStatusMsg{})
}

func (h *Handler) waitToStop() {
	<-h.stopCh
}

func (h *Handler) NotifyStop() {
	h.stopCh <- true
}
