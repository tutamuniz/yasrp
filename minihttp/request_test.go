package minihttp

import (
	"bufio"
	"strings"
	"testing"
)

// TODO - Imrove tests
func TestParseRequest(t *testing.T) {
	testCases := []struct {
		desc     string
		httpdata string
		expected Request
	}{
		{
			desc:     "Simple GET",
			httpdata: "GET / HTTP/1.0\r\nHost: localhost:8080\r\n\r\n",
			expected: Request{Method: "GET", URI: "/", Proto: "HTTP", Version: "1.0"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tC.httpdata))
			req, err := ParseRequest(reader)
			if err != nil {
				t.Errorf("%s", err.Error())
			}

			if req.Method != tC.expected.Method || req.Proto != tC.expected.Proto || req.Version != tC.expected.Version {
				t.Errorf("Error #Invalid FirstLine# expected %v  , Got:%v", tC.expected, req)
			}

		})
	}
}
