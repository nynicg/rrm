package rrm

import (
	"net/http"
	"strings"
	"sync"
)

type restAuth struct {
	method string
	path   string
}

type idmap map[string]struct{}

func (m idmap) Contain(id string) bool {
	_, ok := m[id]
	return ok
}

type RestEnforcer struct {
	fs    []Filter
	auths map[restAuth]idmap
	mu    sync.Mutex
	fsmu  sync.Mutex
}

func (m *RestEnforcer) Enforce(id string, req *http.Request) bool {
	for _, f := range m.fs {
		abort, ok := f(id, req)
		if abort {
			return ok
		}
		if !ok {
			return false
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	path := req.URL.Path
	paths := strings.SplitAfter(path, "/")
	endwithsep := path[len(path)-1] == '/'
	if endwithsep {
		paths = paths[:len(paths)-1]
	}
	for i := 0; i < len(paths)-1; i++ {
		url := strings.Join(paths[:i+1], "")
		if m.find(id, req.Method, url) || m.find(id, req.Method, url+"*") || m.find(id, req.Method, url+":") {
			return true
		}
	}

	if m.find(id, req.Method, path) {
		return true
	}

	if endwithsep && m.find(id, req.Method, path+"*") {
		return true
	}
	return false
}

func (m *RestEnforcer) find(id string, method, path string) bool {
	auth := restAuth{
		method: method,
		path:   path,
	}
	idm, _ := m.auths[auth]
	return idm != nil && idm.Contain(id)
}

func (m *RestEnforcer) Grant(id, method, path string) {
	if method == "*" {
		for _, v := range methods {
			m.Grant(id, v, path)
		}
		return
	}
	a := restAuth{
		method: method,
		path:   path,
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	idm, ok := m.auths[a]
	if !ok {
		idm = make(map[string]struct{})
	}
	idm[id] = struct{}{}
	m.auths[a] = idm
}

func (m *RestEnforcer) AppendFilter(f ...Filter) {
	m.fsmu.Lock()
	defer m.fsmu.Unlock()
	m.fs = append(m.fs, f...)
}

func (m *RestEnforcer) Reset() {
	m.fsmu.Lock()
	m.fs = []Filter{}
	m.fsmu.Unlock()

	m.mu.Lock()
	m.auths = make(map[restAuth]idmap)
	m.mu.Unlock()
}
