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

	"github.com/deadsy/riscv/mem"
	"github.com/deadsy/riscv/rv"
	"github.com/deadsy/riscv/util"
)

//-----------------------------------------------------------------------------

const stackSize = 8 << 10
const heapSize = 256 << 10
const insLimit = 20000

//-----------------------------------------------------------------------------

type testCase struct {
	baseName string
	testName string
	elfFile  string
	sigFile  string
	elfClass elf.Class
}

const sigSuffix = ".signature.output"
const elfSuffix = ".elf"

func getTestCases(testPath string) ([]*testCase, error) {
	x := []*testCase{}
	err := filepath.Walk(testPath, func(path string, info os.FileInfo, err error) error {
		class, err := util.GetELFClass(path)
		if err != nil {
			return nil
		}
		// Do we have a sig file?
		var sigFile string
		if strings.HasSuffix(path, elfSuffix) {
			sigFile = strings.TrimSuffix(path, elfSuffix) + sigSuffix
			if !util.FileExists(sigFile) {
				sigFile = ""
			}
		}
		tc := testCase{
			baseName: testPath,
			testName: strings.TrimPrefix(path, testPath+"/"),
			elfFile:  path,
			sigFile:  sigFile,
			elfClass: class,
		}
		x = append(x, &tc)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return x, nil
}

//-----------------------------------------------------------------------------

// testFixups applies per-test environment tweaks
func testFixups(m *rv.RV, name string) {
	switch name {
	case "rv32uc/rvc",
		"rv32i/I-AUIPC-01",
		"rv32Zifencei/I-FENCE.I-01",
		"rv32ui/fence_i":
		m.Mem.SetAttr(".text.init", mem.AttrRWX)
	case "rv32i/I-SB-01",
		"rv32i/I-SH-01",
		"rv32i/I-SW-01":
		m.Mem.Add(mem.NewSection(".fixup", 0x80001ffc, 4, mem.AttrRW))
	case "rv32mi/ma_addr",
		"rv32i/I-MISALIGN_LDST-01":
		m.Mem.SetAttr(".data", mem.AttrRWM)
	}
}

//-----------------------------------------------------------------------------

// compareSlice32 compares uint32 slices.
func compareSlice32(a, b []uint32) error {
	if len(a) != len(b) {
		return errors.New("len(a) != len(b)")
	}
	for i := range a {
		if a[i] != b[i] {
			return fmt.Errorf("x[%d] %08x (expected) != %08x (actual)", i, a[i], b[i])
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
	return m.Rd32Range(start, (end-start)>>2), nil
}

func checkExit(m *mem.Memory, errIn error) error {
	// we are expecting a breakpoint on write to "tohost".
	e := errIn.(*rv.Error)
	em := e.GetMemError()
	if em == nil {
		return errIn
	}
	if em.Type != mem.ErrWrite|mem.ErrBreak {
		return errIn
	}
	// get the symbol
	addr, err := m.SymbolGetAddress("tohost")
	if err != nil {
		return fmt.Errorf("\"tohost\" symbol not found, %s", errIn)
	}
	// check the address
	if em.Addr != addr {
		return fmt.Errorf("breakpoint not on \"tohost\", %s", errIn)
	}
	// check the exit status
	status, _ := m.Rd32(addr)
	if status != 1 {
		return fmt.Errorf("(%d), %s", status, errIn)
	}
	// looks good
	return nil
}

//-----------------------------------------------------------------------------

func (tc *testCase) Test() error {

	var cpu *rv.RV

	if tc.elfClass == elf.ELFCLASS32 {
		// Setup an RV32 CPU
		// create the ISA
		isa := rv.NewISA()
		err := isa.Add(rv.ISArv32gc)
		if err != nil {
			return err
		}
		// create the cpu
		cpu = rv.NewRV32(isa, mem.NewMem32(0), nil)
	} else {
		// Setup an RV64 CPU
		// create the ISA
		isa := rv.NewISA()
		err := isa.Add(rv.ISArv64gc)
		if err != nil {
			return err
		}
		// create the cpu
		cpu = rv.NewRV64(isa, mem.NewMem64(0), nil)
	}

	// load the elf file
	_, err := cpu.Mem.LoadELF(tc.elfFile, tc.elfClass)
	if err != nil {
		return err
	}

	// add a stack and heap
	cpu.Mem.Add(mem.NewSection("stack", (1<<32)-stackSize, stackSize, mem.AttrRW))
	cpu.Mem.Add(mem.NewSection("heap", 0x80000000, heapSize, mem.AttrRW))

	// Break on the "tohost" write (compliance tests).
	cpu.Mem.AddBreakPointByName("tohost", mem.AttrW)

	// apply per test fixups
	//testFixups(cpu, name)

	// run the emulation
	cpu.Reset()
	var ins int
	for ins = 0; ins < insLimit; ins++ {
		err = cpu.Run()
		if err != nil {
			break
		}
	}
	if ins == insLimit {
		return fmt.Errorf("reached instruction limit")
	}

	// check for a normal exit
	err = checkExit(cpu.Mem, err)
	if err != nil {
		return err
	}

	/*

		// get the test results from memory
		result, err := getResults(m)
		if err != nil {
			return err
		}
		// get the signature results from the file
		sig, err := getSignature(sigFilename)
		if err != nil {
			return err
		}
		// compare the result and signature
		err = compareSlice32(result, sig)
		if err != nil {
			return err
		}

	*/

	return nil
}

//-----------------------------------------------------------------------------

func main() {
	// command line flags
	path := flag.String("p", "test", "path to compliance tests")
	flag.Parse()
	testPath := filepath.Clean(*path)

	// get the test cases
	testCases, err := getTestCases(testPath)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// run the tests
	pass := 0
	fail := 0
	for _, tc := range testCases {
		err := tc.Test()
		fmt.Printf("%-30s ", tc.testName)
		if err != nil {
			fmt.Printf("FAIL %s\n", err)
			fail++
		} else {
			fmt.Printf("PASS\n")
			pass++
		}
	}

	// report aggregate results
	total := pass + fail
	if total != 0 {
		fmt.Printf("result: %d/%d passed (%d failed) score %.2f\n", pass, total, fail, float32(pass)/float32(total))
	} else {
		fmt.Printf("no tests run\n")
	}

	os.Exit(0)
}

//-----------------------------------------------------------------------------
