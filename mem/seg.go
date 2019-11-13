//-----------------------------------------------------------------------------
/*

Memory Segments

*/
//-----------------------------------------------------------------------------

package mem

import "encoding/binary"

//-----------------------------------------------------------------------------

// Exception is a bit mask of memory acess exceptions.
type Exception uint

// Exception values.
const (
	ExcAlign   Exception = 1 << iota // misaligned read/write
	ExcRead                          // can't read this memory
	ExcWrite                         // can't write this memory
	ExcAddress                       // invalid memory address
	ExcExec                          // can't read instructions from this memory
)

// Attribute is a bit mask of memory acess attributes.
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

//-----------------------------------------------------------------------------

// Segment is an interface to a contiguous region of memory.
type Segment interface {
	Rd64(adr uint) (uint64, Exception)
	Rd32(adr uint) (uint32, Exception)
	Rd16(adr uint) (uint16, Exception)
	Rd8(adr uint) (uint8, Exception)
	Wr64(adr uint, val uint64) Exception
	Wr32(adr uint, val uint32) Exception
	Wr16(adr uint, val uint16) Exception
	Wr8(adr uint, val uint8) Exception
	In(adr, size uint) bool
}

//-----------------------------------------------------------------------------

// Chunk is a contiguous chunk of memory.
type Chunk struct {
	attr       Attribute // bitmask of attributes
	start, end uint      // address range
	mem        []uint8   // memory array
}

// NewChunk allocates and returns a memory chunk.
func NewChunk(start, size uint, attr Attribute) *Chunk {
	// allocate the memory and set it to all ones
	mem := make([]uint8, size)
	for i := range mem {
		mem[i] = 0xff
	}
	return &Chunk{
		attr:  attr,
		start: start,
		end:   start + size - 1,
		mem:   mem,
	}
}

// In returns true if the adr, size is entirely within the memory chunk.
func (m *Chunk) In(adr, size uint) bool {
	end := adr + size - 1
	return (adr >= m.start) && (end <= m.end)
}

// Rd64 reads a 64-bit data value from memory.
func (m *Chunk) Rd64(adr uint) (uint64, Exception) {
	return binary.LittleEndian.Uint64(m.mem[adr-m.start:]), 0
}

// Rd32 reads a 32-bit data value from memory.
func (m *Chunk) Rd32(adr uint) (uint32, Exception) {
	return binary.LittleEndian.Uint32(m.mem[adr-m.start:]), 0
}

// Rd16 reads a 16-bit data value from memory.
func (m *Chunk) Rd16(adr uint) (uint16, Exception) {
	return binary.LittleEndian.Uint16(m.mem[adr-m.start:]), 0
}

// Rd8 reads an 8-bit data value from memory.
func (m *Chunk) Rd8(adr uint) (uint8, Exception) {
	return m.mem[adr-m.start], 0
}

// Wr64 writes a 64-bit data value to memory.
func (m *Chunk) Wr64(adr uint, val uint64) Exception {
	binary.LittleEndian.PutUint64(m.mem[adr-m.start:], val)
	return 0
}

// Wr32 writes a 32-bit data value to memory.
func (m *Chunk) Wr32(adr uint, val uint32) Exception {
	binary.LittleEndian.PutUint32(m.mem[adr-m.start:], val)
	return 0
}

// Wr16 writes a 16-bit data value to memory.
func (m *Chunk) Wr16(adr uint, val uint16) Exception {
	binary.LittleEndian.PutUint16(m.mem[adr-m.start:], val)
	return 0
}

// Wr8 writes an 8-bit data value to memory.
func (m *Chunk) Wr8(adr uint, val uint8) Exception {
	m.mem[adr-m.start] = val
	return 0
}

//-----------------------------------------------------------------------------

// Empty is an empty memory region.
type Empty struct {
	attr       Attribute // bitmask of attributes
	start, end uint      // address range
}

// NewEmpty allocates and returns an empty memory region.
func NewEmpty(start, size uint, attr Attribute) *Empty {
	return &Empty{
		attr:  attr,
		start: start,
		end:   start + size - 1,
	}
}

// In returns true if the adr, size is entirely within the empty region.
func (m *Empty) In(adr, size uint) bool {
	end := adr + size - 1
	return (adr >= m.start) && (end <= m.end)
}

// Rd64 reads a 64-bit data value from memory.
func (m *Empty) Rd64(adr uint) (uint64, Exception) {
	return 0, 0
}

// Rd32 reads a 32-bit data value from memory.
func (m *Empty) Rd32(adr uint) (uint32, Exception) {
	return 0, 0
}

// Rd16 reads a 16-bit data value from memory.
func (m *Empty) Rd16(adr uint) (uint16, Exception) {
	return 0, 0
}

// Rd8 reads an 8-bit data value from memory.
func (m *Empty) Rd8(adr uint) (uint8, Exception) {
	return 0, 0
}

// Wr64 writes a 64-bit data value to memory.
func (m *Empty) Wr64(adr uint, val uint64) Exception {
	return 0
}

// Wr32 writes a 32-bit data value to memory.
func (m *Empty) Wr32(adr uint, val uint32) Exception {
	return 0
}

// Wr16 writes a 16-bit data value to memory.
func (m *Empty) Wr16(adr uint, val uint16) Exception {
	return 0
}

// Wr8 writes an 8-bit data value to memory.
func (m *Empty) Wr8(adr uint, val uint8) Exception {
	return 0
}

//-----------------------------------------------------------------------------
