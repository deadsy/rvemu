//-----------------------------------------------------------------------------
/*

RISC-V CPU Definitions

*/
//-----------------------------------------------------------------------------

package cpu

//-----------------------------------------------------------------------------

type ISASet uint32

const ISArv32i = ISASet(1 << 0)
const ISArv32m = ISASet(1 << 1)
const ISArv32a = ISASet(1 << 2)
const ISArv32f = ISASet(1 << 3)
const ISArv32d = ISASet(1 << 4)
const ISArv64i = ISASet(1 << 5)
const ISArv64m = ISASet(1 << 6)
const ISArv64a = ISASet(1 << 7)
const ISArv64f = ISASet(1 << 8)
const ISArv64d = ISASet(1 << 9)

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
	set       ISASet     // isa set to which this instruction belongs
	mneumonic string     // instruction mneumonic
	val       uint32     // value of the fixed bits in the instruction
	mask      uint32     // mask of the fixed bits in the instruction
	decode    decodeType // instruction decode type
}

//-----------------------------------------------------------------------------

// RV32I is a CPU implementing the RV32I ISA
type RV32I struct {
	PC uint32
	X  [32]uint32
}

//-----------------------------------------------------------------------------
