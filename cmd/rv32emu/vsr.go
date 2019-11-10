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

func vsrOpen(m *rv.RV32) {
	fmt.Printf("*** vsrOpen ***\n")
}
func vsrClose(m *rv.RV32) {
	fmt.Printf("*** vsrClose ***\n")
}
func vsrRead(m *rv.RV32) {
	fmt.Printf("*** vsrRead ***\n")
}
func vsrWrite(m *rv.RV32) {
	fmt.Printf("*** vsrWrite ***\n")
}
func vsrArgs(m *rv.RV32) {
	fmt.Printf("*** vsrArgs ***\n")
}
func vsrExit(m *rv.RV32) {
	m.Exit(0)
}

//-----------------------------------------------------------------------------
