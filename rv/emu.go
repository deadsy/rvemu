//-----------------------------------------------------------------------------
/*

RISC-V Disassembler

*/
//-----------------------------------------------------------------------------

package rv

import "fmt"

//-----------------------------------------------------------------------------
// default emulation

func emuNone(m *RV32, ins uint) {
	m.todo = true
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

	if m.illegal {
		return fmt.Errorf("illegal instruction at %08x", m.PC)
	}

	if m.exit {
		return fmt.Errorf("exit at %08x, status %08x", m.PC, m.X[1])
	}

	if m.todo {
		return fmt.Errorf("unimplemented instruction at %08x", m.PC)
	}

	return nil
}

//-----------------------------------------------------------------------------
