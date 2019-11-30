//-----------------------------------------------------------------------------
/*

RISC-V Compliance Testing

*/
//-----------------------------------------------------------------------------

package main

import (
	"debug/elf"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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

// compareSlice32 compares uint32 slices.
func compareSlice32(a, b []uint32) error {
	if len(a) != len(b) {
		return errors.New("len(a) != len(b)")
	}
	for i := range a {
		if a[i] != b[i] {
			return fmt.Errorf("a[%d] != b[%d], %08x != %08x", i, i, a[i], b[i])
		}
	}
	return nil
}

// getSignature reads the signature file.
func getSignature(filename string) ([]uint32, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	x := strings.Split(string(buf), "\n")
	x = x[:len(x)-1]
	sig := make([]uint32, len(x))
	for i := range sig {
		k, err := strconv.ParseUint(x[i], 16, 32)
		if err != nil {
			return nil, err
		}
		sig[i] = uint32(k)
	}
	return sig, nil
}

// getResults gets the test results from memory.
func getResults(m *mem.Memory) ([]uint32, error) {
	start, err := m.SymbolGetAddress("begin_signature")
	if err != nil {
		return nil, err
	}
	end, err := m.SymbolGetAddress("end_signature")
	if err != nil {
		return nil, err
	}
	if start >= end {
		return nil, errors.New("result length <= 0")
	}
	return m.RangeRd32(start, (end-start)>>2), nil
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
	_, err = m.LoadELF(c.elfFilename, elf.ELFCLASS32)
	if err != nil {
		return err
	}

	// create the cpu
	cpu := rv.NewRV32(isa, m, ecall.NewCompliance())

	// run the emulation
	for true {
		err = cpu.Run()
		if err != nil {
			break
		}
	}

	// check for a normal exit
	ex := err.(*rv.Exception)
	if ex.N != rv.ExExit {
		return err
	}

	// get the test results from memory
	result, err := getResults(m)
	if err != nil {
		return err
	}
	// get the signature results from the file
	sig, err := getSignature(c.sigFilename)
	if err != nil {
		return err
	}
	// compare the result and signature
	err = compareSlice32(result, sig)
	if err != nil {
		return err
	}

	return nil
}

func (c *testCase) TestRV64() error {
	fmt.Printf("Testing %s (RV64)\n", c.elfFilename)

	// create the ISA
	isa := rv.NewISA()
	err := isa.Add(rv.ISArv64gc)
	if err != nil {
		return err
	}

	// create the memory
	m := mem.NewMem64(0)
	m.Add(mem.NewSection("stack", (1<<32)-stackSize, stackSize, mem.AttrRW))

	// load the elf file
	_, err = m.LoadELF(c.elfFilename, elf.ELFCLASS64)
	if err != nil {
		return err
	}

	// create the cpu
	cpu := rv.NewRV64(isa, m, ecall.NewCompliance())

	// run the emulation
	for true {
		err = cpu.Run()
		if err != nil {
			break
		}
	}

	// check for a normal exit
	ex := err.(*rv.Exception)
	if ex.N != rv.ExExit {
		return err
	}

	// get the test results from memory
	result, err := getResults(m)
	if err != nil {
		return err
	}
	// get the signature results from the file
	sig, err := getSignature(c.sigFilename)
	if err != nil {
		return err
	}
	// compare the result and signature
	err = compareSlice32(result, sig)
	if err != nil {
		return err
	}

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
			fmt.Printf("FAIL %s\n\n", err)
		} else {
			fmt.Printf("PASS\n\n")
		}
	}

	os.Exit(0)
}

//-----------------------------------------------------------------------------
