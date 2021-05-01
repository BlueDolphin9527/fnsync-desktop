package entity

type AcceptMsg struct {
	PhoneId       string `json:"phoneid"`
	Key           string `json:"key"`
	Oldconnection bool   `json:"oldconnection"`
}
