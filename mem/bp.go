//-----------------------------------------------------------------------------
/*

Memory Break Points

Perform a debug action when specific memory addresses have RWX access.

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
	"sort"
	"strings"
)

//-----------------------------------------------------------------------------

type bpState uint

const (
	sOff  bpState = iota // disabled
	sOn                  // enabled
	sSkip                // skip
)

func (s bpState) String() string {
	return []string{"off", "on", "skip"}[s]
}

//-----------------------------------------------------------------------------

type bpFunc func(bp *BreakPoint) bool

// BreakPoint stores the information for a memory breakpoint.
type BreakPoint struct {
	Name   string    // breakpoint name
	Addr   uint      // address for trigger
	Access Attribute // access for trigger
	alen   uint      // address bit length
	state  bpState   // breakpoint state
	cond   bpFunc    // condition function
}

func (mon *BreakPoint) String() string {
	s := []string{}
	s = append(s, addrStr(mon.Addr, mon.alen))
	s = append(s, mon.Access.String())
	s = append(s, mon.state.String())
	return strings.Join(s, " ")
}

//-----------------------------------------------------------------------------

func (m *Memory) monitor(addr, size uint, access Attribute) {
	bp, ok := m.bp[addr]
	if !ok {
		return
	}
	if access&bp.Access == 0 || bp.state == sOff {
		// no trigger
		return
	}
	if bp.state == sSkip {
		// trigger on the next access
		bp.state = sOn
		return
	}
	// triggered...
	brk := true
	if bp.cond != nil {
		brk = bp.cond(bp)
	}
	// should we break?
	if brk && m.brk == nil {
		m.brk = breakError(addr, access, bp.Name)
	}
}

//-----------------------------------------------------------------------------

// AddBreakPoint adds a break point.
func (m *Memory) AddBreakPoint(name string, addr uint, attr Attribute, cond bpFunc) {
	bp := &BreakPoint{
		Name:   name,
		Addr:   addr,
		Access: attr,
		alen:   m.alen,
		state:  sOn,
		cond:   cond,
	}
	m.bp[addr] = bp
}

// AddBreakPointByName adds a break point by symbol name.
func (m *Memory) AddBreakPointByName(name string, attr Attribute, cond bpFunc) error {
	s := m.SymbolByName(name)
	if s == nil {
		return fmt.Errorf("symbol \"%s\" not found", name)
	}
	m.AddBreakPoint(s.Name, s.Addr, attr, cond)
	return nil
}

//-----------------------------------------------------------------------------

// GetBreak returned (and resets) any pending breakpoint.
func (m *Memory) GetBreak() error {
	err := m.brk
	m.brk = nil
	return err
}

//-----------------------------------------------------------------------------

// sort monitor points by address
type bpByAddr []*BreakPoint

func (a bpByAddr) Len() int           { return len(a) }
func (a bpByAddr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bpByAddr) Less(i, j int) bool { return a[i].Addr < a[j].Addr }

// DisplayBreakPoints displays a string for the memory break points.
func (m *Memory) DisplayBreakPoints() string {
	if len(m.bp) == 0 {
		return "no break points"
	}
	// list of break points
	bpList := []*BreakPoint{}
	for _, v := range m.bp {
		bpList = append(bpList, v)
	}
	// sort by address
	sort.Sort(bpByAddr(bpList))
	// display string
	s := make([]string, len(bpList))
	for i, bp := range bpList {
		s[i] = bp.String()
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
