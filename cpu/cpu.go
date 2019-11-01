//-----------------------------------------------------------------------------
/*

RISC-V CPU Definitions

*/
//-----------------------------------------------------------------------------

package cpu

//-----------------------------------------------------------------------------

type ISASet uint32

const ISArv32i = ISASet(1 << 0) // Integer
const ISArv32m = ISASet(1 << 1) // Integer Multiplication and Division
const ISArv32a = ISASet(1 << 2) // Atomics
const ISArv32f = ISASet(1 << 3) // Single-Precision Floating-Point
const ISArv32d = ISASet(1 << 4) // Double-Precision Floating-Point
const ISArv64i = ISASet(1 << 5) // Integer
const ISArv64m = ISASet(1 << 6) // Integer Multiplication and Division
const ISArv64a = ISASet(1 << 7) // Atomics
const ISArv64f = ISASet(1 << 8) // Single-Precision Floating-Point
const ISArv64d = ISASet(1 << 9) // Double-Precision Floating-Point

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

type insInfo struct {
	set       ISASet     // isa set to which this instruction belongs
	mneumonic string     // instruction mneumonic
	val       uint32     // value of the fixed bits in the instruction
	mask      uint32     // mask of the fixed bits in the instruction
	decode    decodeType // instruction decode type
}

type daFunc32 func(m *RV32, ins uint32) string
type daFunc64 func(m *RV64, ins uint32) string

type decoder struct {
	mask  uint32
	table map[uint32]*decoder
	da32  daFunc32
	da64  daFunc64
}

//-----------------------------------------------------------------------------

// RV32 is a 32-bit RISC-V CPU
type RV32 struct {
	PC uint32
	X  [32]uint32
}

// RV64 is a 64-bit RISC-V CPU
type RV64 struct {
	PC uint64
	X  [32]uint64
}

//-----------------------------------------------------------------------------
