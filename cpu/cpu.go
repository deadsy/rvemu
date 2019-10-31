//-----------------------------------------------------------------------------
/*

RISC-V CPU Definitions

*/
//-----------------------------------------------------------------------------

package cpu

//-----------------------------------------------------------------------------

type isaSet int

const (
	isaNone isaSet = iota
	isaRV32I
	isaRV64I
	isaRV32M
	isaRV64M
	isaRV32A
	isaRV64A
	isaRV32F
	isaRV64F
	isaRV32D
	isaRV64D
)

type decodeType int

const (
	decodeNone decodeType = iota
	decodeR
	decodeI
	decodeS
	decodeB
	decodeU
	decodeJ
	decodeSB
	decodeUJ
	decodeFence
)

//-----------------------------------------------------------------------------

type insDecode struct {
	mneumonic string     // instruction mneumonic
	mask      uint32     // mask of the fixed bits in the instruction
	val       uint32     // value of the fixed bits in the instruction
	decode    decodeType // decode type for the instruction
	blah      string
}

//-----------------------------------------------------------------------------

// RV32I is a CPU implementing the RV32I ISA
type RV32I struct {
	PC uint32
	X  [32]uint32
}

//-----------------------------------------------------------------------------
