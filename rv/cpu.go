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
func bitExtract(x, msb, lsb uint) uint {
	return x & bitMask(msb, lsb)
}

// bitSex sign extends the value using the n-th bit as the sign.
func bitSex(x int, n uint) int {
	mask := 1 << n
	return (x ^ mask) - mask
}

// bitUnsigned extracts an unsigned bit field.
func bitUnsigned(x, msb, lsb uint) uint {
	return bitExtract(x, msb, lsb) >> lsb
}

// bitSigned extracts an signed bit field.
func bitSigned(x, msb, lsb uint) int {
	return bitSex(int(bitUnsigned(x, msb, lsb)), msb-lsb)
}

//-----------------------------------------------------------------------------
// instruction decoding

func decodeR(ins uint) (uint, uint, uint) {
	rs2 := bitUnsigned(ins, 24, 20)
	rs1 := bitUnsigned(ins, 19, 15)
	rd := bitUnsigned(ins, 11, 7)
	return rs2, rs1, rd
}

func decodeR4(ins uint) (uint, uint, uint, uint, uint) {
	rs3 := bitUnsigned(ins, 31, 27)
	rs2 := bitUnsigned(ins, 24, 20)
	rs1 := bitUnsigned(ins, 19, 15)
	rm := bitUnsigned(ins, 14, 12)
	rd := bitUnsigned(ins, 11, 7)
	return rs3, rs2, rs1, rm, rd
}

func decodeI(ins uint) (int, uint, uint) {
	imm := bitSigned(ins, 31, 20) // imm[11:0]
	rs1 := bitUnsigned(ins, 19, 15)
	rd := bitUnsigned(ins, 11, 7)
	return imm, rs1, rd
}

func decodeS(ins uint) (int, uint, uint) {
	imm0 := bitUnsigned(ins, 31, 25) // imm[11:5]
	imm1 := bitUnsigned(ins, 11, 7)  // imm[4:0]
	imm := bitSex(int((imm0<<5)+imm1), 11)
	rs2 := bitUnsigned(ins, 24, 20)
	rs1 := bitUnsigned(ins, 19, 15)
	return imm, rs2, rs1
}

func decodeB(ins uint) (int, uint, uint) {
	imm0 := bitUnsigned(ins, 31, 31) // imm[12]
	imm1 := bitUnsigned(ins, 30, 25) // imm[10:5]
	imm2 := bitUnsigned(ins, 11, 8)  // imm[4:1]
	imm3 := bitUnsigned(ins, 7, 7)   // imm[11]
	imm := bitSex(int((imm0<<12)+(imm1<<5)+(imm2<<1)+(imm3<<11)), 12)
	rs2 := bitUnsigned(ins, 24, 20)
	rs1 := bitUnsigned(ins, 19, 15)
	return imm, rs2, rs1
}

func decodeU(ins uint) (uint, uint) {
	imm := bitUnsigned(ins, 31, 12) // imm[31:12]
	rd := bitUnsigned(ins, 11, 7)
	return imm, rd
}

func decodeJ(ins uint) (int, uint) {
	imm0 := bitUnsigned(ins, 31, 31) // imm[20]
	imm1 := bitUnsigned(ins, 30, 21) // imm[10:1]
	imm2 := bitUnsigned(ins, 20, 20) // imm[11]
	imm3 := bitUnsigned(ins, 19, 12) // imm[19:12]
	imm := bitSex(int((imm0<<20)+(imm1<<1)+(imm2<<11)+(imm3<<12)), 20)
	rd := bitUnsigned(ins, 11, 7)
	return imm, rd
}

func decodeCIa(ins uint) (int, uint) {
	imm0 := bitUnsigned(ins, 12, 12) // imm[5]
	imm1 := bitUnsigned(ins, 6, 2)   // imm[4:0]
	imm := bitSex(int((imm0<<5)+(imm1<<0)), 5)
	rd := bitUnsigned(ins, 11, 7)
	return imm, rd
}

func decodeCIb(ins uint) int {
	imm0 := bitUnsigned(ins, 12, 12) // imm[9]
	imm1 := bitUnsigned(ins, 6, 6)   // imm[4]
	imm2 := bitUnsigned(ins, 5, 5)   // imm[6]
	imm3 := bitUnsigned(ins, 4, 3)   // imm[8:7]
	imm4 := bitUnsigned(ins, 2, 2)   // imm[5]
	imm := bitSex(int((imm0<<9)+(imm1<<4)+(imm2<<6)+(imm3<<7)+(imm4<<5)), 9)
	return imm
}

