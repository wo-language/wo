// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sets

import "iter"

// All returns an iterator over key-value pairs from m.
// The iteration order is not specified and is not guaranteed
// to be the same from one call to the next.
func All[Map ~map[K]V, K comparable, V any](m Map) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Elems returns an iterator over elements in s.
// The iteration order is not specified and is not guaranteed
// to be the same from one call to the next.
func Elems[Map ~map[K]V, K comparable, V any](m Map) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

// Insert adds the key-value pairs from seq to m.
// If a key in seq already exists in m, its value will be overwritten.
// TODO(bran) After set keyword is implemented, this must be compiled differently; and order.go insertFunc must follow
func Insert[Map ~map[K]V, K comparable, V any](m Map, seq iter.Seq2[K, V]) {
	for k, v := range seq {
		m[k] = v
	}
}

// Collect collects key-value pairs from seq into a new map
// and returns it.
func Collect[K comparable, V any](seq iter.Seq2[K, V]) map[K]V {
	m := make(map[K]V)
	Insert(m, seq)
	return m
}
