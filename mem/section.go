//-----------------------------------------------------------------------------
/*

Memory Sections

A memory section is a contiguous region of real memory. The access attributes
control it's usage (read.write, executable, misaligned access). Loading an ELF
file will typically create a number of sections with appropriate attributes.

*/
//-----------------------------------------------------------------------------

package mem

import "encoding/binary"

//-----------------------------------------------------------------------------

// Section is a contiguous region of real memory.
type Section struct {
	name       string    // section name
	attr       Attribute // bitmask of attributes
	start, end uint      // address range
	mem        []uint8   // memory array
}

// NewSection allocates and returns a memory chunk.
func NewSection(name string, start, size uint, attr Attribute) *Section {
	mem := make([]uint8, size)
	return &Section{
		name:  name,
		attr:  attr,
		start: start,
		end:   start + size - 1,
		mem:   mem,
	}
}

// SetAttr sets the attributes for this memory section.
func (m *Section) SetAttr(attr Attribute) {
	m.attr = attr
}

// Info returns the information for this region.
func (m *Section) Info() *RegionInfo {
	return &RegionInfo{
		name:  m.name,
		start: m.start,
		end:   m.end,
		attr:  m.attr,
	}
}

// In returns true if the adr, size is entirely within the memory chunk.
func (m *Section) In(adr, size uint) bool {
	end := adr + size - 1
	return (adr >= m.start) && (end <= m.end)
}

// RdIns reads a 32-bit instruction from memory.
func (m *Section) RdIns(adr uint) (uint, error) {
	return uint(binary.LittleEndian.Uint32(m.mem[adr-m.start:])), rdInsError(adr, m.attr, m.name)
}

// Rd64 reads a 64-bit data value from memory.
func (m *Section) Rd64(adr uint) (uint64, error) {
	return binary.LittleEndian.Uint64(m.mem[adr-m.start:]), rdError(adr, m.attr, m.name, 8)
}

// Rd32 reads a 32-bit data value from memory.
func (m *Section) Rd32(adr uint) (uint32, error) {
	return binary.LittleEndian.Uint32(m.mem[adr-m.start:]), rdError(adr, m.attr, m.name, 4)
}

// Rd16 reads a 16-bit data value from memory.
func (m *Section) Rd16(adr uint) (uint16, error) {
	return binary.LittleEndian.Uint16(m.mem[adr-m.start:]), rdError(adr, m.attr, m.name, 2)
}

// Rd8 reads an 8-bit data value from memory.
func (m *Section) Rd8(adr uint) (uint8, error) {
	return m.mem[adr-m.start], rdError(adr, m.attr, m.name, 1)
}

// Wr64 writes a 64-bit data value to memory.
func (m *Section) Wr64(adr uint, val uint64) error {
	if m.attr&AttrW != 0 {
		binary.LittleEndian.PutUint64(m.mem[adr-m.start:], val)
	}
	return wrError(adr, m.attr, m.name, 8)
}

// Wr32 writes a 32-bit data value to memory.
func (m *Section) Wr32(adr uint, val uint32) error {
	if m.attr&AttrW != 0 {
		binary.LittleEndian.PutUint32(m.mem[adr-m.start:], val)
	}
	return wrError(adr, m.attr, m.name, 4)
}

// Wr16 writes a 16-bit data value to memory.
func (m *Section) Wr16(adr uint, val uint16) error {
	if m.attr&AttrW != 0 {
		binary.LittleEndian.PutUint16(m.mem[adr-m.start:], val)
	}
	return wrError(adr, m.attr, m.name, 2)
}

// Wr8 writes an 8-bit data value to memory.
func (m *Section) Wr8(adr uint, val uint8) error {
	if m.attr&AttrW != 0 {
		m.mem[adr-m.start] = val
	}
	return wrError(adr, m.attr, m.name, 1)
}

//-----------------------------------------------------------------------------
