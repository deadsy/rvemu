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
)

//-----------------------------------------------------------------------------

const historyPath = ".rv32emu_history"

//-----------------------------------------------------------------------------

// userApp is state associated with the user application.
type userApp struct {
	mem *mem.Memory
	cpu *rv.RV32
}

// newUserApp returns a user application.
func newUserApp() (*userApp, error) {

	// create the ISA
	isa := rv.NewISA()
	err := isa.Add(rv.ISArv32gc)
	if err != nil {
		return nil, err
	}

	// create the memory
	m := mem.NewMemory()
	m.Add(mem.NewChunk(0, 256<<10, mem.AttrRWX))
	//m.Add(mem.NewChunk(0x80000000, 64<<10, mem.AttrRW))
	m.Add(mem.NewEmpty(0, 1<<32, 0)) // no access

	// create the cpu
	cpu := rv.NewRV32(isa, m)

	return &userApp{
		mem: m,
		cpu: cpu,
	}, nil
}

//-----------------------------------------------------------------------------

// Put outputs a string to the user application.
func (u *userApp) Put(s string) {
	fmt.Printf("%s", s)
}

//-----------------------------------------------------------------------------

func main() {
	// command line flags
	fname := flag.String("f", "out.bin", "file to load (ELF)")
	flag.Parse()

	// create the application
	app, err := newUserApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// load the file
	status, err := app.mem.LoadELF(*fname, elf.ELFCLASS32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "%s\n", status)

	// create the cli
	c := cli.NewCLI(app)
	c.HistoryLoad(historyPath)
	c.SetRoot(menuRoot)
	c.SetPrompt("rv32> ")

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
