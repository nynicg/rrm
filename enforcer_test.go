package rrm

import (
	"net/http"
	"net/url"
	"testing"
)

var cases map[Case]bool

type Case struct {
	id  string
	req *http.Request
}

func init() {

	cases = map[Case]bool{
		Case{"root", newMock(http.MethodGet, "/")}:                 true,
		Case{"common", newMock(http.MethodPost, "/123")}:           false,
		Case{"common", newMock(http.MethodPost, "/v1/log")}:        true,
		Case{"common", newMock(http.MethodDelete, "/123/123")}:     false,
		Case{"common", newMock(http.MethodDelete, "/v1/file/123")}: true,
		Case{"common", newMock(http.MethodDelete, "/v1/file")}:     false,
		Case{"common", newMock(http.MethodDelete, "/v1/file/")}:    false,
		Case{"icg", newMock(http.MethodOptions, "/123/123")}:       true,
	}
}

func newMock(m string, path string) *http.Request {
	return &http.Request{
		Method: m,
		URL: &url.URL{
			Path: path,
		},
	}
}

func TestRouter_GET(t *testing.T) {
	en := NewStdEnforcer()
	en.Grant("root", "GET", "/*any")
	en.Grant("common", "POST", "/v1/log")
	en.Grant("common", "*", "/v1/file/:id")
	en.AppendFilter(func(id string, req *http.Request) bool {
		return id == "icg" && req.Method == http.MethodOptions
	})
	for k, v := range cases {
		if en.EnforceReq(k.id, k.req) != v {
			t.Log("failed expected:", v, k.id, k.req.Method, k.req.URL.Path)
			t.Fail()
		}
	}
	en.Reset()
	t.Log(en.Enforce("root" ,"GET" ,"/"))
}
