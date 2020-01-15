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
	"fmt"
	"strings"

	"github.com/deadsy/riscv/csr"
)

//-----------------------------------------------------------------------------

// Pagesize is 4KiB
const riscvPageShift = 12
const riscvPageMask = (1 << riscvPageShift) - 1

//-----------------------------------------------------------------------------
// Common PTE functions (the first 9 bits are the same for all PTEs)

func pteIsValid(pte uint) bool {
	// v == 1 and WR != 10
	return (pte&1) == 1 && ((pte>>1)&3) != 2
}

func pteIsPointer(pte uint) bool {
	// XWRV == 0001
	return pte&15 == 1
}

// pteGetUser gets the PTE user flag.
func pteGetUser(pte uint) bool {
	return (pte & (1 << 4 /*U*/)) != 0
}

// pteGetAccess gets the PTE access flag.
func pteGetAccess(pte uint) bool {
	return pte&(1<<6 /*A*/) != 0
}

// pteGetDirty gets the PTE dirty flag.
func pteGetDirty(pte uint) bool {
	return pte&(1<<7 /*D*/) != 0
}

// pteCanRead returns true if the PTE indicates read permission for the page.
func pteCanRead(pte uint) bool {
	return (pte & (1 << 1 /*R*/)) != 0
}

// pteCanWrite returns true if the PTE indicates write permission for the page.
func pteCanWrite(pte uint) bool {
	return (pte & (1 << 2 /*W*/)) != 0
}

// pteCanExec returns true if the PTE indicates execute permission for the page.
func pteCanExec(pte uint) bool {
	return (pte & (1 << 3 /*X*/)) != 0
}

// pteSetRead sets the PTE read bit.
func pteSetRead(pte uint) uint {
	return pte | (1 << 1 /*R*/)
}

// pteSetAccess sets the PTE access bit.
func pteSetAccess(pte uint) uint {
	return pte | (1 << 6 /*A*/)
}

// pteSetDirty sets the PTE dirty bit.
func pteSetDirty(pte uint) uint {
	return pte | (1 << 7 /*D*/)
}

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
func (m *Memory) va2pa(va uint, attr Attribute) (uint, error) {

	// If mstatus.MPRV == 1 then mode = mstatus.MPP
	var mode csr.Mode
	if m.csr.GetMPRV() {
		mode = m.csr.GetMPP()
	} else {
		mode = m.csr.GetMode()
	}

	var pa uint
	var err error

	// get the vm
	vm := m.csr.GetVM()
	if mode == csr.ModeM {
		// machine mode va == pa
		vm = csr.Bare
	}

	// run the va to pa mapping
	switch vm {
	case csr.Bare:
		pa, _, err = m.bare(va, mode, attr, false)
	case csr.SV32:
		pa, _, err = m.sv32(sv32va(va), mode, attr, false)
	case csr.SV39:
		pa, _, err = m.sv39(sv39va(va), mode, attr, false)
	case csr.SV48:
		pa, _, err = m.sv48(sv48va(va), mode, attr, false)
	default:
		err = fmt.Errorf("%s not implmented", vm)
	}

	return pa, err
}

//-----------------------------------------------------------------------------

// PageTableWalk returns a string annotating the va->pa page table walk.
func (m *Memory) PageTableWalk(va uint, mode csr.Mode, attr Attribute) string {

	// get the vm
	vm := m.csr.GetVM()
	if mode == csr.ModeM {
		vm = csr.Bare
	}

	var s []string
	var err error

	// run the va to pa mapping
	switch vm {
	case csr.Bare:
		_, s, err = m.bare(va, mode, attr, true)
	case csr.SV32:
		_, s, err = m.sv32(sv32va(va), mode, attr, true)
	case csr.SV39:
		_, s, err = m.sv39(sv39va(va), mode, attr, true)
	case csr.SV48:
		_, s, err = m.sv48(sv48va(va), mode, attr, true)
	default:
		return fmt.Sprintf("%s not implemented", vm)
	}
	if err != nil {
		s = append(s, err.Error())
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
