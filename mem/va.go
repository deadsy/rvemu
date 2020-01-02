//-----------------------------------------------------------------------------
/*

Virtual Address Translation

SV32: 32-bit VA maps to 34-bit PA
SV39: 39-bit VA maps to 56-bit PA
SV48: 48-bit VA maps to 56-bit PA
SV57: 57-bit VA maps to ?-bit PA
SV64: 64-bit VA maps to ?-bit PA

*/
//-----------------------------------------------------------------------------

package mem

import (
	"errors"
	"fmt"
	"strings"

	"github.com/deadsy/riscv/csr"
)

//-----------------------------------------------------------------------------

const riscvPageShift = 12

//-----------------------------------------------------------------------------

// bare - no translation
func (m *Memory) bare(va uint, mode csr.Mode, attr Attribute, debug bool) (uint, []string, error) {
	dbg := []string{}
	if debug {
		dbg = append(dbg, fmt.Sprintf("va   %08x", va))
		dbg = append(dbg, fmt.Sprintf("satp %s", csr.DisplaySATP(m.csr)))
		dbg = append(dbg, fmt.Sprintf("pa   %08x", va))
	}
	return va, dbg, nil
}

//-----------------------------------------------------------------------------

// va2pa translates a virtual address to a physical address.
func (m *Memory) va2pa(va uint, attr Attribute, debug bool) (uint, []string, error) {

	// If mstatus.MPRV == 1 then mode = mstatus.MPP
	var mode csr.Mode
	if m.csr.GetMPRV() {
		mode = m.csr.GetMPP()
	} else {
		mode = m.csr.GetMode()
	}

	if mode == csr.ModeM {
		// machine mode va == pa
		return m.bare(va, mode, attr, debug)
	}

	switch m.csr.GetVM() {
	case csr.Bare:
		return m.bare(va, mode, attr, debug)
	case csr.SV32:
		return m.sv32(sv32va(va), mode, attr, debug)
	case csr.SV39:
		return 0, nil, nil
	case csr.SV48:
		return 0, nil, nil
	case csr.SV57:
		return 0, nil, nil
	case csr.SV64:
		return 0, nil, nil
	}
	return 0, nil, errors.New("unknown vm mode")
}

//-----------------------------------------------------------------------------

// PageTableWalk returns a string annotating the va->pa page table walk.
func (m *Memory) PageTableWalk(va uint, attr Attribute) string {
	_, s, err := m.va2pa(va, attr, true)
	if err != nil {
		s = append(s, err.Error())
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
