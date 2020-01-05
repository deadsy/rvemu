//-----------------------------------------------------------------------------
/*

Memory Monitor Points

Perform a debug action when specific memory addresses have RWX access.

*/
//-----------------------------------------------------------------------------

package mem

import (
	"fmt"
	"sort"

	"github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------

type mpType uint

const (
	tBreak mpType = iota // break point
	tCall                // call a function
)

func (t mpType) String() string {
	return []string{"break", "call"}[t]
}

//-----------------------------------------------------------------------------

type mpState uint

const (
	sOff  mpState = iota // disabled
	sOn                  // enabled
	sSkip                // skip
)

func (s mpState) String() string {
	return []string{"off", "on", "skip"}[s]
}

//-----------------------------------------------------------------------------

type mPoint struct {
	name   string
	addr   uint
	access Attribute
	typ    mpType
	state  mpState
}

//-----------------------------------------------------------------------------

func (m *Memory) monitor(addr, size uint, access Attribute) {
	mp, ok := m.mp[addr]
	if !ok {
		return
	}
	if access&mp.access == 0 || mp.state == sOff {
		// no trigger
		return
	}
	if mp.state == sSkip {
		// trigger on the next access
		mp.state = sOn
		return
	}

	switch mp.typ {
	case tBreak:
		// set the pending breakpoint
		if m.brk == nil {
			m.brk = breakError(addr, access, mp.name)
		}
	case tCall:
	}
}

//-----------------------------------------------------------------------------

// AddBreakPoint adds a breakpoint
func (m *Memory) AddBreakPoint(name string, addr uint, attr Attribute) {
	m.mp[addr] = &mPoint{name, addr, attr, tBreak, sOn}
}

// AddBreakPointByName adds a breakpoint by symbol name.
func (m *Memory) AddBreakPointByName(name string, attr Attribute) error {
	s := m.SymbolByName(name)
	if s == nil {
		return fmt.Errorf("symbol \"%s\" not found", name)
	}
	m.AddBreakPoint(s.Name, s.Addr, attr)
	return nil
}

// AddCallPoint adds a callpoint.
func (m *Memory) AddCallPoint(name string, addr uint, attr Attribute) {
	m.mp[addr] = &mPoint{name, addr, attr, tCall, sOn}
}

// AddCallPointByName adds a callpoint by symbol name.
func (m *Memory) AddCallPointByName(name string, attr Attribute) error {
	s := m.SymbolByName(name)
	if s == nil {
		return fmt.Errorf("symbol \"%s\" not found", name)
	}
	m.AddCallPoint(s.Name, s.Addr, attr)
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
type mpByAddr []*mPoint

func (a mpByAddr) Len() int           { return len(a) }
func (a mpByAddr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a mpByAddr) Less(i, j int) bool { return a[i].addr < a[j].addr }

// Display a string for the monitor points.
func (m *Memory) DisplayMonitorPoints() string {
	if len(m.mp) == 0 {
		return "no monitor points"
	}
	// list of break points
	mpList := []*mPoint{}
	for _, v := range m.mp {
		mpList = append(mpList, v)
	}
	// sort by address
	sort.Sort(mpByAddr(mpList))
	// display string
	s := make([][]string, len(mpList))
	for i, mp := range mpList {
		s[i] = []string{m.AddrStr(mp.addr), mp.access.String(), mp.typ.String(), mp.state.String()}
	}
	return cli.TableString(s, []int{0, 0, 0, 0}, 1)
}

//-----------------------------------------------------------------------------
