//-----------------------------------------------------------------------------
/*

Virtual Address Translation

*/
//-----------------------------------------------------------------------------

package mem

import (
	"errors"

	"github.com/deadsy/riscv/util"
)

//-----------------------------------------------------------------------------

type vmMode uint

const (
	vmBare vmMode = iota
	vmSV32
	vmSV39
	vmSV48
	vmSV57
	vmSV64
)

// SetSATP is a callback (from CSR) to set the SATP value in the memory subsystem.
func (m *Memory) SetSATP(satp, sxlen uint) {
	m.satp = satp
	if sxlen == 32 {
		// RV32
		m.vm = [2]vmMode{vmBare, vmSV32}[(satp>>31)&1]
	} else {
		// RV64
		m.vm = map[uint]vmMode{0: vmBare, 8: vmSV39, 9: vmSV48, 10: vmSV57, 11: vmSV64}[(satp>>60)&15]
	}
}

//-----------------------------------------------------------------------------

// pteValid returns true if the PTE is valid
func pteValid(pte uint) bool {
	// check V = 1
	if pte&1 == 0 {
		return false
	}
	// check WR != 10
	return (pte>>1)&3 != 2
}

// ptePointer returns true if the PTE points to a next-level page table.
func ptePointer(pte uint) bool {
	// check XWRV != 0001
	return pte&15 == 1
}

//-----------------------------------------------------------------------------

// bare - no translation
func (m *Memory) bare(va uint, attr Attribute) (uint, error) {
	return va, nil
}

// sv32
func (m *Memory) sv32(va uint, attr Attribute) (uint, error) {

	var vpn [2]uint
	vpn[1] = util.RdBits(va, 31, 22)
	vpn[0] = util.RdBits(va, 21, 12)
	pageOffset := util.RdBits(va, 11, 0)

	// 1. Let a be satp.ppn × PAGESIZE, and let i = LEVELS − 1. (For Sv32, PAGESIZE=4096 and LEVELS=2.)
	a := (m.satp & ((1 << 22) - 1)) << 12
	i := 1

	for true {
		// 2. Let pte be the value of the PTE at address a+va.vpn[i]×PTESIZE. (For Sv32, PTESIZE=4.)
		// If accessing pte violates a PMA or PMP check, raise an access exception corresponding to
		// the original access type.
		x, err := m.Rd32Phys(a + (vpn[i] << 2))
		if err != nil {
			return 0, pageError(va, attr)
		}
		pte := uint(x)

		// 3. If pte.v = 0, or if pte.r = 0 and pte.w = 1, stop and raise a page-fault exception corresponding
		// to the original access type.
		if !pteValid(pte) {
			return 0, pageError(va, attr)
		}

		// 4. Otherwise, the PTE is valid. If pte.r = 1 or pte.x = 1, go to step 5. Otherwise, this PTE is a
		// pointer to the next level of the page table. Let i = i − 1. If i < 0, stop and raise a page-fault
		// exception corresponding to the original access type. Otherwise, let a = pte.ppn × PAGESIZE
		// and go to step 2.
		if !ptePointer(pte) {
			break
		}
		i = i - 1
		if i < 0 {
			return 0, pageError(va, attr)
		}
		var ppn [2]uint
		ppn[0] = util.RdBits(pte, 19, 10)
		ppn[1] = util.RdBits(pte, 31, 20)
		a = ppn[0] << 12
	}

	// 5. A leaf PTE has been found. Determine if the requested memory access is allowed by the
	// pte.r, pte.w, pte.x, and pte.u bits, given the current privilege mode and the value of the
	// SUM and MXR fields of the mstatus register. If not, stop and raise a page-fault exception
	// corresponding to the original access type.

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

// va2pa translates a virtual address to a physical address.
func (m *Memory) va2pa(va uint, attr Attribute) (uint, error) {
	switch m.vm {
	case vmBare:
		return m.bare(va, attr)
	case vmSV32:
		return m.sv32(va, attr)
	case vmSV39:
		return 0, nil
	case vmSV48:
		return 0, nil
	case vmSV57:
		return 0, nil
	case vmSV64:
		return 0, nil
	}
	return 0, errors.New("unknown vm mode")
}

//-----------------------------------------------------------------------------
