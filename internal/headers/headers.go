package headers

import (
	"bytes"
	"fmt"
	"strings"
)

func isToken(s string) bool {
	if len(s) == 0 {
		return false
	}

	for i := 0; i < len(s); i++ {
		c := s[i]

		switch {
		case 'A' <= c && c <= 'Z':
		case 'a' <= c && c <= 'z':
		case '0' <= c && c <= '9':
		case strings.ContainsRune("!#$%&'*+-.^_`|~", rune(c)):
		default:
			return false
		}
	}

	return true
}

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

var rn = []byte("\r\n")

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed header")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name), string(value), nil
}

func (h *Headers) Get(name string) (string, bool) {
	v, ok := h.headers[strings.ToLower(name)]
	return v, ok
}

func (h *Headers) Replace(name, value string) {
	name = strings.ToLower(name)
	h.headers[name] = value
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)

	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s,%s", v, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range h.headers {
		cb(n, v)
	}
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}
		if !isToken(name) {
			return 0, false, fmt.Errorf("malformed header name")
		}

		read += idx + len(rn)

		h.Set(name, value)

	}

	return read, done, nil
}
