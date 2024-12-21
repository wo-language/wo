// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/abi"
	"internal/goarch"
	"unsafe"
)

//func setaccess1_fast32(t *settype, h *hset, key uint32) unsafe.Pointer {
//	if raceenabled && h != nil {
//		callerpc := getcallerpc()
//		racereadpc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(setaccess1_fast32))
//	}
//	if h == nil || h.count == 0 {
//		return unsafe.Pointer(&zeroVal[0])
//	}
//	if h.flags&hashWriting != 0 {
//		fatal("concurrent set read and set write")
//	}
//	var b *bset
//	if h.B == 0 {
//		// One-bucket table. No need to hash.
//		b = (*bset)(h.buckets)
//	} else {
//		hash := t.Hasher(noescape(unsafe.Pointer(&key)), uintptr(h.hash0))
//		m := bucketMask(h.B)
//		b = (*bset)(add(h.buckets, (hash&m)*uintptr(t.BucketSize)))
//		if c := h.oldbuckets; c != nil {
//			if !h.sameSizeGrow() {
//				// There used to be half as many buckets; mask down one more power of two.
//				m >>= 1
//			}
//			oldb := (*bset)(add(c, (hash&m)*uintptr(t.BucketSize)))
//			if !evacuatedSet(oldb) {
//				b = oldb
//			}
//		}
//	}
//	for ; b != nil; b = b.overflow(t) {
//		for i, k := uintptr(0), b.keys(); i < abi.MapBucketCount; i, k = i+1, add(k, 4) {
//			if *(*uint32)(k) == key && !isEmpty(b.tophash[i]) {
//				return add(unsafe.Pointer(b), dataOffset+abi.MapBucketCount*4+i*uintptr(t.ValueSize))
//			}
//		}
//	}
//	return unsafe.Pointer(&zeroVal[0])
//}

// setaccess2_fast32 should be an internal detail,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - github.com/ugorji/go/codec
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname setaccess2_fast32
func setaccess2_fast32(t *settype, h *hset, key uint32) bool {
	if raceenabled && h != nil {
		callerpc := getcallerpc()
		racereadpc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(setaccess2_fast32))
	}
	if h == nil || h.count == 0 {
		return false
	}
	if h.flags&hashWriting != 0 {
		fatal("concurrent set read and set write")
	}
	var b *bset
	if h.B == 0 {
		// One-bucket table. No need to hash.
		b = (*bset)(h.buckets)
	} else {
		hash := t.Hasher(noescape(unsafe.Pointer(&key)), uintptr(h.hash0))
		m := bucketMask(h.B)
		b = (*bset)(add(h.buckets, (hash&m)*uintptr(t.BucketSize)))
		if c := h.oldbuckets; c != nil {
			if !h.sameSizeGrow() {
				// There used to be half as many buckets; mask down one more power of two.
				m >>= 1
			}
			oldb := (*bset)(add(c, (hash&m)*uintptr(t.BucketSize)))
			if !evacuatedSet(oldb) {
				b = oldb
			}
		}
	}
	for ; b != nil; b = b.overflow(t) {
		for i, k := uintptr(0), b.keys(); i < abi.MapBucketCount; i, k = i+1, add(k, 4) {
			if *(*uint32)(k) == key && !isEmpty(b.tophash[i]) {
				return true
			}
		}
	}
	return false
}

