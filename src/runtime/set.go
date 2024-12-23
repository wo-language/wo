// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copied from map.go
// The comments are not unreliable
// wo:linkname is not real, it just lacks the linked counterparts

package runtime

// This file contains the implementation of Wo's set type.
//
// A set is just a hash table. The data is arranged
// into an array of buckets. Each bucket contains up to
// 8 elements. The low-order bits of the hash are
// used to select a bucket. Each bucket contains a few
// high-order bits of each hash to distinguish the entries
// within a single bucket.
//
// If more than 8 keysSet hash to a bucket, we chain on
// extra buckets.
//
// When the hashtable grows, we allocate a new array
// of buckets twice as big. Buckets are incrementally
// copied from the old bucket array to the new bucket array.
//
// Set iterators walk through the array of buckets and
// return the keysSet in walk order (bucket #, then overflow
// chain order, then bucket index).  To maintain iteration
// semantics, we never move keysSet within their bucket (if
// we did, keysSet might be returned 0 or 2 times).  When
// growing the table, iterators remain iterating through the
// old table and must check the new table if the bucket
// they are iterating through has been moved ("evacuatedSet")
// to the new table.

// Picking loadFactor: too large and we have lots of overflow
// buckets, too small and we waste a lot of space. I wrote
// a simple program to check some stats for different loads:
// (64-bit, 8 byte keysSet and elems)
//  loadFactor    %overflow  bytes/entry     hitprobe    missprobe
//

//
// %overflow   = percentage of buckets which have an overflow bucket
// bytes/entry = overhead bytes used per key/elem pair
// hitprobe    = # of entries to check when looking up a present key
// missprobe   = # of entries to check when looking up an absent key
//
// Keep in mind this data is for maximally loaded tables, i.e. just
// before the table grows. Typical tables will be somewhat less loaded.

import (
	"internal/abi"
	"internal/goarch"
	"internal/runtime/atomic"
	"runtime/internal/math"
	"unsafe"
)

const (
//bucketCntBits  = abi.MapBucketCountBits
//loadFactorDen  = 2
//loadFactorNum  = loadFactorDen * abi.MapBucketCount * 13 / 16
//emptyRest      = 0
//emptyOne       = 1
//evacuatedX     = 2
//evacuatedY     = 3
//evacuatedEmpty = 4
//minTopHash     = 5

// iterator       = 1
// oldIterator    = 2
// hashWriting    = 4
// sameSizeGrow   = 8
// noCheck		 = 1<<(8*goarch.PtrSize) - 1
)

// A header for a Wo set.
type hset struct {
	// #wo Note: the format of the hset is also encoded in cmd/compile/internal/reflectdata/reflect.go in SetType.
	// #wo Make sure this stays in sync with the compiler's definition.
	// "SetType returns a type interchangeable with runtime.hset."
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuatedSet)

	extra *setextra // optional fields
}

// setextra holds fields that are not present on all maps.
type setextra struct {
	// If both key and elem do not contain pointers and are inline, then we mark bucket
	// type as containing no pointers. This avoids scanning such maps.
	// However, bset.overflow is a pointer. In order to keep overflow buckets
	// alive, we store pointers to all overflow buckets in hset.extra.overflow and hset.extra.oldoverflow.
	// overflow and oldoverflow are only used if key and elem do not contain pointers.
	// overflow contains overflow buckets for hset.buckets.
	// oldoverflow contains overflow buckets for hset.oldbuckets.
	// The indirection allows to store a pointer to the slice in hsetiter.
	overflow    *[]*bset
	oldoverflow *[]*bset

	// nextOverflow holds a pointer to a free overflow bucket.
	nextOverflow *bset
}

// A bucket for a Wo set.
type bset struct {
	// tophash generally contains the top byte of the hash value
	// for each key in this bucket. If tophash[0] < minTopHash,
	// tophash[0] is a bucket evacuation state instead.
	tophash [abi.MapBucketCount]uint8 // [8]uint8
	// Followed by bucketCnt keys and then bucketCnt elems.
	// NOTE: packing all the keys together and then all the elems together makes the
	// code a bit more complicated than alternating key/elem/key/elem/... but it allows
	// us to eliminate padding which would be needed for, e.g., map[int64]int8.
	// Followed by an overflow pointer.
}

// A hash iteration structure.
// If you modify hsetiter, also change cmd/compile/internal/reflectdata/reflect.go
// and reflect/value.go to match the layout of this structure.
type hsetiter struct {
	key         unsafe.Pointer // Must be in first position.  Write nil to indicate iteration end (see cmd/compile/internal/walk/range.go).
	t           *settype
	h           *hset
	buckets     unsafe.Pointer // bucket ptr at hash_iter initialization time
	bptr        *bset          // current bucket
	overflow    *[]*bset       // keeps overflow buckets of hset.buckets alive
	oldoverflow *[]*bset       // keeps overflow buckets of hset.oldbuckets alive
	startBucket uintptr        // bucket iteration started at
	offset      uint8          // intra-bucket offset to start from during iteration (should be big enough to hold bucketCnt-1)
	wrapped     bool           // already wrapped around from end of bucket array to beginning
	B           uint8
	i           uint8
	bucket      uintptr
	checkBucket uintptr
}

func evacuatedSet(b *bset) bool {
	h := b.tophash[0]
	return h > emptyOne && h < minTopHash
}

func (b *bset) overflow(t *settype) *bset {
	return *(**bset)(add(unsafe.Pointer(b), uintptr(t.BucketSize)-goarch.PtrSize))
}

func (b *bset) setoverflow(t *settype, ovf *bset) {
	*(**bset)(add(unsafe.Pointer(b), uintptr(t.BucketSize)-goarch.PtrSize)) = ovf
}

func (b *bset) keys() unsafe.Pointer {
	return add(unsafe.Pointer(b), dataOffset)
}

// incrnoverflow increments h.noverflow.
// noverflow counts the number of overflow buckets.
// This is used to trigger same-size map growth.
// See also tooManyOverflowBuckets.
// To keep hset small, noverflow is a uint16.
// When there are few buckets, noverflow is an exact count.
// When there are many buckets, noverflow is an approximate count.
func (h *hset) incrnoverflow() {
	// We trigger same-size map growth if there are
	// as many overflow buckets as buckets.
	// We need to be able to count to 1<<h.B.
	if h.B < 16 {
		h.noverflow++
		return
	}
	// Increment with probability 1/(1<<(h.B-15)).
	// When we reach 1<<15 - 1, we will have approximately
	// as many overflow buckets as buckets.
	mask := uint32(1)<<(h.B-15) - 1
	// Example: if h.B == 18, then mask == 7,
	// and rand() & 7 == 0 with probability 1/8.
	if uint32(rand())&mask == 0 {
		h.noverflow++
	}
}

