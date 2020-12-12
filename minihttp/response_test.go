package minihttp

import (
	"bufio"
	"strings"
	"testing"
)

// TODO - Improve tests

var header Header = Header{"Host": "www.google.com"}

func TestNewResponse(t *testing.T) {
	r, err := NewResponse(404, header, []byte("Hello World!"))
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if r.Headers["Host"] != "www.google.com" {
		t.Error("Header error")
	}

	if string(r.Body) != "Hello World!" {
		t.Error("Body error")
	}

}
func TestParseResponse(t *testing.T) {
	testCases := []struct {
		desc     string
		httpdata string
		expected Response
	}{
		{
			desc:     "Simple Response",
			httpdata: "HTTP/1.1 200 OK\r\nHost: localhost:8080\r\nContent-Length:10\r\nContent-type: text/html\r\n\r\nABCDEFGHIJ",
			expected: Response{Proto: "HTTP", Version: "1.1", Status: 200, StatusText: "OK", Body: []byte("ABCDEFGHIJ")},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tC.httpdata))
			res, err := ParseResponse(reader)
			if err != nil {
				t.Errorf("%s", err.Error())
			}

			if res.Status != tC.expected.Status ||
				res.Proto != tC.expected.Proto ||
				res.Version != tC.expected.Version ||
				res.StatusText != tC.expected.StatusText {
				t.Errorf("Error #Invalid FirstLine# expected %v  , Got:%v", tC.expected, res)
			}

		})
	}
}
