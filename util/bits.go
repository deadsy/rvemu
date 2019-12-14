//-----------------------------------------------------------------------------
/*

Utilities to Read/Write Bit Fields

*/
//-----------------------------------------------------------------------------

package util

//-----------------------------------------------------------------------------

// bitMask returns a bit mask from the msb to lsb bits.
func bitMask(msb, lsb uint) uint {
	n := msb - lsb + 1
	return ((1 << n) - 1) << lsb
}

func RdBits(x, msb, lsb uint) uint {
	return (x & bitMask(msb, lsb)) >> lsb
}

func MaskBits(x, msb, lsb uint) uint {
	return x & bitMask(msb, lsb)
}

func WrBits(x, val, msb, lsb uint) uint {
	mask := bitMask(msb, lsb)
	val = (val << lsb) & mask
	return (x & ^mask) | val
}

//-----------------------------------------------------------------------------
