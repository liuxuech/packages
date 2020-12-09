package request

import (
	"net/http"
	"net/url"
)

type Option func(*Request)

func WithBaseUrl(baseUrl string) Option {
	return func(r *Request) {
		urlParse, _ := url.Parse(baseUrl)
		r.baseUrl = urlParse
	}
}

// 自定义transport
func (r *Request) WithTransport(transport *http.Transport) {
	r.defaultClient.Transport = transport
}

// 自定义client
func (r *Request) WithClient(client *http.Client) {
	r.defaultClient = client
}

// 设置Header，每次调用Header方法，将会清空之前的Header信息
func (r *Request) Header(header map[string]string) *Request {
	r.header = make(http.Header)
	for k, v := range header {
		r.header.Add(k, v)
	}
	return r
}
