package req_test

import (
	"net/http/httputil"
	"strings"
	"testing"

	"req"
)

func Test(t *testing.T) {
	request, err := req.New(
		req.Options.POST(strings.NewReader("This is the request body.")),
		req.Options.Target("https://www.smarty.com/team"),
		req.Options.Header("x-hello", "world!"),
	)
	if err != nil {
		t.Fatal(err)
	}

	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\n" + string(dump))
}