// setassign_fast32 should be an internal detail.
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname setassign_fast32
//func setassign_fast32(t *settype, h *hset, key uint32) unsafe.Pointer {
//	if h == nil {
//		panic(plainError("assignment to entry in nil set"))
//	}
//	if raceenabled {
//		callerpc := getcallerpc()
//		racewritepc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(setassign_fast32))
//	}
//	if h.flags&hashWriting != 0 {
//		fatal("concurrent set writes")
//	}
//	hash := t.Hasher(noescape(unsafe.Pointer(&key)), uintptr(h.hash0))
//
//	// Set hashWriting after calling t.hasher for consistency with setassign.
//	h.flags ^= hashWriting
//
//	if h.buckets == nil {
//		h.buckets = newobject(t.Bucket) // newarray(t.bucket, 1)
//	}
//
//again:
//	bucket := hash & bucketMask(h.B)
//	if h.growing() {
//		growWorkSet_fast32(t, h, bucket)
//	}
//	b := (*bset)(add(h.buckets, bucket*uintptr(t.BucketSize)))
//
//	var insertb *bset
//	var inserti uintptr
//	var insertk unsafe.Pointer
//
//bucketloop:
//	for {
//		for i := uintptr(0); i < abi.MapBucketCount; i++ {
//			if isEmpty(b.tophash[i]) {
//				if insertb == nil {
//					inserti = i
//					insertb = b
//				}
//				if b.tophash[i] == emptyRest {
//					break bucketloop
//				}
//				continue
//			}
//			k := *((*uint32)(add(unsafe.Pointer(b), dataOffset+i*4)))
//			if k != key {
//				continue
//			}
//			inserti = i
//			insertb = b
//			goto done
//		}
//		ovf := b.overflow(t)
//		if ovf == nil {
//			break
//		}
//		b = ovf
//	}
//
//	// Did not find mapping for key. Allocate new cell & add entry.
//
//	// If we hit the max load factor or we have too many overflow buckets,
//	// and we're not already in the middle of growing, start growing.
//	if !h.growing() && (overLoadFactor(h.count+1, h.B) || tooManyOverflowBuckets(h.noverflow, h.B)) {
//		hashGrowSet(t, h)
//		goto again // Growing the table invalidates everything, so try again
//	}
//
//	if insertb == nil {
//		// The current bucket and all the overflow buckets connected to it are full, allocate a new one.
//		insertb = h.newoverflow(t, b)
//		inserti = 0 // not necessary, but avoids needlessly spilling inserti
//	}
//	insertb.tophash[inserti&(abi.MapBucketCount-1)] = tophash(hash) // mask inserti to avoid bounds checks
//
//	insertk = add(unsafe.Pointer(insertb), dataOffset+inserti*4)
//	// store new key at insert position
//	*(*uint32)(insertk) = key
//
//	h.count++
//
//done:
//	elem := add(unsafe.Pointer(insertb), dataOffset+abi.MapBucketCount*4+inserti*uintptr(t.ValueSize))
//	if h.flags&hashWriting == 0 {
//		fatal("concurrent set writes")
//	}
//	h.flags &^= hashWriting
//	return elem
//}
//
//// setassign_fast32ptr should be an internal detail.
//// Do not access it using linkname.
////
//// Do not remove or change the type signature.
//// See go.dev/issue/67401.
////
////wo:linkname setassign_fast32ptr
//func setassign_fast32ptr(t *settype, h *hset, key unsafe.Pointer) unsafe.Pointer {
//	if h == nil {
//		panic(plainError("assignment to entry in nil set"))
//	}
//	if raceenabled {
//		callerpc := getcallerpc()
//		racewritepc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(setassign_fast32))
//	}
//	if h.flags&hashWriting != 0 {
//		fatal("concurrent set writes")
//	}
//	hash := t.Hasher(noescape(unsafe.Pointer(&key)), uintptr(h.hash0))
//
//	// Set hashWriting after calling t.hasher for consistency with setassign.
//	h.flags ^= hashWriting
//
//	if h.buckets == nil {
//		h.buckets = newobject(t.Bucket) // newarray(t.bucket, 1)
//	}
//
//again:
//	bucket := hash & bucketMask(h.B)
//	if h.growing() {
//		growWorkSet_fast32(t, h, bucket)
//	}
//	b := (*bset)(add(h.buckets, bucket*uintptr(t.BucketSize)))
//
//	var insertb *bset
//	var inserti uintptr
//	var insertk unsafe.Pointer
//
//bucketloop:
//	for {
//		for i := uintptr(0); i < abi.MapBucketCount; i++ {
//			if isEmpty(b.tophash[i]) {
//				if insertb == nil {
//					inserti = i
//					insertb = b
//				}
//				if b.tophash[i] == emptyRest {
//					break bucketloop
//				}
//				continue
//			}
//			k := *((*unsafe.Pointer)(add(unsafe.Pointer(b), dataOffset+i*4)))
//			if k != key {
//				continue
//			}
//			inserti = i
//			insertb = b
//			goto done
//		}
//		ovf := b.overflow(t)
//		if ovf == nil {
//			break
//		}
//		b = ovf
//	}
//
//	// Did not find mapping for key. Allocate new cell & add entry.
//
//	// If we hit the max load factor or we have too many overflow buckets,
//	// and we're not already in the middle of growing, start growing.
//	if !h.growing() && (overLoadFactor(h.count+1, h.B) || tooManyOverflowBuckets(h.noverflow, h.B)) {
//		hashGrowSet(t, h)
//		goto again // Growing the table invalidates everything, so try again
//	}
//
//	if insertb == nil {
//		// The current bucket and all the overflow buckets connected to it are full, allocate a new one.
//		insertb = h.newoverflow(t, b)
//		inserti = 0 // not necessary, but avoids needlessly spilling inserti
//	}
//	insertb.tophash[inserti&(abi.MapBucketCount-1)] = tophash(hash) // mask inserti to avoid bounds checks
//
//	insertk = add(unsafe.Pointer(insertb), dataOffset+inserti*4)
//	// store new key at insert position
//	*(*unsafe.Pointer)(insertk) = key
//
//	h.count++
//
//done:
//	elem := add(unsafe.Pointer(insertb), dataOffset+abi.MapBucketCount*4+inserti*uintptr(t.ValueSize))
//	if h.flags&hashWriting == 0 {
//		fatal("concurrent set writes")
//	}
//	h.flags &^= hashWriting
//	return elem
//}

