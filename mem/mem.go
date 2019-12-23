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

// Symbol is a memory symbol.
type Symbol struct {
	Name string
	Addr uint
	Size uint
}

// sort symbols by address
type symbolByAddr []*Symbol

func (a symbolByAddr) Len() int           { return len(a) }
func (a symbolByAddr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a symbolByAddr) Less(i, j int) bool { return a[i].Addr < a[j].Addr }

//-----------------------------------------------------------------------------

// Memory is emulated target memory.
type Memory struct {
	Entry     uint64             // entry point from ELF
	BP        breakPoints        // break points
	alen      uint               // address bit length
	csr       *csr.State         // CSR state
	region    []Region           // memory regions
	symByAddr map[uint]*Symbol   // symbol table by address
	symByName map[string]*Symbol // symbol table by name
	noMemory  Region             // empty memory region
}

// newMemory returns a memory object.
func newMemory(alen uint, csr *csr.State, empty Attribute) *Memory {
	return &Memory{
		BP:        newBreakPoints(),
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

// AddrStr returns a string for the address.
func (m *Memory) AddrStr(addr uint) string {
	if m.alen == 32 {
		return fmt.Sprintf("%08x", addr)
	}
	return fmt.Sprintf("%016x", addr)
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
	val, err := m.findByAddr(pa, 4).RdIns(pa)
	if err == nil {
		err = m.BP.checkX(pa)
	}
	return val, err
}

// Rd64Phys reads a 64-bit data value from memory.
func (m *Memory) Rd64Phys(pa uint) (uint64, error) {
	val, err := m.findByAddr(pa, 8).Rd64(pa)
	if err == nil {
		err = m.BP.checkR(pa)
	}
	return val, err
}

// Rd32Phys reads a 32-bit data value from memory.
func (m *Memory) Rd32Phys(pa uint) (uint32, error) {
	val, err := m.findByAddr(pa, 4).Rd32(pa)
	if err == nil {
		err = m.BP.checkR(pa)
	}
	return val, err
}

// Rd16Phys reads a 16-bit data value from memory.
func (m *Memory) Rd16Phys(pa uint) (uint16, error) {
	val, err := m.findByAddr(pa, 2).Rd16(pa)
	if err == nil {
		err = m.BP.checkR(pa)
	}
	return val, err
}

// Rd8Phys reads an 8-bit data value from memory.
func (m *Memory) Rd8Phys(pa uint) (uint8, error) {
	val, err := m.findByAddr(pa, 1).Rd8(pa)
	if err == nil {
		err = m.BP.checkR(pa)
	}
	return val, err
}

// Rd32RangePhys reads a range of 32-bit data values from memory.
func (m *Memory) Rd32RangePhys(pa, n uint) []uint32 {
	x := make([]uint32, n)
	for i := uint(0); i < n; i++ {
		x[i], _ = m.Rd32(pa + (i * 4))
	}
	return x
}

//-----------------------------------------------------------------------------
// Virtual Address Read Functions

// RdIns reads a 32-bit instruction from memory.
func (m *Memory) RdIns(va uint) (uint, error) {
	pa, err := m.va2pa(va, AttrRX)
	if err != nil {
		return 0, err
	}
	return m.RdInsPhys(pa)
}

// Rd64 reads a 64-bit data value from memory.
func (m *Memory) Rd64(va uint) (uint64, error) {
	pa, err := m.va2pa(va, AttrR)
	if err != nil {
		return 0, err
	}
	return m.Rd64Phys(pa)
}

// Rd32 reads a 32-bit data value from memory.
func (m *Memory) Rd32(va uint) (uint32, error) {
	pa, err := m.va2pa(va, AttrR)
	if err != nil {
		return 0, err
	}
	return m.Rd32Phys(pa)
}

// Rd16 reads a 16-bit data value from memory.
func (m *Memory) Rd16(va uint) (uint16, error) {
	pa, err := m.va2pa(va, AttrR)
	if err != nil {
		return 0, err
	}
	return m.Rd16Phys(pa)
}

// Rd8 reads an 8-bit data value from memory.
func (m *Memory) Rd8(va uint) (uint8, error) {
	pa, err := m.va2pa(va, AttrR)
	if err != nil {
		return 0, err
	}
	return m.Rd8Phys(pa)
}

//-----------------------------------------------------------------------------
// Physical Address Write Functions

// Wr64Phys writes a 64-bit data value to memory.
func (m *Memory) Wr64Phys(pa uint, val uint64) error {
	err := m.findByAddr(pa, 8).Wr64(pa, val)
	if err == nil {
		err = m.BP.checkW(pa)
	}
	return err
}

// Wr32Phys writes a 32-bit data value to memory.
func (m *Memory) Wr32Phys(pa uint, val uint32) error {
	err := m.findByAddr(pa, 4).Wr32(pa, val)
	if err == nil {
		err = m.BP.checkW(pa)
	}
	return err
}

// Wr16Phys writes a 16-bit data value to memory.
func (m *Memory) Wr16Phys(pa uint, val uint16) error {
	err := m.findByAddr(pa, 2).Wr16(pa, val)
	if err == nil {
		err = m.BP.checkW(pa)
	}
	return err
}

// Wr8Phys writes an 8-bit data value to memory.
func (m *Memory) Wr8Phys(pa uint, val uint8) error {
	err := m.findByAddr(pa, 1).Wr8(pa, val)
	if err == nil {
		err = m.BP.checkW(pa)
	}
	return err
}

//-----------------------------------------------------------------------------
// Virtual Address Write functions

// Wr64 writes a 64-bit data value to memory.
func (m *Memory) Wr64(va uint, val uint64) error {
	pa, err := m.va2pa(va, AttrW)
	if err != nil {
		return err
	}
	return m.Wr64Phys(pa, val)
}

// Wr32 writes a 32-bit data value to memory.
func (m *Memory) Wr32(va uint, val uint32) error {
	pa, err := m.va2pa(va, AttrW)
	if err != nil {
		return err
	}
	return m.Wr32Phys(pa, val)
}

// Wr16 writes a 16-bit data value to memory.
func (m *Memory) Wr16(va uint, val uint16) error {
	pa, err := m.va2pa(va, AttrW)
	if err != nil {
		return err
	}
	return m.Wr16Phys(pa, val)
}

// Wr8 writes an 8-bit data value to memory.
func (m *Memory) Wr8(va uint, val uint8) error {
	pa, err := m.va2pa(va, AttrW)
	if err != nil {
		return err
	}
	return m.Wr8Phys(pa, val)
}

//-----------------------------------------------------------------------------
// symbol functions

// SymbolByAddress returns a symbol for the memory address.
func (m *Memory) SymbolByAddress(adr uint) *Symbol {
	return m.symByAddr[adr]
}

// SymbolByName returns the symbol for a symbol name.
func (m *Memory) SymbolByName(s string) *Symbol {
	return m.symByName[s]
}

// SymbolGetAddress returns the symbol address for a symbol name.
func (m *Memory) SymbolGetAddress(s string) (uint, error) {
	symbol := m.symByName[s]
	if symbol == nil {
		return 0, fmt.Errorf("%s not found", s)
	}
	return symbol.Addr, nil
}

// AddSymbol adds a symbol to the symbol table.
func (m *Memory) AddSymbol(s string, adr, size uint) error {
	if m.findByAddr(adr, size) != nil {
		symbol := Symbol{s, adr, size}
		m.symByAddr[adr] = &symbol
		m.symByName[s] = &symbol
		return nil
	}
	return fmt.Errorf("%s is not in a memory region", s)
}

//-----------------------------------------------------------------------------

// AddBreakPointByName adds a breakpoint by symbol name.
func (m *Memory) AddBreakPointByName(s string, attr Attribute) error {
	sym := m.symByName[s]
	if sym == nil {
		return fmt.Errorf("%s not found", s)
	}
	m.BP.add(&breakPoint{sym.Name, sym.Addr, attr, bpBreak})
	return nil
}

//-----------------------------------------------------------------------------
