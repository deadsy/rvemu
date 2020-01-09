//-----------------------------------------------------------------------------
/*

Multiply Utilities

*/
//-----------------------------------------------------------------------------

package util

//-----------------------------------------------------------------------------

const mask32 = (1 << 32) - 1

// Mulhu1 returns the upper 64-bits of the product of 2 unsigned 64-bit integers.
func Mulhu1(u, v uint64) uint64 {
	ul := uint(u & mask32)
	uh := uint(u >> 32)
	vl := uint(v & mask32)
	vh := uint(v >> 32)
	//x := (uh<<32 + ul) * (vh<<32 + vl)
	//x := ((uh * vh) << 64) + ((uh * vl) << 32) + ((ul * vh) << 32) + (ul * vl)
	x := (uh * vh) + ((uh * vl) >> 32) + ((ul * vh) >> 32)
	return uint64(x)
}

// Mulhu2 returns the upper 64-bits of the product of 2 unsigned 64-bit integers.
func Mulhu2(u, v uint64) uint64 {
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

// Mulhs returns the upper 64-bits of the product of 2 signed 64-bit integers.
func Mulhs(u, v int64) int64 {
	p := Mulhu1(uint64(u), uint64(v))
	t1 := (u >> 63) & v
	t2 := (v >> 63) & u
	return int64(p) - t1 - t2
}

//-----------------------------------------------------------------------------
