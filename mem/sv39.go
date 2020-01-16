//-----------------------------------------------------------------------------
/*

SV39 Virtual Memory Address Translation

39-bit VA maps to 56-bit PA

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"

	"github.com/deadsy/riscv/csr"
	"github.com/deadsy/riscv/util"
)

//-----------------------------------------------------------------------------
// page table entry

type sv39pte uint

func (pte sv39pte) String() string {
	fs := util.FieldSet{
		{"ppn2", 53, 28, util.FmtHex},
		{"ppn1", 27, 19, util.FmtHex},
		{"ppn0", 18, 10, util.FmtHex},
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

func (pte sv39pte) ppn(a, b int) uint {
	hi := [3]uint{18, 27, 53}[a]
	lo := [3]uint{10, 19, 28}[b]
	return util.GetBits(uint(pte), hi, lo)
}

//-----------------------------------------------------------------------------
// virtual address

type sv39va uint

func (va sv39va) String() string {
	fs := util.FieldSet{
		{"vpn2", 38, 30, util.FmtHex},
		{"vpn1", 29, 21, util.FmtHex},
		{"vpn0", 20, 12, util.FmtHex},
		{"ofs", 11, 0, util.FmtHex},
	}
	return fs.Display(uint(va))
}

func (va sv39va) vpn(a, b int) uint {
	hi := [3]uint{20, 29, 38}[a]
	lo := [3]uint{12, 21, 30}[b]
	return util.GetBits(uint(va), hi, lo)
}

func (va sv39va) ofs() uint {
	return uint(va) & riscvPageMask
}

func (va sv39va) pageError(attr Attribute) error {
	return pageError(uint(va), attr)
}

//-----------------------------------------------------------------------------

func (m *Memory) sv39(va sv39va, mode csr.Mode, attr Attribute, debug bool) (uint, []string, error) {
	const levels = 3
	var pteAddr uint
	var pte uint
	dbg := []string{}

	if debug {
		dbg = append(dbg, fmt.Sprintf("%s %s", mode, attr))
		dbg = append(dbg, fmt.Sprintf("va %08x %s", uint(va), va))
		dbg = append(dbg, fmt.Sprintf("satp %s", csr.DisplaySATP(m.csr)))
	}

	// 1. Let baseAddr be satp.ppn × PAGESIZE, and let i = LEVELS − 1. (For SV39, PAGESIZE=4096 and LEVELS=3)
	baseAddr := m.csr.GetPPN() << riscvPageShift
	i := levels - 1

	for true {
		// 2. Let pte be the value of the PTE at address a+va.vpn[i]×PTESIZE. (For SV39, PTESIZE=8)
		// If accessing pte violates a PMA or PMP check, raise an access exception corresponding to
		// the original access type.
		pteAddr = baseAddr + (va.vpn(i, i) << 3)
		x, err := m.Rd64Phys(pteAddr)
		if err != nil {
			return 0, dbg, va.pageError(attr)
		}
		pte = uint(x)

		if debug {
			dbg = append(dbg, fmt.Sprintf("pte%d [%014x] %s", i, pteAddr, sv39pte(pte)))
		}

		// 3. If pte.v = 0, or if pte.r = 0 and pte.w = 1, stop and raise a page-fault exception corresponding
		// to the original access type.
		if !pteIsValid(pte) {
			return 0, dbg, va.pageError(attr)
		}

		// 4. Otherwise, the PTE is valid. If pte.r = 1 or pte.x = 1, go to step 5. Otherwise, this PTE is a
		// pointer to the next level of the page table. Let i = i − 1. If i < 0, stop and raise a page-fault
		// exception corresponding to the original access type. Otherwise, let a = pte.ppn × PAGESIZE
		// and go to step 2.
		if !pteIsPointer(pte) {
			break
		}
		i = i - 1
		if i < 0 {
			return 0, dbg, va.pageError(attr)
		}
		baseAddr = sv39pte(pte).ppn(levels-1, 0) << riscvPageShift
	}

	// 5. A leaf PTE has been found. Determine if the requested memory access is allowed by the
	// pte.r, pte.w, pte.x, and pte.u bits, given the current privilege mode and the value of the
	// SUM and MXR fields of the mstatus register. If not, stop and raise a page-fault exception
	// corresponding to the original access type.

	// If mstatus.MXR == 1 and pte.X == 1 then pte.R = 1
	if m.csr.GetMXR() && pteCanExec(pte) {
		pte = pteSetRead(pte)
	}
	// check the RWX permissions
	if attr&AttrR != 0 && !pteCanRead(pte) {
		return 0, dbg, va.pageError(attr)
	}
	if attr&AttrW != 0 && !pteCanWrite(pte) {
		return 0, dbg, va.pageError(attr)
	}
	if attr&AttrX != 0 && !pteCanExec(pte) {
		return 0, dbg, va.pageError(attr)
	}

	// check user/supervisor mode
	switch mode {
	case csr.ModeU:
		if !pteGetUser(pte) {
			return 0, dbg, va.pageError(attr)
		}
	case csr.ModeS:
		if pteGetUser(pte) {
			if !m.csr.GetSUM() {
				// U == 1 and mstatus.SUM == 0
				return 0, dbg, va.pageError(attr)
			}
			if attr&AttrX != 0 {
				// Irrespective of SUM, the supervisor may not execute code on pages with U=1.
				return 0, dbg, va.pageError(attr)
			}
		}
	}

	// 6. If i > 0 and pte.ppn[i − 1 : 0] != 0, this is a misaligned superpage; stop and raise a page-fault
	// exception corresponding to the original access type.
	if i > 0 && sv39pte(pte).ppn(i-1, 0) != 0 {
		return 0, dbg, va.pageError(attr)
	}

	// 7. If pte.a = 0, or if the memory access is a store and pte.d = 0, either raise a page-fault
	// exception corresponding to the original access type, or:
	// • Set pte.a to 1 and, if the memory access is a store, also set pte.d to 1.
	// • If this access violates a PMA or PMP check, raise an access exception corresponding to
	//   the original access type.
	// • This update and the loading of pte in step 2 must be atomic; in particular, no intervening
	//   store to the PTE may be perceived to have occurred in-between.

	var access, dirty bool
	if !pteGetAccess(pte) {
		access = true
	}
	if attr&AttrW != 0 && !pteGetDirty(pte) {
		dirty = true
	}
	if access || dirty {
		// Note: We may have set the R bit previously, so re-read the pte.
		x, _ := m.Rd64Phys(pteAddr)
		pte := uint(x)
		if access {
			pte = pteSetAccess(pte)
		}
		if dirty {
			pte = pteSetDirty(pte)
		}
		m.Wr64Phys(pteAddr, uint64(pte))
	}

	// 8. The translation is successful. The translated physical address is given as follows:
	// • pa.pgoff = va.pgoff.
	// • If i > 0, then this is a superpage translation and pa.ppn[i − 1 : 0] = va.vpn[i − 1 : 0].
	// • pa.ppn[LEVELS − 1 : i] = pte.ppn[LEVELS − 1 : i].
	pa := sv39pte(pte).ppn(levels-1, i)
	if i == 2 {
		pa = (pa << 18) + va.vpn(1, 0)
	} else if i == 1 {
		pa = (pa << 9) + va.vpn(0, 0)
	}
	pa = (pa << riscvPageShift) + va.ofs()

	if debug {
		dbg = append(dbg, fmt.Sprintf("pa %014x", pa))
	}

	return pa, dbg, nil
}

//-----------------------------------------------------------------------------