func decodeCIc(ins uint) (uint, uint) {
	imm0 := bitUnsigned(ins, 12, 12) // imm[5]
	imm1 := bitUnsigned(ins, 6, 2)   // imm[4:0]
	imm := (imm0 << 5) + (imm1 << 0)
	rd := bitUnsigned(ins, 9, 7) + 8
	return imm, rd
}

func decodeCId(ins uint) (uint, uint) {
	imm0 := bitUnsigned(ins, 12, 12) // imm[5]
	imm1 := bitUnsigned(ins, 6, 2)   // imm[4:0]
	imm := (imm0 << 5) + (imm1 << 0)
	rd := bitUnsigned(ins, 11, 7)
	return imm, rd
}

func decodeCIe(ins uint) (int, uint) {
	imm0 := bitUnsigned(ins, 12, 12) // imm[5]
	imm1 := bitUnsigned(ins, 6, 2)   // imm[4:0]
	imm := bitSex(int((imm0<<5)+(imm1<<0)), 5)
	rd := bitUnsigned(ins, 9, 7) + 8
	return imm, rd
}

func decodeCIf(ins uint) (int, uint) {
	imm0 := bitUnsigned(ins, 12, 12) // imm[17]
	imm1 := bitUnsigned(ins, 6, 2)   // imm[16:12]
	imm := bitSex(int((imm0<<17)+(imm1<<12)), 17) >> 12
	rd := bitUnsigned(ins, 11, 7)
	return imm, rd
}

func decodeCIW(ins uint) (uint, uint) {
	imm0 := bitUnsigned(ins, 12, 11) // imm[5:4]
	imm1 := bitUnsigned(ins, 10, 7)  // imm[9:6]
	imm2 := bitUnsigned(ins, 6, 6)   // imm[2]
	imm3 := bitUnsigned(ins, 5, 5)   // imm[3]
	imm := (imm0 << 4) + (imm1 << 6) + (imm2 << 2) + (imm3 << 3)
	rd := bitUnsigned(ins, 4, 2) + 8
	return imm, rd
}

func decodeCJa(ins uint) uint {
	rs1 := bitUnsigned(ins, 11, 7)
	return rs1
}

func decodeCJb(ins uint) int {
	imm0 := bitUnsigned(ins, 12, 12) // imm[11]
	imm1 := bitUnsigned(ins, 11, 11) // imm[4]
	imm2 := bitUnsigned(ins, 10, 9)  // imm[9:8]
	imm3 := bitUnsigned(ins, 8, 8)   // imm[10]
	imm4 := bitUnsigned(ins, 7, 7)   // imm[6]
	imm5 := bitUnsigned(ins, 6, 6)   // imm[7]
	imm6 := bitUnsigned(ins, 5, 3)   // imm[3:1]
	imm7 := bitUnsigned(ins, 2, 2)   // imm[5]
	imm := bitSex(int((imm0<<11)+(imm1<<4)+(imm2<<8)+(imm3<<10)+(imm4<<6)+(imm5<<7)+(imm6<<1)+(imm7<<5)), 11)
	return imm
}

func decodeCR(ins uint) (uint, uint) {
	rd := bitUnsigned(ins, 11, 7)
	rs := bitUnsigned(ins, 6, 2)
	return rd, rs
}

func decodeCSSa(ins uint) (uint, uint) {
	imm0 := bitUnsigned(ins, 12, 12) // imm[5]
	imm1 := bitUnsigned(ins, 6, 4)   // imm[4:2]
	imm2 := bitUnsigned(ins, 3, 2)   // imm[7:6]
	imm := (imm0 << 5) + (imm1 << 2) + (imm2 << 6)
	rd := bitUnsigned(ins, 11, 7)
	return imm, rd
}

func decodeCSSb(ins uint) (uint, uint) {
	imm0 := bitUnsigned(ins, 12, 9) // imm[5:2]
	imm1 := bitUnsigned(ins, 8, 7)  // imm[7:6]
	imm := (imm0 << 2) + (imm1 << 6)
	rs2 := bitUnsigned(ins, 6, 2)
	return imm, rs2
}

func decodeCB(ins uint) (int, uint) {
	imm0 := bitUnsigned(ins, 12, 12) // imm[8]
	imm1 := bitUnsigned(ins, 11, 10) // imm[4:3]
	imm2 := bitUnsigned(ins, 6, 5)   // imm[7:6]
	imm3 := bitUnsigned(ins, 4, 3)   // imm[2:1]
	imm4 := bitUnsigned(ins, 2, 2)   // imm[5]
	imm := bitSex(int((imm0<<8)+(imm1<<3)+(imm2<<6)+(imm3<<1)+(imm4<<5)), 8)
	rs := bitUnsigned(ins, 9, 7) + 8
	return imm, rs
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
	Symbol(adr uint32) string
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
