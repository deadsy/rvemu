//-----------------------------------------------------------------------------
/*

SV32 Virtual Memory Address Translation

32-bit VA maps to 34-bit PA

*/
//-----------------------------------------------------------------------------

package mem

import (
	"github.com/deadsy/riscv/csr"
	"github.com/deadsy/riscv/util"
)

//-----------------------------------------------------------------------------
// page table entry

type sv32pte uint

func (pte sv32pte) String() string {
	fs := util.FieldSet{
		{"raw", 31, 0, util.FmtHex8},
		{"ppn", 31, 10, util.FmtHex},
		{"ppn1", 31, 20, util.FmtHex},
		{"ppn0", 19, 10, util.FmtHex},
		{"rsw", 9, 8, util.FmtHex},
		{"d", 7, 7, util.FmtHex},
		{"a", 6, 6, util.FmtHex},
		{"g", 5, 5, util.FmtHex},
		{"u", 4, 4, util.FmtHex},
		{"x", 3, 3, util.FmtHex},
		{"w", 2, 2, util.FmtHex},
		{"r", 1, 1, util.FmtHex},
		{"v", 0, 0, util.FmtHex},
	}
	return fs.Display(uint(pte))
}

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

// ppn returns the full physical page number from the PTE.
func (pte sv32pte) ppn() uint {
	return util.RdBits(uint(pte), 31, 10)
}

// ppn1 returns physical page number 1 from the PTE.
func (pte sv32pte) ppn1() uint {
	return util.RdBits(uint(pte), 31, 20)
}

// ppn0 returns physical page number 0 from the PTE.
func (pte sv32pte) ppn0() uint {
	return util.RdBits(uint(pte), 19, 10)
}

// getUser gets the PTE user flag.
func (pte sv32pte) getUser() bool {
	return (pte & (1 << 4 /*U*/)) != 0
}

// getAccess gets the PTE access flag.
func (pte sv32pte) getAccess() bool {
	return pte&(1<<6 /*A*/) != 0
}

// getDirty gets the PTE dirty flag.
func (pte sv32pte) getDirty() bool {
	return pte&(1<<7 /*D*/) != 0
}

// setRead sets the PTE read bit.
func (pte sv32pte) setRead() {
	pte |= (1 << 1 /*R*/)
}

// setAccess sets the PTE access bit.
func (pte sv32pte) setAccess() {
	pte |= (1 << 6 /*A*/)
}

// setAccess sets the PTE dirty bit.
func (pte sv32pte) setDirty() {
	pte |= (1 << 7 /*D*/)
}

// canRead returns true if the PTE indicates read permission for the page.
func (pte sv32pte) canRead() bool {
	return (pte & (1 << 1 /*R*/)) != 0
}

// canWrite returns true if the PTE indicates write permission for the page.
func (pte sv32pte) canWrite() bool {
	return (pte & (1 << 2 /*W*/)) != 0
}

// canExec returns true if the PTE indicates execute permission for the page.
func (pte sv32pte) canExec() bool {
	return (pte & (1 << 3 /*X*/)) != 0
}

//-----------------------------------------------------------------------------
// virtual address

type sv32va uint

func (va sv32va) String() string {
	fs := util.FieldSet{
		{"vpn1", 31, 22, util.FmtHex},
		{"vpn0", 21, 12, util.FmtHex},
		{"ofs", 11, 0, util.FmtHex},
	}
	return fs.Display(uint(va))
}

func (va sv32va) vpn(i int) uint {
	if i == 1 {
		return util.RdBits(uint(va), 31, 22)
	}
	return util.RdBits(uint(va), 21, 12)
}

func (va sv32va) ofs() uint {
	return util.RdBits(uint(va), 11, 0)
}

func (va sv32va) pageError(attr Attribute) error {
	return pageError(uint(va), attr)
}

//-----------------------------------------------------------------------------

