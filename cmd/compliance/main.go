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
const insLimit = 20000

//-----------------------------------------------------------------------------

const elfSuffix = ".elf"
const sigSuffix = ".signature.output"

func getTestCases(testPath string) ([]string, error) {
	x := []string{}
	err := filepath.Walk(testPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, elfSuffix) {
			testName := strings.TrimPrefix(path, testPath+"/")
			testName = strings.TrimSuffix(testName, elfSuffix)
			x = append(x, testName)
		}
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
	case "rv32mi/illegal",
		"rv32mi/shamt":
		m.SetHandler(rv.ErrIllegal)
	case "rv32i/I-EBREAK-01",
		"rv32mi/sbreak":
		m.SetHandler(rv.ErrBreak)
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

//-----------------------------------------------------------------------------

func TestRV32(base, name string) error {

	elfFilename := base + "/" + name + elfSuffix
	sigFilename := base + "/" + name + sigSuffix

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
	_, err = m.LoadELF(elfFilename, elf.ELFCLASS32)
	if err != nil {
		return err
	}

	// create the cpu
	cpu := rv.NewRV32(isa, m, ecall.NewCompliance())

	// apply per test fixups
	testFixups(cpu, name)

	// run the emulation
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
	e := err.(*rv.Error)
	if e.Type != rv.ErrExit {
		return err
	}

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
	err = compareSlice32(sig, result)
	if err != nil {
		return err
	}

	return nil
}

//-----------------------------------------------------------------------------

func TestRV64(base, name string) error {

	elfFilename := base + "/" + name + elfSuffix
	sigFilename := base + "/" + name + sigSuffix

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
	_, err = m.LoadELF(elfFilename, elf.ELFCLASS64)
	if err != nil {
		return err
	}

	// create the cpu
	cpu := rv.NewRV64(isa, m, ecall.NewCompliance())

	// apply per test fixups
	testFixups(cpu, name)

	// run the emulation
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
	e := err.(*rv.Error)
	if e.Type != rv.ErrExit {
		return err
	}

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

	return nil
}

//-----------------------------------------------------------------------------

func Test(base, name string) error {
	if strings.Contains(name, "rv32") {
		return TestRV32(base, name)
	}
	return TestRV64(base, name)
}

//-----------------------------------------------------------------------------

func main() {
	// command line flags
	path := flag.String("p", "test", "path to compliance tests")
	flag.Parse()
	testPath := filepath.Clean(*path)
	testCases, err := getTestCases(testPath)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	pass := 0
	fail := 0

	for _, name := range testCases {
		err := Test(testPath, name)
		testName := fmt.Sprintf("%s/%s", testPath, name)
		fmt.Printf("%-50s ", testName)
		if err != nil {
			fmt.Printf("FAIL %s\n", err)
			fail++
		} else {
			fmt.Printf("PASS\n")
			pass++
		}
	}

	total := pass + fail
	fmt.Printf("result: %d/%d passed (%d failed) score %.2f\n", pass, total, fail, float32(pass)/float32(total))

	os.Exit(0)
}

//-----------------------------------------------------------------------------
