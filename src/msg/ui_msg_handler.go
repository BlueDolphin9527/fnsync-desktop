package msg

import "github.com/cxfksword/fnsync-desktop/entity"

var UIMsgHandler *uiMsgHandler = &uiMsgHandler{
	uiMsgCh: make(chan entity.UIMsg, 10),
}

type uiMsgHandler struct {
	uiMsgCh chan entity.UIMsg
}

func (ui *uiMsgHandler) Start() chan entity.UIMsg {
	return ui.uiMsgCh
}

func (ui *uiMsgHandler) Send(msg entity.UIMsg) {
	ui.uiMsgCh <- msg
}
