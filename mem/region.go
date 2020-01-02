//-----------------------------------------------------------------------------
/*

Memory Region

*/
//-----------------------------------------------------------------------------

package mem

//-----------------------------------------------------------------------------

// Region is an interface to a region of memory.
type Region interface {
	Info() *RegionInfo
	SetAttr(attr Attribute)
	RdIns(adr uint) (uint, error)
	Rd64(adr uint) (uint64, error)
	Rd32(adr uint) (uint32, error)
	Rd16(adr uint) (uint16, error)
	Rd8(adr uint) (uint8, error)
	Wr64(adr uint, val uint64) error
	Wr32(adr uint, val uint32) error
	Wr16(adr uint, val uint16) error
	Wr8(adr uint, val uint8) error
	In(adr, size uint) bool
}

// RegionInfo contains information for the memory region.
type RegionInfo struct {
	name       string
	start, end uint
	attr       Attribute
}

//-----------------------------------------------------------------------------
// sort regions by start address

type regionByStart []*RegionInfo

func (a regionByStart) Len() int           { return len(a) }
func (a regionByStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a regionByStart) Less(i, j int) bool { return a[i].start < a[j].start }

//-----------------------------------------------------------------------------