func setdelete_fast32(t *settype, h *hset, key uint32) {
	if raceenabled && h != nil {
		callerpc := getcallerpc()
		racewritepc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(setdelete_fast32))
	}
	if h == nil || h.count == 0 {
		return
	}
	if h.flags&hashWriting != 0 {
		fatal("concurrent set writes")
	}

	hash := t.Hasher(noescape(unsafe.Pointer(&key)), uintptr(h.hash0))

	// Set hashWriting after calling t.hasher for consistency with setdelete
	h.flags ^= hashWriting

	bucket := hash & bucketMask(h.B)
	if h.growing() {
		growWorkSet_fast32(t, h, bucket)
	}
	b := (*bset)(add(h.buckets, bucket*uintptr(t.BucketSize)))
	bOrig := b
search:
	for ; b != nil; b = b.overflow(t) {
		for i, k := uintptr(0), b.keys(); i < abi.MapBucketCount; i, k = i+1, add(k, 4) {
			if key != *(*uint32)(k) || isEmpty(b.tophash[i]) {
				continue
			}
			// Only clear key if there are pointers in it.
			// This can only happen if pointers are 32 bit
			// wide as 64 bit pointers do not fit into a 32 bit key.
			if goarch.PtrSize == 4 && t.Elem.Pointers() {
				// The key must be a pointer as we checked pointers are
				// 32 bits wide and the key is 32 bits wide also.
				*(*unsafe.Pointer)(k) = nil
			}
			//e := add(unsafe.Pointer(b), dataOffset+abi.MapBucketCount*4+i*uintptr(t.ValueSize))
			//if t.Elem.Pointers() {
			//	memclrHasPointers(e, t.Elem.Size_)
			//} else {
			//	memclrNoHeapPointers(e, t.Elem.Size_)
			//}
			b.tophash[i] = emptyOne
			// If the bucket now ends in a bunch of emptyOne states,
			// change those to emptyRest states.
			if i == abi.MapBucketCount-1 {
				if b.overflow(t) != nil && b.overflow(t).tophash[0] != emptyRest {
					goto notLast
				}
			} else {
				if b.tophash[i+1] != emptyRest {
					goto notLast
				}
			}
			for {
				b.tophash[i] = emptyRest
				if i == 0 {
					if b == bOrig {
						break // beginning of initial bucket, we're done.
					}
					// Find previous bucket, continue at its last entry.
					c := b
					for b = bOrig; b.overflow(t) != c; b = b.overflow(t) {
					}
					i = abi.MapBucketCount - 1
				} else {
					i--
				}
				if b.tophash[i] != emptyOne {
					break
				}
			}
		notLast:
			h.count--
			// Reset the hash seed to make it more difficult for attackers to
			// repeatedly trigger hash collisions. See issue 25237.
			if h.count == 0 {
				h.hash0 = uint32(rand())
			}
			break search
		}
	}

	if h.flags&hashWriting == 0 {
		fatal("concurrent set writes")
	}
	h.flags &^= hashWriting
}

