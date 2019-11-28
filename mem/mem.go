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
	region     []Region           // memory regions
	Entry      uint64             // entry point from ELF
	AddrLength int                // address bit length
	symByAddr  map[uint]*Symbol   // symbol table by address
	symByName  map[string]*Symbol // symbol table by name
	noMemory   Region             // empty memory region
	da         map[uint]string    // reference disassembly

}

// newMemory returns a memory object.
func newMemory(alen int, empty Attribute) *Memory {
	return &Memory{
		AddrLength: alen,
		region:     make([]Region, 0),
		symByAddr:  make(map[uint]*Symbol),
		symByName:  make(map[string]*Symbol),
		noMemory:   newEmpty(empty),
		da:         make(map[uint]string),
	}
}

// NewMem32 returns the memory for 32-bit processor.
func NewMem32(empty Attribute) *Memory {
	return newMemory(32, empty)
}

// NewMem64 returns the memory for 64-bit processor.
func NewMem64(empty Attribute) *Memory {
	return newMemory(64, empty)
}

// Add a memory region to the memory.
func (m *Memory) Add(r Region) {
	m.region = append(m.region, r)
}

// find returns the memory region
func (m *Memory) find(adr, size uint) Region {
	for _, r := range m.region {
		if r.In(adr, size) {
			return r
		}
	}
	return m.noMemory
}

//-----------------------------------------------------------------------------

// RdIns reads a 32-bit instruction from memory.
func (m *Memory) RdIns(adr uint) (uint, Exception) {
	return m.find(adr, 4).RdIns(adr)
}

// Rd64 reads a 64-bit data value from memory.
func (m *Memory) Rd64(adr uint) (uint64, Exception) {
	return m.find(adr, 8).Rd64(adr)
}

// Rd32 reads a 32-bit data value from memory.
func (m *Memory) Rd32(adr uint) (uint32, Exception) {
	return m.find(adr, 4).Rd32(adr)
}

// Rd16 reads a 16-bit data value from memory.
func (m *Memory) Rd16(adr uint) (uint16, Exception) {
	return m.find(adr, 2).Rd16(adr)
}

// Rd8 reads an 8-bit data value from memory.
func (m *Memory) Rd8(adr uint) (uint8, Exception) {
	return m.find(adr, 1).Rd8(adr)
}

// Wr64 writes a 64-bit data value to memory.
func (m *Memory) Wr64(adr uint, val uint64) Exception {
	return m.find(adr, 8).Wr64(adr, val)
}

// Wr32 writes a 32-bit data value to memory.
func (m *Memory) Wr32(adr uint, val uint32) Exception {
	return m.find(adr, 4).Wr32(adr, val)
}

// Wr16 writes a 16-bit data value to memory.
func (m *Memory) Wr16(adr uint, val uint16) Exception {
	return m.find(adr, 2).Wr16(adr, val)
}

// Wr8 writes an 8-bit data value to memory.
func (m *Memory) Wr8(adr uint, val uint8) Exception {
	return m.find(adr, 1).Wr8(adr, val)
}

//-----------------------------------------------------------------------------

// SymbolByAddress returns a symbol for the memory address.
func (m *Memory) SymbolByAddress(adr uint) *Symbol {
	return m.symByAddr[adr]
}

// SymbolByName returns the symbol for a symbol name.
func (m *Memory) SymbolByName(s string) *Symbol {
	return m.symByName[s]
}

// AddSymbol adds a symbol to the symbol table.
func (m *Memory) AddSymbol(s string, adr, size uint) error {
	if m.find(adr, size) != nil {
		symbol := Symbol{s, adr, size}
		m.symByAddr[adr] = &symbol
		m.symByName[s] = &symbol
		return nil
	}
	return fmt.Errorf("%s is not in a memory segment", s)
}

//-----------------------------------------------------------------------------

// Disassembly returns the reference disassembly for the memory address (if there is any).
func (m *Memory) Disassembly(adr uint) string {
	return m.da[adr]
}

// AddDisassembly adds reference diassembly to the disassembly table.
func (m *Memory) AddDisassembly(s string, adr uint) {
	m.da[adr] = s
}

//-----------------------------------------------------------------------------