func (h *hset) newoverflow(t *settype, b *bset) *bset {
	var ovf *bset
	if h.extra != nil && h.extra.nextOverflow != nil {
		// We have preallocated overflow buckets available.
		// See makeSetBucketArray for more details.
		ovf = h.extra.nextOverflow
		if ovf.overflow(t) == nil {
			// We're not at the end of the preallocated overflow buckets. Bump the pointer.
			h.extra.nextOverflow = (*bset)(add(unsafe.Pointer(ovf), uintptr(t.BucketSize)))
		} else {
			// This is the last preallocated overflow bucket.
			// Reset the overflow pointer on this bucket,
			// which was set to a non-nil sentinel value.
			ovf.setoverflow(t, nil)
			h.extra.nextOverflow = nil
		}
	} else {
		ovf = (*bset)(newobject(t.Bucket))
	}
	h.incrnoverflow()
	if !t.Bucket.Pointers() {
		h.createOverflow()
		*h.extra.overflow = append(*h.extra.overflow, ovf)
	}
	b.setoverflow(t, ovf)
	return ovf
}

func (h *hset) createOverflow() {
	if h.extra == nil {
		h.extra = new(setextra)
	}
	if h.extra.overflow == nil {
		h.extra.overflow = new([]*bset)
	}
}

func makeset64(t *settype, hint int64, h *hset) *hset {
	if int64(int(hint)) != hint {
		hint = 0
	}
	return makeset(t, int(hint), h)
}

// makeset_small implements Go map creation for make(map[k]v) and
// make(map[k]v, hint) when hint is known to be at most bucketCnt
// at compile time and the map needs to be allocated on the heap.
//
// makeset_small should be an internal detail.
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname makeset_small
func makeset_small() *hset {
	h := new(hset)
	h.hash0 = uint32(rand())
	return h
}

// makeset implements a set creation for make(set[e], hint).
// If the compiler has determined that the map or the first bucket
// can be created on the stack, h and/or bucket may be non-nil.
// If h != nil, the map can be created directly in h.
// If h.buckets != nil, bucket pointed to can be used as the first bucket.
//
// makeset should be an internal detail.
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname makeset
func makeset(t *settype, hint int, h *hset) *hset {
	mem, overflow := math.MulUintptr(uintptr(hint), t.Bucket.Size_)
	if overflow || mem > maxAlloc {
		hint = 0
	}

	// initialize Hmap
	if h == nil {
		h = new(hset)
	}
	h.hash0 = uint32(rand())

	// Find the size parameter B which will hold the requested # of elements.
	// For hint < 0 overLoadFactor returns false since hint < bucketCnt.
	B := uint8(0)
	for overLoadFactor(hint, B) {
		B++
	}
	h.B = B

	// allocate initial hash table
	// if B == 0, the buckets field is allocated lazily later (in setassign)
	// If hint is large zeroing this memory could take a while.
	if h.B != 0 {
		var nextOverflow *bset
		h.buckets, nextOverflow = makeSetBucketArray(t, h.B, nil)
		if nextOverflow != nil {
			h.extra = new(setextra)
			h.extra.nextOverflow = nextOverflow
		}
	}

	return h
}

// makeSetBucketArray initializes a backing array for map buckets.
// 1<<b is the minimum number of buckets to allocate.
// dirtyalloc should either be nil or a bucket array previously
// allocated by makeSetBucketArray with the same t and b parameters.
// If dirtyalloc is nil a new backing array will be alloced and
// otherwise dirtyalloc will be cleared and reused as backing array.
func makeSetBucketArray(t *settype, b uint8, dirtyalloc unsafe.Pointer) (buckets unsafe.Pointer, nextOverflow *bset) {
	base := bucketShift(b)
	nbuckets := base
	// For small b, overflow buckets are unlikely.
	// Avoid the overhead of the calculation.
	if b >= 4 {
		// Add on the estimated number of overflow buckets
		// required to insert the median number of elements
		// used with this value of b.
		nbuckets += bucketShift(b - 4)
		sz := t.Bucket.Size_ * nbuckets
		up := roundupsize(sz, !t.Bucket.Pointers())
		if up != sz {
			nbuckets = up / t.Bucket.Size_
		}
	}

	if dirtyalloc == nil {
		buckets = newarray(t.Bucket, int(nbuckets))
	} else {
		// dirtyalloc was previously generated by
		// the above newarray(t.Bucket, int(nbuckets))
		// but may not be empty.
		buckets = dirtyalloc
		size := t.Bucket.Size_ * nbuckets
		if t.Bucket.Pointers() {
			memclrHasPointers(buckets, size)
		} else {
			memclrNoHeapPointers(buckets, size)
		}
	}

	if base != nbuckets {
		// We preallocated some overflow buckets.
		// To keep the overhead of tracking these overflow buckets to a minimum,
		// we use the convention that if a preallocated overflow bucket's overflow
		// pointer is nil, then there are more available by bumping the pointer.
		// We need a safe non-nil pointer for the last overflow bucket; just use buckets.
		nextOverflow = (*bset)(add(buckets, base*uintptr(t.BucketSize)))
		last := (*bset)(add(buckets, (nbuckets-1)*uintptr(t.BucketSize)))
		last.setoverflow(t, (*bset)(buckets))
	}
	return buckets, nextOverflow
}

// setaccess1 returns a pointer to h[key], which does not make sense
// for sets, so it and its predecessors are removed.
//
// v := &h[key]
// var v *ValueType = h[key]
//
// Never returns nil, instead
// it will return a reference to the zero object for the elem type if
// the key is not in the map.
// NOTE: The returned pointer may keep the whole map live, so don't
// hold onto it for very long.
func setaccess1(t *settype, h *hset, key unsafe.Pointer) unsafe.Pointer {
	return nil
	//	if raceenabled && h != nil {
	//		callerpc := getcallerpc()
	//		pc := abi.FuncPCABIInternal(setaccess1)
	//		racereadpc(unsafe.Pointer(h), callerpc, pc)
	//		raceReadObjectPC(t.Key, key, callerpc, pc)
	//	}
	//	if msanenabled && h != nil {
	//		msanread(key, t.Key.Size_)
	//	}
	//	if asanenabled && h != nil {
	//		asanread(key, t.Key.Size_)
	//	}
	//	if h == nil || h.count == 0 {
	//		if err := setKeyError(t, key); err != nil {
	//			panic(err) // see issue 23734
	//		}
	//		return unsafe.Pointer(&zeroVal[0])
	//	}
	//	if h.flags&hashWriting != 0 {
	//		fatal("concurrent set read and map write")
	//	}
	//	// shift to the respective bucket
	//	hash := t.Hasher(key, uintptr(h.hash0))
	//	m := bucketMask(h.B)
	//	b := (*bset)(add(h.buckets, (hash&m)*uintptr(t.BucketSize)))
	//	if c := h.oldbuckets; c != nil {
	//		if !h.sameSizeGrow() {
	//			// There used to be half as many buckets; mask down one more power of two.
	//			m >>= 1
	//		}
	//		oldb := (*bset)(add(c, (hash&m)*uintptr(t.BucketSize)))
	//		if !evacuatedSet(oldb) {
	//			b = oldb
	//		}
	//	}
	//	top := tophash(hash)
	//bucketloop:
	//	for ; b != nil; b = b.overflow(t) {
	//		for i := uintptr(0); i < abi.MapBucketCount; i++ {
	//			if b.tophash[i] != top {
	//				if b.tophash[i] == emptyRest {
	//					break bucketloop
	//				}
	//				continue
	//			}
	//			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
	//			if t.IndirectKey() {
	//				k = *((*unsafe.Pointer)(k))
	//			}
	//			if t.Key.Equal(key, k) {
	//				e := add(unsafe.Pointer(b), dataOffset+abi.MapBucketCount*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
	//				if t.IndirectElem() {
	//					e = *((*unsafe.Pointer)(e))
	//				}
	//				return e
	//			}
	//		}
	//	}
	//	return unsafe.Pointer(&zeroVal[0])
}

