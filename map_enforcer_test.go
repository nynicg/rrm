package rrm

import (
	"net/http"
	"net/url"
	"testing"
)

func newMock(m string, path string) *http.Request {
	return &http.Request{
		Method: m,
		URL: &url.URL{
			Path: path,
		},
	}
}

func AssertBool(t *testing.T, result, expect bool) {
	if result != expect {
		t.Errorf("expect:%t ,got:%t", expect, result)
	}
}

func TestRestEnforcer_Enforce(t *testing.T) {
	en := RestEnforcer{
		fs:    nil,
		auths: map[restAuth]idmap{},
	}

	en.Grant("root", "GET", "/123")
	t.Log(en.Enforce("root", newMock("GET", "/")))
	t.Log(en.Enforce("root", newMock("GET", "/123")))
	t.Log(en.Enforce("root", newMock("GET", "/123/")))

	en.Grant("rrm", "GET", "/*any")
	t.Log(en.Enforce("rrm", newMock("GET", "/")))
	t.Log(en.Enforce("rrm", newMock("GET", "/123")))
	t.Log(en.Enforce("rrm", newMock("GET", "/123/")))

	en.Grant("rrm2", "GET", "/:any")
	t.Log(en.Enforce("rrm2", newMock("GET", "/")))
	t.Log(en.Enforce("rrm2", newMock("GET", "/123")))
	t.Log(en.Enforce("rrm2", newMock("GET", "/123/")))
}

func TestRestEnforcer_trimSuffix(t *testing.T){
	t.Log(trimSuffix("/"))
	t.Log(trimSuffix("/:id"))
	t.Log(trimSuffix("/*id"))
	t.Log(trimSuffix("/*"))
	t.Log(trimSuffix("/123"))
	t.Log(trimSuffix("/123*idd"))
	t.Log(trimSuffix("/123:idd"))
	t.Log(trimSuffix("/123/"))
	t.Log(trimSuffix("/123/*idd"))
	t.Log(trimSuffix("/123/:idd"))
	t.Log(trimSuffix("/123/123"))

	t.Log(trimSuffix("/:userID/:repo"))
}
