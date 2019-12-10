//-----------------------------------------------------------------------------
/*

Dump strings for memory.

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

// Display returns a string for a contiguous region of memory.
func (m *Memory) Display(adr, size, width uint) string {
	s := []string{}

	// round down address to 16 byte boundary
	adr &= ^uint(15)

	// round up n to an integral multiple of 16 bytes
	size = (size + 15) & ^uint(15)

	// build the header
	var hdr, fmtx string
	if m.alen == 32 {
		hdr = "addr      "
		fmtx = "%08x  "
	} else {
		hdr = "addr              "
		fmtx = "%016x  "
	}
	switch width {
	case 8:
		hdr += "0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F"
		fmtx += "%s %s"
	case 16:
		hdr += "0    2    4    6    8    A    C    E"
		fmtx += "%s"
	case 32:
		hdr += "0        4        8        C"
		fmtx += "%s"
	case 64:
		hdr += "0                8"
		fmtx += "%s"
	}
	s = append(s, hdr)

	// read and print the data
	for i := 0; i < int(size>>4); i++ {
		if width == 8 {
			// read 16x8 bits per line
			var data [16]string
			var ascii [16]string
			for j := 0; j < 16; j++ {
				x, _ := m.Rd8(adr + uint(j))
				data[j] = fmt.Sprintf("%02x", x)
				if x >= 32 && x <= 126 {
					ascii[j] = fmt.Sprintf("%c", x)
				} else {
					ascii[j] = "."
				}
			}
			dataStr := strings.Join(data[:], " ")
			asciiStr := strings.Join(ascii[:], "")
			s = append(s, fmt.Sprintf(fmtx, adr, dataStr, asciiStr))
		} else if width == 16 {
			// read 8x16 bits per line
			var data [8]string
			for j := 0; j < 8; j++ {
				x, _ := m.Rd16(adr + uint(j*2))
				data[j] = fmt.Sprintf("%04x", x)
			}
			dataStr := strings.Join(data[:], " ")
			s = append(s, fmt.Sprintf(fmtx, adr, dataStr))
		} else if width == 32 {
			// read 4x32 bits per line
			var data [4]string
			for j := 0; j < 4; j++ {
				x, _ := m.Rd32(adr + uint(j*4))
				data[j] = fmt.Sprintf("%08x", x)
			}
			dataStr := strings.Join(data[:], " ")
			s = append(s, fmt.Sprintf(fmtx, adr, dataStr))
		} else if width == 64 {
			// read 2x64 bits per line
			var data [2]string
			for j := 0; j < 2; j++ {
				x, _ := m.Rd64(adr + uint(j*8))
				data[j] = fmt.Sprintf("%016x", x)
			}
			dataStr := strings.Join(data[:], " ")
			s = append(s, fmt.Sprintf(fmtx, adr, dataStr))
		}
		adr += 16
		adr &= (1 << m.alen) - 1
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
