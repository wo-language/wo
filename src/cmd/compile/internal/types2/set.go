// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types2

// A Set represents a set type.
type Set struct {
	elem Type
}

// NewSet returns a new set for the given element type.
func NewSet(elem Type) *Set {
	return &Set{elem: elem}
}

// Elem returns the element type of set m.
func (m *Set) Elem() Type { return m.elem }

func (t *Set) Underlying() Type { return t }
func (t *Set) String() string   { return TypeString(t, nil) }
