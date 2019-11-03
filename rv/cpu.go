//-----------------------------------------------------------------------------
/*

RISC-V CPU Definitions

*/
//-----------------------------------------------------------------------------

package rv

import "fmt"

//-----------------------------------------------------------------------------

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

// Memory is an interface to the memory of the target system.
type Memory interface {
	Read32(adr uint32) uint32
	Write32(adr uint32, val uint32)
}

//-----------------------------------------------------------------------------

type Variant int

const (
	VariantRV32e Variant = iota
	VariantRV32
	VariantRV64
	VariantRV128
)

type RV32e struct {
	PC uint32
	X  [16]uint32
}

type RV32 struct {
	PC uint32
	X  [32]uint32
}

type RV64 struct {
	PC uint64
	X  [32]uint64
}

// NewRV returns a RISC-V CPU
func NewRV(variant Variant, isa *ISA, mem Memory) (*RV, error) {
	cpu := RV{}
	switch variant {
	case VariantRV32e:
		cpu.regs = &RV32e{}
	case VariantRV32:
		cpu.regs = &RV32{}
	case VariantRV64:
		cpu.regs = &RV64{}
	default:
		return nil, fmt.Errorf("unsupported cpu variant %d", variant)
	}
	cpu.Mem = mem
	return &cpu, nil
}

// RV is a RISC-V CPU
type RV struct {
	Mem      Memory      // memory of the target system
	regs     interface{} // cpu registers
	daDecode []linearDecode
}

//-----------------------------------------------------------------------------
