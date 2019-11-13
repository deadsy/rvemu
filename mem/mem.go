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

type memRange struct {
	adr, size uint
}

//-----------------------------------------------------------------------------
// memory segments

const AttrR = 1 << 0 // read
const AttrW = 1 << 0 // write
const AttrX = 1 << 0 // execute

// Segment is a memory segment.
type Segment struct {
	attr       uint    // bitmask of segment attributes
	start, end uint    // address range
	mem        []uint8 // memory array
}

// NewSegment alocates and returns a memory segment.
func NewSegment(start, size uint, attr uint) *Segment {
	// allocate the memory and set it to all ones
	mem := make([]uint8, size)
	for i := range mem {
		mem[i] = 0xff
	}
	return &Segment{
		attr:  attr,
		start: start,
		end:   start + size - 1,
		mem:   mem,
	}
}

// In returns true if the memory region is entirely within the segment.
func (s *Segment) In(adr, size uint) bool {
	end := adr + size - 1
	return (adr >= s.start) && (end <= s.end)
}

//-----------------------------------------------------------------------------

// Memory is emulated target memory.
type Memory struct {
	segment   []*Segment          // memory segments
	symByAddr map[uint]string     // symbol table by address
	symByName map[string]memRange // symbol table by name
	da        map[uint]string     // reference disassembly
}

// NewMemory returns a memory object.
func NewMemory() *Memory {
	return &Memory{
		seg:       make([]*Segment, 0),
		symByAddr: make(map[uint]string),
		symByName: make(map[string]memRange),
		da:        make(map[uint]string),
	}
}

// Add a memory segment to the memory.
func (m *Memory) Add(s *Segment) {
	m.segment = append(m.segment, s)
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
	return m.symByAddr[adr]
}

// AddSymbol adds a symbol to the symbol table.
func (m *Memory) AddSymbol(s string, adr, size uint) error {
	if len(s) == 0 {
		return errors.New("zero length symbol")
	}
	for i := range m.segment {
		if m.segment[i].In(adr, size) {
			m.symByAddr[adr] = s
			m.symByName[s] = memRange{adr, size}
			return nil
		}
	}
	return fmt.Errorf("%s is not in a memory segment", s)
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
