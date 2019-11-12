//-----------------------------------------------------------------------------
/*

Emulated Target Memory

*/
//-----------------------------------------------------------------------------

package mem

import (
	"encoding/binary"
	"errors"
	"fmt"
)

//-----------------------------------------------------------------------------

// Memory is emulated read/write target memory.
type Memory struct {
	base  uint32          // base memory address
	mem   []uint8         // memory array
	st    map[uint]string // symbol table
	da    map[uint]string // reference disassembly
	align bool            // panic on misaligned access
}

// NewMemory returns a target memory object.
func NewMemory(base uint32, size int, align bool) *Memory {
	// allocate the memory and set it to all ones
	mem := make([]uint8, size)
	for i := range mem {
		mem[i] = 0xff
	}
	return &Memory{
		base:  base,
		mem:   mem,
		st:    make(map[uint]string),
		da:    make(map[uint]string),
		align: align,
	}
}

// Rd32 reads a 32-bit data value from memory.
func (m *Memory) Rd32(adr uint32) uint32 {
	if m.align && (adr&3 != 0) {
		panic(fmt.Sprintf("misaligned 32-bit read @ 0x%08x", adr))
	}
	return binary.LittleEndian.Uint32(m.mem[adr-m.base:])
}

// Wr32 writes a 32-bit data value to memory.
func (m *Memory) Wr32(adr uint32, val uint32) {
	if m.align && (adr&3 != 0) {
		panic(fmt.Sprintf("misaligned 32-bit write @ 0x%08x", adr))
	}
	binary.LittleEndian.PutUint32(m.mem[adr-m.base:], val)
}

// Rd16 reads a 16-bit data value from memory.
func (m *Memory) Rd16(adr uint32) uint16 {
	if m.align && (adr&1 != 0) {
		panic(fmt.Sprintf("misaligned 16-bit read @ 0x%08x", adr))
	}
	return binary.LittleEndian.Uint16(m.mem[adr-m.base:])
}

// Wr16 writes a 16-bit data value to memory.
func (m *Memory) Wr16(adr uint32, val uint16) {
	if m.align && (adr&1 != 0) {
		panic(fmt.Sprintf("misaligned 16-bit write @ 0x%08x", adr))
	}
	binary.LittleEndian.PutUint16(m.mem[adr-m.base:], val)
}

// Rd8 reads an 8-bit data value from memory.
func (m *Memory) Rd8(adr uint32) uint8 {
	return m.mem[adr-m.base]
}

// Wr8 writes an 8-bit data value to memory.
func (m *Memory) Wr8(adr uint32, val uint8) {
	m.mem[adr-m.base] = val
}

// Symbol returns a symbol for the memory address (if there is one).
func (m *Memory) Symbol(adr uint) string {
	return m.st[adr]
}

// AddSymbol adds a symbol to the symbol table.
func (m *Memory) AddSymbol(s string, adr, size uint) error {
	if len(s) == 0 {
		return errors.New("zero length symbol")
	}
	// check the symbol is within memory range
	s0 := uint(m.base)
	e0 := uint(m.base) + uint(len(m.mem))
	s1 := adr
	e1 := adr + size
	if s1 >= s0 && e1 <= e0 {
		m.st[adr] = s
		return nil
	}
	return fmt.Errorf("%s is out of memory range %08x-%08x", s, s1, e1)
}

// Disassembly returns the reference disassembly for the memory address (if there is any).
func (m *Memory) Disassembly(adr uint) string {
	return m.da[adr]
}

// AddDisassembly adds reference diassembly to the disassembly table.
func (m *Memory) AddDisassembly(s string, adr uint) {
	m.da[adr] = s
}

//-----------------------------------------------------------------------------
