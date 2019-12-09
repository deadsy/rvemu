//-----------------------------------------------------------------------------
/*

Memory Segments

*/
//-----------------------------------------------------------------------------

package mem

import (
	"encoding/binary"
	"math"
	"strings"
)

//-----------------------------------------------------------------------------

// Error is a bit mask of memory access errors.
type Error uint

// Error values.
const (
	ErrAlign Error = 1 << iota // misaligned read/write
	ErrRead                    // can't read this memory
	ErrWrite                   // can't write this memory
	ErrExec                    // can't read instructions from this memory
	ErrBreak                   // break on memory access
)

func (e Error) String() string {
	s := make([]string, 0)
	if e&ErrAlign != 0 {
		s = append(s, "align")
	}
	if e&ErrRead != 0 {
		s = append(s, "read")
	}
	if e&ErrWrite != 0 {
		s = append(s, "write")
	}
	if e&ErrExec != 0 {
		s = append(s, "exec")
	}
	return strings.Join(s, ",")
}

//-----------------------------------------------------------------------------

// Attribute is a bit mask of memory access attributes.
type Attribute uint

// Attribute values.
const (
	AttrR Attribute = 1 << iota // read
	AttrW                       // write
	AttrX                       // execute
)

// AttrRW = read/write
const AttrRW = AttrR | AttrW

// AttrRX = read/execute
const AttrRX = AttrR | AttrX

// AttrRWX = read/write/execute
const AttrRWX = AttrR | AttrW | AttrX

func (a Attribute) String() string {
	s := make([]string, 3)
	for i := range s {
		s[i] = "-"
	}
	if a&AttrR != 0 {
		s[0] = "r"
	}
	if a&AttrW != 0 {
		s[1] = "w"
	}
	if a&AttrX != 0 {
		s[2] = "x"
	}
	return strings.Join(s, "")
}

//-----------------------------------------------------------------------------
// memory access errors

func wrError(adr uint, attr Attribute, align uint) Error {
	var err Error
	if attr&AttrW == 0 {
		err |= ErrWrite
	}
	if adr&(align-1) != 0 {
		err |= ErrAlign
	}
	return err
}

func rdError(adr uint, attr Attribute, align uint) Error {
	var err Error
	if attr&AttrR == 0 {
		err |= ErrRead
	}
	if adr&(align-1) != 0 {
		err |= ErrAlign
	}
	return err
}

func rdInsError(adr uint, attr Attribute) Error {
	// rv32c has mixed 32/16 bit instruction streams so
	// we allow 32-bit reads on 2 byte address boundaries.
	err := rdError(adr, attr, 2)
	if attr&AttrX == 0 {
		err |= ErrExec
	}
	return err
}

//-----------------------------------------------------------------------------

// RegionInfo contains information for the memory region.
type RegionInfo struct {
	name       string
	start, end uint
	attr       Attribute
}

// sort regions by start address
type regionByStart []*RegionInfo

func (a regionByStart) Len() int           { return len(a) }
func (a regionByStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a regionByStart) Less(i, j int) bool { return a[i].start < a[j].start }

// Region is an interface to a contiguous region of memory.
type Region interface {
	Info() *RegionInfo
	SetAttr(attr Attribute)
	RdIns(adr uint) (uint, Error)
	Rd64(adr uint) (uint64, Error)
	Rd32(adr uint) (uint32, Error)
	Rd16(adr uint) (uint16, Error)
	Rd8(adr uint) (uint8, Error)
	Wr64(adr uint, val uint64) Error
	Wr32(adr uint, val uint32) Error
	Wr16(adr uint, val uint16) Error
	Wr8(adr uint, val uint8) Error
	In(adr, size uint) bool
}

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
	// allocate the memory and set it to all ones
	mem := make([]uint8, size)
	for i := range mem {
		mem[i] = 0xff
	}
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
func (m *Section) RdIns(adr uint) (uint, Error) {
	return uint(binary.LittleEndian.Uint32(m.mem[adr-m.start:])), rdInsError(adr, m.attr)
}

// Rd64 reads a 64-bit data value from memory.
func (m *Section) Rd64(adr uint) (uint64, Error) {
	return binary.LittleEndian.Uint64(m.mem[adr-m.start:]), rdError(adr, m.attr, 8)
}

