//-----------------------------------------------------------------------------
/*

Multiply High Utilities

We need the upper 64-bits of the product of 64-bit operands.
Go does not have a 128-bit integer type, so we work this out
using 64-bit operations.

*/
//-----------------------------------------------------------------------------

package rv

//-----------------------------------------------------------------------------

const mask32 = (1 << 32) - 1

// mulhuu returns the upper 64-bits of the product of 2 unsigned 64-bit integers.
func mulhuu(u, v uint64) uint64 {
	ul := u & mask32
	vl := v & mask32
	uh := u >> 32
	vh := v >> 32
	w0 := ul * vl
	t := uh*vl + (w0 >> 32)
	w1 := t & mask32
	w2 := t >> 32
	w1 = ul*vh + w1
	return uh*vh + w2 + (w1 >> 32)
}

// mulhss returns the upper 64-bits of the product of 2 signed 64-bit integers.
func mulhss(u, v int64) int64 {
	p := mulhuu(uint64(u), uint64(v))
	t1 := (u >> 63) & v
	t2 := (v >> 63) & u
	return int64(p) - t1 - t2
}

// mulhsu returns the upper 64-bits of the product of signed and unsigned 64-bit integers.
func mulhsu(u int64, v uint64) int64 {
	p := mulhuu(uint64(u), v)
	t1 := (u >> 63) & int64(v)
	return int64(p) - t1
}

//-----------------------------------------------------------------------------
