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

	"github.com/deadsy/riscv/csr"
	"github.com/deadsy/riscv/host"
	"github.com/deadsy/riscv/mem"
	"github.com/deadsy/riscv/rv"
	"github.com/deadsy/riscv/util"
)

//-----------------------------------------------------------------------------

const heapSize = 1 << 20
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

// Fixups applies per-test environment tweaks
func (tc *testCase) Fixups(m *rv.RV) {
	switch tc.testName {
	case "rv32uc/rvc.elf",
		"rv32i/I-AUIPC-01.elf",
		"rv32Zifencei/I-FENCE.I-01.elf",
		"rv32ui/fence_i.elf",
		"rv32ui-p-fence_i",
		"rv64ui-p-fence_i":
		m.Mem.SetAttr(".text.init", mem.AttrRWX)
	case "rv32mi/ma_addr.elf",
		"rv32mi-p-ma_addr",
		"rv64mi-p-ma_addr":
		m.Mem.SetAttr(".data", mem.AttrRWM)
	}
}

//-----------------------------------------------------------------------------

// compareSlice compares uint slices.
func compareSlice(a, b []uint) error {
	if len(a) != len(b) {
		return fmt.Errorf("len(a) %d != len(b) %d", len(a), len(b))
	}
	for i := range a {
		if a[i] != b[i] {
			return fmt.Errorf("x[%d] %x (expected) != %x (actual)", i, a[i], b[i])
		}
	}
	return nil
}

// getSignature reads the signature file.
func getSignature(filename string) ([]uint, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	x := strings.Split(string(buf), "\n")
	x = x[:len(x)-1]
	sig := make([]uint, len(x))
	for i := range sig {
		k, err := strconv.ParseUint(x[i], 16, 32)
		if err != nil {
			return nil, err
		}
		sig[i] = uint(k)
	}
	return sig, nil
}

// getResults gets the test results from memory.
func getResults(m *mem.Memory) ([]uint, error) {
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
	return m.RdBuf(start, (end-start)>>2, 32, false), nil
}

func checkExit(tohost *host.Host, errIn error) error {
	// we are expecting a breakpoint on write to "tohost".
	e := errIn.(*rv.Error)
	em := e.GetMemError()
	if em == nil {
		return errIn
	}
	if em.Type != mem.ErrWrite|mem.ErrBreak {
		return errIn
	}
	// check the exit status
	if !tohost.Passed() {
		return fmt.Errorf("%s", tohost)
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
		isa := rv.NewISA(csr.IsaExtS | csr.IsaExtU)
		err := isa.Add(rv.ISArv32gc)
		if err != nil {
			return err
		}
		// create the cpu
		csr := csr.NewState(32, isa.GetExtensions())
		mem := mem.NewMem32(csr, 0)
		cpu = rv.NewRV32(isa, mem, csr)
	} else {
		// Setup an RV64 CPU
		// create the ISA
		isa := rv.NewISA(csr.IsaExtS | csr.IsaExtU)
		err := isa.Add(rv.ISArv64gc)
		if err != nil {
			return err
		}
		// create the cpu
		csr := csr.NewState(64, isa.GetExtensions())
		mem := mem.NewMem64(csr, 0)
		cpu = rv.NewRV64(isa, mem, csr)
	}

	// load the elf file
	_, err := cpu.Mem.LoadELF(tc.elfFile, tc.elfClass)
	if err != nil {
		return err
	}

	// add a heap
	cpu.Mem.Add(mem.NewSection("heap", 0x80000000, heapSize, mem.AttrRW))

	// Callback on the "tohost" write (compliance tests).
	var tohost *host.Host
	sym := cpu.Mem.SymbolByName("tohost")
	if sym != nil {
		tohost = host.NewHost(cpu.Mem, sym.Addr)
		if sym.Size == 8 {
			// trap on a write to the most significant word
			fn := func(bp *mem.BreakPoint) bool { return tohost.To64(bp) }
			cpu.Mem.AddBreakPoint(sym.Name, sym.Addr+4, mem.AttrW, fn)
		} else {
			// 32-bit variable
			fn := func(bp *mem.BreakPoint) bool { return tohost.To32(bp) }
			cpu.Mem.AddBreakPoint(sym.Name, sym.Addr, mem.AttrW, fn)
		}
	}

	// apply per test fixups
	tc.Fixups(cpu)

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
	err = checkExit(tohost, err)
	if err != nil {
		return err
	}

	// do we have a signature file?
	if tc.sigFile == "" {
		return nil
	}

	// get the test results from memory
	result, err := getResults(cpu.Mem)
	if err != nil {
		return err
	}
	// get the signature results from the file
	sig, err := getSignature(tc.sigFile)
	if err != nil {
		return err
	}
	// compare the result and signature
	err = compareSlice(sig, result)
	if err != nil {
		return err
	}

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
		fmt.Printf("%-80s ", tc.elfFile)
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
