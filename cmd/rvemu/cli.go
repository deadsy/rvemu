//-----------------------------------------------------------------------------
/*

RISC-V RV32 Emulator CLI

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
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
// memory functions

var helpMemDisplay = []cli.Help{
	{"<adr> [len]", "address (hex) - default is 0"},
	{"", "length (hex) - default is 0x40"},
}

var cmdMemDisplay8 = cli.Leaf{
	Descr: "display memory (8-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 8))
	},
}

var cmdMemDisplay16 = cli.Leaf{
	Descr: "display memory (16-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 16))
	},
}

var cmdMemDisplay32 = cli.Leaf{
	Descr: "display memory (32-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 32))
	},
}

var cmdMemDisplay64 = cli.Leaf{
	Descr: "display memory (64-bits)",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*emuApp).mem
		c.User.Put(m.Display(adr, size, 64))
	},
}

// memDisplayMenu submenu items
var memDisplayMenu = cli.Menu{
	{"b", cmdMemDisplay8, helpMemDisplay},
	{"h", cmdMemDisplay16, helpMemDisplay},
	{"w", cmdMemDisplay32, helpMemDisplay},
	{"d", cmdMemDisplay64, helpMemDisplay},
}

//-----------------------------------------------------------------------------

var helpGo = []cli.Help{
	{"<adr>", "address (hex) - default is PC"},
}

var cmdGo = cli.Leaf{
	Descr: "run the emulation (no tracing)",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).cpu
		adr, err := util.AddrArg(uint(m.PC), maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m.PC = uint64(adr)
		for true {
			err := m.Run()
			if err != nil {
				c.User.Put(fmt.Sprintf("%s\n", err))
				break
			}
		}
	},
}

var cmdTrace = cli.Leaf{
	Descr: "run the emulation (with tracing)",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).cpu
		adr, err := util.AddrArg(uint(m.PC), maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m.PC = uint64(adr)
		for true {
			s := m.Disassemble(uint(m.PC))
			err := m.Run()
			c.User.Put(fmt.Sprintf("%s\n", s))
			if err != nil {
				c.User.Put(fmt.Sprintf("%s\n", err))
				break
			}
		}
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
// breakpoints

var helpBreakpointAdd = []cli.Help{}

var helpBreakpointSetClr = []cli.Help{}

var cmdBreakpointAdd = cli.Leaf{
	Descr: "add a breakpoint",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdBreakpointSet = cli.Leaf{
	Descr: "set a breakpoint",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdBreakpointClr = cli.Leaf{
	Descr: "clear a breakpoint",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdBreakpointShow = cli.Leaf{
	Descr: "show the breakpoints",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*emuApp).mem
		c.User.Put(fmt.Sprintf("%s\n", m.BP.Display(32)))
	},
}

// memBreakpointMenu submenu items
var memBreakpointMenu = cli.Menu{
	{"add", cmdBreakpointAdd, helpBreakpointAdd},
	{"clr", cmdBreakpointClr, helpBreakpointSetClr},
	{"set", cmdBreakpointSet, helpBreakpointSetClr},
	{"show", cmdBreakpointShow},
}

//-----------------------------------------------------------------------------

// root menu
var menuRoot = cli.Menu{
	{"bp", memBreakpointMenu, "breakpoint functions"},
	{"csr", cmdCSR},
	{"da", cmdDisassemble, helpDisassemble},
	{"exit", cmdExit},
	{"go", cmdGo, helpGo},
	{"help", cmdHelp},
	{"history", cmdHistory, cli.HistoryHelp},
	{"map", cmdMap},
	{"md", memDisplayMenu, "memory display"},
	{"rf", cmdFloatRegisters},
	{"ri", cmdIntRegisters},
	{"reset", cmdReset},
	{"step", cmdStep, helpGo},
	{"sym", cmdSymbol},
	{"trace", cmdTrace, helpGo},
}

//-----------------------------------------------------------------------------
