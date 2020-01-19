//-----------------------------------------------------------------------------
/*

Emulated Target Memory

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"

	"github.com/deadsy/riscv/csr"
)

//-----------------------------------------------------------------------------

// Memory is emulated target memory.
type Memory struct {
	Entry     uint64               // entry point from ELF
	brk       error                // pending breakpoint
	bp        map[uint]*BreakPoint // break points
	alen      uint                 // address bit length
	csr       *csr.State           // CSR state
	region    []Region             // memory regions
	symByAddr map[uint]*Symbol     // symbol table by address
	symByName map[string]*Symbol   // symbol table by name
	noMemory  Region               // empty memory region
}

// newMemory returns a memory object.
func newMemory(alen uint, csr *csr.State, empty Attribute) *Memory {
	return &Memory{
		bp:        make(map[uint]*BreakPoint),
		alen:      alen,
		csr:       csr,
		region:    make([]Region, 0),
		symByAddr: make(map[uint]*Symbol),
		symByName: make(map[string]*Symbol),
		noMemory:  newEmpty(empty),
	}
}

// NewMem32 returns memory with a 32-bit address bus.
func NewMem32(csr *csr.State, empty Attribute) *Memory {
	return newMemory(32, csr, empty)
}

// NewMem64 returns memory with a 64-bit address bus.
func NewMem64(csr *csr.State, empty Attribute) *Memory {
	return newMemory(64, csr, empty)
}

func addrStr(addr, alen uint) string {
	if alen == 32 {
		return fmt.Sprintf("%08x", addr)
	}
	return fmt.Sprintf("%016x", addr)
}

// AddrStr returns a string for the address.
func (m *Memory) AddrStr(addr uint) string {
	return addrStr(addr, m.alen)
}

// Add a memory region to the memory.
func (m *Memory) Add(r Region) {
	m.region = append(m.region, r)
}

// findByName returns the memory region by name.
func (m *Memory) findByName(name string) Region {
	for _, r := range m.region {
		if r.Info().name == name {
			return r
		}
	}
	return nil
}

// findByAddr returns the memory region by address.
func (m *Memory) findByAddr(addr, size uint) Region {
	for _, r := range m.region {
		if r.In(addr, size) {
			return r
		}
	}
	return m.noMemory
}

//-----------------------------------------------------------------------------

// SetAttr sets the attributes of a named region.
func (m *Memory) SetAttr(name string, attr Attribute) error {
	r := m.findByName(name)
	if r == nil {
		return fmt.Errorf("no region named \"%s\"", name)
	}
	r.SetAttr(attr)
	return nil
}

// GetSectionName returns the name of the memory section containing the address.
func (m *Memory) GetSectionName(adr uint) string {
	return m.findByAddr(adr, 1).Info().name
}

//-----------------------------------------------------------------------------
// Physical Address Read Functions

// RdInsPhys reads a 32-bit instruction from memory.
func (m *Memory) RdInsPhys(pa uint) (uint, error) {
	return m.findByAddr(pa, 4).RdIns(pa)
}

// Rd64Phys reads a 64-bit data value from memory.
func (m *Memory) Rd64Phys(pa uint) (uint64, error) {
	return m.findByAddr(pa, 8).Rd64(pa)
}

// Rd32Phys reads a 32-bit data value from memory.
func (m *Memory) Rd32Phys(pa uint) (uint32, error) {
	return m.findByAddr(pa, 4).Rd32(pa)
}

// Rd16Phys reads a 16-bit data value from memory.
func (m *Memory) Rd16Phys(pa uint) (uint16, error) {
	return m.findByAddr(pa, 2).Rd16(pa)
}

// Rd8Phys reads an 8-bit data value from memory.
func (m *Memory) Rd8Phys(pa uint) (uint8, error) {
	return m.findByAddr(pa, 1).Rd8(pa)
}

//-----------------------------------------------------------------------------
// Virtual Address Read Functions

// RdIns reads a 32-bit instruction from memory.
func (m *Memory) RdIns(va uint) (uint, error) {
	pa, err := m.va2pa(va, AttrX)
	if err != nil {
		return 0, err
	}
	val, err := m.RdInsPhys(pa)
	m.monitor(pa, 4, AttrX)
	return val, err
}

// Rd64 reads a 64-bit data value from memory.
func (m *Memory) Rd64(va uint) (uint64, error) {
	pa, err := m.va2pa(va, AttrR)
	if err != nil {
		return 0, err
	}
	val, err := m.Rd64Phys(pa)
	m.monitor(pa, 8, AttrR)
	return val, err
}

// Rd32 reads a 32-bit data value from memory.
func (m *Memory) Rd32(va uint) (uint32, error) {
	pa, err := m.va2pa(va, AttrR)
	if err != nil {
		return 0, err
	}
	val, err := m.Rd32Phys(pa)
	m.monitor(pa, 4, AttrR)
	return val, err
}

// Rd16 reads a 16-bit data value from memory.
func (m *Memory) Rd16(va uint) (uint16, error) {
	pa, err := m.va2pa(va, AttrR)
	if err != nil {
		return 0, err
	}
	val, err := m.Rd16Phys(pa)
	m.monitor(pa, 2, AttrR)
	return val, err
}

// Rd8 reads an 8-bit data value from memory.
func (m *Memory) Rd8(va uint) (uint8, error) {
	pa, err := m.va2pa(va, AttrR)
	if err != nil {
		return 0, err
	}
	val, err := m.Rd8Phys(pa)
	m.monitor(pa, 1, AttrR)
	return val, err
}

//-----------------------------------------------------------------------------
// Physical Address Write Functions

// Wr64Phys writes a 64-bit data value to memory.
func (m *Memory) Wr64Phys(pa uint, val uint64) error {
	return m.findByAddr(pa, 8).Wr64(pa, val)
}

// Wr32Phys writes a 32-bit data value to memory.
func (m *Memory) Wr32Phys(pa uint, val uint32) error {
	return m.findByAddr(pa, 4).Wr32(pa, val)
}

// Wr16Phys writes a 16-bit data value to memory.
func (m *Memory) Wr16Phys(pa uint, val uint16) error {
	return m.findByAddr(pa, 2).Wr16(pa, val)
}

// Wr8Phys writes an 8-bit data value to memory.
func (m *Memory) Wr8Phys(pa uint, val uint8) error {
	return m.findByAddr(pa, 1).Wr8(pa, val)
}

//-----------------------------------------------------------------------------
// Virtual Address Write functions

// Wr64 writes a 64-bit data value to memory.
func (m *Memory) Wr64(va uint, val uint64) error {
	pa, err := m.va2pa(va, AttrW)
	if err != nil {
		return err
	}
	err = m.Wr64Phys(pa, val)
	m.monitor(pa, 8, AttrW)
	return err
}

// Wr32 writes a 32-bit data value to memory.
func (m *Memory) Wr32(va uint, val uint32) error {
	pa, err := m.va2pa(va, AttrW)
	if err != nil {
		return err
	}
	err = m.Wr32Phys(pa, val)
	m.monitor(pa, 4, AttrW)
	return err
}

// Wr16 writes a 16-bit data value to memory.
func (m *Memory) Wr16(va uint, val uint16) error {
	pa, err := m.va2pa(va, AttrW)
	if err != nil {
		return err
	}
	err = m.Wr16Phys(pa, val)
	m.monitor(pa, 2, AttrW)
	return err
}

// Wr8 writes an 8-bit data value to memory.
func (m *Memory) Wr8(va uint, val uint8) error {
	pa, err := m.va2pa(va, AttrW)
	if err != nil {
		return err
	}
	err = m.Wr8Phys(pa, val)
	m.monitor(pa, 1, AttrW)
	return err
}

//-----------------------------------------------------------------------------

// RdBuf reads a buffer of data from memory.
func (m *Memory) RdBuf(addr, n, width uint, vm bool) []uint {
	buf := make([]uint, n)
	for i := range buf {
		pa := addr + (uint(i) * (width >> 3))
		if vm {
			pa, _ = m.va2pa(pa, AttrR)
		}
		switch width {
		case 8:
			val, _ := m.Rd8Phys(pa)
			buf[i] = uint(val)
		case 16:
			val, _ := m.Rd16Phys(pa)
			buf[i] = uint(val)
		case 32:
			val, _ := m.Rd32Phys(pa)
			buf[i] = uint(val)
		case 64:
			val, _ := m.Rd64Phys(pa)
			buf[i] = uint(val)
		default:
			panic("bad width")
		}
	}
	return buf
}

//-----------------------------------------------------------------------------
