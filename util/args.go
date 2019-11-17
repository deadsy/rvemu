//-----------------------------------------------------------------------------
/*

Utilities for CLI Argument Processing

*/
//-----------------------------------------------------------------------------

package util

import (
	cli "github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------

// AddressArg converts an address argument to an address value.
func AddressArg(defAdr, maxAdr uint, args []string) (uint, error) {
	err := cli.CheckArgc(args, []int{0, 1})
	if err != nil {
		return 0, err
	}
	// address
	adr := defAdr
	if len(args) >= 1 {
		adr, err = cli.UintArg(args[0], [2]uint{0, maxAdr}, 16)
		if err != nil {
			return 0, err
		}
	}
	return adr, nil
}

//-----------------------------------------------------------------------------
