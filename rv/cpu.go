//-----------------------------------------------------------------------------
/*

RISC-V CPU Definitions

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
	"math"
	"strings"

	"github.com/deadsy/riscv/csr"
	"github.com/deadsy/riscv/mem"
)

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
func bitUnsigned(x, msb, lsb, shift uint) uint {
	return (bitExtract(x, msb, lsb) >> lsb) << shift
}

// bitSigned extracts an signed bit field.
func bitSigned(x, msb, lsb uint) int {
	return bitSex(int(bitUnsigned(x, msb, lsb, 0)), msb-lsb)
}

//-----------------------------------------------------------------------------
// instruction decoding

func decodeR(ins uint) (uint, uint, uint, uint) {
	rs2 := bitUnsigned(ins, 24, 20, 0)
	rs1 := bitUnsigned(ins, 19, 15, 0)
	rm := bitUnsigned(ins, 14, 12, 0)
	rd := bitUnsigned(ins, 11, 7, 0)
	return rs2, rs1, rm, rd
}

func decodeR4(ins uint) (uint, uint, uint, uint, uint) {
	rs3 := bitUnsigned(ins, 31, 27, 0)
	rs2 := bitUnsigned(ins, 24, 20, 0)
	rs1 := bitUnsigned(ins, 19, 15, 0)
	rm := bitUnsigned(ins, 14, 12, 0)
	rd := bitUnsigned(ins, 11, 7, 0)
	return rs3, rs2, rs1, rm, rd
}

func decodeIa(ins uint) (int, uint, uint) {
	imm := bitSigned(ins, 31, 20) // imm[11:0]
	rs1 := bitUnsigned(ins, 19, 15, 0)
	rd := bitUnsigned(ins, 11, 7, 0)
	return imm, rs1, rd
}

func decodeIb(ins uint) (uint, uint, uint) {
	csr := bitUnsigned(ins, 31, 20, 0)
	rs1 := bitUnsigned(ins, 19, 15, 0)
	rd := bitUnsigned(ins, 11, 7, 0)
	return csr, rs1, rd
}

func decodeIc(ins uint) (uint, uint, uint) {
	shamt := bitUnsigned(ins, 25, 20, 0)
	rs1 := bitUnsigned(ins, 19, 15, 0)
	rd := bitUnsigned(ins, 11, 7, 0)
	return shamt, rs1, rd
}

func decodeS(ins uint) (int, uint, uint) {
	uimm := bitUnsigned(ins, 31, 25, 5) // imm[11:5]
	uimm += bitUnsigned(ins, 11, 7, 0)  // imm[4:0]
	imm := bitSex(int(uimm), 11)
	rs2 := bitUnsigned(ins, 24, 20, 0)
	rs1 := bitUnsigned(ins, 19, 15, 0)
	return imm, rs2, rs1
}

func decodeB(ins uint) (int, uint, uint) {
	uimm := bitUnsigned(ins, 31, 31, 12) // imm[12]
	uimm += bitUnsigned(ins, 30, 25, 5)  // imm[10:5]
	uimm += bitUnsigned(ins, 11, 8, 1)   // imm[4:1]
	uimm += bitUnsigned(ins, 7, 7, 11)   // imm[11]
	imm := bitSex(int(uimm), 12)
	rs2 := bitUnsigned(ins, 24, 20, 0)
	rs1 := bitUnsigned(ins, 19, 15, 0)
	return imm, rs2, rs1
}

func decodeU(ins uint) (int, uint) {
	uimm := bitUnsigned(ins, 31, 12, 0) // imm[31:12]
	imm := bitSex(int(uimm), 19)
	rd := bitUnsigned(ins, 11, 7, 0)
	return imm, rd
}

func decodeJ(ins uint) (int, uint) {
	uimm := bitUnsigned(ins, 31, 31, 20) // imm[20]
	uimm += bitUnsigned(ins, 30, 21, 1)  // imm[10:1]
	uimm += bitUnsigned(ins, 20, 20, 11) // imm[11]
	uimm += bitUnsigned(ins, 19, 12, 12) // imm[19:12]
	imm := bitSex(int(uimm), 20)
	rd := bitUnsigned(ins, 11, 7, 0)
	return imm, rd
}

func decodeCIa(ins uint) (int, uint) {
	uimm := bitUnsigned(ins, 12, 12, 5) // imm[5]
	uimm += bitUnsigned(ins, 6, 2, 0)   // imm[4:0]
	imm := bitSex(int(uimm), 5)
	rd := bitUnsigned(ins, 11, 7, 0)
	return imm, rd
}

func decodeCIb(ins uint) int {
	uimm := bitUnsigned(ins, 12, 12, 9) // imm[9]
	uimm += bitUnsigned(ins, 6, 6, 4)   // imm[4]
	uimm += bitUnsigned(ins, 5, 5, 6)   // imm[6]
	uimm += bitUnsigned(ins, 4, 3, 7)   // imm[8:7]
	uimm += bitUnsigned(ins, 2, 2, 5)   // imm[5]
	imm := bitSex(int(uimm), 9)
	return imm
}

func decodeCIc(ins uint) (uint, uint) {
	uimm := bitUnsigned(ins, 12, 12, 5) // imm[5]
	uimm += bitUnsigned(ins, 6, 2, 0)   // imm[4:0]
	rd := bitUnsigned(ins, 9, 7, 0) + 8
	return uimm, rd
}

func decodeCId(ins uint) (uint, uint) {
	uimm := bitUnsigned(ins, 12, 12, 5) // uimm[5]
	uimm += bitUnsigned(ins, 6, 2, 0)   // uimm[4:0]
	rd := bitUnsigned(ins, 11, 7, 0)
	return uimm, rd
}

func decodeCIe(ins uint) (int, uint) {
	uimm := bitUnsigned(ins, 12, 12, 5) // imm[5]
	uimm += bitUnsigned(ins, 6, 2, 0)   // imm[4:0]
	imm := bitSex(int(uimm), 5)
	rd := bitUnsigned(ins, 9, 7, 0) + 8
	return imm, rd
}

func decodeCIf(ins uint) (int, uint) {
	uimm := bitUnsigned(ins, 12, 12, 17) // imm[17]
	uimm += bitUnsigned(ins, 6, 2, 12)   // imm[16:12]
	imm := bitSex(int(uimm), 17) >> 12
	rd := bitUnsigned(ins, 11, 7, 0)
	return imm, rd
}

func decodeCIg(ins uint) (uint, uint) {
	uimm := bitUnsigned(ins, 12, 12, 5) // imm[5]
	uimm += bitUnsigned(ins, 6, 5, 3)   // imm[4:3]
	uimm += bitUnsigned(ins, 4, 2, 6)   // imm[8:6]
	rd := bitUnsigned(ins, 11, 7, 0)
	return uimm, rd
}

func decodeCIW(ins uint) (uint, uint) {
	uimm := bitUnsigned(ins, 12, 11, 4) // imm[5:4]
	uimm += bitUnsigned(ins, 10, 7, 6)  // imm[9:6]
	uimm += bitUnsigned(ins, 6, 6, 2)   // imm[2]
	uimm += bitUnsigned(ins, 5, 5, 3)   // imm[3]
	rd := bitUnsigned(ins, 4, 2, 0) + 8
	return uimm, rd
}

func decodeCJ(ins uint) int {
	uimm := bitUnsigned(ins, 12, 12, 11) // imm[11]
	uimm += bitUnsigned(ins, 11, 11, 4)  // imm[4]
	uimm += bitUnsigned(ins, 10, 9, 8)   // imm[9:8]
	uimm += bitUnsigned(ins, 8, 8, 10)   // imm[10]
	uimm += bitUnsigned(ins, 7, 7, 6)    // imm[6]
	uimm += bitUnsigned(ins, 6, 6, 7)    // imm[7]
	uimm += bitUnsigned(ins, 5, 3, 1)    // imm[3:1]
	uimm += bitUnsigned(ins, 2, 2, 5)    // imm[5]
	imm := bitSex(int(uimm), 11)
	return imm
}

func decodeCR(ins uint) (uint, uint) {
	rd := bitUnsigned(ins, 11, 7, 0)
	rs := bitUnsigned(ins, 6, 2, 0)
	return rd, rs
}

func decodeCRa(ins uint) (uint, uint) {
	rd := bitUnsigned(ins, 9, 7, 0) + 8
	rs := bitUnsigned(ins, 4, 2, 0) + 8
	return rd, rs
}

func decodeCS(ins uint) (uint, uint, uint) {
	uimm := bitUnsigned(ins, 12, 10, 3) // imm[5:3]
	uimm += bitUnsigned(ins, 6, 6, 2)   // imm[2]
	uimm += bitUnsigned(ins, 5, 5, 6)   // imm[6]
	rs1 := bitUnsigned(ins, 9, 7, 0) + 8
	rs2 := bitUnsigned(ins, 4, 2, 0) + 8
	return uimm, rs1, rs2
}

func decodeCSa(ins uint) (uint, uint, uint) {
	uimm := bitUnsigned(ins, 12, 10, 3) // imm[5:3]
	uimm += bitUnsigned(ins, 6, 5, 6)   // imm[7:6]
	rs1 := bitUnsigned(ins, 9, 7, 0) + 8
	rs2 := bitUnsigned(ins, 4, 2, 0) + 8
	return uimm, rs1, rs2
}

func decodeCSSa(ins uint) (uint, uint) {
	uimm := bitUnsigned(ins, 12, 12, 5) // imm[5]
	uimm += bitUnsigned(ins, 6, 4, 2)   // imm[4:2]
	uimm += bitUnsigned(ins, 3, 2, 6)   // imm[7:6]
	rd := bitUnsigned(ins, 11, 7, 0)
	return uimm, rd
}

func decodeCSSb(ins uint) (uint, uint) {
	uimm := bitUnsigned(ins, 12, 9, 2) // imm[5:2]
	uimm += bitUnsigned(ins, 8, 7, 6)  // imm[7:6]
	rs2 := bitUnsigned(ins, 6, 2, 0)
	return uimm, rs2
}

func decodeCSSc(ins uint) (uint, uint) {
	uimm := bitUnsigned(ins, 12, 10, 3) // imm[5:3]
	uimm += bitUnsigned(ins, 9, 7, 6)   // imm[8:6]
	rs2 := bitUnsigned(ins, 6, 2, 0)
	return uimm, rs2
}

func decodeCB(ins uint) (int, uint) {
	uimm := bitUnsigned(ins, 12, 12, 8) // imm[8]
	uimm += bitUnsigned(ins, 11, 10, 3) // imm[4:3]
	uimm += bitUnsigned(ins, 6, 5, 6)   // imm[7:6]
	uimm += bitUnsigned(ins, 4, 3, 1)   // imm[2:1]
	uimm += bitUnsigned(ins, 2, 2, 5)   // imm[5]
	imm := bitSex(int(uimm), 8)
	rs := bitUnsigned(ins, 9, 7, 0) + 8
	return imm, rs
}

//-----------------------------------------------------------------------------

// Ecall provides pluggable ecall functions.
type Ecall interface {
	Call(m *RV)
}

//-----------------------------------------------------------------------------
// memory exceptions

type memoryException struct {
	addr  uint          // memory address that caused the exception
	ex    mem.Exception // exception bitmap
	sname string        // name of the section containing the exception address
}

func (e memoryException) String() string {
	// TODO - fix the 64-bit address case
	return fmt.Sprintf("%s @ %08x (%s)", e.ex, e.addr, e.sname)
}

//-----------------------------------------------------------------------------
// CSR exceptions

type csrException struct {
	reg uint          // CSR causing the exception
	ex  csr.Exception // exception bitmap
}

func (e csrException) String() string {
	// TODO - fix the 64-bit address case
	return fmt.Sprintf("%s: %s", csr.Name(e.reg), e.ex)
}

//-----------------------------------------------------------------------------

// Exception numbers.
const (
	ExNone    = iota // normal (0)
	ExIllegal        // illegal instruction
	ExExit           // exit from emulation
	ExTodo           // unimplemented instruction
	ExMemory         // memory exception
	ExCSR            // CSR exception
	ExEcall          // unrecognised ecall
	ExBreak          // debug break point
	ExBadReg         // bad register number (rv32e)
	ExStuck          // stuck program counter
)

// Exception is a general emulation exception.
type Exception struct {
	N    int             // exception number
	alen int             // address length
	pc   uint            // program counter at which exception occurrred
	mem  memoryException // memory exception info
	csr  csrException    // csr exception info
}

func (e *Exception) Error() string {

	pcStr := ""
	if e.alen == 32 {
		pcStr = fmt.Sprintf("%08x", e.pc)
	} else {
		pcStr = fmt.Sprintf("%016x", e.pc)
	}

	switch e.N {
	case ExNone:
		return ""
	case ExIllegal:
		return "illegal instruction at PC " + pcStr
	case ExExit:
		return "exit at PC " + pcStr
	case ExTodo:
		return "unimplemented instruction at PC " + pcStr
	case ExMemory:
		return fmt.Sprintf("memory exception at PC %s, %s", pcStr, e.mem)
	case ExCSR:
		return fmt.Sprintf("csr exception at PC %s, %s", pcStr, e.csr)
	case ExEcall:
		return "unrecognized ecall at PC " + pcStr
	case ExBreak:
		return "breakpoint at PC " + pcStr
	}

	return "unknown exception at PC " + pcStr
}

//-----------------------------------------------------------------------------

func intRegString(reg []uint, pc, xlen uint) string {
	fmtx := "%08x"
	if xlen == 64 {
		fmtx = "%016x"
	}
	s := make([]string, len(reg)+1)
	for i := 0; i < len(reg); i++ {
		regStr := fmt.Sprintf("x%d", i)
		valStr := "0"
		if reg[i] != 0 {
			valStr = fmt.Sprintf(fmtx, reg[i])
		}
		s[i] = fmt.Sprintf("%-4s %-4s %s", regStr, abiXName[i], valStr)
	}
	s[len(reg)] = fmt.Sprintf("%-9s "+fmtx, "pc", pc)
	return strings.Join(s, "\n")
}

func floatRegString(reg []uint64) string {
	s := make([]string, len(reg))
	for i := 0; i < len(reg); i++ {
		regStr := fmt.Sprintf("f%d", i)
		valStr := "0"
		if reg[i] != 0 {
			valStr = fmt.Sprintf("%016x", reg[i])
		}
		f32 := math.Float32frombits(uint32(reg[i]))
		s[i] = fmt.Sprintf("%-4s %-4s %-16s %f", regStr, abiFName[i], valStr, f32)
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
