//-----------------------------------------------------------------------------
/*

to/from host functions

Compliance testing communicates with the "host" via reads/writes to memory
locations with symbols "tohost" and "fromhost". This code intercepts this
acess via the breakpoint mechanism and records interesting state.

*/
//-----------------------------------------------------------------------------

package host

import "github.com/deadsy/riscv/mem"

//-----------------------------------------------------------------------------

type Host struct {
	addr uint   // base address of "tohost" symbol
	text []rune // characters written "tohost"
}

func NewHost(m *mem.Memory) *Host {
	sym := m.SymbolByName("tohost")
	if sym == nil {
		return nil
	}
	return &Host{
		addr: sym.Addr,
		text: make([]rune, 0, 128),
	}
}

// The risc-v compliance tests use this as the upper word in "tohost" when
// they want to send a character to the host.
const charMagic = 0x01010000

// To64 intercepts writes to a 64-bit tohost memory location.
func (h *Host) To64(m *mem.Memory, bp *mem.BreakPoint) bool {
	// break by default unless we consume the write
	brk := true
	val, _ := m.Rd64Phys(h.addr)
	// Is this a character write?
	if (val >> 32) == charMagic {
		h.text = append(h.text, rune(val&0xff))
		// signal write consumption to the risc-v test.
		m.Wr64Phys(h.addr, 0)
		brk = false
	}
	return brk
}

// To32 intercepts writes to a 32-bit tohost memory location.
func (h *Host) To32(m *mem.Memory, bp *mem.BreakPoint) bool {
	return true
}

// GetText returns the string of characters written "tohost".
func (h *Host) GetText() string {
	if h.text == nil {
		return ""
	}
	return string(h.text)
}

//-----------------------------------------------------------------------------
