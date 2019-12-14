//-----------------------------------------------------------------------------
/*

Emulated Target Memory

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
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
	region    []Region           // memory regions
	symByAddr map[uint]*Symbol   // symbol table by address
	symByName map[string]*Symbol // symbol table by name
	noMemory  Region             // empty memory region
}

// newMemory returns a memory object.
func newMemory(alen uint, empty Attribute) *Memory {
	return &Memory{
		alen:      alen,
		BP:        newBreakPoints(),
		region:    make([]Region, 0),
		symByAddr: make(map[uint]*Symbol),
		symByName: make(map[string]*Symbol),
		noMemory:  newEmpty(empty),
	}
}

// NewMem32 returns memory with a 32-bit address bus.
func NewMem32(empty Attribute) *Memory {
	return newMemory(32, empty)
}

// NewMem64 returns memory with a 64-bit address bus.
func NewMem64(empty Attribute) *Memory {
	return newMemory(64, empty)
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
// read functions

// RdIns reads a 32-bit instruction from memory.
func (m *Memory) RdIns(adr uint) (uint, error) {
	val, err := m.findByAddr(adr, 4).RdIns(adr)
	if err == nil {
		err = m.BP.checkX(adr)
	}
	return val, err
}

// Rd64 reads a 64-bit data value from memory.
func (m *Memory) Rd64(adr uint) (uint64, error) {
	val, err := m.findByAddr(adr, 8).Rd64(adr)
	if err == nil {
		err = m.BP.checkR(adr)
	}
	return val, err
}

// Rd32 reads a 32-bit data value from memory.
func (m *Memory) Rd32(adr uint) (uint32, error) {
	val, err := m.findByAddr(adr, 4).Rd32(adr)
	if err == nil {
		err = m.BP.checkR(adr)
	}
	return val, err
}

// Rd16 reads a 16-bit data value from memory.
func (m *Memory) Rd16(adr uint) (uint16, error) {
	val, err := m.findByAddr(adr, 2).Rd16(adr)
	if err == nil {
		err = m.BP.checkR(adr)
	}
	return val, err
}

// Rd8 reads an 8-bit data value from memory.
func (m *Memory) Rd8(adr uint) (uint8, error) {
	val, err := m.findByAddr(adr, 1).Rd8(adr)
	if err == nil {
		err = m.BP.checkR(adr)
	}
	return val, err
}

// Rd32Range reads a range of 32-bit data values from memory.
func (m *Memory) Rd32Range(adr, n uint) []uint32 {
	x := make([]uint32, n)
	for i := uint(0); i < n; i++ {
		x[i], _ = m.Rd32(adr + (i * 4))
	}
	return x
}

//-----------------------------------------------------------------------------
// write functions

// Wr64 writes a 64-bit data value to memory.
func (m *Memory) Wr64(adr uint, val uint64) error {
	err := m.findByAddr(adr, 8).Wr64(adr, val)
	if err == nil {
		err = m.BP.checkW(adr)
	}
	return err
}

// Wr32 writes a 32-bit data value to memory.
func (m *Memory) Wr32(adr uint, val uint32) error {
	err := m.findByAddr(adr, 4).Wr32(adr, val)
	if err == nil {
		err = m.BP.checkW(adr)
	}
	return err
}

// Wr16 writes a 16-bit data value to memory.
func (m *Memory) Wr16(adr uint, val uint16) error {
	err := m.findByAddr(adr, 2).Wr16(adr, val)
	if err == nil {
		err = m.BP.checkW(adr)
	}
	return err
}

// Wr8 writes an 8-bit data value to memory.
func (m *Memory) Wr8(adr uint, val uint8) error {
	err := m.findByAddr(adr, 1).Wr8(adr, val)
	if err == nil {
		err = m.BP.checkW(adr)
	}
	return err
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

// Add a breakpoint by symbol name.
func (m *Memory) AddBreakPointByName(s string, attr Attribute) error {
	sym := m.symByName[s]
	if sym == nil {
		return fmt.Errorf("%s not found", s)
	}
	m.BP.add(&breakPoint{sym.Name, sym.Addr, attr, bpBreak})
	return nil
}

//-----------------------------------------------------------------------------
