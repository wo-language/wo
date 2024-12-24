// Code generated by "stringer -type=Op -trimprefix=O node.go"; DO NOT EDIT.

package ir

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OXXX-0]
	_ = x[ONAME-1]
	_ = x[ONONAME-2]
	_ = x[OTYPE-3]
	_ = x[OLITERAL-4]
	_ = x[ONIL-5]
	_ = x[OADD-6]
	_ = x[OSUB-7]
	_ = x[OOR-8]
	_ = x[OXOR-9]
	_ = x[OADDSTR-10]
	_ = x[OADDR-11]
	_ = x[OANDAND-12]
	_ = x[OAPPEND-13]
	_ = x[OBYTES2STR-14]
	_ = x[OBYTES2STRTMP-15]
	_ = x[ORUNES2STR-16]
	_ = x[OSTR2BYTES-17]
	_ = x[OSTR2BYTESTMP-18]
	_ = x[OSTR2RUNES-19]
	_ = x[OSLICE2ARR-20]
	_ = x[OSLICE2ARRPTR-21]
	_ = x[OAS-22]
	_ = x[OAS2-23]
	_ = x[OAS2DOTTYPE-24]
	_ = x[OAS2FUNC-25]
	_ = x[OAS2MAPR-26]
	_ = x[OAS2RECV-27]
	_ = x[OASOP-28]
	_ = x[OCALL-29]
	_ = x[OCALLFUNC-30]
	_ = x[OCALLMETH-31]
	_ = x[OCALLINTER-32]
	_ = x[OCAP-33]
	_ = x[OCLEAR-34]
	_ = x[OCLOSE-35]
	_ = x[OCLOSURE-36]
	_ = x[OCOMPLIT-37]
	_ = x[OMAPLIT-38]
	_ = x[OSTRUCTLIT-39]
	_ = x[OARRAYLIT-40]
	_ = x[OSLICELIT-41]
	_ = x[OPTRLIT-42]
	_ = x[OCONV-43]
	_ = x[OCONVIFACE-44]
	_ = x[OCONVNOP-45]
	_ = x[OCOPY-46]
	_ = x[ODCL-47]
	_ = x[ODCLFUNC-48]
	_ = x[ODELETE-49]
	_ = x[ODOT-50]
	_ = x[ODOTPTR-51]
	_ = x[ODOTMETH-52]
	_ = x[ODOTINTER-53]
	_ = x[OXDOT-54]
	_ = x[ODOTTYPE-55]
	_ = x[ODOTTYPE2-56]
	_ = x[OEQ-57]
	_ = x[ONE-58]
	_ = x[OLT-59]
	_ = x[OLE-60]
	_ = x[OGE-61]
	_ = x[OGT-62]
	_ = x[ODEREF-63]
	_ = x[OINDEX-64]
	_ = x[OINDEXMAP-65]
	_ = x[OKEY-66]
	_ = x[OSTRUCTKEY-67]
	_ = x[OLEN-68]
	_ = x[OMAKE-69]
	_ = x[OMAKECHAN-70]
	_ = x[OMAKEMAP-71]
	_ = x[OMAKESET-72]
	_ = x[OMAKESLICE-73]
	_ = x[OMAKESLICECOPY-74]
	_ = x[OMUL-75]
	_ = x[ODIV-76]
	_ = x[OMOD-77]
	_ = x[OLSH-78]
	_ = x[ORSH-79]
	_ = x[OAND-80]
	_ = x[OANDNOT-81]
	_ = x[ONEW-82]
	_ = x[ONOT-83]
	_ = x[OBITNOT-84]
	_ = x[OPLUS-85]
	_ = x[ONEG-86]
	_ = x[OOROR-87]
	_ = x[OPANIC-88]
	_ = x[OPRINT-89]
	_ = x[OPRINTLN-90]
	_ = x[OPAREN-91]
	_ = x[OSEND-92]
	_ = x[OSLICE-93]
	_ = x[OSLICEARR-94]
	_ = x[OSLICESTR-95]
	_ = x[OSLICE3-96]
	_ = x[OSLICE3ARR-97]
	_ = x[OSLICEHEADER-98]
	_ = x[OSTRINGHEADER-99]
	_ = x[ORECOVER-100]
	_ = x[ORECOVERFP-101]
	_ = x[ORECV-102]
	_ = x[ORUNESTR-103]
	_ = x[OSELRECV2-104]
	_ = x[OMIN-105]
	_ = x[OMAX-106]
	_ = x[OREAL-107]
	_ = x[OIMAG-108]
	_ = x[OCOMPLEX-109]
	_ = x[OUNSAFEADD-110]
	_ = x[OUNSAFESLICE-111]
	_ = x[OUNSAFESLICEDATA-112]
	_ = x[OUNSAFESTRING-113]
	_ = x[OUNSAFESTRINGDATA-114]
	_ = x[OMETHEXPR-115]
	_ = x[OMETHVALUE-116]
	_ = x[OBLOCK-117]
	_ = x[OBREAK-118]
	_ = x[OCASE-119]
	_ = x[OCONTINUE-120]
	_ = x[ODEFER-121]
	_ = x[OFALL-122]
	_ = x[OFOR-123]
	_ = x[OGOTO-124]
	_ = x[OIF-125]
	_ = x[OLABEL-126]
	_ = x[OGO-127]
	_ = x[ORANGE-128]
	_ = x[ORETURN-129]
	_ = x[OSELECT-130]
	_ = x[OSWITCH-131]
	_ = x[OTYPESW-132]
	_ = x[OINLCALL-133]
	_ = x[OMAKEFACE-134]
	_ = x[OITAB-135]
	_ = x[OIDATA-136]
	_ = x[OSPTR-137]
	_ = x[OCFUNC-138]
	_ = x[OCHECKNIL-139]
	_ = x[ORESULT-140]
	_ = x[OINLMARK-141]
	_ = x[OLINKSYMOFFSET-142]
	_ = x[OJUMPTABLE-143]
	_ = x[OINTERFACESWITCH-144]
	_ = x[ODYNAMICDOTTYPE-145]
	_ = x[ODYNAMICDOTTYPE2-146]
	_ = x[ODYNAMICTYPE-147]
	_ = x[OTAILCALL-148]
	_ = x[OGETG-149]
	_ = x[OGETCALLERPC-150]
	_ = x[OGETCALLERSP-151]
	_ = x[OEND-152]
}

