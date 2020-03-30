package task

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
	counter  *Stats
}

// Send req
func (feed *Feed) Send() {
	fmt.Println("Send: " + feed.Req.URL)
	jsonValue, _ := json.Marshal(feed.Req)
	for i := 0; i < 3; i++ {
		resp, err := http.Post(feed.Receiver, "application/json", bytes.NewBuffer(jsonValue))
		if err == nil && resp.StatusCode == http.StatusOK {
			feed.counter.OnSuccess()
			fmt.Println(err.Error())
			break
		} else {
			feed.counter.OnFail()
		}
	}
}
