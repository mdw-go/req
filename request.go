// Package req is a simple example of a single-purpose library
// whose API is based on Dave Cheney's excellent blog post
// entitled "Functional options for friendly APIs".
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
// It also incorporates a namespacing technique oft-used in
// library modules build by the developers at Smarty.
// https://github.com/smarty
package req

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

func New(method, target string, options ...option) (*http.Request, error) {
	c := newConfig(options)
	request, err := http.NewRequest(method, target, c.body)
	if err != nil {
		return request, err
	}
	request.Close = c.close
	if c.ctx != nil {
		request = request.WithContext(c.ctx)
	}
	request.URL.RawQuery = c.query.Encode()
	for key, values := range c.headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	return request, err
}

type cfg struct {
	query   url.Values
	headers http.Header
	body    io.Reader
	ctx     context.Context
	close   bool
}

func newConfig(options []option) (result cfg) {
	result.query = make(url.Values)
	result.headers = make(http.Header)
	for _, opt := range options {
		opt(&result)
	}
	return result
}

type (
	option  func(*cfg)
	options struct{}
)

var Options options

func (options) Body(r io.Reader) option            { return func(c *cfg) { c.body = r } }
func (options) Header(key, value string) option    { return func(c *cfg) { c.headers.Add(key, value) } }
func (options) Query(key, value string) option     { return func(c *cfg) { c.query.Add(key, value) } }
func (options) Context(ctx context.Context) option { return func(c *cfg) { c.ctx = ctx } }
func (options) Close(b bool) option                { return func(c *cfg) { c.close = b } }
