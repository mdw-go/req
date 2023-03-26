// Package req is a simple example of a single-purpose library
// whose API is based on Dave Cheney's excellent blog post
// entitled "Functional options for friendly APIs".
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
// It also incorporates a namespacing technique oft-used in
// library modules build by the developers at Smarty.
// https://github.com/smarty
package req

import (
	"io"
	"net/http"
)

func New(options ...option) (*http.Request, error) {
	c := newConfig(options)
	request, err := http.NewRequest(c.method, c.target, c.body)
	if err != nil {
		return request, err
	}
	for key, values := range c.headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	return request, err
}

type config struct {
	method  string
	target  string
	headers http.Header
	body    io.Reader
}

func newConfig(options []option) (result config) {
	result.headers = make(http.Header)
	for _, opt := range options {
		opt(&result)
	}
	return result
}

type option func(*config)

var Options singleton

type singleton struct{}

func (singleton) Method(method string) option {
	return func(c *config) { c.method = method }
}
func (singleton) GET() option {
	return Options.Method(http.MethodGet)
}
func (singleton) POST(body io.Reader) option {
	return func(c *config) {
		Options.Method(http.MethodPost)(c)
		c.body = body
	}
}
func (singleton) PUT(body io.Reader) option {
	return func(c *config) {
		Options.Method(http.MethodPut)(c)
		c.body = body
	}
}
func (singleton) DELETE() option {
	return Options.Method(http.MethodDelete)
}
func (singleton) Target(t string) option {
	return func(config *config) { config.target = t }
}
func (singleton) Header(key, value string) option {
	return func(c *config) { c.headers.Set(key, value) }
}
