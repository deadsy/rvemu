//-----------------------------------------------------------------------------
/*

RISC-V RV32 Disassembler

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
func (m *RV32) Disassemble(ins uint32) *Disassembly {

	return nil
}

//-----------------------------------------------------------------------------

type daFunc func(m *RV32, ins uint32) (*Disassembly, error)

var daTable00 = map[uint32]daFunc{}

func da00(m *RV32, ins uint32) (*Disassembly, error) {
	mask := uint32(0x7f)
	if da, ok := daTable00[ins&mask]; ok {
		return da(m, ins & ^mask)
	}
	return nil, nil
}

//-----------------------------------------------------------------------------
