//-----------------------------------------------------------------------------
/*

Emulated Target Memory

*/
//-----------------------------------------------------------------------------

package mem

import (
	"errors"
	"fmt"
)

//-----------------------------------------------------------------------------

// Range is a memory range.
type Range struct {
	Addr uint
	Size uint
}

//-----------------------------------------------------------------------------

// Memory is emulated target memory.
type Memory struct {
	region     []Region          // memory regions
	Entry      uint64            // entry point from ELF
	AddrLength int               // address bit length
	symByAddr  map[uint]string   // symbol table by address
	symByName  map[string]*Range // symbol table by name
	da         map[uint]string   // reference disassembly
}

// newMemory returns a memory object.
func newMemory(alen int) *Memory {
	return &Memory{
		AddrLength: alen,
		region:     make([]Region, 0),
		symByAddr:  make(map[uint]string),
		symByName:  make(map[string]*Range),
		da:         make(map[uint]string),
	}
}

// NewMem32 returns the memory for 32-bit processor.
func NewMem32() *Memory {
	return newMemory(32)
}

// NewMem64 returns the memory for 64-bit processor.
func NewMem64() *Memory {
	return newMemory(64)
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
	return nil
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
func (m *Memory) SymbolByAddress(adr uint) string {
	return m.symByAddr[adr]
}

// SymbolByName returns the memory range for a symbol.
func (m *Memory) SymbolByName(s string) *Range {
	return m.symByName[s]
}

// AddSymbol adds a symbol to the symbol table.
func (m *Memory) AddSymbol(s string, adr, size uint) error {
	if len(s) == 0 {
		return errors.New("zero length symbol")
	}
	if m.find(adr, size) != nil {
		m.symByAddr[adr] = s
		m.symByName[s] = &Range{adr, size}
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
