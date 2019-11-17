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

// AddrArg converts an address argument to an address value.
func AddrArg(defAdr, maxAdr uint, args []string) (uint, error) {
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

// MemArg converts memory arguments to an (address, size) tuple.
func MemArg(defAdr, maxAdr uint, args []string) (uint, uint, error) {
	err := cli.CheckArgc(args, []int{0, 1, 2})
	if err != nil {
		return 0, 0, err
	}
	// address
	adr := defAdr
	if len(args) >= 1 {
		adr, err = cli.UintArg(args[0], [2]uint{0, maxAdr}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	// size
	size := uint(0x80) // default size
	if len(args) >= 2 {
		size, err = cli.UintArg(args[1], [2]uint{1, 0x100000000}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	return adr, size, nil
}

//-----------------------------------------------------------------------------
