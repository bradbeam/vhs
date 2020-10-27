package httpx

import (
	"bufio"
	"net/http"
	"strings"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestNewResponset(t *testing.T) {
	cases := []struct {
		desc        string
		b           *bufio.Reader
		r           *Response
		cID         string
		eID         string
		errContains string
	}{
		{
			desc:        "EOF",
			b:           bufio.NewReader(strings.NewReader("")),
			errContains: "EOF",
		},
		{
			desc:        "malformed",
			b:           bufio.NewReader(strings.NewReader("AICHTEETEEPEE/1.1 200 OK")),
			errContains: "malformed HTTP version",
		},
		{
			desc: "success",
			cID:  "111",
			eID:  "111",
			b:    bufio.NewReader(strings.NewReader("HTTP/1.1 204 No Content\r\n\r\n")),
			r: &Response{
				ConnectionID:  "111",
				ExchangeID:    "111",
				Status:        "204 No Content",
				StatusCode:    http.StatusNoContent,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Cookies:       []*http.Cookie{},
				Body:          "",
				ContentLength: 0,
			},
		},
		{
			desc: "success with a cookie",
			cID:  "111",
			eID:  "111",
			b:    bufio.NewReader(strings.NewReader("HTTP/1.1 204 No Content\r\nLocation: /111.html\r\nSet-Cookie: grault=foo\r\n\r\n")),
			r: &Response{
				ConnectionID: "111",
				ExchangeID:   "111",
				Status:       "204 No Content",
				StatusCode:   http.StatusNoContent,
				Proto:        "HTTP/1.1",
				ProtoMajor:   1,
				ProtoMinor:   1,
				Header: http.Header{
					"Set-Cookie": {"grault=foo"},
					"Location":   {"/111.html"},
				},
				Cookies:       []*http.Cookie{{Name: "grault", Value: "foo", Raw: "grault=foo"}},
				Body:          "",
				ContentLength: 0,
				Location:      "/111.html",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			r, err := NewResponse(c.b, c.cID, c.eID)
			if c.errContains != "" {
				assert.ErrorContains(t, err, c.errContains)
			} else {
				assert.NilError(t, err)
				assert.DeepEqual(t, r, c.r)
				now := time.Now()
				r.SetCreated(now)
				assert.Equal(t, now, r.Created)
				r.SetSessionID("111")
				assert.Equal(t, "111", r.SessionID)
			}
		})
	}
}
