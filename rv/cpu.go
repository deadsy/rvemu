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

// RV is a RISC-V CPU
type RV struct {
	Mem   Memory // memory of the target system
	X     [32]uint64
	PC    uint64
	xlen  uint // register bit length
	nregs uint // number of registers
	isa   *ISA // ISA implemented for the CPU
}

// NewRV returns a RISC-V CPU
func NewRV(variant Variant, isa *ISA, mem Memory) (*RV, error) {
	cpu := RV{
		Mem: mem,
		isa: isa,
	}

	switch variant {
	case VariantRV32e:
		cpu.nregs = 16
		cpu.xlen = 32
	case VariantRV32:
		cpu.nregs = 32
		cpu.xlen = 32
	case VariantRV64:
		cpu.nregs = 32
		cpu.xlen = 64
	default:
		return nil, fmt.Errorf("unsupported cpu variant %d", variant)
	}

	return &cpu, nil
}

//-----------------------------------------------------------------------------
