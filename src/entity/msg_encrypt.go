package entity

import "github.com/cxfksword/fnsync-desktop/utils"

type EncryptMsg struct {
	Data string `json:"data"`
	Hash string `json:"hash"`
	Time int64  `json:"time"`
}

func NewEncryptMsg(data []byte) *EncryptMsg {
	return &EncryptMsg{
		Data: string(data),
		Hash: utils.HashSHA256(data),
		Time: utils.UnixMicro(),
	}
}
