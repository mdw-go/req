package req_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/mdw-go/req"
	"github.com/mdw-go/req/internal/should"
)

func TestSuite(t *testing.T) {
	should.Run(&Suite{T: should.New(t)}, should.Options.UnitTests())
}

type Suite struct{ *should.T }

func (this *Suite) TestAwkwardHTTPAPI() {
	request, err := http.NewRequest("", "", nil)
	this.So(err, should.BeNil)
	this.So(request.Method, should.Equal, http.MethodGet)
	this.So(request.URL.String(), should.BeEmpty)
	this.So(request.Header, should.BeEmpty)
	this.So(request.Body, should.BeNil)
}
func (this *Suite) TestSimple() {
	request, err := req.New(http.MethodHead, "/")
	this.So(err, should.BeNil)
	this.So(request.Method, should.Equal, http.MethodHead)
	this.So(request.URL.String(), should.Equal, "/")
	this.So(request.Header, should.BeEmpty)
	this.So(request.Body, should.BeNil)
}
func (this *Suite) Test() {
	request, err := req.New(http.MethodPost, "https://www.smarty.com/team",
		req.Options.Query("hello", "world"),
		req.Options.Body(strings.NewReader("This is the request body.")),
		req.Options.Header("x-hello", "world!"),
		req.Options.Context(context.WithValue(context.Background(), "context", "testing")),
		req.Options.Close(true),
	)
	this.So(err, should.BeNil)
	this.So(request.Method, should.Equal, http.MethodPost)
	this.So(request.URL.String(), should.Equal, "https://www.smarty.com/team?hello=world")
	this.So(readString(request.Body), should.Equal, "This is the request body.")
	this.So(request.Header.Get("X-Hello"), should.Equal, "world!")
	this.So(request.Context().Value("context"), should.Equal, "testing")
	this.So(request.Close, should.BeTrue)

	this.Print("\n" + string(justValue(httputil.DumpRequest(request, true))))
}

func justValue[T any](t T, _ error) T { return t }
func readString(body io.ReadCloser) string {
	defer func() { _ = body.Close() }()
	return string(justValue(io.ReadAll(body)))
}
