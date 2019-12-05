//-----------------------------------------------------------------------------
/*

Memory Breakpoints

Return a break exception when specific memory addresses have RWX access.

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
	"sort"

	"github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------

type bpState uint

const (
	bpOff   bpState = iota // disabled
	bpBreak                // break when hit
	bpSkip                 // skip once, break next time
)

var bpStateStr = map[bpState]string{
	bpOff:   "off",
	bpBreak: "brk",
	bpSkip:  "skip",
}

func (s bpState) String() string {
	return bpStateStr[s]
}

//-----------------------------------------------------------------------------

type breakPoint struct {
	addr   uint
	access Attribute
	state  bpState
}

type breakPoints map[uint]*breakPoint

func newBreakPoints() breakPoints {
	return make(map[uint]*breakPoint)
}

// Add a breakpoint
func (b breakPoints) Add(bp *breakPoint) {
	b[bp.addr] = bp
}

// Remove a breakpoint
func (b breakPoints) Remove(addr uint) {
	delete(b, addr)
}

// Set a breakpoint
func (b breakPoints) Set(addr uint) {
	if bp, ok := b[addr]; ok {
		bp.state = bpBreak
	}
}

// Clr a breakpoint
func (b breakPoints) Clr(addr uint) {
	if bp, ok := b[addr]; ok {
		bp.state = bpOff
	}
}

//-----------------------------------------------------------------------------

func (b breakPoints) check(addr uint, access Attribute) Exception {
	if bp, ok := b[addr]; ok {
		if access&bp.access != 0 {
			if bp.state == bpBreak {
				// skip so we don't immediately re-break.
				bp.state = bpSkip
				return ExBreak
			}
			if bp.state == bpSkip {
				// break on the next access.
				bp.state = bpBreak
			}
		}
	}
	return 0
}

func (b breakPoints) checkR(addr uint) Exception {
	ex := b.check(addr, AttrR)
	if ex != 0 {
		return ex | ExRead
	}
	return 0
}

func (b breakPoints) checkW(addr uint) Exception {
	ex := b.check(addr, AttrW)
	if ex != 0 {
		return ex | ExWrite
	}
	return 0
}

func (b breakPoints) checkX(addr uint) Exception {
	ex := b.check(addr, AttrX)
	if ex != 0 {
		return ex | ExExec
	}
	return 0
}

//-----------------------------------------------------------------------------

// sort breakpoints by address
type bpByAddr []*breakPoint

func (a bpByAddr) Len() int           { return len(a) }
func (a bpByAddr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bpByAddr) Less(i, j int) bool { return a[i].addr < a[j].addr }

// Display a string for the breakpoints
func (b breakPoints) Display(alen uint) string {
	if len(b) == 0 {
		return "no breakpoints"
	}
	fmtx := "%08x"
	if alen == 64 {
		fmtx = "%016x"
	}
	// list of break points
	bpList := []*breakPoint{}
	for _, v := range b {
		bpList = append(bpList, v)
	}
	// sort by address
	sort.Sort(bpByAddr(bpList))
	// display string
	s := make([][]string, len(bpList))
	for i, bp := range bpList {
		addrStr := fmt.Sprintf(fmtx, bp.addr)
		s[i] = []string{addrStr, bp.access.String(), bp.state.String()}
	}
	return cli.TableString(s, []int{0, 0, 0}, 1)
}

//-----------------------------------------------------------------------------
