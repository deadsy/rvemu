//-----------------------------------------------------------------------------
/*

RISC-V Emulator

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"strings"

	cli "github.com/deadsy/go-cli"
)

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

// memArgs converts memory arguments to an (address, size) tuple.
func memArgs(args []string) (uint32, uint, error) {
	err := cli.CheckArgc(args, []int{0, 1, 2})
	if err != nil {
		return 0, 0, err
	}
	// address
	adr := 0
	if len(args) >= 1 {
		adr, err = cli.IntArg(args[0], [2]int{0, 0xffffffff}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	// size
	size := 0x80 // default size
	if len(args) >= 2 {
		size, err = cli.IntArg(args[1], [2]int{1, 0x100000000}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	return uint32(adr), uint(size), nil
}

var cmdMemDisplay = cli.Leaf{
	Descr: "display memory",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := memArgs(args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		// round down address to 16 byte boundary
		adr &= ^uint32(15)
		// round up n to an integral multiple of 16 bytes
		size = (size + 15) & ^uint(15)
		// print the header
		c.User.Put("addr  0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F\n")
		// read and print the data
		for i := 0; i < int(size>>4); i++ {
			// read 16 bytes per line
			var data [16]string
			var ascii [16]string
			for j := 0; j < 16; j++ {
				x := c.User.(*userApp).mem.Rd8(adr + uint32(j))
				data[j] = fmt.Sprintf("%02x", x)
				if x >= 32 && x <= 126 {
					ascii[j] = fmt.Sprintf("%c", x)
				} else {
					ascii[j] = "."
				}
			}
			dataStr := strings.Join(data[:], " ")
			asciiStr := strings.Join(ascii[:], "")
			c.User.Put(fmt.Sprintf("%04x  %s  %s\n", adr, dataStr, asciiStr))
			adr += 16
		}
	},
}

//-----------------------------------------------------------------------------

var helpGo = []cli.Help{
	{"<adr>", "address (hex) - default is PC"},
}

// goArgs converts go arguments to an address value.
func goArgs(pc uint32, args []string) (uint32, error) {
	err := cli.CheckArgc(args, []int{0, 1})
	if err != nil {
		return 0, err
	}
	// address
	adr := int(pc)
	if len(args) >= 1 {
		adr, err = cli.IntArg(args[0], [2]int{0, 0xffffffff}, 16)
		if err != nil {
			return 0, err
		}
	}
	return uint32(adr), nil
}

var cmdGo = cli.Leaf{
	Descr: "run the emulation (no tracing)",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*userApp).cpu
		adr, err := goArgs(m.PC, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m.PC = adr
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
		adr, err := goArgs(m.PC, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m.PC = adr
		for true {
			s := m.Disassemble(m.PC)
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
		adr, err := goArgs(m.PC, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		m.PC = adr
		s := m.Disassemble(m.PC)
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

// daArgs converts disassembly arguments to an (address, size) tuple.
func daArgs(pc uint32, args []string) (uint32, int, error) {
	err := cli.CheckArgc(args, []int{0, 1, 2})
	if err != nil {
		return 0, 0, err
	}
	// address
	adr := int(pc) // default address
	if len(args) >= 1 {
		adr, err = cli.IntArg(args[0], [2]int{0, 0xffffffff}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	// size
	size := 64 // default size
	if len(args) >= 2 {
		size, err = cli.IntArg(args[1], [2]int{1, 2048}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	return uint32(adr), size, nil
}

var cmdDisassemble = cli.Leaf{
	Descr: "disassemble memory",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*userApp).cpu
		adr, size, err := daArgs(m.PC, args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		for size > 0 {
			da := m.Disassemble(adr)
			c.User.Put(fmt.Sprintf("%s\n", da))
			adr += uint32(da.Length)
			size -= da.Length
		}
	},
}

//-----------------------------------------------------------------------------

var cmdRegisters = cli.Leaf{
	Descr: "display cpu registers",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*userApp).cpu
		c.User.Put(fmt.Sprintf("%s\n", m.Dump()))
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

// root menu
var menuRoot = cli.Menu{
	{"da", cmdDisassemble, helpDisassemble},
	{"exit", cmdExit},
	{"go", cmdGo, helpGo},
	{"help", cmdHelp},
	{"history", cmdHistory, cli.HistoryHelp},
	{"md", cmdMemDisplay, helpMemDisplay},
	{"regs", cmdRegisters},
	{"reset", cmdReset},
	{"step", cmdStep, helpGo},
	{"trace", cmdTrace, helpGo},
}

//-----------------------------------------------------------------------------
