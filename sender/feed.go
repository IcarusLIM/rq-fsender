package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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
	fmt.Println("Send: " + feed.Req.URL)
	jsonValue, _ := json.Marshal(feed.Req)
	resp, err := http.Post(feed.Receiver, "application/json", bytes.NewBuffer(jsonValue))
	fmt.Println(resp)
	fmt.Println(err)
}