// setaccess2 should be an internal detail.
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
// ok := h[key]
//
// returns whether it contains the key
//
//wo:linkname setaccess2
func setaccess2(t *settype, h *hset, key unsafe.Pointer) bool {
	if raceenabled && h != nil {
		callerpc := getcallerpc()
		pc := abi.FuncPCABIInternal(setaccess2)
		racereadpc(unsafe.Pointer(h), callerpc, pc)
		raceReadObjectPC(t.Elem, key, callerpc, pc)
	}
	if msanenabled && h != nil {
		msanread(key, t.Elem.Size_)
	}
	if asanenabled && h != nil {
		asanread(key, t.Elem.Size_)
	}
	if h == nil || h.count == 0 {
		if err := setKeyError(t, key); err != nil {
			panic(err) // see issue 23734
		}
		return false
	}
	if h.flags&hashWriting != 0 {
		fatal("concurrent set read and set write")
	}
	hash := t.Hasher(key, uintptr(h.hash0))
	m := bucketMask(h.B)
	b := (*bset)(add(h.buckets, (hash&m)*uintptr(t.BucketSize)))
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
	top := tophash(hash)
bucketloop:
	for ; b != nil; b = b.overflow(t) {
		for i := uintptr(0); i < abi.MapBucketCount; i++ {
			if b.tophash[i] != top {
				if b.tophash[i] == emptyRest {
					break bucketloop
				}
				continue
			}
			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
			if t.IndirectKey() {
				k = *((*unsafe.Pointer)(k))
			}
			if t.Elem.Equal(key, k) {
				//e := add(unsafe.Pointer(b), dataOffset+abi.MapBucketCount*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
				//if t.IndirectElem() {
				//	e = *((*unsafe.Pointer)(e))
				//}
				return true
			}
		}
	}
	return false
}

// returns the "key" (usually called elem in sets). Used by set iterator.
func setaccessK(t *settype, h *hset, key unsafe.Pointer) unsafe.Pointer {
	if h == nil || h.count == 0 {
		return nil
	}
	hash := t.Hasher(key, uintptr(h.hash0))
	m := bucketMask(h.B)
	b := (*bset)(add(h.buckets, (hash&m)*uintptr(t.BucketSize)))
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
	top := tophash(hash)
bucketloop:
	for ; b != nil; b = b.overflow(t) {
		for i := uintptr(0); i < abi.MapBucketCount; i++ {
			if b.tophash[i] != top {
				if b.tophash[i] == emptyRest {
					break bucketloop
				}
				continue
			}
			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
			if t.IndirectKey() {
				k = *((*unsafe.Pointer)(k))
			}
			if t.Elem.Equal(key, k) {
				//e := add(unsafe.Pointer(b), dataOffset+abi.MapBucketCount*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
				//if t.IndirectElem() {
				//	e = *((*unsafe.Pointer)(e))
				//}
				return k
			}
		}
	}
	return nil
}

//func setaccess1_fat(t *settype, h *hset, key, zero unsafe.Pointer) unsafe.Pointer {
//	e := setaccess1(t, h, key)
//	if e == unsafe.Pointer(&zeroVal[0]) {
//		return zero
//	}
//	return e
//}

//func setaccess2_fat(t *settype, h *hset, key, zero unsafe.Pointer) (unsafe.Pointer, bool) {
//	e := setaccess1(t, h, key)
//	if e == unsafe.Pointer(&zeroVal[0]) {
//		return zero, false
//	}
//	return e, true
//}

// Like setaccess, but allocates a slot for the key if it is not present in the map.
//
// m[k] =
//
// setassign should be an internal detail.
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname setassign
func setassign(t *settype, h *hset, key unsafe.Pointer) unsafe.Pointer { // TODO
	if h == nil {
		panic(plainError("assignment to entry in nil set"))
	}
	if raceenabled {
		callerpc := getcallerpc()
		pc := abi.FuncPCABIInternal(setassign)
		racewritepc(unsafe.Pointer(h), callerpc, pc)
		raceReadObjectPC(t.Elem, key, callerpc, pc)
	}
	if msanenabled {
		msanread(key, t.Elem.Size_)
	}
	if asanenabled {
		asanread(key, t.Elem.Size_)
	}
	if h.flags&hashWriting != 0 {
		fatal("concurrent set writes")
	}
	hash := t.Hasher(key, uintptr(h.hash0))

	// Set hashWriting after calling t.hasher, since t.hasher may panic,
	// in which case we have not actually done a write.
	h.flags ^= hashWriting

	if h.buckets == nil {
		h.buckets = newobject(t.Bucket) // newarray(t.Bucket, 1)
	}

again:
	bucket := hash & bucketMask(h.B)
	if h.growing() {
		growWorkSet(t, h, bucket)
	}
	b := (*bset)(add(h.buckets, bucket*uintptr(t.BucketSize)))
	top := tophash(hash)

	var inserti *uint8
	var insertk unsafe.Pointer
	var elem unsafe.Pointer
bucketloop:
	for {
		for i := uintptr(0); i < abi.MapBucketCount; i++ {
			if b.tophash[i] != top {
				if isEmpty(b.tophash[i]) && inserti == nil {
					inserti = &b.tophash[i]
					insertk = add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
					elem = add(unsafe.Pointer(b), dataOffset+abi.MapBucketCount*uintptr(t.KeySize))
				}
				if b.tophash[i] == emptyRest {
					break bucketloop
				}
				continue
			}
			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
			if t.IndirectKey() {
				k = *((*unsafe.Pointer)(k))
			}
			if !t.Elem.Equal(key, k) {
				continue
			}
			// already have a mapping for key. Update it.
			if t.NeedKeyUpdate() {
				typedmemmove(t.Elem, k, key)
			}
			elem = add(unsafe.Pointer(b), dataOffset+abi.MapBucketCount*uintptr(t.KeySize))
			goto done
		}
		ovf := b.overflow(t)
		if ovf == nil {
			break
		}
		b = ovf
	}

	// Did not find mapping for key. Allocate new cell & add entry.

	// If we hit the max load factor or we have too many overflow buckets,
	// and we're not already in the middle of growing, start growing.
	if !h.growing() && (overLoadFactor(h.count+1, h.B) || tooManyOverflowBuckets(h.noverflow, h.B)) {
		hashGrowSet(t, h)
		goto again // Growing the table invalidates everything, so try again
	}

	if inserti == nil {
		// The current bucket and all the overflow buckets connected to it are full, allocate a new one.
		newb := h.newoverflow(t, b)
		inserti = &newb.tophash[0]
		insertk = add(unsafe.Pointer(newb), dataOffset)
		elem = add(insertk, abi.MapBucketCount*uintptr(t.KeySize))
	}

	// store new key/elem at insert position
	if t.IndirectKey() {
		kmem := newobject(t.Elem)
		*(*unsafe.Pointer)(insertk) = kmem
		insertk = kmem
	}
	//if t.IndirectElem() {
	//	vmem := newobject(t.Elem)
	//	*(*unsafe.Pointer)(elem) = vmem
	//}
	typedmemmove(t.Elem, insertk, key)
	*inserti = top
	h.count++

done:
	if h.flags&hashWriting == 0 {
		fatal("concurrent set writes")
	}
	h.flags &^= hashWriting
	//if t.IndirectElem() {
	//	elem = *((*unsafe.Pointer)(elem))
	//}
	return elem
}

