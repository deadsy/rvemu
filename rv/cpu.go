//-----------------------------------------------------------------------------
/*

RISC-V CPU Definitions

*/
//-----------------------------------------------------------------------------

package rv

//-----------------------------------------------------------------------------
// Bitfield Operations

// bitMask returns a bit mask from the msb to lsb bits.
func bitMask(msb, lsb uint) uint {
	n := msb - lsb + 1
	return ((1 << n) - 1) << lsb
}

// bitExtract extracts a bit field from a value (no shifting).
func bitExtract(x uint32, msb, lsb uint) uint {
	return uint(x) & bitMask(msb, lsb)
}

// bitSex sign extends the value using the n-th bit as the sign.
func bitSex(x int, n uint) int {
	mask := 1 << n
	return (x ^ mask) - mask
}

// bitUnsigned extracts an unsigned bit field.
func bitUnsigned(x uint32, msb, lsb uint) uint {
	return bitExtract(x, msb, lsb) >> lsb
}

// bitSigned extracts an signed bit field.
func bitSigned(x uint32, msb, lsb uint) int {
	return bitSex(int(bitUnsigned(x, msb, lsb)), msb-lsb)
}

//-----------------------------------------------------------------------------
// instruction decoding

func decodeR(ins uint32) (uint, uint, uint) {
	rs2 := bitUnsigned(ins, 24, 20)
	rs1 := bitUnsigned(ins, 19, 15)
	rd := bitUnsigned(ins, 11, 7)
	return rs2, rs1, rd
}

func decodeR4(ins uint32) (uint, uint, uint, uint, uint) {
	rs3 := bitUnsigned(ins, 31, 27)
	rs2 := bitUnsigned(ins, 24, 20)
	rs1 := bitUnsigned(ins, 19, 15)
	rm := bitUnsigned(ins, 14, 12)
	rd := bitUnsigned(ins, 11, 7)
	return rs3, rs2, rs1, rm, rd
}

func decodeI(ins uint32) (int, uint, uint) {
	imm := bitSigned(ins, 31, 20) // imm[11:0]
	rs1 := bitUnsigned(ins, 19, 15)
	rd := bitUnsigned(ins, 11, 7)
	return imm, rs1, rd
}

func decodeS(ins uint32) (int, uint, uint) {
	imm0 := bitUnsigned(ins, 31, 25) // imm[11:5]
	imm1 := bitUnsigned(ins, 11, 7)  // imm[4:0]
	x := int((imm0 << 5) + imm1)
	imm := bitSex(x, 11)
	rs2 := bitUnsigned(ins, 24, 20)
	rs1 := bitUnsigned(ins, 19, 15)
	return imm, rs2, rs1
}

func decodeB(ins uint32) (int, uint, uint) {
	imm0 := bitUnsigned(ins, 31, 31) // imm[12]
	imm1 := bitUnsigned(ins, 30, 25) // imm[10:5]
	imm2 := bitUnsigned(ins, 11, 8)  // imm[4:1]
	imm3 := bitUnsigned(ins, 7, 7)   // imm[11]
	x := int((imm0 << 12) + (imm1 << 5) + (imm2 << 1) + (imm3 << 11))
	imm := bitSex(x, 12)
	rs2 := bitUnsigned(ins, 24, 20)
	rs1 := bitUnsigned(ins, 19, 15)
	return imm, rs2, rs1
}

func decodeU(ins uint32) (uint, uint) {
	imm := bitUnsigned(ins, 31, 12) // imm[31:12]
	rd := bitUnsigned(ins, 11, 7)
	return imm, rd
}

func decodeJ(ins uint32) (int, uint) {
	imm0 := bitUnsigned(ins, 31, 31) // imm[20]
	imm1 := bitUnsigned(ins, 30, 21) // imm[10:1]
	imm2 := bitUnsigned(ins, 20, 20) // imm[11]
	imm3 := bitUnsigned(ins, 19, 12) // imm[19:12]
	x := int((imm0 << 20) + (imm1 << 1) + (imm2 << 11) + (imm3 << 12))
	imm := bitSex(x, 20)
	rd := bitUnsigned(ins, 11, 7)
	return imm, rd
}

//-----------------------------------------------------------------------------

// Memory is an interface to the memory of the target system.
type Memory interface {
	Rd32(adr uint32) uint32
	Wr32(adr uint32, val uint32)
	Rd16(adr uint32) uint16
	Wr16(adr uint32, val uint16)
	Rd8(adr uint32) uint8
	Wr8(adr uint32, val uint8)
}

//-----------------------------------------------------------------------------

// RV is a RISC-V CPU
type RV struct {
	Mem   Memory // memory of the target system
	X     [32]uint64
	PC    uint64
	xlen  uint // register bit length
	nregs uint // number of registers
	isa   *ISA // ISA implemented for the CPU
}

// newRV returns a RISC-V CPU
func newRV(isa *ISA, mem Memory, xlen, nregs uint) *RV {
	return &RV{
		Mem:   mem,
		xlen:  xlen,
		nregs: nregs,
		isa:   isa,
	}
}

// NewRV32 returns a 32-bit RISC-V CPU
func NewRV32(isa *ISA, mem Memory) *RV {
	return newRV(isa, mem, 32, 32)
}

// NewRV32e returns a 32-bit embedded RISC-V CPU
func NewRV32e(isa *ISA, mem Memory) *RV {
	return newRV(isa, mem, 32, 16)
}

// NewRV64 returns a 64-bit RISC-V CPU
func NewRV64(isa *ISA, mem Memory) *RV {
	return newRV(isa, mem, 64, 32)
}

//-----------------------------------------------------------------------------
