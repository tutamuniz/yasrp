package minihttp

const DefaultProtocol = "HTTP"
const DefaultVersion = "1.1"
const DefaultHTTPPort = "80"
const DefaultHTTPSPort = "443"
const CRLF = "\r\n"

type StatusCode map[int]string

var SupportedStatusCode StatusCode = map[int]string{
	200: "OK",
	404: "Not found",
}

func (s StatusCode) GetText(code int) string {
	var text string
	if text, found := s[code]; found {
		return text
	}
	return text
}
