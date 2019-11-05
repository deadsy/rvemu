//-----------------------------------------------------------------------------
/*

RISC-V Disassembler

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
	"strings"
	"unsafe"
)

//-----------------------------------------------------------------------------

var abiXName = [32]string{
	"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2",
	"s0", "s1", "a0", "a1", "a2", "a3", "a4", "a5",
	"a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7",
	"s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6",
}

var abiFName = [32]string{
	"ft0", "ft1", "ft2", "ft3", "ft4", "ft5", "ft6", "ft7",
	"fs0", "fs1", "fa0", "fa1", "fa2", "fa3", "fa4", "fa5",
	"fa6", "fa7", "fs2", "fs3", "fs4", "fs5", "fs6", "fs7",
	"fs8", "fs9", "fs10", "fs11", "ft8", "ft9", "ft10", "ft11",
}

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
// default decode

func daNone(name string, adr, ins uint32) (string, string) {
	return name, "TODO"
}

//-----------------------------------------------------------------------------
// Type I Decodes

func daTypeI(ins uint32) (int, string, string) {
	imm := bitSigned(ins, 31, 20)
	rs1 := abiXName[bitUnsigned(ins, 19, 15)]
	rd := abiXName[bitUnsigned(ins, 11, 7)]
	return imm, rs1, rd
}

// default
func daTypeIa(name string, adr, ins uint32) (string, string) {
	imm, rs1, rd := daTypeI(ins)
	return fmt.Sprintf("%s %s,%s,%d", name, rd, rs1, imm), ""
}

// addi
func daTypeIb(name string, adr, ins uint32) (string, string) {
	imm, rs1, rd := daTypeI(ins)
	if rs1 == "zero" {
		return fmt.Sprintf("li %s,%d", rd, imm), ""
	}
	if imm == 0 {
		return fmt.Sprintf("mv %s,%s", rd, rs1), ""
	}
	return fmt.Sprintf("%s %s,%s,%d", name, rd, rs1, imm), ""
}

// lb, lh, lw, lbu
func daTypeIc(name string, adr, ins uint32) (string, string) {
	imm, rs1, rd := daTypeI(ins)
	return fmt.Sprintf("%s %s,%d(%s)", name, rd, imm, rs1), ""
}

//-----------------------------------------------------------------------------
// Type U Decodes

func daTypeU(ins uint32) (uint, string) {
	imm := bitUnsigned(ins, 31, 12)
	rd := abiXName[bitUnsigned(ins, 11, 7)]
	return imm, rd
}

// default
func daTypeUa(name string, adr, ins uint32) (string, string) {
	imm, rd := daTypeU(ins)
	return fmt.Sprintf("%s %s,0x%x", name, rd, imm), ""
}

//-----------------------------------------------------------------------------
// Type S Decodes

func daTypeS(ins uint32) (int, string, string) {
	x := (bitUnsigned(ins, 31, 25) << 5) + bitUnsigned(ins, 11, 7)
	imm := bitSex(int(x), 11)
	rs2 := abiXName[bitUnsigned(ins, 24, 20)]
	rs1 := abiXName[bitUnsigned(ins, 19, 15)]
	return imm, rs2, rs1
}

// default
func daTypeSa(name string, adr, ins uint32) (string, string) {
	imm, rs2, rs1 := daTypeS(ins)
	return fmt.Sprintf("%s %s,%d(%s)", name, rs2, imm, rs1), ""
}

//-----------------------------------------------------------------------------
// Type R Decodes

func daTypeR(ins uint32) (string, string, string) {
	rs2 := abiXName[bitUnsigned(ins, 24, 20)]
	rs1 := abiXName[bitUnsigned(ins, 19, 15)]
	rd := abiXName[bitUnsigned(ins, 11, 7)]
	return rs2, rs1, rd
}

// default
func daTypeRa(name string, adr, ins uint32) (string, string) {
	rs2, rs1, rd := daTypeR(ins)
	return fmt.Sprintf("%s %s,%s,%s", name, rd, rs1, rs2), ""
}

//-----------------------------------------------------------------------------

// SymbolTable maps an address to a symbol.
type SymbolTable map[uint32]string

// Disassembly returns the result of the disassembler call.
type Disassembly struct {
	Dump        string // address and memory bytes
	Symbol      string // symbol for the address (if any)
	Instruction string // instruction decode
	Comment     string // useful comment
	N           int    // length in bytes of decode
}

func (da *Disassembly) String() string {
	s := make([]string, 2)
	s[0] = fmt.Sprintf("%-16s %8s %-13s", da.Dump, da.Symbol, da.Instruction)
	if da.Comment != "" {
		s[1] = fmt.Sprintf(" ; %s", da.Comment)
	}
	return strings.Join(s, "")
}

//-----------------------------------------------------------------------------

func daDump(adr, ins uint32) string {
	return fmt.Sprintf("%08x: %08x", adr, ins)
}

func daSymbol(adr uint32, st SymbolTable) string {
	if st != nil {
		return st[adr]
	}
	return ""
}

// daInstruction returns the disassembly and comment for the instruction.
func (isa *ISA) daInstruction(adr, ins uint32) (string, string) {
	ii := isa.lookup(ins)
	if ii != nil {
		return ii.da(ii.name, adr, ins)
	}
	return "?", "unknown"
}

//-----------------------------------------------------------------------------

// Disassemble a RISC-V instruction at the address.
func (m *RV) Disassemble(adr uint32, st SymbolTable) *Disassembly {

	ins := m.Mem.Read32(adr)

	instruction, comment := m.isa.daInstruction(adr, ins)

	return &Disassembly{
		Dump:        daDump(adr, ins),
		Symbol:      daSymbol(adr, st),
		Instruction: instruction,
		Comment:     comment,
		N:           int(unsafe.Sizeof(ins)),
	}
}

//-----------------------------------------------------------------------------
