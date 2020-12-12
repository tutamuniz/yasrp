package minihttp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	Proto   string
	URI     string
	Version string
	Headers Header
	Body    []byte
}

//func NewRequest(method, )

// ParseRequest RFC7230
func ParseRequest(r *bufio.Reader) (*Request, error) {
	log.Printf("Parsing request\n")
	req := &Request{}
	firstLine, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}

	fields := strings.Split(string(firstLine), " ")
	req.Method, req.URI, req.Proto = fields[0], fields[1], fields[2] // Check Valid method

	if !isValidMethod(req.Method) {
		return nil, ErrNotImplementedMethod
	}
	fields = strings.Split(req.Proto, "/")
	req.Proto, req.Version = fields[0], fields[1]

	req.Headers = make(Header)
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
		req.Headers[key] = value

	}

	if clen, ok := req.Headers["Content-Length"]; ok {
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
			fmt.Println(inBytes, n)
		}
		req.Body = bodyBuff.Bytes()

	}
	return req, nil
}

// ToBytes - convert request content into byte slice
func (r Request) ToBytes() []byte {
	var outReq bytes.Buffer

	outReq.WriteString(fmt.Sprintf("%s %s %s/%s\r\n", r.Method,
		r.URI,
		r.Proto,
		r.Version))
	for key, value := range r.Headers {
		l := fmt.Sprintf("%s: %s\r\n", key, value)
		outReq.WriteString(l)
	}
	outReq.WriteString("\r\n")
	if len(r.Body) > 0 {
		outReq.WriteString(string(r.Body))
	}
	return outReq.Bytes()
}
