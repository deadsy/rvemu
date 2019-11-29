//-----------------------------------------------------------------------------
/*

RISC-V Compliance Testing

*/
//-----------------------------------------------------------------------------

package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deadsy/riscv/ecall"
	"github.com/deadsy/riscv/mem"
	"github.com/deadsy/riscv/rv"
)

//-----------------------------------------------------------------------------

const stackSize = 8 << 10

//-----------------------------------------------------------------------------

type testCase struct {
	xlen        int
	elfFilename string
	sigFilename string
}

func getTestCases(path string) ([]*testCase, error) {
	cases := []*testCase{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".elf") {
			base := strings.TrimSuffix(path, ".elf")
			xlen := 32
			if strings.Contains(base, "rv64") {
				xlen = 64
			}
			tc := testCase{
				xlen:        xlen,
				elfFilename: base + ".elf",
				sigFilename: base + ".signature.output",
			}
			cases = append(cases, &tc)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cases, nil
}

//-----------------------------------------------------------------------------

func (c *testCase) TestRV32() error {
	fmt.Printf("Testing %s (RV32)\n", c.elfFilename)

	// create the ISA
	isa := rv.NewISA()
	err := isa.Add(rv.ISArv32gc)
	if err != nil {
		return err
	}

	// create the memory
	m := mem.NewMem32(0)
	m.Add(mem.NewSection("stack", (1<<32)-stackSize, stackSize, mem.AttrRW))

	// load the elf file
	status, err := m.LoadELF(c.elfFilename, elf.ELFCLASS32)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", status)

	// create the cpu
	cpu := rv.NewRV32(isa, m, ecall.NewCompliance())
	cpu.Reset()

	for true {
		err := cpu.Run()
		if err != nil {
			fmt.Printf("%s\n", err)
			break
		}
	}

	return nil
}

func (c *testCase) TestRV64() error {
	fmt.Printf("Testing %s (RV64)\n", c.elfFilename)
	return nil
}

func (c *testCase) Test() error {
	if c.xlen == 32 {
		return c.TestRV32()
	}
	return c.TestRV64()
}

//-----------------------------------------------------------------------------

func main() {
	// command line flags
	path := flag.String("p", "test", "path to compliance tests")
	flag.Parse()

	cases, err := getTestCases(*path)

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	for i := range cases {
		err := cases[i].Test()
		if err != nil {
			fmt.Printf("%s\n", err)
			break
		} else {
			fmt.Printf("PASS\n")
		}
	}

	os.Exit(0)
}

//-----------------------------------------------------------------------------
