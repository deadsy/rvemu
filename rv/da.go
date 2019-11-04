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

func bfUnsigned(val uint32, n, shift uint) uint {
	return uint((val >> shift) & ((1 << n) - 1))
}

func bfSigned(val uint32, n, shift uint) int {
	x := int(bfUnsigned(val, n, shift))
	if x&(1<<(n-1)) != 0 {
		x -= 1 << n
	}
	return x
}

//-----------------------------------------------------------------------------

func daNone(mneumonic string, adr, ins uint32) (string, string) {
	return mneumonic, "TODO"
}

func daTypeI(mneumonic string, adr, ins uint32) (string, string) {

	imm := bfSigned(ins, 12, 20)
	rs1 := abiXName[bfUnsigned(ins, 5, 15)]
	rd := abiXName[bfUnsigned(ins, 5, 7)]

	if mneumonic == "addi" && rs1 == "zero" {
		return fmt.Sprintf("li %s,%d", rd, imm), ""
	}

	if mneumonic == "addi" && imm == 0 {
		return fmt.Sprintf("mv %s,%s", rd, rs1), ""
	}

	if mneumonic == "lbu" || mneumonic == "lw" {
		return fmt.Sprintf("%s %s,%d(%s)", mneumonic, rd, imm, rs1), ""
	}

	return fmt.Sprintf("%s %s,%s,%d", mneumonic, rd, rs1, imm), ""
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
		return ii.decode.da(ii.mneumonic, adr, ins)
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
