//-----------------------------------------------------------------------------
/*

Display strings for memory.

*/
//-----------------------------------------------------------------------------

package mem

import (
	"encoding/binary"
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

// Display returns a string for a contiguous region of memory.
func (m *Memory) Display(adr, size, width uint, vm bool) string {

	fmtx := ""
	if m.alen == 32 {
		fmtx = "%08x %s %s"
	} else {
		fmtx = "%016x %s %s"
	}

	s := []string{}

	// round down address to 16 byte boundary
	adr &= ^uint(15)

	// round up n to an integral multiple of 16 bytes
	size = (size + 15) & ^uint(15)

	// read and print the data
	for i := 0; i < int(size>>4); i++ {

		// read 16 bytes per line
		var data [16]uint8
		var ascii [16]string
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

		// create the data string
		dataStr := ""
		switch width {
		case 8:
			var xStr [16]string
			for j := range xStr {
				xStr[j] = fmt.Sprintf("%02x", data[j])
			}
			dataStr = strings.Join(xStr[:], " ")
		case 16:
			var xStr [8]string
			for j := range xStr {
				val := binary.LittleEndian.Uint16(data[2*j:])
				xStr[j] = fmt.Sprintf("%04x", val)
			}
			dataStr = strings.Join(xStr[:], " ")
		case 32:
			var xStr [4]string
			for j := range xStr {
				val := binary.LittleEndian.Uint32(data[4*j:])
				xStr[j] = fmt.Sprintf("%08x", val)
			}
			dataStr = strings.Join(xStr[:], " ")
		case 64:
			var xStr [2]string
			for j := range xStr {
				val := binary.LittleEndian.Uint64(data[8*j:])
				xStr[j] = fmt.Sprintf("%016x", val)
			}
			dataStr = strings.Join(xStr[:], " ")
		}

		asciiStr := strings.Join(ascii[:], "")
		s = append(s, fmt.Sprintf(fmtx, adr, dataStr, asciiStr))

		adr += 16
		adr &= (1 << m.alen) - 1
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
