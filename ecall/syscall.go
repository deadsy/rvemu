//-----------------------------------------------------------------------------
/*

RISC-V Linux System Calls

See:

linux/include/uapi/asm-generic/unistd.h
glibc/sysdeps/unix/sysv/linux/riscv/sysdep.h

glibc passes the syscall number in a7.
Syscall arguments are passed in a0..a6.
The syscall return value is passed in a0.

*/
//-----------------------------------------------------------------------------

package ecall

import "github.com/deadsy/riscv/rv"

//-----------------------------------------------------------------------------
// system calls

func sc_close(m *rv.RV) {
	m.SetBreak()
}

func sc_fstat(m *rv.RV) {
	m.SetBreak()
}

func sc_exit(m *rv.RV) {
	m.Exit(0)
}

//-----------------------------------------------------------------------------

type scFunc func(m *rv.RV)

type scEntry struct {
	name string
	sc   scFunc
}

var scTable = map[uint]scEntry{
	57: {"close", sc_close},
	80: {"fstat", sc_fstat},
	93: {"exit", sc_exit},
}

func scLookup(n uint) *scEntry {
	if e, ok := scTable[n]; ok {
		return &e
	}
	return nil
}

//-----------------------------------------------------------------------------

// Syscall is a syscall ecall object.
type Syscall struct {
}

// NewSyscall returns a syscall ecall object.
func NewSyscall() *Syscall {
	return &Syscall{}
}

// Call is an ecall handler.
func (sc *Syscall) Call(m *rv.RV) {
	n := uint(m.X[rv.RegA7])
	e := scLookup(n)
	if e != nil {
		e.sc(m)
	}
}

//-----------------------------------------------------------------------------
