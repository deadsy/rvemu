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
	segment   []Segment         // memory segments
	Entry     uint64            // entry point from ELF
	symByAddr map[uint]string   // symbol table by address
	symByName map[string]*Range // symbol table by name
	da        map[uint]string   // reference disassembly
}

// NewMemory returns a memory object.
func NewMemory() *Memory {
	return &Memory{
		segment:   make([]Segment, 0),
		symByAddr: make(map[uint]string),
		symByName: make(map[string]*Range),
		da:        make(map[uint]string),
	}
}

// Add a memory segment to the memory.
func (m *Memory) Add(s Segment) {
	m.segment = append(m.segment, s)
}

// find returns the segment
func (m *Memory) find(adr, size uint) Segment {
	for _, s := range m.segment {
		if s.In(adr, size) {
			return s
		}
	}
	// It's expected there will be a catch-all empty
	// memory segment defined.
	panic("where's the empty memory segment?")
}

//-----------------------------------------------------------------------------

// RdIns reads a 32-bit instruction from memory.
func (m *Memory) RdIns(adr uint) (uint, Exception) {
	return m.find(adr, 4).RdIns(adr)
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
	for i := range m.segment {
		if m.segment[i].In(adr, size) {
			m.symByAddr[adr] = s
			m.symByName[s] = &Range{adr, size}
			return nil
		}
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