const _Op_name = "XXXNAMENONAMETYPELITERALNILADDSUBORXORADDSTRADDRANDANDAPPENDBYTES2STRBYTES2STRTMPRUNES2STRSTR2BYTESSTR2BYTESTMPSTR2RUNESSLICE2ARRSLICE2ARRPTRASAS2AS2DOTTYPEAS2FUNCAS2MAPRAS2RECVASOPCALLCALLFUNCCALLMETHCALLINTERCAPCLEARCLOSECLOSURECOMPLITMAPLITSTRUCTLITARRAYLITSLICELITPTRLITCONVCONVIFACECONVNOPCOPYDCLDCLFUNCDELETEDOTDOTPTRDOTMETHDOTINTERXDOTDOTTYPEDOTTYPE2EQNELTLEGEGTDEREFINDEXINDEXMAPKEYSTRUCTKEYLENMAKEMAKECHANMAKEMAPMAKESETMAKESLICEMAKESLICECOPYMULDIVMODLSHRSHANDANDNOTNEWNOTBITNOTPLUSNEGORORPANICPRINTPRINTLNPARENSENDSLICESLICEARRSLICESTRSLICE3SLICE3ARRSLICEHEADERSTRINGHEADERRECOVERRECOVERFPRECVRUNESTRSELRECV2MINMAXREALIMAGCOMPLEXUNSAFEADDUNSAFESLICEUNSAFESLICEDATAUNSAFESTRINGUNSAFESTRINGDATAMETHEXPRMETHVALUEBLOCKBREAKCASECONTINUEDEFERFALLFORGOTOIFLABELGORANGERETURNSELECTSWITCHTYPESWINLCALLMAKEFACEITABIDATASPTRCFUNCCHECKNILRESULTINLMARKLINKSYMOFFSETJUMPTABLEINTERFACESWITCHDYNAMICDOTTYPEDYNAMICDOTTYPE2DYNAMICTYPETAILCALLGETGGETCALLERPCGETCALLERSPEND"

var _Op_index = [...]uint16{0, 3, 7, 13, 17, 24, 27, 30, 33, 35, 38, 44, 48, 54, 60, 69, 81, 90, 99, 111, 120, 129, 141, 143, 146, 156, 163, 170, 177, 181, 185, 193, 201, 210, 213, 218, 223, 230, 237, 243, 252, 260, 268, 274, 278, 287, 294, 298, 301, 308, 314, 317, 323, 330, 338, 342, 349, 357, 359, 361, 363, 365, 367, 369, 374, 379, 387, 390, 399, 402, 406, 414, 421, 428, 437, 450, 453, 456, 459, 462, 465, 468, 474, 477, 480, 486, 490, 493, 497, 502, 507, 514, 519, 523, 528, 536, 544, 550, 559, 570, 582, 589, 598, 602, 609, 617, 620, 623, 627, 631, 638, 647, 658, 673, 685, 701, 709, 718, 723, 728, 732, 740, 745, 749, 752, 756, 758, 763, 765, 770, 776, 782, 788, 794, 801, 809, 813, 818, 822, 827, 835, 841, 848, 861, 870, 885, 899, 914, 925, 933, 937, 948, 959, 962}

func (i Op) String() string {
	if i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}