// setdelete should be an internal detail,
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname setdelete
func setdelete(t *settype, h *hset, key unsafe.Pointer) {
	if raceenabled && h != nil {
		callerpc := getcallerpc()
		pc := abi.FuncPCABIInternal(setdelete)
		racewritepc(unsafe.Pointer(h), callerpc, pc)
		raceReadObjectPC(t.Elem, key, callerpc, pc)
	}
	if msanenabled && h != nil {
		msanread(key, t.Elem.Size_)
	}
	if asanenabled && h != nil {
		asanread(key, t.Elem.Size_)
	}
	if h == nil || h.count == 0 {
		if err := setKeyError(t, key); err != nil {
			panic(err) // see issue 23734
		}
		return
	}
	if h.flags&hashWriting != 0 {
		fatal("concurrent set writes")
	}

	hash := t.Hasher(key, uintptr(h.hash0))

	// Set hashWriting after calling t.hasher, since t.hasher may panic,
	// in which case we have not actually done a write (delete).
	h.flags ^= hashWriting

	bucket := hash & bucketMask(h.B)
	if h.growing() {
		growWorkSet(t, h, bucket)
	}
	b := (*bset)(add(h.buckets, bucket*uintptr(t.BucketSize)))
	bOrig := b
	top := tophash(hash)
search:
	for ; b != nil; b = b.overflow(t) {
		for i := uintptr(0); i < abi.MapBucketCount; i++ {
			if b.tophash[i] != top {
				if b.tophash[i] == emptyRest {
					break search
				}
				continue
			}
			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
			k2 := k
			if t.IndirectKey() {
				k2 = *((*unsafe.Pointer)(k2))
			}
			if !t.Elem.Equal(key, k2) {
				continue
			}
			// Only clear key if there are pointers in it.
			if t.IndirectKey() {
				*(*unsafe.Pointer)(k) = nil
			} else if t.Elem.Pointers() {
				memclrHasPointers(k, t.Elem.Size_)
			}
			//e := add(unsafe.Pointer(b), dataOffset+abi.MapBucketCount*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
			//if t.IndirectElem() {
			//	*(*unsafe.Pointer)(e) = nil
			//} else if t.Elem.Pointers() {
			//	memclrHasPointers(e, t.Elem.Size_)
			//} else {
			//	memclrNoHeapPointers(e, t.Elem.Size_)
			//}
			b.tophash[i] = emptyOne
			// If the bucket now ends in a bunch of emptyOne states,
			// change those to emptyRest states.
			// It would be nice to make this a separate function, but
			// for loops are not currently inlineable.
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

// setiterinit initializes the hsetiter struct used for ranging over maps.
// The hsetiter struct pointed to by 'it' is allocated on the stack
// by the compilers order pass or on the heap by reflect_setiterinit.
// Both need to have zeroed hsetiter since the struct contains pointers.
//
// setiterinit should be an internal detail,
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname setiterinit
func setiterinit(t *settype, h *hset, it *hsetiter) {
	if raceenabled && h != nil {
		callerpc := getcallerpc()
		racereadpc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(setiterinit))
	}

	it.t = t
	if h == nil || h.count == 0 {
		return
	}

	if unsafe.Sizeof(hsetiter{})/goarch.PtrSize != 12 {
		throw("hash_iter size incorrect") // see cmd/compile/internal/reflectdata/reflect.go
	}
	it.h = h

	// grab snapshot of bucket state
	it.B = h.B
	it.buckets = h.buckets
	if !t.Bucket.Pointers() {
		// Allocate the current slice and remember pointers to both current and old.
		// This preserves all relevant overflow buckets alive even if
		// the table grows and/or overflow buckets are added to the table
		// while we are iterating.
		h.createOverflow()
		it.overflow = h.extra.overflow
		it.oldoverflow = h.extra.oldoverflow
	}

	// decide where to start
	r := uintptr(rand())
	it.startBucket = r & bucketMask(h.B)
	it.offset = uint8(r >> h.B & (abi.MapBucketCount - 1))

	// iterator state
	it.bucket = it.startBucket

	// Remember we have an iterator.
	// Can run concurrently with another setiterinit().
	if old := h.flags; old&(iterator|oldIterator) != iterator|oldIterator {
		atomic.Or8(&h.flags, iterator|oldIterator)
	}

	setiternext(it)
}

// setiternext should be an internal detail,
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname setiternext
func setiternext(it *hsetiter) {
	h := it.h
	if raceenabled {
		callerpc := getcallerpc()
		racereadpc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(setiternext))
	}
	if h.flags&hashWriting != 0 {
		fatal("concurrent set iteration and set write")
	}
	t := it.t
	bucket := it.bucket
	b := it.bptr
	i := it.i
	checkBucket := it.checkBucket

next:
	if b == nil {
		if bucket == it.startBucket && it.wrapped {
			// end of iteration
			it.key = nil
			//it.elem = nil
			return
		}
		if h.growing() && it.B == h.B {
			// Iterator was started in the middle of a grow, and the grow isn't done yet.
			// If the bucket we're looking at hasn't been filled in yet (i.e. the old
			// bucket hasn't been evacuatedSet) then we need to iterate through the old
			// bucket and only return the ones that will be migrated to this bucket.
			oldbucket := bucket & it.h.oldbucketmask()
			b = (*bset)(add(h.oldbuckets, oldbucket*uintptr(t.BucketSize)))
			if !evacuatedSet(b) {
				checkBucket = bucket
			} else {
				b = (*bset)(add(it.buckets, bucket*uintptr(t.BucketSize)))
				checkBucket = noCheck
			}
		} else {
			b = (*bset)(add(it.buckets, bucket*uintptr(t.BucketSize)))
			checkBucket = noCheck
		}
		bucket++
		if bucket == bucketShift(it.B) {
			bucket = 0
			it.wrapped = true
		}
		i = 0
	}
	for ; i < abi.MapBucketCount; i++ {
		offi := (i + it.offset) & (abi.MapBucketCount - 1)
		if isEmpty(b.tophash[offi]) || b.tophash[offi] == evacuatedEmpty {
			// TODO: emptyRest is hard to use here, as we start iterating
			// in the middle of a bucket. It's feasible, just tricky.
			continue
		}
		k := add(unsafe.Pointer(b), dataOffset+uintptr(offi)*uintptr(t.KeySize))
		if t.IndirectKey() {
			k = *((*unsafe.Pointer)(k))
		}
		// #wo  e := add(unsafe.Pointer(b), dataOffset+abi.MapBucketCount*uintptr(t.KeySize)+uintptr(offi)*uintptr(t.ValueSize))
		if checkBucket != noCheck && !h.sameSizeGrow() {
			// Special case: iterator was started during a grow to a larger size
			// and the grow is not done yet. We're working on a bucket whose
			// oldbucket has not been evacuatedSet yet. Or at least, it wasn't
			// evacuatedSet when we started the bucket. So we're iterating
			// through the oldbucket, skipping any keysSet that will go
			// to the other new bucket (each oldbucket expands to two
			// buckets during a grow).
			if t.ReflexiveKey() || t.Elem.Equal(k, k) {
				// If the item in the oldbucket is not destined for
				// the current new bucket in the iteration, skip it.
				hash := t.Hasher(k, uintptr(h.hash0))
				if hash&bucketMask(it.B) != checkBucket {
					continue
				}
			} else {
				// Hash isn't repeatable if k != k (NaNs).  We need a
				// repeatable and randomish choice of which direction
				// to send NaNs during evacuation. We'll use the low
				// bit of tophash to decide which way NaNs go.
				// NOTE: this case is why we need two evacuateSet tophash
				// valuesSet, evacuatedX and evacuatedY, that differ in
				// their low bit.
				if checkBucket>>(it.B-1) != uintptr(b.tophash[offi]&1) {
					continue
				}
			}
		}
		if (b.tophash[offi] != evacuatedX && b.tophash[offi] != evacuatedY) ||
			!(t.ReflexiveKey() || t.Elem.Equal(k, k)) {
			// This is the golden data, we can return it.
			// OR
			// key!=key, so the entry can't be deleted or updated, so we can just return it.
			// That's lucky for us because when key!=key we can't look it up successfully.
			it.key = k
		} else {
			// The hash table has grown since the iterator was started.
			// The golden data for this key is now somewhere else.
			// Check the current hash table for the data.
			// This code handles the case where the key
			// has been deleted, updated, or deleted and reinserted.
			// NOTE: we need to regrab the key as it has potentially been
			// updated to an equal() but not identical key (e.g. +0.0 vs -0.0).
			rk := setaccessK(t, h, k)
			if rk == nil {
				continue // key has been deleted
			}
			it.key = rk
		}
		it.bucket = bucket
		if it.bptr != b { // avoid unnecessary write barrier; see issue 14921
			it.bptr = b
		}
		it.i = i + 1
		it.checkBucket = checkBucket
		return
	}
	b = b.overflow(t)
	i = 0
	goto next
}

// setclear deletes all keysSet from a map.
// It is called by the compiler.
//
// setclear should be an internal detail,
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname setclear
func setclear(t *settype, h *hset) {
	if raceenabled && h != nil {
		callerpc := getcallerpc()
		pc := abi.FuncPCABIInternal(setclear)
		racewritepc(unsafe.Pointer(h), callerpc, pc)
	}

	if h == nil || h.count == 0 {
		return
	}

	if h.flags&hashWriting != 0 {
		fatal("concurrent set writes")
	}

	h.flags ^= hashWriting

	// Mark buckets empty, so existing iterators can be terminated, see issue #59411.
	markBucketsEmpty := func(bucket unsafe.Pointer, mask uintptr) {
		for i := uintptr(0); i <= mask; i++ {
			b := (*bset)(add(bucket, i*uintptr(t.BucketSize)))
			for ; b != nil; b = b.overflow(t) {
				for i := uintptr(0); i < abi.MapBucketCount; i++ {
					b.tophash[i] = emptyRest
				}
			}
		}
	}
	markBucketsEmpty(h.buckets, bucketMask(h.B))
	if oldBuckets := h.oldbuckets; oldBuckets != nil {
		markBucketsEmpty(oldBuckets, h.oldbucketmask())
	}

	h.flags &^= sameSizeGrow
	h.oldbuckets = nil
	h.nevacuate = 0
	h.noverflow = 0
	h.count = 0

	// Reset the hash seed to make it more difficult for attackers to
	// repeatedly trigger hash collisions. See issue 25237.
	h.hash0 = uint32(rand())

	// Keep the setextra allocation but clear any extra information.
	if h.extra != nil {
		*h.extra = setextra{}
	}

	// makeSetBucketArray clears the memory pointed to by h.buckets
	// and recovers any overflow buckets by generating them
	// as if h.buckets was newly alloced.
	_, nextOverflow := makeSetBucketArray(t, h.B, h.buckets)
	if nextOverflow != nil {
		// If overflow buckets are created then h.extra
		// will have been allocated during initial bucket creation.
		h.extra.nextOverflow = nextOverflow
	}

	if h.flags&hashWriting == 0 {
		fatal("concurrent set writes")
	}
	h.flags &^= hashWriting
}

func hashGrowSet(t *settype, h *hset) {
	// If we've hit the load factor, get bigger.
	// Otherwise, there are too many overflow buckets,
	// so keep the same number of buckets and "grow" laterally.
	bigger := uint8(1)
	if !overLoadFactor(h.count+1, h.B) {
		bigger = 0
		h.flags |= sameSizeGrow
	}
	oldbuckets := h.buckets
	newbuckets, nextOverflow := makeSetBucketArray(t, h.B+bigger, nil)

	flags := h.flags &^ (iterator | oldIterator)
	if h.flags&iterator != 0 {
		flags |= oldIterator
	}
	// commit the grow (atomic wrt gc)
	h.B += bigger
	h.flags = flags
	h.oldbuckets = oldbuckets
	h.buckets = newbuckets
	h.nevacuate = 0
	h.noverflow = 0

	if h.extra != nil && h.extra.overflow != nil {
		// Promote current overflow buckets to the old generation.
		if h.extra.oldoverflow != nil {
			throw("oldoverflow is not nil")
		}
		h.extra.oldoverflow = h.extra.overflow
		h.extra.overflow = nil
	}
	if nextOverflow != nil {
		if h.extra == nil {
			h.extra = new(setextra)
		}
		h.extra.nextOverflow = nextOverflow
	}

	// the actual copying of the hash table data is done incrementally
	// by growWorkSet() and evacuateSet().
}

// growing reports whether h is growing. The growth may be to the same size or bigger.
func (h *hset) growing() bool {
	return h.oldbuckets != nil
}

// sameSizeGrow reports whether the current growth is to a map of the same size.
func (h *hset) sameSizeGrow() bool {
	return h.flags&sameSizeGrow != 0
}

//wo:linkname sameSizeGrowForIssue69110TestSet
func sameSizeGrowForIssue69110TestSet(h *hset) bool {
	return h.sameSizeGrow()
}

// noldbuckets calculates the number of buckets prior to the current map growth.
func (h *hset) noldbuckets() uintptr {
	oldB := h.B
	if !h.sameSizeGrow() {
		oldB--
	}
	return bucketShift(oldB)
}

// oldbucketmask provides a mask that can be applied to calculate n % noldbuckets().
func (h *hset) oldbucketmask() uintptr {
	return h.noldbuckets() - 1
}

func growWorkSet(t *settype, h *hset, bucket uintptr) {
	// make sure we evacuateSet the oldbucket corresponding
	// to the bucket we're about to use
	evacuateSet(t, h, bucket&h.oldbucketmask())

	// evacuateSet one more oldbucket to make progress on growing
	if h.growing() {
		evacuateSet(t, h, h.nevacuate)
	}
}

func bucketEvacuatedSet(t *settype, h *hset, bucket uintptr) bool {
	b := (*bset)(add(h.oldbuckets, bucket*uintptr(t.BucketSize)))
	return evacuatedSet(b)
}

// evacDst is an evacuation destination.
type evacDstSet struct {
	b *bset          // current destination bucket
	i int            // elem index into b
	k unsafe.Pointer // pointer to current key storage
}

func evacuateSet(t *settype, h *hset, oldbucket uintptr) {
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
			// Otherwise, GC can see bad pointers.
			y := &xy[1]
			y.b = (*bset)(add(h.buckets, (oldbucket+newbit)*uintptr(t.BucketSize)))
			y.k = add(unsafe.Pointer(y.b), dataOffset)
		}

		for ; b != nil; b = b.overflow(t) {
			k := add(unsafe.Pointer(b), dataOffset)
			for i := 0; i < abi.MapBucketCount; i, k = i+1, add(k, uintptr(t.KeySize)) {
				top := b.tophash[i]
				if isEmpty(top) {
					b.tophash[i] = evacuatedEmpty
					continue
				}
				if top < minTopHash {
					throw("bad set state")
				}
				k2 := k
				if t.IndirectKey() {
					k2 = *((*unsafe.Pointer)(k2))
				}
				var useY uint8
				if !h.sameSizeGrow() {
					// Compute hash to make our evacuation decision (whether we need
					// to send this key/elem to bucket x or bucket y).
					hash := t.Hasher(k2, uintptr(h.hash0))
					if h.flags&iterator != 0 && !t.ReflexiveKey() && !t.Elem.Equal(k2, k2) {
						// If key != key (NaNs), then the hash could be (and probably
						// will be) entirely different from the old hash. Moreover,
						// it isn't reproducible. Reproducibility is required in the
						// presence of iterators, as our evacuation decision must
						// match whatever decision the iterator made.
						// Fortunately, we have the freedom to send these keysSet either
						// way. Also, tophash is meaningless for these kinds of keysSet.
						// We let the low bit of tophash drive the evacuation decision.
						// We recompute a new random tophash for the next level so
						// these keysSet will get evenly distributed across all buckets
						// after multiple grows.
						useY = top & 1
						top = tophash(hash)
					} else {
						if hash&newbit != 0 {
							useY = 1
						}
					}
				}

				if evacuatedX+1 != evacuatedY || evacuatedX^1 != evacuatedY {
					throw("bad evacuatedN")
				}

				b.tophash[i] = evacuatedX + useY // evacuatedX + 1 == evacuatedY
				dst := &xy[useY]                 // evacuation destination

				if dst.i == abi.MapBucketCount {
					dst.b = h.newoverflow(t, dst.b)
					dst.i = 0
					dst.k = add(unsafe.Pointer(dst.b), dataOffset)
				}
				dst.b.tophash[dst.i&(abi.MapBucketCount-1)] = top // mask dst.i as an optimization, to avoid a bounds check
				if t.IndirectKey() {
					*(*unsafe.Pointer)(dst.k) = k2 // copy pointer
				} else {
					typedmemmove(t.Elem, dst.k, k) // copy elem
				}
				//if t.IndirectElem() {
				//	*(*unsafe.Pointer)(dst.e) = *(*unsafe.Pointer)(e)
				//} else {
				//	typedmemmove(t.Elem, dst.e, e)
				//}
				dst.i++
				// These updates might push these pointers past the end of the
				// key or elem arrays.  That's ok, as we have the overflow pointer
				// at the end of the bucket to protect against pointing past the
				// end of the bucket.
				dst.k = add(dst.k, uintptr(t.KeySize))
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

func advanceEvacuationMarkSet(h *hset, t *settype, newbit uintptr) {
	h.nevacuate++
	// Experiments suggest that 1024 is overkill by at least an order of magnitude.
	// Put it in there as a safeguard anyway, to ensure O(1) behavior.
	stop := h.nevacuate + 1024
	if stop > newbit {
		stop = newbit
	}
	for h.nevacuate != stop && bucketEvacuatedSet(t, h, h.nevacuate) {
		h.nevacuate++
	}
	if h.nevacuate == newbit { // newbit == # of oldbuckets
		// Growing is all done. Free old main bucket array.
		h.oldbuckets = nil
		// Can discard old overflow buckets as well.
		// If they are still referenced by an iterator,
		// then the iterator holds a pointers to the slice.
		if h.extra != nil {
			h.extra.oldoverflow = nil
		}
		h.flags &^= sameSizeGrow
	}
}

// Reflect stubs. Called from ../reflect/asm_*.s

// reflect_makeset is for package reflect.
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname reflect_makeset reflect.makemap
func reflect_makeset(t *settype, cap int) *hset {
	// Check invariants and reflects math.
	if t.Elem.Equal == nil {
		throw("runtime.reflect_makeset: unsupported set key type")
	}
	if t.Elem.Size_ > abi.MapMaxKeyBytes && (!t.IndirectKey() || t.KeySize != uint8(goarch.PtrSize)) ||
		t.Elem.Size_ <= abi.MapMaxKeyBytes && (t.IndirectKey() || t.KeySize != uint8(t.Elem.Size_)) {
		throw("key size wrong")
	}
	//if t.Elem.Size_ > abi.MapMaxElemBytes && (!t.IndirectElem() || 0 != uint8(goarch.PtrSize)) ||
	//	t.Elem.Size_ <= abi.MapMaxElemBytes && (t.IndirectElem() || 0 != uint8(t.Elem.Size_)) {
	//	throw("elem size wrong")
	//}
	if t.Elem.Align_ > abi.MapBucketCount {
		throw("key align too big")
	}
	//if t.Elem.Align_ > abi.MapBucketCount {
	//	throw("elem align too big")
	//}
	if t.Elem.Size_%uintptr(t.Elem.Align_) != 0 {
		throw("key size not a multiple of key align")
	}
	//if t.Elem.Size_%uintptr(t.Elem.Align_) != 0 {
	//	throw("elem size not a multiple of elem align")
	//}
	if abi.MapBucketCount < 8 {
		throw("bucketsize too small for proper alignment")
	}
	if dataOffset%uintptr(t.Elem.Align_) != 0 {
		throw("need padding in bucket (key)")
	}
	//if dataOffset%uintptr(t.Elem.Align_) != 0 {
	//	throw("need padding in bucket (elem)")
	//}

	return makeset(t, cap, nil)
}

// reflect_setaccess is for package reflect.
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname reflect_setaccess reflect.setaccess
func reflect_setaccess(t *settype, h *hset, key unsafe.Pointer) bool {
	return setaccess2(t, h, key)
}

//wo:linkname reflect_setaccess_faststr reflect.setaccess_faststr
func reflect_setaccess_faststr(t *settype, h *hset, key string) bool {
	return setaccess2_faststr(t, h, key)
}

// reflect_setassign is for package reflect,
// Do not access it using linkname.

// Do not remove or change the type signature.
//
//wo:linkname reflect_setassign reflect.setassign0
func reflect_setassign(t *settype, h *hset, key unsafe.Pointer, elem unsafe.Pointer) {
	p := setassign(t, h, key)
	typedmemmove(t.Elem, p, elem)
}

//wo:linkname reflect_setassign_faststr reflect.setassign_faststr0
//func reflect_setassign_faststr(t *settype, h *hset, key string, elem unsafe.Pointer) {
//	p := setassign_faststr(t, h, key)
//	typedmemmove(t.Elem, p, elem)
//}

//wo:linkname reflect_setdelete reflect.setdelete
func reflect_setdelete(t *settype, h *hset, key unsafe.Pointer) {
	setdelete(t, h, key)
}

//wo:linkname reflect_setdelete_faststr reflect.setdelete_faststr
func reflect_setdelete_faststr(t *settype, h *hset, key string) {
	setdelete_faststr(t, h, key)
}

// reflect_setiterinit is for package reflect.
// Do not use it as a linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname reflect_setiterinit reflect.setiterinit
func reflect_setiterinit(t *settype, h *hset, it *hsetiter) {
	setiterinit(t, h, it)
}

// reflect_setiternext is for package reflect.
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname reflect_setiternext reflect.setiternext
func reflect_setiternext(it *hsetiter) {
	setiternext(it)
}

// reflect_setiterkey is for package reflect.
// Do not access it using linkname.

// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname reflect_setiterkey reflect.setiterkey
func reflect_setiterkey(it *hsetiter) unsafe.Pointer {
	return it.key
}

// reflect_setiterelem is for package reflect.
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname reflect_setiterelem reflect.setiterelem
//func reflect_setiterelem(it *hsetiter) unsafe.Pointer {
//	return it.elem
//}

// reflect_setlen is for package reflect.
// Do not access it using linkname.
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//wo:linkname reflect_setlen reflect.setlen
func reflect_setlen(h *hset) int {
	if h == nil {
		return 0
	}
	if raceenabled {
		callerpc := getcallerpc()
		racereadpc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(reflect_setlen))
	}
	return h.count
}

//wo:linkname reflect_setclear reflect.setclear
func reflect_setclear(t *settype, h *hset) {
	setclear(t, h)
}

//wo:linkname reflectlite_setlen internal/reflectlite.maplen
func reflectlite_setlen(h *hset) int {
	if h == nil {
		return 0
	}
	if raceenabled {
		callerpc := getcallerpc()
		racereadpc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(reflect_setlen))
	}
	return h.count
}

// setinitnoop is a no-op function known the Go linker; if a given global
// map (of the right size) is determined to be dead, the linker will
// rewrite the relocation (from the package init func) from the outlined
// map init function to this symbol. Defined in assembly so as to avoid
// complications with instrumentation (coverage, etc).
func setinitnoop()

// setclone for implementing maps.Clone
//
//wo:linkname setclone maps.clone
func setclone(m any) any {
	e := efaceOf(&m)
	e.data = unsafe.Pointer(setclone2((*settype)(unsafe.Pointer(e._type)), (*hset)(e.data)))
	return m
}

// moveToBset moves a bucket from src to dst. It returns the destination bucket or new destination bucket if it overflows
// and the pos that the next key/value will be written, if pos == bucketCnt means needs to written in overflow bucket.
func moveToBset(t *settype, h *hset, dst *bset, pos int, src *bset) (*bset, int) {
	for i := 0; i < abi.MapBucketCount; i++ {
		if isEmpty(src.tophash[i]) {
			continue
		}

		for ; pos < abi.MapBucketCount; pos++ {
			if isEmpty(dst.tophash[pos]) {
				break
			}
		}

		if pos == abi.MapBucketCount {
			dst = h.newoverflow(t, dst)
			pos = 0
		}

		srcK := add(unsafe.Pointer(src), dataOffset+uintptr(i)*uintptr(t.KeySize))
		dstK := add(unsafe.Pointer(dst), dataOffset+uintptr(pos)*uintptr(t.KeySize))

		dst.tophash[pos] = src.tophash[i]
		if t.IndirectKey() {
			srcK = *(*unsafe.Pointer)(srcK)
			if t.NeedKeyUpdate() {
				kStore := newobject(t.Elem)
				typedmemmove(t.Elem, kStore, srcK)
				srcK = kStore
			}
			// Note: if NeedKeyUpdate is false, then the memory
			// used to store the key is immutable, so we can share
			// it between the original map and its clone.
			*(*unsafe.Pointer)(dstK) = srcK
		} else {
			typedmemmove(t.Elem, dstK, srcK)
		}
		//if t.IndirectElem() {
		//	srcEle = *(*unsafe.Pointer)(srcEle)
		//	eStore := newobject(t.Elem)
		//	typedmemmove(t.Elem, eStore, srcEle)
		//	*(*unsafe.Pointer)(dstEle) = eStore
		//} else {
		//	typedmemmove(t.Elem, dstEle, srcEle)
		//}
		pos++
		h.count++
	}
	return dst, pos
}

func setclone2(t *settype, src *hset) *hset {
	hint := src.count
	if overLoadFactor(hint, src.B) {
		// Note: in rare cases (e.g. during a same-sized grow) the map
		// can be overloaded. Make sure we don't allocate a destination
		// bucket array larger than the source bucket array.
		// This will cause the cloned map to be overloaded also,
		// but that's better than crashing. See issue 69110.
		hint = int(loadFactorNum * (bucketShift(src.B) / loadFactorDen))
	}
	dst := makeset(t, hint, nil)
	dst.hash0 = src.hash0
	dst.nevacuate = 0
	// flags do not need to be copied here, just like a new map has no flags.

	if src.count == 0 {
		return dst
	}

	if src.flags&hashWriting != 0 {
		fatal("concurrent set clone and set write")
	}

	if src.B == 0 && !(t.IndirectKey() && t.NeedKeyUpdate()) {
		// Quick copy for small maps.
		dst.buckets = newobject(t.Bucket)
		dst.count = src.count
		typedmemmove(t.Bucket, dst.buckets, src.buckets)
		return dst
	}

	if dst.B == 0 {
		dst.buckets = newobject(t.Bucket)
	}
	dstArraySize := int(bucketShift(dst.B))
	srcArraySize := int(bucketShift(src.B))
	for i := 0; i < dstArraySize; i++ {
		dstBmap := (*bset)(add(dst.buckets, uintptr(i*int(t.BucketSize))))
		pos := 0
		for j := 0; j < srcArraySize; j += dstArraySize {
			srcBmap := (*bset)(add(src.buckets, uintptr((i+j)*int(t.BucketSize))))
			for srcBmap != nil {
				dstBmap, pos = moveToBset(t, dst, dstBmap, pos, srcBmap)
				srcBmap = srcBmap.overflow(t)
			}
		}
	}

	if src.oldbuckets == nil {
		return dst
	}

	oldB := src.B
	srcOldbuckets := src.oldbuckets
	if !src.sameSizeGrow() {
		oldB--
	}
	oldSrcArraySize := int(bucketShift(oldB))

	for i := 0; i < oldSrcArraySize; i++ {
		srcBset := (*bset)(add(srcOldbuckets, uintptr(i*int(t.BucketSize))))
		if evacuatedSet(srcBset) {
			continue
		}

		if oldB >= dst.B { // main bucket bits in dst is less than oldB bits in src
			dstBmap := (*bset)(add(dst.buckets, (uintptr(i)&bucketMask(dst.B))*uintptr(t.BucketSize)))
			for dstBmap.overflow(t) != nil {
				dstBmap = dstBmap.overflow(t)
			}
			pos := 0
			for srcBset != nil {
				dstBmap, pos = moveToBset(t, dst, dstBmap, pos, srcBset)
				srcBset = srcBset.overflow(t)
			}
			continue
		}

		// oldB < dst.B, so a single source bucket may go to multiple destination buckets.
		// Process entries one at a time.
		for srcBset != nil {
			// move from oldBlucket to new bucket
			for i := uintptr(0); i < abi.MapBucketCount; i++ {
				if isEmpty(srcBset.tophash[i]) {
					continue
				}

				if src.flags&hashWriting != 0 {
					fatal("concurrent set clone and set write")
				}

				srcK := add(unsafe.Pointer(srcBset), dataOffset+i*uintptr(t.KeySize))
				if t.IndirectKey() {
					srcK = *((*unsafe.Pointer)(srcK))
				}

				//srcEle := add(unsafe.Pointer(srcBset), dataOffset+abi.MapBucketCount*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
				//if t.IndirectElem() {
				//	srcEle = *((*unsafe.Pointer)(srcEle))
				//}
				//dstEle := setassign(t, dst, srcK)
				//typedmemmove(t.Elem, dstEle, srcEle)
			}
			srcBset = srcBset.overflow(t)
		}
	}
	return dst
}

// keysSet for implementing maps.keysSet
//
//wo:linkname keysSet maps.keysSet
func keysSet(m any, p unsafe.Pointer) {
	e := efaceOf(&m)
	t := (*settype)(unsafe.Pointer(e._type))
	h := (*hset)(e.data)

	if h == nil || h.count == 0 {
		return
	}
	s := (*slice)(p)
	r := int(rand())
	offset := uint8(r >> h.B & (abi.MapBucketCount - 1))
	if h.B == 0 {
		copyKeysSet(t, h, (*bset)(h.buckets), s, offset)
		return
	}
	arraySize := int(bucketShift(h.B))
	buckets := h.buckets
	for i := 0; i < arraySize; i++ {
		bucket := (i + r) & (arraySize - 1)
		b := (*bset)(add(buckets, uintptr(bucket)*uintptr(t.BucketSize)))
		copyKeysSet(t, h, b, s, offset)
	}

	if h.growing() {
		oldArraySize := int(h.noldbuckets())
		for i := 0; i < oldArraySize; i++ {
			bucket := (i + r) & (oldArraySize - 1)
			b := (*bset)(add(h.oldbuckets, uintptr(bucket)*uintptr(t.BucketSize)))
			if evacuatedSet(b) {
				continue
			}
			copyKeysSet(t, h, b, s, offset)
		}
	}
	return
}

func copyKeysSet(t *settype, h *hset, b *bset, s *slice, offset uint8) {
	for b != nil {
		for i := uintptr(0); i < abi.MapBucketCount; i++ {
			offi := (i + uintptr(offset)) & (abi.MapBucketCount - 1)
			if isEmpty(b.tophash[offi]) {
				continue
			}
			if h.flags&hashWriting != 0 {
				fatal("concurrent set read and set write")
			}
			k := add(unsafe.Pointer(b), dataOffset+offi*uintptr(t.KeySize))
			if t.IndirectKey() {
				k = *((*unsafe.Pointer)(k))
			}
			if s.len >= s.cap {
				fatal("concurrent set read and set write")
			}
			typedmemmove(t.Elem, add(s.array, uintptr(s.len)*uintptr(t.Elem.Size())), k)
			s.len++
		}
		b = b.overflow(t)
	}
}

// valuesSet for implementing maps.valuesSet
//
//wo:linkname valuesSet maps.valuesSet
func valuesSet(m any, p unsafe.Pointer) {
	// TODO ensure this function is removed deeply
	//e := efaceOf(&m)
	//t := (*settype)(unsafe.Pointer(e._type))
	//h := (*hset)(e.data)
	//if h == nil || h.count == 0 {
	//	return
	//}
	//s := (*slice)(p)
	//r := int(rand())
	//offset := uint8(r >> h.B & (abi.MapBucketCount - 1))
	//if h.B == 0 {
	//	copyValuesSet(t, h, (*bset)(h.buckets), s, offset)
	//	return
	//}
	//arraySize := int(bucketShift(h.B))
	//buckets := h.buckets
	//for i := 0; i < arraySize; i++ {
	//	bucket := (i + r) & (arraySize - 1)
	//	b := (*bset)(add(buckets, uintptr(bucket)*uintptr(t.BucketSize)))
	//	copyValuesSet(t, h, b, s, offset)
	//}
	//
	//if h.growing() {
	//	oldArraySize := int(h.noldbuckets())
	//	for i := 0; i < oldArraySize; i++ {
	//		bucket := (i + r) & (oldArraySize - 1)
	//		b := (*bset)(add(h.oldbuckets, uintptr(bucket)*uintptr(t.BucketSize)))
	//		if evacuatedSet(b) {
	//			continue
	//		}
	//		copyValuesSet(t, h, b, s, offset)
	//	}
	//}
	//return
}

func copyValuesSet(t *settype, h *hset, b *bset, s *slice, offset uint8) {
	// TODO ensure this function is removed deeply
	//for b != nil {
	//	for i := uintptr(0); i < abi.MapBucketCount; i++ {
	//		offi := (i + uintptr(offset)) & (abi.MapBucketCount - 1)
	//		if isEmpty(b.tophash[offi]) {
	//			continue
	//		}
	//
	//		if h.flags&hashWriting != 0 {
	//			fatal("concurrent set read and map write")
	//		}
	//
	//		//ele := add(unsafe.Pointer(b), dataOffset+abi.MapBucketCount*uintptr(t.KeySize))
	//		//if t.IndirectElem() {
	//		//	ele = *((*unsafe.Pointer)(ele))
	//		//}
	//		//if s.len >= s.cap {
	//		//	fatal("concurrent set read and map write")
	//		//}
	//		//typedmemmove(t.Elem, add(s.array, uintptr(s.len)*uintptr(t.Elem.Size())), ele)
	//		s.len++
	//	}
	//	b = b.overflow(t)
	//}
}
