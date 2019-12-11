//-----------------------------------------------------------------------------
/*

RISC-V Compliance Testing

https://github.com/riscv/riscv-compliance

The complicance tests exit the test by setting status in the gp register
and using the ecall instruction. At that point the non-emulation code can
read the output from memory and compare it with the known good signature
value for the test.

*/
//-----------------------------------------------------------------------------

package ecall

import "github.com/deadsy/riscv/rv"

//-----------------------------------------------------------------------------

// Compliance is a compliance ecall object.
type Compliance struct {
}

// NewCompliance returns a compliance ecall object.
func NewCompliance() *Compliance {
	return &Compliance{}
}

// Call is an ecall.
func (c *Compliance) Call(m *rv.RV) error {
	return m.Exit(m.X[rv.RegGp])
}

//-----------------------------------------------------------------------------
