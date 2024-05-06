// This file is modified from the interner found in Prometheus. The original
// license follows:
//
// Copyright 2019 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Inspired / copied / modified from https://gitlab.com/cznic/strutil/blob/master/strutil.go,
// which is MIT licensed, so:
//
// Copyright (c) 2014 The strutil Authors. All rights reserved.

package discovery

import (
	"sync"
	"sync/atomic"
)

// Iterner is a string interner.
type Interner interface {
	// Intern will intern an input string, returning the interned string as
	// a result.
	Intern(string) string

	// Release removes an interned string from interner.
	Release(string)
}

func NewInterner() Interner {
	return &pool{
		pool: map[string]*entry{},
	}
}

type pool struct {
	mtx  sync.RWMutex
	pool map[string]*entry
}

type entry struct {
	// Keep all 64bit atomically accessed variables at the top of this struct.
	// See https://golang.org/pkg/sync/atomic/#pkg-note-BUG for more info.
	refs int64

	s string
}

func (p *pool) Intern(s string) string {
	if s == "" {
		return ""
	}

	p.mtx.RLock()
	interned, ok := p.pool[s]
	p.mtx.RUnlock()
	if ok {
		atomic.AddInt64(&interned.refs, 1)
		return interned.s
	}
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if interned, ok := p.pool[s]; ok {
		atomic.AddInt64(&interned.refs, 1)
		return interned.s
	}

	p.pool[s] = &entry{
		s:    s,
		refs: 1,
	}
	return s
}

func (p *pool) Release(s string) {
	p.mtx.RLock()
	interned, ok := p.pool[s]
	p.mtx.RUnlock()

	if !ok {
		// TODO(rfratto): no reference releases metric?
		return
	}

	refs := atomic.AddInt64(&interned.refs, -1)
	if refs > 0 {
		return
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()
	if atomic.LoadInt64(&interned.refs) != 0 {
		return
	}
	delete(p.pool, s)
}
