//-----------------------------------------------------------------------------
/*

Virtual Address Translation

*/
//-----------------------------------------------------------------------------

package mem

import (
	"errors"

	"github.com/deadsy/riscv/csr"
	"github.com/deadsy/riscv/util"
)

//-----------------------------------------------------------------------------

// bare - no translation
func (m *Memory) bare(va uint, mode csr.Mode, attr Attribute) (uint, error) {
	return va, nil
}

//-----------------------------------------------------------------------------
// SV32: 32-bit VA maps to 34-bit PA

type sv32pte uint

// isPointer returns true if the PTE points to a next-level page table.
func (pte sv32pte) isPointer() bool {
	// XWRV == 0001
	return pte&15 == 1
}

// isValid returns true if the PTE is valid.
func (pte sv32pte) isValid() bool {
	// v == 1 and WR != 10
	return (pte&1) == 1 && ((pte>>1)&3) != 2
}

// ppn returns the physical page number from the PTE.
func (pte sv32pte) ppn(n int) uint {
	if n == 0 {
		return util.RdBits(uint(pte), 19, 10)
	}
	return util.RdBits(uint(pte), 31, 20)
}

func (pte sv32pte) canRead() bool {
	return (pte & (1 << 1 /*R*/)) != 0
}

func (pte sv32pte) setRead() {
	pte |= (1 << 1 /*R*/)
}

func (pte sv32pte) canWrite() bool {
	return (pte & (1 << 2 /*W*/)) != 0
}

func (pte sv32pte) canExec() bool {
	return (pte & (1 << 3 /*X*/)) != 0
}

func (pte sv32pte) userMode() bool {
	return (pte & (1 << 4 /*U*/)) != 0
}

func (m *Memory) sv32(va uint, mode csr.Mode, attr Attribute) (uint, error) {
	var pte sv32pte

	var vpn [2]uint
	vpn[1] = util.RdBits(va, 31, 22)
	vpn[0] = util.RdBits(va, 21, 12)
	pageOffset := util.RdBits(va, 11, 0)

	// 1. Let a be satp.ppn × PAGESIZE, and let i = LEVELS − 1. (For Sv32, PAGESIZE=4096 and LEVELS=2.)
	a := m.csr.GetPPN() << 12
	i := 1

	for true {
		// 2. Let pte be the value of the PTE at address a+va.vpn[i]×PTESIZE. (For Sv32, PTESIZE=4.)
		// If accessing pte violates a PMA or PMP check, raise an access exception corresponding to
		// the original access type.
		x, err := m.Rd32Phys(a + (vpn[i] << 2))
		if err != nil {
			return 0, pageError(va, attr)
		}
		pte = sv32pte(x)

		// 3. If pte.v = 0, or if pte.r = 0 and pte.w = 1, stop and raise a page-fault exception corresponding
		// to the original access type.
		if !pte.isValid() {
			return 0, pageError(va, attr)
		}

		// 4. Otherwise, the PTE is valid. If pte.r = 1 or pte.x = 1, go to step 5. Otherwise, this PTE is a
		// pointer to the next level of the page table. Let i = i − 1. If i < 0, stop and raise a page-fault
		// exception corresponding to the original access type. Otherwise, let a = pte.ppn × PAGESIZE
		// and go to step 2.
		if !pte.isPointer() {
			break
		}
		i = i - 1
		if i < 0 {
			return 0, pageError(va, attr)
		}

		a = pte.ppn(0) << 12
	}

	// 5. A leaf PTE has been found. Determine if the requested memory access is allowed by the
	// pte.r, pte.w, pte.x, and pte.u bits, given the current privilege mode and the value of the
	// SUM and MXR fields of the mstatus register. If not, stop and raise a page-fault exception
	// corresponding to the original access type.

	// If mstatus.MXR == 1 and pte.X == 1 then pte.R = 1
	if m.csr.GetMXR() && pte.canExec() {
		pte.setRead()
	}
	// check the RWX permissions
	if attr&AttrR != 0 && pte.canRead() {
		return 0, pageError(va, attr)
	}
	if attr&AttrW != 0 && pte.canWrite() {
		return 0, pageError(va, attr)
	}
	if attr&AttrX != 0 && pte.canExec() {
		return 0, pageError(va, attr)
	}
	// check user/supervisor mode
	switch mode {
	case csr.ModeU:
		if !pte.userMode() {
			return 0, pageError(va, attr)
		}
	case csr.ModeS:
		if pte.userMode() {
			if !m.csr.GetSUM() {
				// U == 1 and mstatus.SUM == 0
				return 0, pageError(va, attr)
			}
			if attr&AttrX != 0 {
				// Irrespective of SUM, the supervisor may not execute code on pages with U=1.
				return 0, pageError(va, attr)
			}
		}
	}

	// 6. If i > 0 and pte.ppn[i − 1 : 0] != 0, this is a misaligned superpage; stop and raise a page-fault
	// exception corresponding to the original access type.

	// 7. If pte.a = 0, or if the memory access is a store and pte.d = 0, either raise a page-fault
	// exception corresponding to the original access type, or:
	// • Set pte.a to 1 and, if the memory access is a store, also set pte.d to 1.
	// • If this access violates a PMA or PMP check, raise an access exception corresponding to
	//   the original access type.
	// • This update and the loading of pte in step 2 must be atomic; in particular, no intervening
	//   store to the PTE may be perceived to have occurred in-between.

	// 8. The translation is successful. The translated physical address is given as follows:
	// • pa.pgoff = va.pgoff.
	// • If i > 0, then this is a superpage translation and pa.ppn[i − 1 : 0] = va.vpn[i − 1 : 0].
	// • pa.ppn[LEVELS − 1 : i] = pte.ppn[LEVELS − 1 : i].

	_ = pageOffset

	return 0, nil
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
		return m.sv32(va, mode, attr)
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
