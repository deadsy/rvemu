//-----------------------------------------------------------------------------
/*

RISC-V Emulator

*/
//-----------------------------------------------------------------------------

package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/riscv/mem"
	"github.com/deadsy/riscv/rv"
)

//-----------------------------------------------------------------------------

const historyPath = "history.txt"

//-----------------------------------------------------------------------------

// userApp is state associated with the user application.
type userApp struct {
	mem *mem.Memory
	cpu *rv.RV32
}

// newUserApp returns a user application.
func newUserApp() (*userApp, error) {

	// create the ISA
	isa := rv.NewISA("rv32g")
	err := isa.Add(rv.ISArv32i, rv.ISArv32m, rv.ISArv32a, rv.ISArv32f, rv.ISArv32d, rv.ISArv32c)
	if err != nil {
		return nil, err
	}

	// create the memory
	mem := mem.NewMemory(0, 256<<10, false)

	// create the cpu
	cpu := rv.NewRV32(isa, mem)

	return &userApp{
		mem: mem,
		cpu: cpu,
	}, nil
}

// loadRaw loads a raw binary file.
func (u *userApp) loadRaw(filename string, x []uint8) (string, error) {
	// copy the code to the load address
	var loadAdr uint32
	for i, v := range x {
		u.mem.Wr8(loadAdr+uint32(i), v)
	}
	endAdr := loadAdr + uint32(len(x)) - 1
	return fmt.Sprintf("%s code %08x-%08x", filename, loadAdr, endAdr), nil
}

func (u *userApp) loadFile(filename string) (string, error) {
	// get the file contents
	x, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return u.loadRaw(filename, x)
}

//-----------------------------------------------------------------------------

// loadELF loads an ELF file.
func (u *userApp) loadELF(filename string) (string, error) {

	f, err := elf.Open(filename)
	if err != nil {
		return "", err
	}

	status := fmt.Sprintf("%+v", f)

	f.Close()

	return status, nil
}

//-----------------------------------------------------------------------------

// Put outputs a string to the user application.
func (u *userApp) Put(s string) {
	fmt.Printf("%s", s)
}

//-----------------------------------------------------------------------------

func main() {
	// command line flags
	fname := flag.String("f", "out.bin", "file to load (raw binary)")
	flag.Parse()

	// create the application
	app, err := newUserApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// load the file
	//status, err := app.loadFile(*fname)
	status, err := app.loadELF(*fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "%s\n", status)

	// create the cli
	c := cli.NewCLI(app)
	c.HistoryLoad(historyPath)
	c.SetRoot(menuRoot)
	c.SetPrompt("emu> ")

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
