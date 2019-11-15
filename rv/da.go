//-----------------------------------------------------------------------------
/*

RISC-V Disassembler

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
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

const regZero = 0 // zero
const regRa = 1   // return address
const regSp = 2   // stack pointer
const regGp = 3   // global pointer
const regTp = 3   // thread pointer
const regFp = 8   // frame pointer

//-----------------------------------------------------------------------------
// default decode

func daNone(name string, pc uint32, ins uint) string {
	return fmt.Sprintf("%s TODO", name)
}

//-----------------------------------------------------------------------------
// Type I Decodes

func daTypeIa(name string, pc uint32, ins uint) string {
	imm, rs1, rd := decodeIa(ins)
	return fmt.Sprintf("%s %s,%s,%d", name, abiXName[rd], abiXName[rs1], imm)
}

func daTypeIb(name string, pc uint32, ins uint) string {
	imm, rs1, rd := decodeIa(ins)
	if rs1 == 0 {
		return fmt.Sprintf("li %s,%d", abiXName[rd], imm)
	}
	if imm == 0 {
		return fmt.Sprintf("mv %s,%s", abiXName[rd], abiXName[rs1])
	}
	return fmt.Sprintf("%s %s,%s,%d", name, abiXName[rd], abiXName[rs1], imm)
}

func daTypeIc(name string, pc uint32, ins uint) string {
	imm, rs1, rd := decodeIa(ins)
	return fmt.Sprintf("%s %s,%d(%s)", name, abiXName[rd], imm, abiXName[rs1])
}

func daTypeId(name string, pc uint32, ins uint) string {
	imm, rs1, rd := decodeIa(ins)
	return fmt.Sprintf("%s %s,%s,0x%x", name, abiXName[rd], abiXName[rs1], imm)
}

func daTypeIe(name string, pc uint32, ins uint) string {
	imm, rs1, rd := decodeIa(ins)
	if imm == 0 && rd == 0 && rs1 == 1 {
		return "ret"
	}
	if rd == 1 {
		if imm == 0 {
			return fmt.Sprintf("%s %s", name, abiXName[rs1])
		}
		return fmt.Sprintf("%s %s,%d", name, abiXName[rs1], imm)
	}
	return fmt.Sprintf("%s %s,%s,%d", name, abiXName[rd], abiXName[rs1], imm)
}

func daTypeIf(name string, pc uint32, ins uint) string {
	imm, rs1, rd := decodeIa(ins)
	if imm == -1 {
		return fmt.Sprintf("not %s,%s", abiXName[rd], abiXName[rs1])
	}
	return fmt.Sprintf("%s %s,%s,%d", name, abiXName[rd], abiXName[rs1], imm)
}

func daTypeIg(name string, pc uint32, ins uint) string {
	imm, rs1, rd := decodeIa(ins)
	return fmt.Sprintf("%s %s,%d(%s)", name, abiFName[rd], imm, abiXName[rs1])
}

func daTypeIh(name string, pc uint32, ins uint) string {
	csr, rs1, rd := decodeIb(ins)
	if csr == csrFCSR {
		if rd == 0 {
			return fmt.Sprintf("fscsr %s", abiXName[rs1])
		}
		return fmt.Sprintf("fscsr %s,%s", abiXName[rd], abiXName[rs1])
	}
	if name == "csrrs" && rs1 == 0 {
		return fmt.Sprintf("csrr %s,%s", abiXName[rd], csrName(csr))
	}
	return fmt.Sprintf("%s %s,%s,%s", name, abiXName[rd], csrName(csr), abiXName[rs1])
}

//-----------------------------------------------------------------------------
// Type U Decodes

func daTypeUa(name string, pc uint32, ins uint) string {
	imm, rd := decodeU(ins)
	return fmt.Sprintf("%s %s,0x%x", name, abiXName[rd], uint(imm)&0xfffff)
}

//-----------------------------------------------------------------------------
// Type S Decodes

func daTypeSa(name string, pc uint32, ins uint) string {
	imm, rs2, rs1 := decodeS(ins)
	return fmt.Sprintf("%s %s,%d(%s)", name, abiXName[rs2], imm, abiXName[rs1])
}

func daTypeSb(name string, pc uint32, ins uint) string {
	imm, rs2, rs1 := decodeS(ins)
	return fmt.Sprintf("%s %s,%d(%s)", name, abiFName[rs2], imm, abiXName[rs1])
}

//-----------------------------------------------------------------------------
// Type R Decodes

func daTypeRa(name string, pc uint32, ins uint) string {
	rs2, rs1, rd := decodeR(ins)
	if name == "sub" && rs1 == 0 {
		return fmt.Sprintf("neg %s,%s", abiXName[rd], abiXName[rs2])
	}
	return fmt.Sprintf("%s %s,%s,%s", name, abiXName[rd], abiXName[rs1], abiXName[rs2])
}

func daTypeRb(name string, pc uint32, ins uint) string {
	rs2, rs1, rd := decodeR(ins)
	if rs2 == 0 {
		return fmt.Sprintf("%s %s,(%s)", name, abiXName[rd], abiXName[rs1])
	}
	return fmt.Sprintf("%s %s,%s,(%s)", name, abiXName[rd], abiXName[rs2], abiXName[rs1])
}

func daTypeRc(name string, pc uint32, ins uint) string {
	rs2, rs1, rd := decodeR(ins)
	return fmt.Sprintf("%s %s,%s,%s", name, abiFName[rd], abiFName[rs1], abiFName[rs2])
}

func daTypeRd(name string, pc uint32, ins uint) string {
	_, rs1, rd := decodeR(ins)
	return fmt.Sprintf("%s %s,%s", name, abiXName[rd], abiFName[rs1])
}

func daTypeRe(name string, pc uint32, ins uint) string {
	_, rs1, rd := decodeR(ins)
	return fmt.Sprintf("%s %s,%s", name, abiFName[rd], abiXName[rs1])
}

//-----------------------------------------------------------------------------
// Type R4 Decodes

func daTypeR4a(name string, pc uint32, ins uint) string {
	rs3, rs2, rs1, _, rd := decodeR4(ins)
	return fmt.Sprintf("%s %s,%s,%s,%s", name, abiFName[rd], abiFName[rs1], abiFName[rs2], abiFName[rs3])
}

//-----------------------------------------------------------------------------
// Type B Decodes

func daTypeBa(name string, pc uint32, ins uint) string {
	imm, rs2, rs1 := decodeB(ins)
	adr := int(pc) + imm

	if rs2 == 0 {
		switch name {
		case "bge", "beq", "bne", "blt":
			return fmt.Sprintf("%sz %s,%x", name, abiXName[rs1], adr)
		}
	}

	return fmt.Sprintf("%s %s,%s,%x", name, abiXName[rs1], abiXName[rs2], adr)
}

//-----------------------------------------------------------------------------
// Type J Decodes

func daTypeJa(name string, pc uint32, ins uint) string {
	imm, rd := decodeJ(ins)
	if rd == 0 {
		return fmt.Sprintf("j %x", int(pc)+imm)
	}
	return fmt.Sprintf("%s %s,%x", name, abiXName[rd], int(pc)+imm)
}

//-----------------------------------------------------------------------------
// Type CI Decodes

func daTypeCIa(name string, pc uint32, ins uint) string {
	imm, rd := decodeCIa(ins)
	return fmt.Sprintf("%s %s,%d", name, abiXName[rd], imm)
}

func daTypeCIb(name string, pc uint32, ins uint) string {
	imm := decodeCIb(ins)
	return fmt.Sprintf("%s sp,sp,%d", name, imm)
}

func daTypeCIc(name string, pc uint32, ins uint) string {
	imm, rd := decodeCIa(ins)
	return fmt.Sprintf("%s %s,%s,%d", name, abiXName[rd], abiXName[rd], imm)
}

func daTypeCId(name string, pc uint32, ins uint) string {
	imm, rd := decodeCIc(ins)
	return fmt.Sprintf("%s %s,%s,0x%x", name, abiXName[rd], abiXName[rd], imm)
}

func daTypeCIe(name string, pc uint32, ins uint) string {
	uimm, rd := decodeCId(ins)
	return fmt.Sprintf("%s %s,%s,0x%x", name, abiXName[rd], abiXName[rd], uimm)
}

func daTypeCIf(name string, pc uint32, ins uint) string {
	imm, rd := decodeCIe(ins)
	return fmt.Sprintf("%s %s,%s,%d", name, abiXName[rd], abiXName[rd], imm)
}

func daTypeCIg(name string, pc uint32, ins uint) string {
	imm, rd := decodeCIf(ins)
	return fmt.Sprintf("%s %s,0x%x", name, abiXName[rd], imm)
}

//-----------------------------------------------------------------------------
// Type CIW Decodes

func daTypeCIWa(name string, pc uint32, ins uint) string {
	return "?"
}

func daTypeCIWb(name string, pc uint32, ins uint) string {
	uimm, rd := decodeCIW(ins)
	return fmt.Sprintf("%s %s,sp,%d", name, abiXName[rd], uimm)
}

//-----------------------------------------------------------------------------
// Type CJ Decodes

func daTypeCJa(name string, pc uint32, ins uint) string {
	rs1 := decodeCJa(ins)
	if rs1 == 1 {
		return "ret"
	}
	return fmt.Sprintf("%s %s", name, abiXName[rs1])
}

func daTypeCJb(name string, pc uint32, ins uint) string {
	imm := decodeCJb(ins)
	return fmt.Sprintf("%s %x", name, int(pc)+imm)
}

func daTypeCJc(name string, pc uint32, ins uint) string {
	imm := decodeCJb(ins)
	return fmt.Sprintf("%s ra,%x", name, int(pc)+imm)
}

//-----------------------------------------------------------------------------
// Type CR Decodes

func daTypeCRa(name string, pc uint32, ins uint) string {
	rd, rs := decodeCR(ins)
	return fmt.Sprintf("%s %s,%s", name, abiXName[rd], abiXName[rs])
}

func daTypeCRb(name string, pc uint32, ins uint) string {
	rd, rs := decodeCR(ins)
	return fmt.Sprintf("%s %s,%s,%s", name, abiXName[rd], abiXName[rd], abiXName[rs])
}

//-----------------------------------------------------------------------------
// Type CSS Decodes

func daTypeCSSa(name string, pc uint32, ins uint) string {
	imm, rd := decodeCSSa(ins)
	return fmt.Sprintf("%s %s,%d(sp)", name, abiXName[rd], imm)
}

func daTypeCSSb(name string, pc uint32, ins uint) string {
	imm, rs2 := decodeCSSb(ins)
	return fmt.Sprintf("%s %s,%d(sp)", name, abiXName[rs2], imm)
}

//-----------------------------------------------------------------------------
// Type CB Decodes

func daTypeCBa(name string, pc uint32, ins uint) string {
	imm, rs := decodeCB(ins)
	return fmt.Sprintf("%s %s,%x", name, abiXName[rs], int(pc)+imm)
}

//-----------------------------------------------------------------------------

// SymbolTable maps an address to a symbol.
type SymbolTable map[uint32]string

// Disassembly returns the result of the disassembler call.
type Disassembly struct {
	Dump     string // address and memory bytes
	Symbol   string // symbol for the address (if any)
	Assembly string // assembly instructions
	Length   int    // length in bytes of decode
}

func (da *Disassembly) String() string {
	return fmt.Sprintf("%-16s %8s %-13s", da.Dump, da.Symbol, da.Assembly)
}

//-----------------------------------------------------------------------------

func daDump32(pc uint32, ins uint) string {
	return fmt.Sprintf("%08x: %08x", pc, uint32(ins))
}

func daDump16(pc uint32, ins uint) string {
	return fmt.Sprintf("%08x: %04x    ", pc, uint16(ins))
}

func daSymbol(adr uint32, st SymbolTable) string {
	if st != nil {
		return st[adr]
	}
	return ""
}

// daInstruction returns the disassembly for a 16/32-bit instruction.
func (isa *ISA) daInstruction(pc uint32, ins uint) string {
	im := isa.lookup(ins)
	if im != nil {
		return im.defn.da(im.name, pc, ins)
	}
	return "?"
}

//-----------------------------------------------------------------------------

// Disassemble a RISC-V instruction at the address.
func (m *RV32) Disassemble(adr uint32) *Disassembly {
	ins, _ := m.Mem.RdIns(uint(adr))
	var da Disassembly
	da.Symbol = m.Mem.Symbol(uint(adr))
	if ins&3 == 3 {
		da.Dump = daDump32(adr, ins)
		da.Assembly = m.isa.daInstruction(adr, ins)
		da.Length = 4
	} else {
		da.Dump = daDump16(adr, ins)
		da.Assembly = m.isa.daInstruction(adr, ins)
		da.Length = 2
	}
	return &da
}

//-----------------------------------------------------------------------------
