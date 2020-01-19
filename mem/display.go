//-----------------------------------------------------------------------------
/*

Display strings for memory.

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
	"sort"
	"strings"

	"github.com/deadsy/go-cli"
	"github.com/deadsy/riscv/util"
)

//-----------------------------------------------------------------------------

// Map returns a memory map display string.
func (m *Memory) Map() string {
	if len(m.region) == 0 {
		return "no map"
	}
	// list of regions
	regions := []*RegionInfo{}
	for _, r := range m.region {
		regions = append(regions, r.Info())
	}
	// sort by start address
	sort.Sort(regionByStart(regions))
	// display string
	s := make([][]string, len(regions))
	for i, r := range regions {
		addrStr := fmt.Sprintf("%s %s", m.AddrStr(r.start), m.AddrStr(r.end))
		attrStr := r.attr.String()
		sizeStr := fmt.Sprintf("(%d bytes)", r.end-r.start+1)
		s[i] = []string{r.name, addrStr, attrStr, sizeStr}
	}
	return cli.TableString(s, []int{0, 0, 0, 0}, 1)
}

//-----------------------------------------------------------------------------

// Symbols returns an address sorted string of memory symbols.
func (m *Memory) Symbols() string {
	if len(m.symByName) == 0 {
		return "no symbols"
	}
	// list of symbols
	symbols := []*Symbol{}
	for _, v := range m.symByName {
		symbols = append(symbols, v)
	}
	// sort by address
	sort.Sort(symbolByAddr(symbols))
	// display string
	s := make([][]string, len(symbols))
	for i, se := range symbols {
		addrStr := fmt.Sprintf("%s", m.AddrStr(se.Addr))
		sizeStr := fmt.Sprintf("(%d)", se.Size)
		s[i] = []string{addrStr, sizeStr, util.GreenString(se.Name)}
	}
	return cli.TableString(s, []int{0, 0, 0}, 1)
}

//-----------------------------------------------------------------------------

const bytesPerLine = 32 // must be a power of 2

// Display returns a string for a contiguous region of memory.
func (m *Memory) Display(adr, size, width uint, vm bool) string {
	s := []string{}

	fmtLine := fmt.Sprintf("%%0%dx %%s %%s", [2]int{16, 8}[util.BoolToInt(m.alen == 32)])
	fmtData := fmt.Sprintf("%%0%dx", width>>2)

	// round down address to width alignment
	adr &= ^uint((width >> 3) - 1)
	// round up size to an integral multiple of bytesPerLine bytes
	size = (size + bytesPerLine - 1) & ^uint(bytesPerLine-1)

	// read and print the data
	for i := 0; i < int(size/bytesPerLine); i++ {

		// read bytesPerLine bytes
		buf := m.RdBuf(adr, bytesPerLine/(width>>3), width, vm)

		// create the data string
		xStr := make([]string, len(buf))
		for j := range xStr {
			xStr[j] = fmt.Sprintf(fmtData, buf[j])
		}
		dataStr := strings.Join(xStr[:], " ")

		// create the ascii string
		var data [bytesPerLine]uint8
		var ascii [bytesPerLine]string
		for j := range data {
			if vm {
				data[j], _ = m.Rd8(adr + uint(j))
			} else {
				data[j], _ = m.Rd8Phys(adr + uint(j))
			}
			if data[j] >= 32 && data[j] <= 126 {
				ascii[j] = fmt.Sprintf("%c", data[j])
			} else {
				ascii[j] = "."
			}
		}
		asciiStr := strings.Join(ascii[:], "")

		s = append(s, fmt.Sprintf(fmtLine, adr, dataStr, asciiStr))
		adr += bytesPerLine
		adr &= (1 << m.alen) - 1
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
