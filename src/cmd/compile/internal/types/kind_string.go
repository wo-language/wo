// Code generated by "stringer -type Kind -trimprefix T type.go"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Txxx-0]
	_ = x[TINT8-1]
	_ = x[TUINT8-2]
	_ = x[TINT16-3]
	_ = x[TUINT16-4]
	_ = x[TINT32-5]
	_ = x[TUINT32-6]
	_ = x[TINT64-7]
	_ = x[TUINT64-8]
	_ = x[TINT-9]
	_ = x[TUINT-10]
	_ = x[TUINTPTR-11]
	_ = x[TCOMPLEX64-12]
	_ = x[TCOMPLEX128-13]
	_ = x[TFLOAT32-14]
	_ = x[TFLOAT64-15]
	_ = x[TBOOL-16]
	_ = x[TPTR-17]
	_ = x[TFUNC-18]
	_ = x[TSLICE-19]
	_ = x[TARRAY-20]
	_ = x[TSTRUCT-21]
	_ = x[TCHAN-22]
	_ = x[TMAP-23]
	_ = x[TSET-24]
	_ = x[TINTER-25]
	_ = x[TFORW-26]
	_ = x[TANY-27]
	_ = x[TSTRING-28]
	_ = x[TUNSAFEPTR-29]
	_ = x[TIDEAL-30]
	_ = x[TNIL-31]
	_ = x[TBLANK-32]
	_ = x[TFUNCARGS-33]
	_ = x[TCHANARGS-34]
	_ = x[TSSA-35]
	_ = x[TTUPLE-36]
	_ = x[TRESULTS-37]
	_ = x[NTYPE-38]
}

const _Kind_name = "xxxINT8UINT8INT16UINT16INT32UINT32INT64UINT64INTUINTUINTPTRCOMPLEX64COMPLEX128FLOAT32FLOAT64BOOLPTRFUNCSLICEARRAYSTRUCTCHANMAPSETINTERFORWANYSTRINGUNSAFEPTRIDEALNILBLANKFUNCARGSCHANARGSSSATUPLERESULTSNTYPE"

var _Kind_index = [...]uint8{0, 3, 7, 12, 17, 23, 28, 34, 39, 45, 48, 52, 59, 68, 78, 85, 92, 96, 99, 103, 108, 113, 119, 123, 126, 129, 134, 138, 141, 147, 156, 161, 164, 169, 177, 185, 188, 193, 200, 205 }

func (i Kind) String() string {
	if i >= Kind(len(_Kind_index)-1) {
		return "Kind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}
