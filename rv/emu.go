//-----------------------------------------------------------------------------
/*

RISC-V Emulator

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
	"math"
)

//-----------------------------------------------------------------------------
// default emulation

func emuNone(m *RV32, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv32i

func emuLUI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuAUIPC(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuJAL(m *RV32, ins uint) {
	imm, rd := decodeJ(ins)
	m.wrX(rd, m.PC+4)
	m.PC = uint32(int(m.PC) + int(imm))
}

func emuJALR(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuBEQ(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuBNE(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuBLT(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuBGE(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuBLTU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuBGEU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuLB(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuLH(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuLW(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuLBU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuLHU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSB(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSH(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSW(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuADDI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSLTI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSLTIU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuXORI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuORI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuANDI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSLLI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSRLI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSRAI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuADD(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSUB(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSLL(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSLT(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSLTU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuXOR(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSRL(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuSRA(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuOR(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuAND(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuFENCE(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuFENCExI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuECALL(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuEBREAK(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuCSRRW(m *RV32, ins uint) {
	csr, rs1, rd := decodeIb(ins)
	t := m.rdCSR(csr)
	m.wrCSR(csr, m.X[rs1])
	m.wrX(rd, t)
	m.PC += 4
}

func emuCSRRS(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuCSRRC(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuCSRRWI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuCSRRSI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emuCSRRCI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv32f

func emuFMVxWxX(m *RV32, ins uint) {
	_, rs1, rd := decodeR(ins)
	m.F[rd] = math.Float32frombits(m.X[rs1])
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv32c

func emuCxLI(m *RV32, ins uint) {
	imm, rd := decodeCIa(ins)
	m.wrX(rd, uint32(imm))
	m.PC += 2
}

//-----------------------------------------------------------------------------

// Reset the RV32 CPU.
func (m *RV32) Reset() {
}

// Run the RV32 CPU for a single instruction.
func (m *RV32) Run() error {

	// normal instructions
	ins, _ := m.Mem.RdIns(uint(m.PC))
	im := m.isa.lookup(ins)
	im.defn.emu32(m, ins)

	if m.flag != 0 {
		if m.flag&flagIllegal != 0 {
			return fmt.Errorf("illegal instruction at %08x", m.PC)
		}
		if m.flag&flagExit != 0 {
			return fmt.Errorf("exit at %08x, status %08x", m.PC, m.X[1])
		}
		if m.flag&flagTodo != 0 {
			return fmt.Errorf("unimplemented instruction at %08x", m.PC)
		}
	}

	return nil
}

//-----------------------------------------------------------------------------
