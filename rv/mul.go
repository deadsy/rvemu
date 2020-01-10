//-----------------------------------------------------------------------------
/*

Multiply Utilities

*/
//-----------------------------------------------------------------------------

package rv

//-----------------------------------------------------------------------------

const mask32 = (1 << 32) - 1

// mulhu64 returns the upper 64-bits of the product of 2 unsigned 64-bit integers.
func mulhu64(u, v uint64) uint64 {
	ul := uint(u & mask32)
	vl := uint(v & mask32)
	uh := uint(u >> 32)
	vh := uint(v >> 32)

	w0 := ul * vl
	t := uh*vl + (w0 >> 32)
	w1 := t & mask32
	w2 := t >> 32
	w1 = ul*vh + w1
	return uint64(uh*vh + w2 + (w1 >> 32))
}

// mulhs64 returns the upper 64-bits of the product of 2 signed 64-bit integers.
func mulhs64(u, v int64) int64 {
	p := mulhu64(uint64(u), uint64(v))
	t1 := (u >> 63) & v
	t2 := (v >> 63) & u
	return int64(p) - t1 - t2
}

//-----------------------------------------------------------------------------
