//-----------------------------------------------------------------------------
/*

Int128/Uint128 Types

*/
//-----------------------------------------------------------------------------

package big

import (
	"math"
	"math/bits"
)

//-----------------------------------------------------------------------------

type Uint128 struct {
	Lo, Hi uint64
}

func Uint128FromUint(lo uint64) Uint128 {
	return Uint128{lo, 0}
}

func (u Uint128) Mul(v Uint128) Uint128 {
	hi, lo := bits.Mul64(u.Lo, v.Lo)
	hi += u.Hi*v.Lo + u.Lo*v.Hi
	return Uint128{lo, hi}
}

//-----------------------------------------------------------------------------

type Int128 struct {
	Lo, Hi uint64
}

func Int128FromInt(lo int64) Int128 {
	var hi uint64
	if lo < 0 {
		hi = math.MaxUint64
	}
	return Int128{uint64(lo), hi}
}

func Int128FromUint(lo uint64) Int128 {
	return Int128{lo, 0}
}

func (u Int128) Mul(v Int128) Int128 {
	hi, lo := bits.Mul64(u.Lo, v.Lo)
	hi += u.Hi*v.Lo + u.Lo*v.Hi
	return Int128{lo, hi}
}

//-----------------------------------------------------------------------------
