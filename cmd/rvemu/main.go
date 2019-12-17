//-----------------------------------------------------------------------------
/*

RISC-V 32-bit Emulator

*/
//-----------------------------------------------------------------------------

package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"os"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/riscv/mem"
	"github.com/deadsy/riscv/rv"
	"github.com/deadsy/riscv/util"
)

//-----------------------------------------------------------------------------

const historyPath = ".rvemu_history"
const stackSize = 8 << 10

//-----------------------------------------------------------------------------

// emuApp is state associated with the emulator application.
type emuApp struct {
	mem      *mem.Memory
	cpu      *rv.RV
	elfClass elf.Class
	prompt   string
}

// newEmu32 returns a 32-bit emulator.
func newEmu32() (*emuApp, error) {
	// 32-bit ISA
	isa := rv.NewISA()
	err := isa.Add(rv.ISArv32gc)
	if err != nil {
		return nil, err
	}
	// 32-bit memory systems
	m := mem.NewMem32(0)
	// add a stack to the memory
	m.Add(mem.NewSection("stack", (1<<32)-stackSize, stackSize, mem.AttrRW))

	return &emuApp{
		mem:      m,
		cpu:      rv.NewRV32(isa, m, nil),
		elfClass: elf.ELFCLASS32,
		prompt:   "rv32> ",
	}, nil
}

// newEmu64 returns a 64-bit emulator.
func newEmu64() (*emuApp, error) {
	// 64-bit ISA
	isa := rv.NewISA()
	err := isa.Add(rv.ISArv64gc)
	if err != nil {
		return nil, err
	}
	// 64-bit memory system
	m := mem.NewMem64(0)
	// add a stack to the memory
	m.Add(mem.NewSection("stack", (1<<32)-stackSize, stackSize, mem.AttrRW))

	return &emuApp{
		mem:      m,
		cpu:      rv.NewRV64(isa, m, nil),
		elfClass: elf.ELFCLASS64,
		prompt:   "rv64> ",
	}, nil
}

//-----------------------------------------------------------------------------

// Put outputs a string to the user application.
func (u *emuApp) Put(s string) {
	fmt.Printf("%s", s)
}

//-----------------------------------------------------------------------------

func main() {
	// command line flags
	fname := flag.String("f", "out.bin", "file to load (ELF)")
	flag.Parse()

	elfClass, err := util.GetELFClass(*fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// create the application
	var app *emuApp
	switch elfClass {
	case elf.ELFCLASS32:
		app, err = newEmu32()
	case elf.ELFCLASS64:
		app, err = newEmu64()
	default:
		fmt.Fprintf(os.Stderr, "ELF class %d is not supported\n", elfClass)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// load the file
	status, err := app.mem.LoadELF(*fname, app.elfClass)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "%s\n", status)

	// Break on the "tohost" write (compliance tests).
	app.mem.AddBreakPointByName("tohost", mem.AttrW)

	// create the cli
	c := cli.NewCLI(app)
	c.HistoryLoad(historyPath)
	c.SetRoot(menuRoot)
	c.SetPrompt(app.prompt)

	// reset the cpu
	app.cpu.Reset()

	// run the cli
	for c.Running() {
		c.Run()
	}

	// exit
	c.HistorySave(historyPath)
	os.Exit(0)
}

//-----------------------------------------------------------------------------