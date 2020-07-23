package entity

type Req struct {
	url string
}

func NewReq(url string) *Req {
	return &Req{url: url}
}
