//-----------------------------------------------------------------------------
/*

Utilities for Memory Display

*/
//-----------------------------------------------------------------------------

package util

import (
	"fmt"
	"strings"

	"github.com/deadsy/riscv/mem"
)

//-----------------------------------------------------------------------------

func MemDisplay(m *mem.Memory, adr, size, xlen uint) string {
	s := make([]string, 0)
	// round down address to 16 byte boundary
	adr &= ^uint(15)
	// round up n to an integral multiple of 16 bytes
	size = (size + 15) & ^uint(15)
	// print the header
	if xlen == 32 {
		s = append(s, "addr      0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F")
	} else {
		s = append(s, "addr              0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F")
	}
	// read and print the data
	for i := 0; i < int(size>>4); i++ {
		// read 16 bytes per line
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
		if xlen == 32 {
			s = append(s, fmt.Sprintf("%08x  %s  %s", adr, dataStr, asciiStr))
		} else {
			s = append(s, fmt.Sprintf("%016x  %s  %s", adr, dataStr, asciiStr))
		}
		adr += 16
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
