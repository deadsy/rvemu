//-----------------------------------------------------------------------------
/*

RISC-V RV32I Disassembler

*/
//-----------------------------------------------------------------------------

package cpu

//-----------------------------------------------------------------------------

// Disassembly returns the result of the disassembler call.
type Disassembly struct {
	Symbol      string // symbol for the address (if any)
	Instruction string // instruction decode
	Comment     string // useful comment
}

//-----------------------------------------------------------------------------

// Disassemble disassembles a RV32I RISC-V instruction.
func (m *RV32I) Disassemble(ins uint32) *Disassembly {

	return nil
}

//-----------------------------------------------------------------------------
