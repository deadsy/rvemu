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
// default decode

func daNone(name string, adr, ins uint32) (string, string) {
	return name, "TODO"
}

//-----------------------------------------------------------------------------
// Type I Decodes

func daTypeIa(name string, adr, ins uint32) (string, string) {
	imm, rs1, rd := decodeI(ins)
	return fmt.Sprintf("%s %s,%s,%d", name, abiXName[rd], abiXName[rs1], imm), ""
}

func daTypeIb(name string, adr, ins uint32) (string, string) {
	imm, rs1, rd := decodeI(ins)
	if rs1 == 0 {
		return fmt.Sprintf("li %s,%d", abiXName[rd], imm), ""
	}
	if imm == 0 {
		return fmt.Sprintf("mv %s,%s", abiXName[rd], abiXName[rs1]), ""
	}
	return fmt.Sprintf("%s %s,%s,%d", name, abiXName[rd], abiXName[rs1], imm), ""
}

func daTypeIc(name string, adr, ins uint32) (string, string) {
	imm, rs1, rd := decodeI(ins)
	return fmt.Sprintf("%s %s,%d(%s)", name, abiXName[rd], imm, abiXName[rs1]), ""
}

func daTypeId(name string, adr, ins uint32) (string, string) {
	imm, rs1, rd := decodeI(ins)
	return fmt.Sprintf("%s %s,%s,0x%x", name, abiXName[rd], abiXName[rs1], imm), ""
}

func daTypeIe(name string, adr, ins uint32) (string, string) {
	imm, rs1, rd := decodeI(ins)
	if imm == 0 && rd == 0 && rs1 == 1 {
		return "ret", ""
	}
	if rd == 1 {
		if imm == 0 {
			return fmt.Sprintf("%s %s", name, abiXName[rs1]), ""
		}
		return fmt.Sprintf("%s %s,%d", name, abiXName[rs1], imm), ""
	}
	return fmt.Sprintf("%s %s,%s,%d", name, abiXName[rd], abiXName[rs1], imm), ""
}

func daTypeIf(name string, adr, ins uint32) (string, string) {
	imm, rs1, rd := decodeI(ins)
	if imm == -1 {
		return fmt.Sprintf("not %s,%s", abiXName[rd], abiXName[rs1]), ""
	}
	return fmt.Sprintf("%s %s,%s,%d", name, abiXName[rd], abiXName[rs1], imm), ""
}

func daTypeIg(name string, adr, ins uint32) (string, string) {
	imm, rs1, rd := decodeI(ins)
	return fmt.Sprintf("%s %s,%d(%s)", name, abiFName[rd], imm, abiXName[rs1]), ""
}

//-----------------------------------------------------------------------------
// Type U Decodes

func daTypeUa(name string, adr, ins uint32) (string, string) {
	imm, rd := decodeU(ins)
	return fmt.Sprintf("%s %s,0x%x", name, abiXName[rd], imm), ""
}

//-----------------------------------------------------------------------------
// Type S Decodes

func daTypeSa(name string, adr, ins uint32) (string, string) {
	imm, rs2, rs1 := decodeS(ins)
	return fmt.Sprintf("%s %s,%d(%s)", name, abiXName[rs2], imm, abiXName[rs1]), ""
}

func daTypeSb(name string, adr, ins uint32) (string, string) {
	imm, rs2, rs1 := decodeS(ins)
	return fmt.Sprintf("%s %s,%d(%s)", name, abiFName[rs2], imm, abiXName[rs1]), ""
}

//-----------------------------------------------------------------------------
// Type R Decodes

func daTypeRa(name string, adr, ins uint32) (string, string) {
	rs2, rs1, rd := decodeR(ins)
	return fmt.Sprintf("%s %s,%s,%s", name, abiXName[rd], abiXName[rs1], abiXName[rs2]), ""
}

func daTypeRb(name string, adr, ins uint32) (string, string) {
	rs2, rs1, rd := decodeR(ins)
	if rs2 == 0 {
		return fmt.Sprintf("%s %s,(%s)", name, abiXName[rd], abiXName[rs1]), ""
	}
	return fmt.Sprintf("%s %s,%s,(%s)", name, abiXName[rd], abiXName[rs2], abiXName[rs1]), ""
}

//-----------------------------------------------------------------------------
// Type R4 Decodes

func daTypeR4a(name string, adr, ins uint32) (string, string) {
	rs3, rs2, rs1, _, rd := decodeR4(ins)
	return fmt.Sprintf("%s %s,%s,%s,%s", name, abiFName[rd], abiFName[rs1], abiFName[rs2], abiFName[rs3]), ""
}

//-----------------------------------------------------------------------------
// Type B Decodes

func daTypeBa(name string, pc, ins uint32) (string, string) {
	imm, rs2, rs1 := decodeB(ins)
	adr := int(pc) + imm

	if rs2 == 0 {
		switch name {
		case "bge", "beq", "bne", "blt":
			return fmt.Sprintf("%sz %s,%x", name, abiXName[rs1], adr), ""
		}
	}

	return fmt.Sprintf("%s %s,%s,%x", name, abiXName[rs1], abiXName[rs2], adr), ""
}

//-----------------------------------------------------------------------------
// Type J Decodes

func daTypeJa(name string, pc, ins uint32) (string, string) {
	imm, rd := decodeJ(ins)
	if rd == 0 {
		return fmt.Sprintf("j %x", int(pc)+imm), ""
	}
	return fmt.Sprintf("%s %s,%x", name, abiXName[rd], int(pc)+imm), ""
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

	ins := m.Mem.Rd32(adr)

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
