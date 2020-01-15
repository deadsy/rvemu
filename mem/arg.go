//-----------------------------------------------------------------------------
/*

Memory Arguments

*/
//-----------------------------------------------------------------------------

package mem

import (
	"errors"
	"strconv"
)

//-----------------------------------------------------------------------------

// AddrArg converts an address argument to an address value.
func (m *Memory) AddrArg(arg string) (uint, error) {

	x, err := strconv.ParseUint(arg, 16, 64)
	if err != nil {
		return 0, errors.New("invalid address")
	}
	addr := uint(x)

	// check the limits
	maxAddr := uint((1 << m.alen) - 1)
	if addr > maxAddr {
		return 0, errors.New("invalid address, out of range")
	}

	return addr, nil
}

//-----------------------------------------------------------------------------
