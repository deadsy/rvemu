//-----------------------------------------------------------------------------
/*

RISC-V 32-bit Emulator CLI

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

var cmdMemDisplay = cli.Leaf{
	Descr: "display memory",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := util.MemArg(0, maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m := c.User.(*userApp).mem
		c.User.Put(m.Display(adr, size))
	},
}

//-----------------------------------------------------------------------------

var helpGo = []cli.Help{
	{"<adr>", "address (hex) - default is PC"},
}

var cmdGo = cli.Leaf{
	Descr: "run the emulation (no tracing)",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*userApp).cpu
		adr, err := util.AddrArg(uint(m.PC), maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m.PC = uint32(adr)
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
		m := c.User.(*userApp).cpu
		adr, err := util.AddrArg(uint(m.PC), maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m.PC = uint32(adr)
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
		m := c.User.(*userApp).cpu
		adr, err := util.AddrArg(uint(m.PC), maxAdr, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m.PC = uint32(adr)
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
		m := c.User.(*userApp).cpu
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

var cmdRegisters = cli.Leaf{
	Descr: "display registers",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*userApp).cpu
		c.User.Put(fmt.Sprintf("%s\n", m.IRegs()))
	},
}

//-----------------------------------------------------------------------------

var cmdReset = cli.Leaf{
	Descr: "reset the cpu",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*userApp).cpu
		m.Reset()
	},
}

//-----------------------------------------------------------------------------

var cmdSymbol = cli.Leaf{
	Descr: "display the symbol table",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*userApp).mem
		c.User.Put(fmt.Sprintf("%s\n", m.Symbols()))
	},
}

//-----------------------------------------------------------------------------

var cmdCSR = cli.Leaf{
	Descr: "display the control and status registers",
	F: func(c *cli.CLI, args []string) {
		csr := c.User.(*userApp).cpu.CSR
		c.User.Put(fmt.Sprintf("%s\n", csr.Display()))
	},
}

//-----------------------------------------------------------------------------

// root menu
var menuRoot = cli.Menu{
	{"csr", cmdCSR},
	{"da", cmdDisassemble, helpDisassemble},
	{"exit", cmdExit},
	{"go", cmdGo, helpGo},
	{"help", cmdHelp},
	{"history", cmdHistory, cli.HistoryHelp},
	{"reg", cmdRegisters},
	{"md", cmdMemDisplay, helpMemDisplay},
	{"reset", cmdReset},
	{"step", cmdStep, helpGo},
	{"sym", cmdSymbol},
	{"trace", cmdTrace, helpGo},
}

//-----------------------------------------------------------------------------
