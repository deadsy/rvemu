//-----------------------------------------------------------------------------
/*

to/from host functions

Compliance testing communicates with the "host" via reads/writes to memory
locations with symbols "tohost" and "fromhost". This code intercepts this
acess via the breakpoint mechanism and records interesting state.

*/
//-----------------------------------------------------------------------------

package host

import (
	"fmt"
	"strings"

	"github.com/deadsy/riscv/mem"
)

//-----------------------------------------------------------------------------

// Host stores the state of writes to the "tohost" location.
type Host struct {
	mem    *mem.Memory // memory sub-system
	addr   uint        // base address of "tohost" symbol
	status uint
	text   []rune // characters written "tohost"
}

// NewHost returns a Host empty structure.
func NewHost(mem *mem.Memory, addr uint) *Host {
	return &Host{
		mem:  mem,
		addr: addr,
		text: make([]rune, 0, 128),
	}
}

func (h *Host) String() string {
	s := []string{}
	s = append(s, fmt.Sprintf("(%d)", h.status))
	if h.text != nil {
		s = append(s, strings.TrimSpace(string(h.text)))
	}
	return strings.Join(s, " ")
}

// Passed returns if the compliance test has passed.
func (h *Host) Passed() bool {
	return h.status == 1
}

//-----------------------------------------------------------------------------

// The risc-v compliance tests use this as the upper word in "tohost" when
// they want to send a character to the host.
const charMagic = 0x01010000

// To64 intercepts writes to a 64-bit tohost memory location.
func (h *Host) To64(bp *mem.BreakPoint) bool {
	// get the tohost value
	val, _ := h.mem.Rd64Phys(h.addr)
	// Is this a character write?
	if (val >> 32) == charMagic {
		h.text = append(h.text, rune(val&0xff))
		// signal write consumption to the risc-v test.
		h.mem.Wr64Phys(h.addr, 0)
		// no break
		return false
	}
	h.status = uint(val)
	// break
	return true
}

//-----------------------------------------------------------------------------

// To32 intercepts writes to a 32-bit tohost memory location.
func (h *Host) To32(bp *mem.BreakPoint) bool {
	// get the tohost value
	val, _ := h.mem.Rd32Phys(h.addr)
	h.status = uint(val)
	// break
	return true
}

//-----------------------------------------------------------------------------
