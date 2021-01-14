package rrm

import (
	"net/http"
)

var methods []string

func init() {
	methods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}
}

type Filter func(id string, req *http.Request) (abort bool, ok bool)

type Enforcer interface {
	Enforce(id string, req *http.Request) bool
	Grant(id, method, path string)
	AppendFilter(f ...Filter)
	Reset()
}
