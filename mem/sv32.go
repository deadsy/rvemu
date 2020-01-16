//-----------------------------------------------------------------------------
/*

SV32 Virtual Memory Address Translation

32-bit VA maps to 34-bit PA

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

type sv32pte uint

func (pte sv32pte) String() string {
	fs := util.FieldSet{
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

func (pte sv32pte) ppn(a, b int) uint {
	hi := [2]uint{19, 31}[a]
	lo := [2]uint{10, 20}[b]
	return util.GetBits(uint(pte), hi, lo)
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

func (va sv32va) vpn(a, b int) uint {
	hi := [2]uint{21, 31}[a]
	lo := [2]uint{12, 22}[b]
	return util.GetBits(uint(va), hi, lo)
}

func (va sv32va) ofs() uint {
	return uint(va) & riscvPageMask
}

func (va sv32va) pageError(attr Attribute) error {
	return pageError(uint(va), attr)
}

//-----------------------------------------------------------------------------

func (m *Memory) sv32(va sv32va, mode csr.Mode, attr Attribute, debug bool) (uint, []string, error) {
	const levels = 2
	var pteAddr uint
	var pte uint
	dbg := []string{}

	if debug {
		dbg = append(dbg, fmt.Sprintf("%s %s", mode, attr))
		dbg = append(dbg, fmt.Sprintf("va %08x %s", uint(va), va))
		dbg = append(dbg, fmt.Sprintf("satp %s", csr.DisplaySATP(m.csr)))
	}

	// 1. Let baseAddr be satp.ppn × PAGESIZE, and let i = LEVELS − 1. (For Sv32, PAGESIZE=4096 and LEVELS=2.)
	baseAddr := m.csr.GetPPN() << riscvPageShift
	i := levels - 1

	for true {
		// 2. Let pte be the value of the PTE at address a+va.vpn[i]×PTESIZE. (For SV32, PTESIZE=4.)
		// If accessing pte violates a PMA or PMP check, raise an access exception corresponding to
		// the original access type.
		pteAddr = baseAddr + (va.vpn(i, i) << 2)
		x, err := m.Rd32Phys(pteAddr)
		if err != nil {
			return 0, dbg, va.pageError(attr)
		}
		pte = uint(x)

		if debug {
			dbg = append(dbg, fmt.Sprintf("pte%d [%09x] %s", i, pteAddr, sv32pte(pte)))
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
		baseAddr = sv32pte(pte).ppn(levels-1, 0) << riscvPageShift
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
	if i > 0 && sv32pte(pte).ppn(0, 0) != 0 {
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
		x, _ := m.Rd32Phys(pteAddr)
		pte := uint(x)
		if access {
			pte = pteSetAccess(pte)
		}
		if dirty {
			pte = pteSetDirty(pte)
		}
		m.Wr32Phys(pteAddr, uint32(pte))
	}

	// 8. The translation is successful. The translated physical address is given as follows:
	// • pa.pgoff = va.pgoff.
	// • If i > 0, then this is a superpage translation and pa.ppn[i − 1 : 0] = va.vpn[i − 1 : 0].
	// • pa.ppn[LEVELS − 1 : i] = pte.ppn[LEVELS − 1 : i].
	pa := sv32pte(pte).ppn(levels-1, i)
	if i == 1 {
		pa = (pa << 10) + va.vpn(0, 0)
	}
	pa = (pa << riscvPageShift) + va.ofs()

	if debug {
		dbg = append(dbg, fmt.Sprintf("pa %09x", pa))
	}

	return pa, dbg, nil
}

//-----------------------------------------------------------------------------
