package rrm

import (
	"net/http"
	"sync"
)

type Handle struct {
	m map[string]struct{}
	sync.Mutex
}

func newHandle() *Handle {
	return &Handle{
		m: make(map[string]struct{}),
	}
}

func (s *Handle) Contain(id string) bool {
	s.Lock()
	defer s.Unlock()
	_, ok := s.m[id]
	return ok
}

func (s *Handle) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[string]struct{})
}

// Add return whether key exists before Add()
func (s *Handle) Add(id string) bool {
	s.Lock()
	defer s.Unlock()
	_, ok := s.m[id]
	if !ok {
		s.m[id] = struct{}{}
	}
	return ok
}

// Remove return whether key exists before Remove()
func (s *Handle) Remove(id string) bool {
	s.Lock()
	defer s.Unlock()
	_, ok := s.m[id]
	if ok {
		delete(s.m, id)
	}
	return ok
}

type FilterFunc func(id string, req *http.Request) bool

type Enforcer interface {
	Enforce(id, method, path string) bool
	EnforceReq(id string, req *http.Request) bool
	Grant(id, method, path string)
	AppendFilter(f ...FilterFunc)
	Reset()
}

type StdEnforcer struct {
	router  *router
	filters []FilterFunc
}

func (s *StdEnforcer) AppendFilter(f ...FilterFunc) {
	s.filters = append(s.filters, f...)
}

func NewStdEnforcer() Enforcer {
	return &StdEnforcer{
		router: newRouter(),
	}
}

func (s *StdEnforcer) Grant(id string, method string, path string) {
	h, _, rd := s.router.lookup(method, path)
	if rd {
		h, _, _ = s.router.lookup(method, path+"/")
	}
	if h == nil {
		h = newHandle()
		h.Add(id)
		s.router.Handle(method, path, h)
	} else {
		h.Add(id)
	}
}

// Reset clear all permission
func (s *StdEnforcer) Reset() {
	s.router = newRouter()
}

// Enforce will not call filter func chain.
func (s *StdEnforcer) Enforce(id string, method, path string) bool {
	h, _, rd := s.router.lookup(method, path)
	if rd {
		h, _, _ = s.router.lookup(method, path+"/")
	}
	if h == nil {
		return false
	}
	return h.Contain(id)
}

func (s *StdEnforcer) EnforceReq(id string, req *http.Request) bool {
	for _, flt := range s.filters {
		if flt(id, req) {
			return true
		}
	}
	return s.Enforce(id, req.Method, req.URL.Path)
}
