//-----------------------------------------------------------------------------
/*

RISC-V CPU Definitions

*/
//-----------------------------------------------------------------------------

package cpu

//-----------------------------------------------------------------------------

// ISAModule is the numeric identifier of an ISA sub-module.
type ISAModule uint32

// Identifiers for ISA sub-modules.
const (
	ISArv32i ISAModule = (1 << iota) // Integer
	ISArv32m                         // Integer Multiplication and Division
	ISArv32a                         // Atomics
	ISArv32f                         // Single-Precision Floating-Point
	ISArv32d                         // Double-Precision Floating-Point
	ISArv64i                         // Integer
	ISArv64m                         // Integer Multiplication and Division
	ISArv64a                         // Atomics
	ISArv64f                         // Single-Precision Floating-Point
	ISArv64d                         // Double-Precision Floating-Point
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

type insInfo struct {
	module    ISAModule  // isa module to which this instruction belongs
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
