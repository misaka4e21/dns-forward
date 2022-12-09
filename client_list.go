// Copyright 2018 The Rind Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net"
	"sync"
)

type clientList struct {
	sync.RWMutex
	data map[string][]net.UDPAddr
}

func (b *clientList) get(key string) ([]net.UDPAddr, bool) {
	b.RLock()
	val, ok := b.data[key]
	b.RUnlock()
	return val, ok
}

func (b *clientList) set(key string, addr net.UDPAddr) {
	b.Lock()
	b.data[key] = append(b.data[key], addr)
	b.Unlock()
}

func (b *clientList) remove(key string) bool {
	b.Lock()
	_, ok := b.data[key]
	delete(b.data, key)
	b.Unlock()
	return ok
}
