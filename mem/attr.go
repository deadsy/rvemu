//-----------------------------------------------------------------------------
/*

Memory Attributes

*/
//-----------------------------------------------------------------------------

package mem

import (
	"strings"
)

//-----------------------------------------------------------------------------

// Attribute is a bit mask of memory access attributes.
type Attribute uint

// Attribute values.
const (
	AttrR Attribute = 1 << iota // read
	AttrW                       // write
	AttrX                       // execute
	AttrM                       // misaligned access
)

// AttrRW = read/write
const AttrRW = AttrR | AttrW

// AttrRWM = read/write/misaligned
const AttrRWM = AttrR | AttrW | AttrM

// AttrRX = read/execute
const AttrRX = AttrR | AttrX

// AttrRWX = read/write/execute
const AttrRWX = AttrR | AttrW | AttrX

func (a Attribute) String() string {
	s := make([]string, 4)
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
	if a&AttrM != 0 {
		s[3] = "m"
	}
	return strings.Join(s, "")
}

//-----------------------------------------------------------------------------
