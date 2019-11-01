//-----------------------------------------------------------------------------
/*

RISC-V Disassembler

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

type daFunc func(m *RV32, ins uint32) (*Disassembly, error)

type linearDecode struct {
	val  uint32
	mask uint32
	fn   daFunc
}

//-----------------------------------------------------------------------------

// Disassemble disassembles a RISC-V instruction.
func (m *RV32) Disassemble(ins uint32) (*Disassembly, error) {
	for i := range m.daDecode {
		d := &m.daDecode[i]
		if ins&d.mask == d.val {
			return d.fn(m, ins)
		}
	}
	return nil, nil
}

//-----------------------------------------------------------------------------
