package entity

type QRCode struct {
	Base64Data string   `json:"base64data"`
	Data       HelloMsg `json:"data"`
}
