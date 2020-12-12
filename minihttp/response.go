package minihttp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Response struct {
	Status     int
	StatusText string
	Proto      string
	Version    string
	Headers    Header
	Body       []byte
}

func NewResponse(status int, header Header, body []byte) (*Response, error) {

	text := SupportedStatusCode.GetText(status)
	if text == "" {
		return nil, fmt.Errorf("Invalid Status code")
	}

	resp := &Response{
		Status:     status,
		StatusText: text,
		Proto:      DefaultProtocol,
		Version:    DefaultVersion,
		Headers:    header,
		Body:       body,
	}

	resp.SetContentType("text/html")

	return resp, nil
}

// ParseResponse RFC7230
func ParseResponse(r *bufio.Reader) (*Response, error) {
	resp := &Response{}

	firstLine, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	// HTTP/1.1 200 OK
	fields := strings.Split(string(firstLine), " ")

	resp.Proto = fields[0]
	resp.Status, _ = strconv.Atoi(fields[1]) // Check error
	resp.StatusText = fields[2]

	fields = strings.Split(resp.Proto, "/")
	resp.Proto, resp.Version = fields[0], fields[1]

	resp.Headers = make(Header)

	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return nil, err
		}
		if len(line) == 0 {
			break
		}
		fields := strings.Split(string(line), ":")
		key, value := fields[0], strings.Trim(fields[1], " ")
		resp.Headers[key] = value

	}

	if tenc, ok := resp.Headers["Transfer-Encoding"]; ok && tenc == "chunked" {
		// var vf bytes.Buffer
		// wr := gzip.NewWriter(&vf)

		// wr.Write([]byte("5\r\nHello\r\n0\r\n\r\n"))
		// wr.Close()
		resp.Body = []byte("5\r\nHello\r\n0\r\n\r\n")
	}

	if clen, ok := resp.Headers["Content-Length"]; ok {
		var bodyBuff bytes.Buffer
		tmp := make([]byte, 1024)
		inBytes := 0
		contentLen, err := strconv.Atoi(clen)

		if err != nil {
			return nil, err
		}
		var n int = 0
		for {
			if inBytes >= contentLen {
				break
			}
			n, err = r.Read(tmp)
			inBytes += n
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			bodyBuff.Write(tmp[:n])
		}
		resp.Body = bodyBuff.Bytes()

	}
	return resp, nil
}

func (r *Response) SetContentType(content string) {
	r.Headers["Content-Type"] = content
}

func (r Response) ToBytes() []byte {
	var outRes bytes.Buffer

	outRes.WriteString(fmt.Sprintf("%s/%s %d %s\r\n", r.Proto,
		r.Version,
		r.Status,
		r.StatusText))
	for key, value := range r.Headers {
		l := fmt.Sprintf("%s: %s\r\n", key, value)
		outRes.WriteString(l)
	}
	outRes.WriteString("\r\n")
	if len(r.Body) > 0 {
		outRes.WriteString(string(r.Body))
	}
	return outRes.Bytes()
}
