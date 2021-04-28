package msg

import "github.com/cxfksword/fnsync-desktop/entity"

var UIMsgHandler *uiMsgHandler = &uiMsgHandler{
	uiMsgCh: make(chan entity.UIMsg, 10),
}

type uiMsgHandler struct {
	uiMsgCh chan entity.UIMsg
}

func (ui *uiMsgHandler) StartMsgLoop() {
	for {
		msg := <-ui.uiMsgCh

		msg.Execute()
	}
}

func (ui *uiMsgHandler) Send(msg entity.UIMsg) {
	ui.uiMsgCh <- msg
}
