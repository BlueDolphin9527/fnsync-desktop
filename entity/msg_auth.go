package entity

type AuthMsg struct {
	PhoneId    string `json:"phoneid"`
	WaitAccept bool   `json:"wait_accept"`
	PhoneName  string `json:"phonename"`
	Key        string `json:"key"`
}