func (m *Memory) sv32(va sv32va, mode csr.Mode, attr Attribute) (uint, error) {
	var pte sv32pte
	var pteAddr uint

	//fmt.Printf("va %s attr %s\n", va, attr)

	// 1. Let a be satp.ppn × PAGESIZE, and let i = LEVELS − 1. (For Sv32, PAGESIZE=4096 and LEVELS=2.)
	a := m.csr.GetPPN() << riscvPageShift
	i := 1

	for true {
		// 2. Let pte be the value of the PTE at address a+va.vpn[i]×PTESIZE. (For Sv32, PTESIZE=4.)
		// If accessing pte violates a PMA or PMP check, raise an access exception corresponding to
		// the original access type.
		pteAddr = a + (va.vpn(i) << 2)
		x, err := m.Rd32Phys(pteAddr)
		if err != nil {
			return 0, va.pageError(attr)
		}
		pte = sv32pte(x)

		//fmt.Printf("%d addr %08x pte %s\n", i, pteAddr, pte)

		// 3. If pte.v = 0, or if pte.r = 0 and pte.w = 1, stop and raise a page-fault exception corresponding
		// to the original access type.
		if !pte.isValid() {
			return 0, va.pageError(attr)
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
			return 0, va.pageError(attr)
		}
		a = pte.ppn() << riscvPageShift
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
	if attr&AttrR != 0 && !pte.canRead() {
		return 0, va.pageError(attr)
	}
	if attr&AttrW != 0 && !pte.canWrite() {
		return 0, va.pageError(attr)
	}
	if attr&AttrX != 0 && !pte.canExec() {
		return 0, va.pageError(attr)
	}

	// check user/supervisor mode
	switch mode {
	case csr.ModeU:
		if !pte.getUser() {
			return 0, va.pageError(attr)
		}
	case csr.ModeS:
		if pte.getUser() {
			if !m.csr.GetSUM() {
				// U == 1 and mstatus.SUM == 0
				return 0, va.pageError(attr)
			}
			if attr&AttrX != 0 {
				// Irrespective of SUM, the supervisor may not execute code on pages with U=1.
				return 0, va.pageError(attr)
			}
		}
	}

	// 6. If i > 0 and pte.ppn[i − 1 : 0] != 0, this is a misaligned superpage; stop and raise a page-fault
	// exception corresponding to the original access type.
	if i > 0 && pte.ppn0() != 0 {
		return 0, va.pageError(attr)
	}

	// 7. If pte.a = 0, or if the memory access is a store and pte.d = 0, either raise a page-fault
	// exception corresponding to the original access type, or:
	// • Set pte.a to 1 and, if the memory access is a store, also set pte.d to 1.
	// • If this access violates a PMA or PMP check, raise an access exception corresponding to
	//   the original access type.
	// • This update and the loading of pte in step 2 must be atomic; in particular, no intervening
	//   store to the PTE may be perceived to have occurred in-between.

	var pteUpdate bool
	if !pte.getAccess() {
		pte.setAccess()
		pteUpdate = true
	}
	if attr&AttrW != 0 && !pte.getDirty() {
		pte.setDirty()
		pteUpdate = true
	}
	if pteUpdate {
		err := m.Wr32Phys(pteAddr, uint32(pte))
		if err != nil {
			return 0, va.pageError(attr)
		}
	}

	// 8. The translation is successful. The translated physical address is given as follows:
	// • pa.pgoff = va.pgoff.
	// • If i > 0, then this is a superpage translation and pa.ppn[i − 1 : 0] = va.vpn[i − 1 : 0].
	// • pa.ppn[LEVELS − 1 : i] = pte.ppn[LEVELS − 1 : i].
	pa := va.ofs()
	if i > 0 {
		pa += pte.ppn1() << 22
		pa += va.vpn(0) << riscvPageShift
	} else {
		pa += pte.ppn() << riscvPageShift
	}

	//fmt.Printf("pa %09x\n", pa)
	return pa, nil
}

//-----------------------------------------------------------------------------
