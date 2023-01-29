package message

type Message struct {
	Key        string `json:"key" form:"key"`
	Expiration string `json:"expiration" form:"expiration"`
	Body       string `json:"body" form:"body"`
}
