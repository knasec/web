package web

type Responce struct {
	Status int    `json:"status"`
	Body   []byte `json:"body"`
}
