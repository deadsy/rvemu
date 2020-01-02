//-----------------------------------------------------------------------------
/*

Memory Errors

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
	"strings"

	"github.com/deadsy/riscv/csr"
)

//-----------------------------------------------------------------------------

// Error is a memory acccess error.
type Error struct {
	Type uint   // bitmap of memory errors
	Ex   int    // riscv exception
	Addr uint   // memory address causing the error
	Name string // section name for the address
}

// Memory error bits.
const (
	ErrAlign = 1 << iota // misaligned read/write
	ErrRead              // can't read this memory
	ErrWrite             // can't write this memory
	ErrExec              // can't read instructions from this memory
	ErrPage              // error with page table translation
	ErrBreak             // break on memory access
)

func (e *Error) Error() string {
	s := make([]string, 0)
	if e.Type&ErrAlign != 0 {
		s = append(s, "align")
	}
	if e.Type&ErrRead != 0 {
		s = append(s, "read")
	}
	if e.Type&ErrWrite != 0 {
		s = append(s, "write")
	}
	if e.Type&ErrExec != 0 {
		s = append(s, "exec")
	}
	if e.Type&ErrPage != 0 {
		s = append(s, "page")
	}
	if e.Type&ErrBreak != 0 {
		s = append(s, "break")
	}
	errStr := strings.Join(s, ",")
	return fmt.Sprintf("%s @ %08x (%s)", errStr, e.Addr, e.Name)
}

//-----------------------------------------------------------------------------

func pageError(va uint, attr Attribute) error {
	n := uint(ErrPage)
	ex := -1
	// The attribute is what the cpu was trying to do when the page error occured.
	// It's sense is inverted from the other error cases.
	if attr&AttrR != 0 {
		n |= ErrRead
		ex = csr.ExLoadPageFault
	}
	if attr&AttrW != 0 {
		n |= ErrWrite
		ex = csr.ExStorePageFault
	}
	if attr&AttrX != 0 {
		n |= ErrExec
		ex = csr.ExInsPageFault
	}
	return &Error{n, ex, va, ""}
}

func wrError(addr uint, attr Attribute, name string, align uint) error {
	var n uint
	ex := -1
	if attr&AttrW == 0 {
		n |= ErrWrite
		ex = csr.ExStoreAccessFault
	}
	if (attr&AttrM == 0) && (addr&(align-1) != 0) {
		n |= ErrAlign
		ex = csr.ExStoreAddrMisaligned
	}
	if n != 0 {
		return &Error{n, ex, addr, name}
	}
	return nil
}

func rdError(addr uint, attr Attribute, name string, align uint) error {
	var n uint
	ex := -1
	if attr&AttrR == 0 {
		n |= ErrRead
		ex = csr.ExLoadAccessFault
	}
	if (attr&AttrM == 0) && (addr&(align-1) != 0) {
		n |= ErrAlign
		ex = csr.ExLoadAddrMisaligned
	}
	if n != 0 {
		return &Error{n, ex, addr, name}
	}
	return nil
}

func rdInsError(addr uint, attr Attribute, name string) error {
	// rv32c has mixed 32/16 bit instruction streams so
	// we allow 32-bit reads on 2 byte address boundaries.
	var n uint
	ex := -1
	if attr&AttrX == 0 {
		n |= ErrExec
		ex = csr.ExInsAccessFault
	}
	if (attr&AttrM == 0) && (addr&1 != 0) {
		n |= ErrAlign
		ex = csr.ExInsAddrMisaligned
	}
	if n != 0 {
		return &Error{n, ex, addr, name}
	}
	return nil
}

//-----------------------------------------------------------------------------
