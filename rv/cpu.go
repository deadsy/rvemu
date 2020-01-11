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

func decodeId(ins uint) (uint, uint) {
	rs2 := bitUnsigned(ins, 24, 20, 0)
	rs1 := bitUnsigned(ins, 19, 15, 0)
	return rs2, rs1
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
	Call(m *RV) error
}

//-----------------------------------------------------------------------------

// Emulation error types.
const (
	ErrIllegal = (1 << iota) // illegal instruction
	ErrMemory                // memory exception
	ErrEcall                 // ecall exception
	ErrEbreak                // ebreak exception
	ErrCSR                   // CSR exception
	ErrTodo                  // unimplemented instruction
	ErrStuck                 // stuck program counter
	//ErrExit                  // exit from emulation
)

// Error is a general emulation error.
type Error struct {
	Type uint   // error type
	alen uint   // address length
	ins  uint   // illegal instruction value
	pc   uint64 // program counter at which error occurrred
	err  error  // sub error
}

func (e *Error) Error() string {
	pcStr := ""
	if e.alen == 32 {
		pcStr = fmt.Sprintf("%08x", e.pc)
	} else {
		pcStr = fmt.Sprintf("%016x", e.pc)
	}
	switch e.Type {
	case ErrIllegal:
		return "illegal instruction at PC " + pcStr
	case ErrMemory:
		return fmt.Sprintf("memory exception at PC %s, %s", pcStr, e.err)
	case ErrEcall:
		return "ecall exception at PC " + pcStr
	case ErrEbreak:
		return "ebreak exception at PC " + pcStr
	case ErrCSR:
		return fmt.Sprintf("csr exception at PC %s, %s", pcStr, e.err)
	//case ErrExit:
	//	return "exit at PC " + pcStr
	case ErrTodo:
		return "unimplemented instruction at PC " + pcStr
	case ErrStuck:
		return "stuck at PC " + pcStr
	}
	return "unknown exception at PC " + pcStr
}

// GetMemError returns a memory error from the general CPU error.
func (e *Error) GetMemError() *mem.Error {
	if e.Type != ErrMemory {
		return nil
	}
	return e.err.(*mem.Error)
}

// GetCSRError returns a CSR error from the general CPU error.
func (e *Error) GetCSRError() *csr.Error {
	if e.Type != ErrCSR {
		return nil
	}
	return e.err.(*csr.Error)
}

// errIllegal returns the error for an illegal instruction exception.
func (m *RV) errIllegal(ins uint) error {
	return &Error{
		Type: ErrIllegal,
		ins:  ins,
		alen: m.xlen,
		pc:   m.PC,
	}
}

// errEcall returns the error for an environment call exception.
func (m *RV) errEcall() error {
	return &Error{
		Type: ErrEcall,
		alen: m.xlen,
		pc:   m.PC,
	}
}

// errEbreak returns the error for an environment break exception.
func (m *RV) errEbreak() error {
	return &Error{
		Type: ErrEbreak,
		alen: m.xlen,
		pc:   m.PC,
	}
}

// errMemory returns the error for a memory exception.
func (m *RV) errMemory(err error) error {
	return &Error{
		Type: ErrMemory,
		alen: m.xlen,
		pc:   m.PC,
		err:  err,
	}
}

// errCSR returns the error for CSR access exception.
func (m *RV) errCSR(err error, ins uint) error {
	return &Error{
		Type: ErrCSR,
		ins:  ins,
		alen: m.xlen,
		pc:   m.PC,
		err:  err,
	}
}

func (m *RV) errStuckPC() error {
	return &Error{
		Type: ErrStuck,
		alen: m.xlen,
		pc:   m.PC,
	}
}

func (m *RV) errTodo() error {
	return &Error{
		Type: ErrTodo,
		alen: m.xlen,
		pc:   m.PC,
	}
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
