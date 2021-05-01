package entity

type Msg struct {
	Text    string `json:"text,omitempty"`
	MsgType string `json:"msgtype"`
}
