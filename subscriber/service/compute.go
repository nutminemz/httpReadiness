package service

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// FillProto :: prevent input without protocal
func FillProto(url string) (out string) {
	if strings.Contains(url, "http") {
		return url
	}
	return ("http://" + url)
}

// FetchHTTP :: get target URL return fail if timeout
func FetchHTTP(msg *redis.Message) string {
	url := FillProto(msg.Payload)
	req, err := http.NewRequest("GET", url, nil)
	resp, err := newNetClient().Do(req)
	if err != nil {
		if resp == nil {
			log.Println(url, " unreachable")
			return "fail"
		}
	}
	return "success"
}

// SetResult :: increment redis value for counting result
func SetResult(c context.Context, r *redis.Client, re string) {
	if re == "success" {
		r.Incr(c, "success")
	} else {
		_ = r.Incr(c, "fail")
	}
}

var (
	once      sync.Once
	netClient *http.Client
)

// newNetClient :: makeing new net/http connection
func newNetClient() *http.Client {
	once.Do(func() {
		var netTransport = &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 15 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   8 * time.Second,
			ExpectContinueTimeout: 8 * time.Second,
			ResponseHeaderTimeout: 8 * time.Second,
			ReadBufferSize:        8,
		}
		netClient = &http.Client{
			Timeout:   time.Second * 15,
			Transport: netTransport,
		}
	})

	return netClient
}
