//-----------------------------------------------------------------------------
/*

RISC-V POSIX System Calls

See:

linux/include/uapi/asm-generic/unistd.h
glibc/sysdeps/unix/sysv/linux/riscv/sysdep.h

glibc passes the syscall number in a7.
Syscall arguments are passed in a0..a6.
The syscall return value is passed in a0.

*/
//-----------------------------------------------------------------------------

package rv

//-----------------------------------------------------------------------------

// scFunc32 is 32-bit system call function
type scFunc32 func(m *RV32, s *scEntry)

// scFunc64 is a 64-bit system call function
type scFunc64 func(m *RV64, s *scEntry)

type scEntry struct {
	name string
	sc32 scFunc32
	sc64 scFunc64
}

var scTable = map[int]scEntry{
	57: {"close", sc32_close, sc64_close},
	80: {"fstat", sc32_fstat, sc64_fstat},
	93: {"exit", sc32_exit, sc64_exit},
}

func scLookup(n int) *scEntry {
	if s, ok := scTable[n]; ok {
		return &s
	}
	return nil
}

//-----------------------------------------------------------------------------
// 32-bit system calls

func sc32_close(m *RV32, s *scEntry) {
	m.flag |= flagBreak
}

func sc32_fstat(m *RV32, s *scEntry) {
	m.flag |= flagBreak
}

func sc32_exit(m *RV32, s *scEntry) {
	m.X[regA0] = 0
	m.flag |= flagExit
}

//-----------------------------------------------------------------------------
// 64-bit system calls

func sc64_close(m *RV64, s *scEntry) {
	m.flag |= flagBreak
}

func sc64_fstat(m *RV64, s *scEntry) {
	m.flag |= flagBreak
}

func sc64_exit(m *RV64, s *scEntry) {
	m.X[regA0] = 0
	m.flag |= flagExit
}

//-----------------------------------------------------------------------------
