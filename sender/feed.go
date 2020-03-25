package sender

import "fmt"

// Request TODO
type Request struct {
	URL      string
	Method   string
	Meta     map[string]interface{}
	Headers  map[string]string
	Cookies  map[string]string
	Body     []byte
	Encoding string
}

// Feed TODO
type Feed struct {
	Receiver string
	Req      Request
}

func (feed *Feed) send() {
	fmt.Print(feed.Req.URL)
}
