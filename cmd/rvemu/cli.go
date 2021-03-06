//-----------------------------------------------------------------------------
/*

RISC-V RV32/RV64 Emulator CLI

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/riscv/csr"
	"github.com/deadsy/riscv/mem"
	"github.com/deadsy/riscv/util"
)

//-----------------------------------------------------------------------------

const maxAdr = (1 << 32) - 1

//-----------------------------------------------------------------------------
// cli related leaf functions

var cmdHelp = cli.Leaf{
	Descr: "general help",
	F: func(c *cli.CLI, args []string) {
		c.GeneralHelp()
	},
}

var cmdHistory = cli.Leaf{
	Descr: "command history",
	F: func(c *cli.CLI, args []string) {
		c.SetLine(c.DisplayHistory(args))
	},
}

var cmdExit = cli.Leaf{
	Descr: "exit application",
	F: func(c *cli.CLI, args []string) {
		c.Exit()
	},
}

//-----------------------------------------------------------------------------
// virtual memory functions

var helpMemDisplay = []cli.Help{
	{"<adr> [len]", "address (hex) - default is 0"},
	{"", "length (hex) - default is 0x40"},
}

var cmdVmDisplay8 = cli.Leaf{
	Descr: "display virtual memory (8-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 8, true))
	},
}

var cmdVmDisplay16 = cli.Leaf{
	Descr: "display virtual memory (16-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 16, true))
	},
}

var cmdVmDisplay32 = cli.Leaf{
	Descr: "display virtual memory (32-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 32, true))
	},
}

var cmdVmDisplay64 = cli.Leaf{
	Descr: "display virtual memory (64-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 64, true))
	},
}

// memDisplayMenu submenu items
var memDisplayVm = cli.Menu{
	{"b", cmdVmDisplay8, helpMemDisplay},
	{"h", cmdVmDisplay16, helpMemDisplay},
	{"w", cmdVmDisplay32, helpMemDisplay},
	{"d", cmdVmDisplay64, helpMemDisplay},
}

//-----------------------------------------------------------------------------
// physical memory functions

var cmdPmDisplay8 = cli.Leaf{
	Descr: "display physical memory (8-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 8, false))
	},
}

var cmdPmDisplay16 = cli.Leaf{
	Descr: "display physical memory (16-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 16, false))
	},
}

var cmdPmDisplay32 = cli.Leaf{
	Descr: "display physical memory (32-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 32, false))
	},
}

var cmdPmDisplay64 = cli.Leaf{
	Descr: "display physical memory (64-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 64, false))
	},
}

// memDisplayMenu submenu items
var memDisplayPm = cli.Menu{
	{"b", cmdPmDisplay8, helpMemDisplay},
	{"h", cmdPmDisplay16, helpMemDisplay},
	{"w", cmdPmDisplay32, helpMemDisplay},
	{"d", cmdPmDisplay64, helpMemDisplay},
}

//-----------------------------------------------------------------------------

var helpGo = []cli.Help{
	{"<adr>", "address (hex) - default is PC"},
}

func goLoop(c *cli.CLI) bool {
	m := c.User.(*emuApp).cpu
	err := m.Run()
	if err != nil {
		c.User.Put(fmt.Sprintf("%s\r\n", err))
		return true
	}
	return false
}

var cmdGo = cli.Leaf{
	Descr: "run the emulation (ctrl-d to stop)",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).cpu
		adr, err := util.AddrArg(uint(m.PC), maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m.PC = uint64(adr)
		c.Loop(func() bool { return goLoop(c) }, cli.KeycodeCtrlD)
	},
}

func traceLoop(c *cli.CLI) bool {
	m := c.User.(*emuApp).cpu
	s := m.Disassemble(uint(m.PC))
	err := m.Run()
	c.User.Put(fmt.Sprintf("%s\r\n", s))
	if err != nil {
		c.User.Put(fmt.Sprintf("%s\r\n", err))
		return true
	}
	return false
}

var cmdTrace = cli.Leaf{
	Descr: "run the emulation with tracing (ctrl-d to stop)",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).cpu
		adr, err := util.AddrArg(uint(m.PC), maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m.PC = uint64(adr)
		c.Loop(func() bool { return traceLoop(c) }, cli.KeycodeCtrlD)
	},
}

var cmdStep = cli.Leaf{
	Descr: "single step the emulation",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).cpu
		adr, err := util.AddrArg(uint(m.PC), maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m.PC = uint64(adr)
		s := m.Disassemble(adr)
		err = m.Run()
		c.User.Put(fmt.Sprintf("%s\n", s))
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
		}
	},
}

//-----------------------------------------------------------------------------

var helpDisassemble = []cli.Help{
	{"[adr] [len]", "address (hex) - default is current pc"},
	{"", "length (hex) - default is 0x10"},
}

var cmdDisassemble = cli.Leaf{
	Descr: "disassemble memory",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).cpu
		adr, size, err := util.MemArg(uint(m.PC), maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		n := int(size)
		for n > 0 {
			da := m.Disassemble(adr)
			c.User.Put(fmt.Sprintf("%s\n", da))
			adr += da.Length
			n -= int(da.Length)
		}
	},
}

//-----------------------------------------------------------------------------

var cmdFloatRegisters = cli.Leaf{
	Descr: "display float registers",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).cpu
		c.User.Put(fmt.Sprintf("%s\n", m.FloatRegs()))
	},
}

var cmdIntRegisters = cli.Leaf{
	Descr: "display integer registers",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).cpu
		c.User.Put(fmt.Sprintf("%s\n", m.IntRegs()))
	},
}

//-----------------------------------------------------------------------------

var cmdReset = cli.Leaf{
	Descr: "reset the cpu",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).cpu
		m.Reset()
	},
}

//-----------------------------------------------------------------------------

var cmdSymbol = cli.Leaf{
	Descr: "display the symbol table",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).mem
		c.User.Put(fmt.Sprintf("%s\n", m.Symbols()))
	},
}

//-----------------------------------------------------------------------------

var cmdCSR = cli.Leaf{
	Descr: "display the control and status registers",
	F: func(c *cli.CLI, args []string) {
		csr := c.User.(*emuApp).cpu.CSR
		c.User.Put(fmt.Sprintf("%s\n", csr.Display()))
	},
}

//-----------------------------------------------------------------------------

var cmdMap = cli.Leaf{
	Descr: "display the memory map",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).mem
		c.User.Put(fmt.Sprintf("%s\n", m.Map()))
	},
}

//-----------------------------------------------------------------------------
// memory monitors

var helpBreakPointAdd = []cli.Help{}

var helpBreakPointSetClr = []cli.Help{}

var cmdBreakPointAdd = cli.Leaf{
	Descr: "add a monitor point",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdBreakPointSet = cli.Leaf{
	Descr: "set a monitor point",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdBreakPointClr = cli.Leaf{
	Descr: "clear a monitor point",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdBreakPointShow = cli.Leaf{
	Descr: "show the monitor points",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).mem
		c.User.Put(fmt.Sprintf("%s\n", m.DisplayBreakPoints()))
	},
}

// memBreakPointMenu submenu items
var memBreakPointMenu = cli.Menu{
	{"add", cmdBreakPointAdd, helpBreakPointAdd},
	{"clr", cmdBreakPointClr, helpBreakPointSetClr},
	{"set", cmdBreakPointSet, helpBreakPointSetClr},
	{"show", cmdBreakPointShow},
}

//-----------------------------------------------------------------------------

var helpPageTable = []cli.Help{
	{"<va> <mode> <attr>", "address (hex) - default is PC"},
	{"", "mode (u/s/m) - cpu mode - default is s(upervisor)"},
	{"", "attr (rwxm) - access attribute - default is r(ead)"},
}

var cmdPageTable = cli.Leaf{
	Descr: "display a page table walk",
	F: func(c *cli.CLI, args []string) {
		var err error

		cpu := c.User.(*emuApp).cpu
		m := c.User.(*emuApp).mem
		addr := uint(cpu.PC)
		mode := csr.ModeS
		attr := mem.AttrR

		if len(args) >= 1 {
			addr, err = m.AddrArg(args[0])
			if err != nil {
				c.User.Put(fmt.Sprintf("%s\n", err))
				return
			}
		}

		if len(args) >= 2 {
			mode, err = csr.ModeArg(args[1])
			if err != nil {
				c.User.Put(fmt.Sprintf("%s\n", err))
				return
			}
		}

		if len(args) >= 3 {
			attr, err = mem.AttrArg(args[2])
			if err != nil {
				c.User.Put(fmt.Sprintf("%s\n", err))
				return
			}
		}

		c.User.Put(fmt.Sprintf("%s\n", m.PageTableWalk(addr, mode, attr)))
	},
}

//-----------------------------------------------------------------------------

var cmdHost = cli.Leaf{
	Descr: "display host state",
	F: func(c *cli.CLI, args []string) {
		host := c.User.(*emuApp).host
		c.User.Put(fmt.Sprintf("%s\n", host))
	},
}

//-----------------------------------------------------------------------------

var cmdErrors = cli.Leaf{
	Descr: "display emulation errors",
	F: func(c *cli.CLI, args []string) {
		cpu := c.User.(*emuApp).cpu
		c.User.Put(fmt.Sprintf("%s\n", cpu.DisplayErrorBuffer()))
	},
}

//-----------------------------------------------------------------------------

// root menu
var menuRoot = cli.Menu{
	{"csr", cmdCSR},
	{"da", cmdDisassemble, helpDisassemble},
	{"errors", cmdErrors},
	{"exit", cmdExit},
	{"go", cmdGo, helpGo},
	{"help", cmdHelp},
	{"history", cmdHistory, cli.HistoryHelp},
	{"host", cmdHost},
	{"map", cmdMap},
	{"mm", memBreakPointMenu, "memory monitor functions"},
	{"pm", memDisplayPm, "physical memory menu"},
	{"pt", cmdPageTable, helpPageTable},
	{"rf", cmdFloatRegisters},
	{"ri", cmdIntRegisters},
	{"reset", cmdReset},
	{"step", cmdStep, helpGo},
	{"sym", cmdSymbol},
	{"trace", cmdTrace, helpGo},
	{"vm", memDisplayVm, "virtual memory menu"},
}

//-----------------------------------------------------------------------------