func growWorkSet_fast32(t *settype, h *hset, bucket uintptr) {
	// make sure we evacuate the oldbucket corresponding
	// to the bucket we're about to use
	evacuateSet_fast32(t, h, bucket&h.oldbucketmask())

	// evacuate one more oldbucket to make progress on growing
	if h.growing() {
		evacuateSet_fast32(t, h, h.nevacuate)
	}
}

func evacuateSet_fast32(t *settype, h *hset, oldbucket uintptr) {
	b := (*bset)(add(h.oldbuckets, oldbucket*uintptr(t.BucketSize)))
	newbit := h.noldbuckets()
	if !evacuatedSet(b) {
		// TODO: reuse overflow buckets instead of using new ones, if there
		// is no iterator using the old buckets.  (If !oldIterator.)

		// xy contains the x and y (low and high) evacuation destinations.
		var xy [2]evacDstSet
		x := &xy[0]
		x.b = (*bset)(add(h.buckets, oldbucket*uintptr(t.BucketSize)))
		x.k = add(unsafe.Pointer(x.b), dataOffset)

		if !h.sameSizeGrow() {
			// Only calculate y pointers if we're growing bigger.
			// Otherwise GC can see bad pointers.
			y := &xy[1]
			y.b = (*bset)(add(h.buckets, (oldbucket+newbit)*uintptr(t.BucketSize)))
			y.k = add(unsafe.Pointer(y.b), dataOffset)
		}

		for ; b != nil; b = b.overflow(t) {
			k := add(unsafe.Pointer(b), dataOffset)
			for i := 0; i < abi.MapBucketCount; i, k = i+1, add(k, 4) {
				top := b.tophash[i]
				if isEmpty(top) {
					b.tophash[i] = evacuatedEmpty
					continue
				}
				if top < minTopHash {
					throw("bad set state")
				}
				var useY uint8
				if !h.sameSizeGrow() {
					// Compute hash to make our evacuation decision (whether we need
					// to send this key/elem to bucket x or bucket y).
					hash := t.Hasher(k, uintptr(h.hash0))
					if hash&newbit != 0 {
						useY = 1
					}
				}

				b.tophash[i] = evacuatedX + useY // evacuatedX + 1 == evacuatedY, enforced in makeset
				dst := &xy[useY]                 // evacuation destination

				if dst.i == abi.MapBucketCount {
					dst.b = h.newoverflow(t, dst.b)
					dst.i = 0
					dst.k = add(unsafe.Pointer(dst.b), dataOffset)
				}
				dst.b.tophash[dst.i&(abi.MapBucketCount-1)] = top // mask dst.i as an optimization, to avoid a bounds check

				// Copy key.
				if goarch.PtrSize == 4 && t.Elem.Pointers() && writeBarrier.enabled {
					// Write with a write barrier.
					*(*unsafe.Pointer)(dst.k) = *(*unsafe.Pointer)(k)
				} else {
					*(*uint32)(dst.k) = *(*uint32)(k)
				}

				dst.i++
				// These updates might push these pointers past the end of the
				// key or elem arrays.  That's ok, as we have the overflow pointer
				// at the end of the bucket to protect against pointing past the
				// end of the bucket.
				dst.k = add(dst.k, 4)
			}
		}
		// Unlink the overflow buckets & clear key/elem to help GC.
		if h.flags&oldIterator == 0 && t.Bucket.Pointers() {
			b := add(h.oldbuckets, oldbucket*uintptr(t.BucketSize))
			// Preserve b.tophash because the evacuation
			// state is maintained there.
			ptr := add(b, dataOffset)
			n := uintptr(t.BucketSize) - dataOffset
			memclrHasPointers(ptr, n)
		}
	}

	if oldbucket == h.nevacuate {
		advanceEvacuationMarkSet(h, t, newbit)
	}
}
