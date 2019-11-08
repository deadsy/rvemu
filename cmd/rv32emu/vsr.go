//-----------------------------------------------------------------------------
/*

Virtual Subroutines

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/deadsy/riscv/rv"
)

//-----------------------------------------------------------------------------

func vsrOpen(m *rv.RV) {
	fmt.Printf("*** vsrOpen ***\n")
}
func vsrClose(m *rv.RV) {
	fmt.Printf("*** vsrClose ***\n")
}
func vsrRead(m *rv.RV) {
	fmt.Printf("*** vsrRead ***\n")
}
func vsrWrite(m *rv.RV) {
	fmt.Printf("*** vsrWrite ***\n")
}
func vsrArgs(m *rv.RV) {
	fmt.Printf("*** vsrArgs ***\n")
}
func vsrExit(m *rv.RV) {
	m.Exit(0)
}

//-----------------------------------------------------------------------------
