//-----------------------------------------------------------------------------
/*

Virtual Address Translation

*/
//-----------------------------------------------------------------------------

package mem

import (
	"errors"

	"github.com/deadsy/riscv/csr"
)

//-----------------------------------------------------------------------------

const riscvPageShift = 12

//-----------------------------------------------------------------------------

// bare - no translation
func (m *Memory) bare(va uint, mode csr.Mode, attr Attribute) (uint, error) {
	return va, nil
}

//-----------------------------------------------------------------------------
// SV39: 39-bit VA maps to 56-bit PA

//-----------------------------------------------------------------------------
// SV48: 48-bit VA maps to 56-bit PA

//-----------------------------------------------------------------------------
// SV57: 57-bit VA maps to ?-bit PA

//-----------------------------------------------------------------------------
// SV64: 64-bit VA maps to ?-bit PA

//-----------------------------------------------------------------------------

// va2pa translates a virtual address to a physical address.
func (m *Memory) va2pa(va uint, attr Attribute) (uint, error) {

	// If mstatus.MPRV == 1 then mode = mstatus.MPP
	var mode csr.Mode
	if m.csr.GetMPRV() {
		mode = m.csr.GetMPP()
	} else {
		mode = m.csr.GetMode()
	}

	if mode == csr.ModeM {
		// machine mode va == pa
		return m.bare(va, mode, attr)
	}

	switch m.csr.GetVM() {
	case csr.Bare:
		return m.bare(va, mode, attr)
	case csr.SV32:
		return m.sv32(sv32va(va), mode, attr)
	case csr.SV39:
		return 0, nil
	case csr.SV48:
		return 0, nil
	case csr.SV57:
		return 0, nil
	case csr.SV64:
		return 0, nil
	}
	return 0, errors.New("unknown vm mode")
}

//-----------------------------------------------------------------------------