// Rd32 reads a 32-bit data value from memory.
func (m *Section) Rd32(adr uint) (uint32, Error) {
	return binary.LittleEndian.Uint32(m.mem[adr-m.start:]), rdError(adr, m.attr, 4)
}

// Rd16 reads a 16-bit data value from memory.
func (m *Section) Rd16(adr uint) (uint16, Error) {
	return binary.LittleEndian.Uint16(m.mem[adr-m.start:]), rdError(adr, m.attr, 2)
}

// Rd8 reads an 8-bit data value from memory.
func (m *Section) Rd8(adr uint) (uint8, Error) {
	return m.mem[adr-m.start], rdError(adr, m.attr, 1)
}

// Wr64 writes a 64-bit data value to memory.
func (m *Section) Wr64(adr uint, val uint64) Error {
	if m.attr&AttrW != 0 {
		binary.LittleEndian.PutUint64(m.mem[adr-m.start:], val)
	}
	return wrError(adr, m.attr, 8)
}

// Wr32 writes a 32-bit data value to memory.
func (m *Section) Wr32(adr uint, val uint32) Error {
	if m.attr&AttrW != 0 {
		binary.LittleEndian.PutUint32(m.mem[adr-m.start:], val)
	}
	return wrError(adr, m.attr, 4)
}

// Wr16 writes a 16-bit data value to memory.
func (m *Section) Wr16(adr uint, val uint16) Error {
	if m.attr&AttrW != 0 {
		binary.LittleEndian.PutUint16(m.mem[adr-m.start:], val)
	}
	return wrError(adr, m.attr, 2)
}

// Wr8 writes an 8-bit data value to memory.
func (m *Section) Wr8(adr uint, val uint8) Error {
	if m.attr&AttrW != 0 {
		m.mem[adr-m.start] = val
	}
	return wrError(adr, m.attr, 1)
}

//-----------------------------------------------------------------------------
// If a memory access does not correspond to a defined memory Region the
// empty memory region will be used.

// empty memory region.
type empty struct {
	attr Attribute // bitmask of attributes
}

// newEmpty allocates and returns the empty memory region.
func newEmpty(attr Attribute) *empty {
	return &empty{
		attr: attr,
	}
}

// SetAttr sets the attributes for the empty region.
func (m *empty) SetAttr(attr Attribute) {
	m.attr = attr
}

// Info returns the information for the empty region.
func (m *empty) Info() *RegionInfo {
	return &RegionInfo{
		name: "empty",
		attr: m.attr,
	}
}

// In returns true if the adr, size is entirely within the empty region.
func (m *empty) In(adr, size uint) bool {
	return true
}

// RdIns reads a 32-bit instruction from memory.
func (m *empty) RdIns(adr uint) (uint, Error) {
	return math.MaxUint32, rdInsError(adr, m.attr)
}

// Rd64 reads a 64-bit data value from memory.
func (m *empty) Rd64(adr uint) (uint64, Error) {
	return math.MaxUint64, rdError(adr, m.attr, 8)
}

// Rd32 reads a 32-bit data value from memory.
func (m *empty) Rd32(adr uint) (uint32, Error) {
	return math.MaxUint32, rdError(adr, m.attr, 4)
}

// Rd16 reads a 16-bit data value from memory.
func (m *empty) Rd16(adr uint) (uint16, Error) {
	return math.MaxUint16, rdError(adr, m.attr, 2)
}

// Rd8 reads an 8-bit data value from memory.
func (m *empty) Rd8(adr uint) (uint8, Error) {
	return math.MaxUint8, rdError(adr, m.attr, 1)
}

// Wr64 writes a 64-bit data value to memory.
func (m *empty) Wr64(adr uint, val uint64) Error {
	return wrError(adr, m.attr, 8)
}

// Wr32 writes a 32-bit data value to memory.
func (m *empty) Wr32(adr uint, val uint32) Error {
	return wrError(adr, m.attr, 4)
}

// Wr16 writes a 16-bit data value to memory.
func (m *empty) Wr16(adr uint, val uint16) Error {
	return wrError(adr, m.attr, 2)
}

// Wr8 writes an 8-bit data value to memory.
func (m *empty) Wr8(adr uint, val uint8) Error {
	return wrError(adr, m.attr, 1)
}

//-----------------------------------------------------------------------------
