//-----------------------------------------------------------------------------
/*

RISC-V Control and Status Register Definitions

*/
//-----------------------------------------------------------------------------

package csr

import (
	"errors"
	"fmt"
	"strings"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/riscv/util"
)

//-----------------------------------------------------------------------------

// Error is a CSR access error.
type Error struct {
	reg uint // csr register number
	n   uint // bitmap of csr access errors
}

// CSR error bits.
const (
	ErrTodo      = 1 << iota // csr not implemented
	ErrPrivilege             // insufficient privilege
	ErrReadOnly              // trying to write a read-only register
	ErrNoRead                // no read function (todo)
	ErrNoWrite               // no write function (todo)
)

func (e *Error) Error() string {
	s := []string{}
	if e.n&ErrTodo != 0 {
		s = append(s, "not implemented")
	}
	if e.n&ErrPrivilege != 0 {
		s = append(s, "insufficient privilege")
	}
	if e.n&ErrReadOnly != 0 {
		s = append(s, "read only")
	}
	if e.n&ErrNoRead != 0 {
		s = append(s, "no read function")
	}
	if e.n&ErrNoWrite != 0 {
		s = append(s, "no write function")
	}
	return fmt.Sprintf("%s: %s", Name(e.reg), strings.Join(s, ","))
}

//-----------------------------------------------------------------------------

// Mode is the processor mode (user/supervisor/machine).
type Mode uint

// Modes
const (
	ModeU Mode = 0 // user
	ModeS      = 1 // supervisor
	ModeM      = 3 // machine
)

func (m Mode) String() string {
	return [4]string{"user", "supervisor", "?", "machine"}[m]
}

// ModeArg converts a mode argument string to a mode value.
func ModeArg(arg string) (Mode, error) {
	arg = strings.ToLower(arg)
	mode, ok := map[string]Mode{"u": ModeU, "s": ModeS, "m": ModeM}[arg]
	if !ok {
		return 0, fmt.Errorf("mode \"%s\" is not valid", arg)
	}
	return mode, nil
}

// GetMode returns the current processor mode.
func (s *State) GetMode() Mode {
	return s.mode
}

// setMode sets the current processor mode.
func (s *State) setMode(mode Mode) {
	s.mode = mode
}

// hasMode returns true if the mode is supported.
func (s *State) hasMode(mode Mode) bool {
	switch mode {
	case ModeU:
		return (s.misa & IsaExtU) != 0
	case ModeS:
		return (s.misa & IsaExtS) != 0
	case ModeM:
		return true
	}
	return false
}

func (s *State) getMinMode() Mode {
	if s.hasMode(ModeU) {
		return ModeU
	}
	return ModeM
}

// getNextMode returns the target mode for an exception/interrupt.
func (s *State) getNextMode(code uint, isInterrupt bool) Mode {
	var sMask, mMask uint
	if isInterrupt {
		mMask = s.mideleg
		sMask = s.sideleg
	} else {
		mMask = s.medeleg
		sMask = s.sedeleg
	}
	var nextMode Mode
	// get target mode implied by delegation registers
	if mMask&(1<<code) == 0 {
		nextMode = ModeM
	} else if sMask&(1<<code) == 0 {
		nextMode = ModeS
	} else {
		nextMode = ModeU
	}
	// exception cannot be taken to lower-privilege mode
	if nextMode < s.mode {
		return s.mode
	}
	return nextMode
}

//-----------------------------------------------------------------------------

// Register numbers for specific CSRs.
const (
	FFLAGS  = 0x001
	FRM     = 0x002
	FCSR    = 0x003
	SSTATUS = 0x100
	SEDELEG = 0x102
	SIDELEG = 0x103
	MSTATUS = 0x300
	MEDELEG = 0x302
	MIDELEG = 0x303
	MTVEC   = 0x305
	MEPC    = 0x341
	MCAUSE  = 0x342
	MTVAL   = 0x343
)

//-----------------------------------------------------------------------------

// ICode is a RISC-V interrupt code.
type ICode uint

// Interrupt Codes
const (
	IntUserSoftware       ICode = 0  // User software interrupt
	IntSupervisorSoftware ICode = 1  // Supervisor software interrupt
	IntMachineSoftware    ICode = 3  // Machine software interrupt
	IntUserTimer          ICode = 4  // User timer interrupt
	IntSupervisorTimer    ICode = 5  // Supervisor timer interrupt
	IntMachineTimer       ICode = 7  // Machine timer interrupt
	IntUserExternal       ICode = 8  // User external interrupt
	IntSupervisorExternal ICode = 9  // Supervisor external interrupt
	IntMachineExternal    ICode = 11 // Machine external interrupt
)

func (n ICode) String() string {
	return [16]string{
		"IntUserSoftware",       // 0
		"IntSupervisorSoftware", // 1
		"IntUnknown(2)",         // 2
		"IntMachineSoftware",    // 3
		"IntUserTimer",          // 4
		"IntSupervisorTimer",    // 5
		"IntUnknown(6)",         // 6
		"IntMachineTimer",       // 7
		"IntUserExternal",       // 8
		"IntSupervisorExternal", // 9
		"IntUnknown(10)",        // 10
		"IntMachineExternal",    // 11
		"IntUnknown(12)",        // 12
		"IntUnknown(13)",        // 13
		"IntUnknown(14)",        // 14
		"IntUnknown(15)",        // 15
	}[n]
}

// ECode is a RISC-V exception code.
type ECode uint

// Exception Codes
const (
	ExInsAddrMisaligned         ECode = 0  // Instruction address misaligned
	ExInsAccessFault            ECode = 1  // Instruction access fault
	ExInsIllegal                ECode = 2  // Illegal instruction
	ExBreakpoint                ECode = 3  // Breakpoint
	ExLoadAddrMisaligned        ECode = 4  // Load address misaligned
	ExLoadAccessFault           ECode = 5  // Load access fault
	ExStoreAddrMisaligned       ECode = 6  // Store/AMO address misaligned
	ExStoreAccessFault          ECode = 7  // Store/AMO access fault
	ExEnvCallFromUserMode       ECode = 8  // Environment call from U-mode
	ExEnvCallFromSupervisorMode ECode = 9  // Environment call from S-mode
	ExUnknown10                 ECode = 10 // unknown 10
	ExEnvCallFromMachineMode    ECode = 11 // Environment call from M-mode
	ExInsPageFault              ECode = 12 // Instruction page fault
	ExLoadPageFault             ECode = 13 // Load page fault
	ExUnknown14                 ECode = 14 // unknown 14
	ExStorePageFault            ECode = 15 // Store/AMO page fault
)

func (n ECode) String() string {
	return [16]string{
		"ExInsAddrMisaligned",         // 0
		"ExInsAccessFault",            // 1
		"ExInsIllegal",                // 2
		"ExBreakpoint",                // 3
		"ExLoadAddrMisaligned",        // 4
		"ExLoadAccessFault",           // 5
		"ExStoreAddrMisaligned",       // 6
		"ExStoreAccessFault",          // 7
		"ExEnvCallFromUserMode",       // 8
		"ExEnvCallFromSupervisorMode", // 9
		"ExUnknown(10)",               // 10
		"ExEnvCallFromMachineMode",    // 11
		"ExInsPageFault",              // 12
		"ExLoadPageFault",             // 13
		"ExUnknown(14)",               // 14
		"ExStorePageFault",            // 15
	}[n]
}

//-----------------------------------------------------------------------------

// wrIgnore is a no-op write function.
func wrIgnore(s *State, val uint) {
}

// rdZero always reads the CSR as zero.
func rdZero(s *State) uint {
	return 0
}

//-----------------------------------------------------------------------------
// u/s/m cause register

func (s *State) setCause(ecode uint, isInterrupt bool, mode Mode) {
	cause := ecode
	if isInterrupt {
		cause |= 1 << (s.mxlen - 1)
	}
	switch mode {
	case ModeU:
		s.ucause = cause
	case ModeS:
		s.scause = cause
	case ModeM:
		s.mcause = cause
	}
}

func rdMCAUSE(s *State) uint {
	return s.mcause
}

func rdSCAUSE(s *State) uint {
	return s.scause
}

func rdUCAUSE(s *State) uint {
	return s.ucause
}

//-----------------------------------------------------------------------------
// machine isa register

// ISA Extension Bitmap
const (
	IsaExtA = (1 << iota) // Atomic extension
	IsaExtB               // Tentatively reserved for Bit-Manipulation extension
	IsaExtC               // Compressed extension
	IsaExtD               // Double-precision floating-point extension
	IsaExtE               // RV32E base ISA
	IsaExtF               // Single-precision floating-point extension
	IsaExtG               // Additional standard extensions present
	IsaExtH               // Hypervisor extension
	IsaExtI               // RV32I/64I/128I base ISA
	IsaExtJ               // Tentatively reserved for Dynamically Translated Languages extension
	IsaExtK               // Reserved
	IsaExtL               // Tentatively reserved for Decimal Floating-Point extension
	IsaExtM               // Integer Multiply/Divide extension
	IsaExtN               // User-level interrupts supported
	IsaExtO               // Reserved
	IsaExtP               // Tentatively reserved for Packed-SIMD extension
	IsaExtQ               // Quad-precision floating-point extension
	IsaExtR               // Reserved
	IsaExtS               // Supervisor mode implemented
	IsaExtT               // Tentatively reserved for Transactional Memory extension
	IsaExtU               // User mode implemented
	IsaExtV               // Tentatively reserved for Vector extension
	IsaExtW               // Reserved
	IsaExtX               // Non-standard extensions present
	IsaExtY               // Reserved
	IsaExtZ               // Reserved
)

func mxl(xlen uint) uint {
	return map[uint]uint{32: 1, 64: 2, 128: 3}[xlen]
}

func initMISA(s *State, ext uint) {
	s.misa = (mxl(s.xlen) << (s.mxlen - 2)) | ext
	s.ialign = []uint{32, 16}[util.BoolToInt(s.misa&IsaExtC != 0)]
}

func rdMISA(s *State) uint {
	return s.misa
}

func wrMISA(s *State, val uint) {

	// From the Spec:
	// 1. Writing misa may increase IALIGN, e.g., by disabling the “C” extension. If an instruction that
	// would write misa increases IALIGN, and the subsequent instruction’s address is not IALIGN-bit
	// aligned, the write to misa is suppressed, leaving misa unchanged.

	// 2. If a write to misa causes MXLEN to change, the position of
	// MXL moves to the most-significant two bits of misa at the new width.

	// TODO I don't easily know the PC alignment in CSR, so I'm ignoring MISA writes for now.
	// s.misa = val
	// s.ialign = []uint{32,16}[util.BoolToInt(misa&IsaExtC != 0)]
}

func fmtMXL(x uint) string {
	return []string{"?", "32", "64", "128"}[x]
}

func fmtExtensions(x uint) string {
	s := []rune{}
	for i := 0; i < 26; i++ {
		if x&1 != 0 {
			s = append(s, 'a'+rune(i))
		}
		x >>= 1
	}
	if len(s) != 0 {
		return fmt.Sprintf("\"%s\"", string(s))
	}
	return "none"
}

func displayMISA(s *State) string {
	fs := util.FieldSet{
		{"mxl", s.mxlen - 1, s.mxlen - 2, fmtMXL},
		{"extensions", 25, 0, fmtExtensions},
	}
	return fs.Display(s.misa)
}

//-----------------------------------------------------------------------------
// u/s/m exception program counter

func wrUEPC(s *State, val uint) {
	s.uepc = val & ^uint(1)
}

func wrSEPC(s *State, val uint) {
	s.sepc = val & ^uint(1)
}

func wrMEPC(s *State, val uint) {
	s.mepc = val & ^uint(1)
}

func rdUEPC(s *State) uint {
	if s.ialign == 32 {
		return s.uepc & ^uint(3)
	}
	return s.uepc
}

func rdSEPC(s *State) uint {
	if s.ialign == 32 {
		return s.sepc & ^uint(3)
	}
	return s.sepc
}

func rdMEPC(s *State) uint {
	if s.ialign == 32 {
		return s.mepc & ^uint(3)
	}
	return s.mepc
}

func (s *State) setEPC(pc uint64, mode Mode) {
	epc := uint(pc & ^uint64(1))
	switch mode {
	case ModeU:
		s.uepc = epc
	case ModeS:
		s.sepc = epc
	case ModeM:
		s.mepc = epc
	}
}

//-----------------------------------------------------------------------------
// u/s/m trap value register

func wrUTVAL(s *State, val uint) {
	s.utval = val
}

func wrSTVAL(s *State, val uint) {
	s.stval = val
}

func wrMTVAL(s *State, val uint) {
	s.mtval = val
}

func rdUTVAL(s *State) uint {
	return s.utval
}

func rdSTVAL(s *State) uint {
	return s.stval
}

func rdMTVAL(s *State) uint {
	return s.mtval
}

func (s *State) setTrapValue(val uint, mode Mode) {
	switch mode {
	case ModeU:
		s.utval = val
	case ModeS:
		s.stval = val
	case ModeM:
		s.mtval = val
	}
}

//-----------------------------------------------------------------------------
// u/s/m trap vector

func wrUTVEC(s *State, val uint) {
	s.utvec = val
}

func rdUTVEC(s *State) uint {
	return s.utvec
}

func wrSTVEC(s *State, val uint) {
	s.stvec = val
}

func rdSTVEC(s *State) uint {
	return s.stvec
}

func wrMTVEC(s *State, val uint) {
	s.mtvec = val
}

func rdMTVEC(s *State) uint {
	return s.mtvec
}

func fmtBase(x uint) string {
	return fmt.Sprintf("%x", x<<2)
}

func fmtTrapVectorMode(x uint) string {
	m := map[uint]string{0: "direct", 1: "vectored"}
	return util.DisplayEnum(x, m, "reserved")
}

func displaySTVEC(s *State) string {
	fs := util.FieldSet{
		{"base", s.sxlen - 1, 2, fmtBase},
		{"mode", 1, 0, fmtTrapVectorMode},
	}
	return fs.Display(s.stvec)
}

func displayMTVEC(s *State) string {
	fs := util.FieldSet{
		{"base", s.mxlen - 1, 2, fmtBase},
		{"mode", 1, 0, fmtTrapVectorMode},
	}
	return fs.Display(s.mtvec)
}

// getTrapVector returns the base/mode of a u/s/m trap vector.
func (s *State) getTrapVector(mode Mode) (uint, uint) {
	var tvec, msb uint
	switch mode {
	case ModeU:
		tvec = s.utvec
		msb = s.uxlen - 1
	case ModeS:
		tvec = s.stvec
		msb = s.sxlen - 1
	case ModeM:
		tvec = s.mtvec
		msb = s.mxlen - 1
	}
	return util.MaskBits(tvec, msb, 2), util.MaskBits(tvec, 1, 0)
}

//-----------------------------------------------------------------------------
// u/s/m scratch

func wrUSCRATCH(s *State, val uint) {
	s.uscratch = val
}

func wrSSCRATCH(s *State, val uint) {
	s.sscratch = val
}

func wrMSCRATCH(s *State, val uint) {
	s.mscratch = val
}

func rdUSCRATCH(s *State) uint {
	return s.uscratch
}

func rdSSCRATCH(s *State) uint {
	return s.sscratch
}

func rdMSCRATCH(s *State) uint {
	return s.mscratch
}

//-----------------------------------------------------------------------------
// fcsr

const frmMask = uint(7 << 5)
const fflagsMask = uint(31 << 0)
const fcsrMask = frmMask | fflagsMask

func wrFCSR(s *State, val uint) {
	s.fcsr = val & fcsrMask
}

func rdFCSR(s *State) uint {
	return s.fcsr
}

func wrFRM(s *State, val uint) {
	s.fcsr &= ^frmMask
	s.fcsr |= (val & 7) << 5
}

func rdFRM(s *State) uint {
	return (s.fcsr & frmMask) >> 5
}

func wrFFLAGS(s *State, val uint) {
	s.fcsr &= ^fflagsMask
	s.fcsr |= val & fflagsMask
}

func rdFFLAGS(s *State) uint {
	return s.fcsr & fflagsMask
}

//-----------------------------------------------------------------------------
// u/s/m status

// xsState is the state of the float/user extension (XS/FS bitfields in mstatus)
type xsState uint

const (
	xsOff   xsState = 0 // extension is off
	xsInit  xsState = 1 // extension is in the initial state
	xsClean xsState = 2 // extension state not modified since context switch
	xsDirty xsState = 3 // extension state modified since context switch
)

func (s xsState) String() string {
	return []string{"off", "init", "clean", "dirty"}[s]
}

// IsFloatOff returns true if the floating point has been disabled in mstatus.fs.
func (s *State) IsFloatOff() bool {
	return s.mstatusRdFS() == uint(xsOff)
}

const tsrMask = (1 << 22)
const twMask = (1 << 21)
const tvmMask = (1 << 20)
const mxrMask = (1 << 19)
const sumMask = (1 << 18)
const mprvMask = (1 << 17)
const xsMask = (3 << 15)
const fsMask = (3 << 13)
const mppMask = (3 << 11)
const sppMask = (1 << 8)
const mpieMask = (1 << 7)
const spieMask = (1 << 5)
const upieMask = (1 << 4)
const mieMask = (1 << 3)
const sieMask = (1 << 1)
const uieMask = (1 << 0)
const uxlMask = (3 << 32)
const sxlMask = (3 << 34)

type mStatus struct {
	val      uint // u/s/m status CSR
	wpriMask uint // read WPRI fields as 0
	uMask    uint // bits seen in user mode
	sMask    uint // bits seen in supervisor mode
}

func (m *mStatus) init(mxlen uint) {

	// mark the initial state for FS and XS
	m.val = (uint(xsInit) << 15 /*XS*/) | (uint(xsInit) << 13 /*FS*/)

	// sxlen == uxlen == mxlen
	if mxlen == 64 {
		m.val |= (2 << 32 /*UXL*/) | (2 << 34 /*SXL*/)
	}

	// set up the access masks
	m.uMask = uieMask | upieMask
	m.sMask = uieMask | upieMask | sppMask | spieMask | sieMask | sumMask | xsMask | fsMask | mxrMask
	m.wpriMask = util.BitMask(2, 2) | util.BitMask(6, 6) | util.BitMask(10, 9) | util.BitMask(30, 23)
	if mxlen == 32 {
		m.sMask |= (1 << 31 /*SD*/)
	} else {
		m.uMask |= uxlMask
		m.sMask |= uxlMask | sxlMask | (1 << 63 /*SD*/)
		m.wpriMask |= (util.BitMask(31, 31) | util.BitMask(62, 36))
	}
}

func (m *mStatus) wr(x uint, mode Mode) {
	x &= ^m.wpriMask

	// UXL is WARL
	if x&uxlMask == 0 {
		// preserve the existing value
		x |= m.val & uxlMask
	}

	// SXL is WARL
	if x&sxlMask == 0 {
		// preserve the existing value
		x |= m.val & sxlMask
	}

	switch mode {
	case ModeU:
		m.val = (m.val & ^m.uMask) | (x & m.uMask)
	case ModeS:
		m.val = (m.val & ^m.sMask) | (x & m.sMask)
	case ModeM:
		m.val = x
	}
}

func (m *mStatus) rd(mode Mode) uint {
	switch mode {
	case ModeU:
		return m.val & m.uMask
	case ModeS:
		return m.val & m.sMask
	case ModeM:
		return m.val
	}
	return 0
}

func wrUSTATUS(s *State, x uint) {
	s.mstatus.wr(x, ModeU)
}

func rdUSTATUS(s *State) uint {
	return s.mstatus.rd(ModeU)
}

func wrSSTATUS(s *State, x uint) {
	s.mstatus.wr(x, ModeS)
}

func rdSSTATUS(s *State) uint {
	return s.mstatus.rd(ModeS)
}

func wrMSTATUS(s *State, x uint) {
	s.mstatus.wr(x, ModeM)
}

func rdMSTATUS(s *State) uint {
	return s.mstatus.rd(ModeM)
}

// GetMPRV returns the MPRV bit of mstatus.
func (s *State) GetMPRV() bool {
	return s.mstatus.rd(ModeM)&mprvMask != 0
}

// GetSUM returns the sum bit of mstatus.
func (s *State) GetSUM() bool {
	return s.mstatus.rd(ModeM)&sumMask != 0
}

// GetMXR returns the MXR bit of mstatus.
func (s *State) GetMXR() bool {
	return s.mstatus.rd(ModeM)&mxrMask != 0
}

// GetMPP returns the MPP bits of mstatus.
func (s *State) GetMPP() Mode {
	return Mode((s.mstatus.rd(ModeM) >> 11 /*MPP*/) & 3)
}

func fmtXS(x uint) string {
	return xsState(x).String()
}

func displaySSTATUS(s *State) string {
	var fs util.FieldSet
	if s.sxlen == 32 {
		// RV32
		fs = util.FieldSet{
			{"sd", 31, 31, util.FmtDec},
			{"mxr", 19, 19, util.FmtDec},
			{"sum", 18, 18, util.FmtDec},
			{"xs", 16, 15, fmtXS},
			{"fs", 14, 13, fmtXS},
			{"spp", 8, 8, util.FmtDec},
			{"spie", 5, 5, util.FmtDec},
			{"upie", 4, 4, util.FmtDec},
			{"sie", 1, 1, util.FmtDec},
			{"uie", 0, 0, util.FmtDec},
		}
	} else {
		// RV64
		fs = util.FieldSet{
			{"sd", s.mxlen - 1, s.mxlen - 1, util.FmtDec},
			{"sxl", 35, 34, util.FmtDec},
			{"uxl", 33, 32, util.FmtDec},
			{"mxr", 19, 19, util.FmtDec},
			{"sum", 18, 18, util.FmtDec},
			{"xs", 16, 15, fmtXS},
			{"fs", 14, 13, fmtXS},
			{"spp", 8, 8, util.FmtDec},
			{"spie", 5, 5, util.FmtDec},
			{"upie", 4, 4, util.FmtDec},
			{"sie", 1, 1, util.FmtDec},
			{"uie", 0, 0, util.FmtDec},
		}
	}
	return fs.Display(s.mstatus.rd(ModeS))
}

func displayMSTATUS(s *State) string {
	var fs util.FieldSet
	if s.sxlen == 32 {
		// RV32
		fs = util.FieldSet{
			{"sd", 31, 31, util.FmtDec},
			{"tsr", 22, 22, util.FmtDec},
			{"tw", 21, 21, util.FmtDec},
			{"tvm", 20, 20, util.FmtDec},
			{"mxr", 19, 19, util.FmtDec},
			{"sum", 18, 18, util.FmtDec},
			{"mprv", 17, 17, util.FmtDec},
			{"xs", 16, 15, fmtXS},
			{"fs", 14, 13, fmtXS},
			{"mpp", 12, 11, util.FmtDec},
			{"spp", 8, 8, util.FmtDec},
			{"mpie", 7, 7, util.FmtDec},
			{"spie", 5, 5, util.FmtDec},
			{"upie", 4, 4, util.FmtDec},
			{"mie", 3, 3, util.FmtDec},
			{"sie", 1, 1, util.FmtDec},
			{"uie", 0, 0, util.FmtDec},
		}
	} else {
		// RV64
		fs = util.FieldSet{
			{"sd", s.mxlen - 1, s.mxlen - 1, util.FmtDec},
			{"sxl", 35, 34, util.FmtDec},
			{"uxl", 33, 32, util.FmtDec},
			{"tsr", 22, 22, util.FmtDec},
			{"tw", 21, 21, util.FmtDec},
			{"tvm", 20, 20, util.FmtDec},
			{"mxr", 19, 19, util.FmtDec},
			{"sum", 18, 18, util.FmtDec},
			{"mprv", 17, 17, util.FmtDec},
			{"xs", 16, 15, fmtXS},
			{"fs", 14, 13, fmtXS},
			{"mpp", 12, 11, util.FmtDec},
			{"spp", 8, 8, util.FmtDec},
			{"mpie", 7, 7, util.FmtDec},
			{"spie", 5, 5, util.FmtDec},
			{"upie", 4, 4, util.FmtDec},
			{"mie", 3, 3, util.FmtDec},
			{"sie", 1, 1, util.FmtDec},
			{"uie", 0, 0, util.FmtDec},
		}
	}
	return fs.Display(s.mstatus.rd(ModeM))
}

func (s *State) mstatusRdMPP() uint {
	return util.GetBits(s.mstatus.rd(ModeM), 12, 11)
}

func (s *State) mstatusWrMPP(x uint) {
	s.mstatus.wr(util.SetBits(s.mstatus.val, x, 12, 11), ModeM)
}

func (s *State) mstatusRdSPP() uint {
	return util.GetBits(s.mstatus.rd(ModeM), 8, 8)
}

func (s *State) mstatusWrSPP(x uint) {
	s.mstatus.wr(util.SetBits(s.mstatus.val, x, 8, 8), ModeM)
}

func (s *State) mstatusRdMPIE() uint {
	return util.GetBits(s.mstatus.rd(ModeM), 7, 7)
}

func (s *State) mstatusWrMPIE(x uint) {
	s.mstatus.wr(util.SetBits(s.mstatus.val, x, 7, 7), ModeM)
}

func (s *State) mstatusRdSPIE() uint {
	return util.GetBits(s.mstatus.rd(ModeM), 5, 5)
}

func (s *State) mstatusWrSPIE(x uint) {
	s.mstatus.wr(util.SetBits(s.mstatus.val, x, 5, 5), ModeM)
}

func (s *State) mstatusRdUPIE() uint {
	return util.GetBits(s.mstatus.rd(ModeM), 4, 4)
}

func (s *State) mstatusWrUPIE(x uint) {
	s.mstatus.wr(util.SetBits(s.mstatus.val, x, 4, 4), ModeM)
}

func (s *State) mstatusRdMIE() uint {
	return util.GetBits(s.mstatus.rd(ModeM), 3, 3)
}

func (s *State) mstatusWrMIE(x uint) {
	s.mstatus.wr(util.SetBits(s.mstatus.val, x, 3, 3), ModeM)
}

func (s *State) mstatusRdSIE() uint {
	return util.GetBits(s.mstatus.rd(ModeM), 1, 1)
}

func (s *State) mstatusWrSIE(x uint) {
	s.mstatus.wr(util.SetBits(s.mstatus.val, x, 1, 1), ModeM)
}

func (s *State) mstatusRdUIE() uint {
	return util.GetBits(s.mstatus.rd(ModeM), 0, 0)
}

func (s *State) mstatusWrUIE(x uint) {
	s.mstatus.wr(util.SetBits(s.mstatus.val, x, 0, 0), ModeM)
}

func (s *State) mstatusRdFS() uint {
	return util.GetBits(s.mstatus.rd(ModeM), 14, 13)
}

func (s *State) mstatusWrFS(x uint) {
	s.mstatus.wr(util.SetBits(s.mstatus.val, x, 14, 13), ModeM)
}

func (s *State) updateMSTATUS(mode Mode) {
	switch mode {
	case ModeU:
		s.mstatusWrUPIE(s.mstatusRdUIE())
		s.mstatusWrUIE(0)
	case ModeS:
		s.mstatusWrSPIE(s.mstatusRdSIE())
		s.mstatusWrSIE(0)
	case ModeM:
		s.mstatusWrMPIE(s.mstatusRdMIE())
		s.mstatusWrMIE(0)
	}
}

//-----------------------------------------------------------------------------
// s/m exception/interrupt delegation registers

func wrMEDELEG(s *State, x uint) {
	s.medeleg = x
}

func rdMEDELEG(s *State) uint {
	return s.medeleg
}

func wrMIDELEG(s *State, x uint) {
	s.mideleg = x
}

func rdMIDELEG(s *State) uint {
	return s.mideleg
}

func wrSEDELEG(s *State, x uint) {
	s.sedeleg = x
}

func rdSEDELEG(s *State) uint {
	return s.sedeleg
}

func wrSIDELEG(s *State, x uint) {
	s.sideleg = x
}

func rdSIDELEG(s *State) uint {
	return s.sideleg
}

//-----------------------------------------------------------------------------
// u/s/m interrupt enable

func rdUIE(s *State) uint {
	return s.mie // TODO mask
}

func rdSIE(s *State) uint {
	return s.mie // TODO mask
}

func rdMIE(s *State) uint {
	return s.mie
}

func wrUIE(s *State, x uint) {
	s.mie = x // TODO mask
}

func wrSIE(s *State, x uint) {
	s.mie = x // TODO mask
}

func wrMIE(s *State, x uint) {
	s.mie = x
}

//-----------------------------------------------------------------------------
// u/s/m interrupt pending

func rdUIP(s *State) uint {
	return s.mip // TODO mask
}

func rdSIP(s *State) uint {
	return s.mip // TODO mask
}

func rdMIP(s *State) uint {
	return s.mip
}

func wrUIP(s *State, x uint) {
	s.mip = x // TODO mask
}

func wrSIP(s *State, x uint) {
	s.mip = x // TODO mask
}

func wrMIP(s *State, x uint) {
	s.mip = x
}

//-----------------------------------------------------------------------------
// supervisor address translation and protection

// VM is the virtual memory mode.
type VM uint

// Virtual memory mode.
const (
	Bare VM = iota
	SV32
	SV39
	SV48
	SV57
	SV64
)

func (vm VM) String() string {
	return []string{"bare", "sv32", "sv39", "sv48", "sv57", "sv64"}[vm]
}

// GetVM returns the VM mode set in the SATP.
func (s *State) GetVM() VM {
	return s.vm
}

// GetPPN returns the physical page number set in the SATP.
func (s *State) GetPPN() uint {
	return s.ppn
}

func wrSATP(s *State, x uint) {
	s.satp = x
	// cache the vm and ppn
	if s.sxlen == 32 {
		// RV32
		s.vm = [2]VM{Bare, SV32}[util.GetBits(s.satp, 31, 31)]
		s.ppn = util.GetBits(s.satp, 21, 0)
	} else {
		// RV64
		s.vm = map[uint]VM{0: Bare, 8: SV39, 9: SV48, 10: SV57, 11: SV64}[util.GetBits(s.satp, 63, 60)]
		s.ppn = util.GetBits(s.satp, 43, 0)
	}
}

func rdSATP(s *State) uint {
	return s.satp
}

func fmtMode32(x uint) string {
	s := [2]string{"bare", "sv32"}[x]
	return fmt.Sprintf("%s(%d)", s, x)
}

func fmtMode64(x uint) string {
	m := map[uint]string{0: "bare", 8: "sv39", 9: "sv48", 10: "sv57", 11: "sv64"}
	return util.DisplayEnum(x, m, "reserved")
}

// DisplaySATP returns a display string for the SATP register.
func DisplaySATP(s *State) string {
	var fs util.FieldSet
	if s.sxlen == 32 {
		// RV32
		fs = util.FieldSet{
			{"mode", 31, 31, fmtMode32},
			{"asid", 30, 22, util.FmtHex},
			{"ppn", 21, 0, util.FmtHex},
		}
	} else {
		// RV64
		fs = util.FieldSet{
			{"mode", 63, 60, fmtMode64},
			{"asid", 59, 44, util.FmtHex},
			{"ppn", 43, 0, util.FmtHex},
		}
	}
	return fs.Display(s.satp)
}

//-----------------------------------------------------------------------------
// mcycle

func rdMCYCLE(s *State) uint {
	if s.mxlen == 32 {
		return uint(uint32(s.mcycle))
	}
	return uint(s.mcycle)
}

func rdMCYCLEH(s *State) uint {
	return uint(s.mcycle >> 32)
}

// IncClockCycles increments the CSR clock cycle counter.
func (s *State) IncClockCycles(n uint) {
	s.mcycle += uint64(n)
}

//-----------------------------------------------------------------------------
// minstret

func rdMINSTRET(s *State) uint {
	if s.mxlen == 32 {
		return uint(uint32(s.minstret))
	}
	return uint(s.minstret)
}

func rdMINSTRETH(s *State) uint {
	return uint(s.minstret >> 32)
}

// IncInstructions increments the CSR instructions retired counter.
func (s *State) IncInstructions() {
	s.minstret++
}

//-----------------------------------------------------------------------------

type wrFunc func(s *State, val uint)
type rdFunc func(s *State) uint
type displayFunc func(s *State) string

type csrDefn struct {
	name    string      // name of CSR
	wr      wrFunc      // write function for CSR
	rd      rdFunc      // read function for CSR
	display displayFunc // display function for CSR
}

var lookup = map[uint]csrDefn{
	// User CSRs 0x000 - 0x0ff (read/write)
	0x000: {"ustatus", nil, nil, nil},
	0x001: {"fflags", wrFFLAGS, rdFFLAGS, nil},
	0x002: {"frm", wrFRM, rdFRM, nil},
	0x003: {"fcsr", wrFCSR, rdFCSR, nil},
	0x004: {"uie", wrUIE, rdUIE, nil},
	0x005: {"utvec", wrUTVEC, rdUTVEC, nil},
	0x040: {"uscratch", wrUSCRATCH, rdUSCRATCH, nil},
	0x041: {"uepc", nil, rdUEPC, nil},
	0x042: {"ucause", nil, rdUCAUSE, nil},
	0x043: {"utval", wrUTVAL, rdUTVAL, nil},
	0x044: {"uip", wrUIP, rdUIP, nil},
	// User CSRs 0xc00 - 0xc7f (read only)
	0xc00: {"cycle", nil, rdMCYCLE, nil},
	0xc01: {"time", nil, nil, nil},
	0xc02: {"instret", nil, rdMINSTRET, nil},
	0xc03: {"hpmcounter3", nil, nil, nil},
	0xc04: {"hpmcounter4", nil, nil, nil},
	0xc05: {"hpmcounter5", nil, nil, nil},
	0xc06: {"hpmcounter6", nil, nil, nil},
	0xc07: {"hpmcounter7", nil, nil, nil},
	0xc08: {"hpmcounter8", nil, nil, nil},
	0xc09: {"hpmcounter9", nil, nil, nil},
	0xc0a: {"hpmcounter10", nil, nil, nil},
	0xc0b: {"hpmcounter11", nil, nil, nil},
	0xc0c: {"hpmcounter12", nil, nil, nil},
	0xc0d: {"hpmcounter13", nil, nil, nil},
	0xc0e: {"hpmcounter14", nil, nil, nil},
	0xc0f: {"hpmcounter15", nil, nil, nil},
	0xc10: {"hpmcounter16", nil, nil, nil},
	0xc11: {"hpmcounter17", nil, nil, nil},
	0xc12: {"hpmcounter18", nil, nil, nil},
	0xc13: {"hpmcounter19", nil, nil, nil},
	0xc14: {"hpmcounter20", nil, nil, nil},
	0xc15: {"hpmcounter21", nil, nil, nil},
	0xc16: {"hpmcounter22", nil, nil, nil},
	0xc17: {"hpmcounter23", nil, nil, nil},
	0xc18: {"hpmcounter24", nil, nil, nil},
	0xc19: {"hpmcounter25", nil, nil, nil},
	0xc1a: {"hpmcounter26", nil, nil, nil},
	0xc1b: {"hpmcounter27", nil, nil, nil},
	0xc1c: {"hpmcounter28", nil, nil, nil},
	0xc1d: {"hpmcounter29", nil, nil, nil},
	0xc1e: {"hpmcounter30", nil, nil, nil},
	0xc1f: {"hpmcounter31", nil, nil, nil},
	// User CSRs 0xc80 - 0xcbf (read only)
	0xc80: {"cycleh", nil, rdMCYCLEH, nil},
	0xc81: {"timeh", nil, nil, nil},
	0xc82: {"instreth", nil, rdMINSTRETH, nil},
	0xc83: {"hpmcounter3h", nil, nil, nil},
	0xc84: {"hpmcounter4h", nil, nil, nil},
	0xc85: {"hpmcounter5h", nil, nil, nil},
	0xc86: {"hpmcounter6h", nil, nil, nil},
	0xc87: {"hpmcounter7h", nil, nil, nil},
	0xc88: {"hpmcounter8h", nil, nil, nil},
	0xc89: {"hpmcounter9h", nil, nil, nil},
	0xc8a: {"hpmcounter10h", nil, nil, nil},
	0xc8b: {"hpmcounter11h", nil, nil, nil},
	0xc8c: {"hpmcounter12h", nil, nil, nil},
	0xc8d: {"hpmcounter13h", nil, nil, nil},
	0xc8e: {"hpmcounter14h", nil, nil, nil},
	0xc8f: {"hpmcounter15h", nil, nil, nil},
	0xc90: {"hpmcounter16h", nil, nil, nil},
	0xc91: {"hpmcounter17h", nil, nil, nil},
	0xc92: {"hpmcounter18h", nil, nil, nil},
	0xc93: {"hpmcounter19h", nil, nil, nil},
	0xc94: {"hpmcounter20h", nil, nil, nil},
	0xc95: {"hpmcounter21h", nil, nil, nil},
	0xc96: {"hpmcounter22h", nil, nil, nil},
	0xc97: {"hpmcounter23h", nil, nil, nil},
	0xc98: {"hpmcounter24h", nil, nil, nil},
	0xc99: {"hpmcounter25h", nil, nil, nil},
	0xc9a: {"hpmcounter26h", nil, nil, nil},
	0xc9b: {"hpmcounter27h", nil, nil, nil},
	0xc9c: {"hpmcounter28h", nil, nil, nil},
	0xc9d: {"hpmcounter29h", nil, nil, nil},
	0xc9e: {"hpmcounter30h", nil, nil, nil},
	0xc9f: {"hpmcounter31h", nil, nil, nil},
	// Supervisor CSRs 0x100 - 0x1ff (read/write)
	0x100: {"sstatus", wrSSTATUS, rdSSTATUS, displaySSTATUS},
	0x102: {"sedeleg", wrSEDELEG, rdSEDELEG, nil},
	0x103: {"sideleg", wrSIDELEG, rdSIDELEG, nil},
	0x104: {"sie", wrSIE, rdSIE, nil},
	0x105: {"stvec", wrSTVEC, rdSTVEC, displaySTVEC},
	0x106: {"scounteren", nil, nil, nil},
	0x140: {"sscratch", wrSSCRATCH, rdSSCRATCH, nil},
	0x141: {"sepc", wrSEPC, rdSEPC, nil},
	0x142: {"scause", nil, rdSCAUSE, nil},
	0x143: {"stval", wrSTVAL, rdSTVAL, nil},
	0x144: {"sip", wrSIP, rdSIP, nil},
	0x180: {"satp", wrSATP, rdSATP, DisplaySATP},
	// Machine CSRs 0xf00 - 0xf7f (read only)
	0xf11: {"mvendorid", nil, rdZero, nil},
	0xf12: {"marchid", nil, rdZero, nil},
	0xf13: {"mimpid", nil, rdZero, nil},
	0xf14: {"mhartid", nil, rdZero, nil},
	// Machine CSRs 0x300 - 0x3ff (read/write)
	0x300: {"mstatus", wrMSTATUS, rdMSTATUS, displayMSTATUS},
	0x301: {"misa", wrMISA, rdMISA, displayMISA},
	0x302: {"medeleg", wrMEDELEG, rdMEDELEG, nil},
	0x303: {"mideleg", wrMIDELEG, rdMIDELEG, nil},
	0x304: {"mie", wrMIE, rdMIE, nil},
	0x305: {"mtvec", wrMTVEC, rdMTVEC, displayMTVEC},
	0x306: {"mcounteren", nil, nil, nil},
	0x320: {"mucounteren", nil, nil, nil},
	0x321: {"mscounteren", nil, nil, nil},
	0x322: {"mhcounteren", nil, nil, nil},
	0x323: {"mhpmevent3", nil, nil, nil},
	0x324: {"mhpmevent4", nil, nil, nil},
	0x325: {"mhpmevent5", nil, nil, nil},
	0x326: {"mhpmevent6", nil, nil, nil},
	0x327: {"mhpmevent7", nil, nil, nil},
	0x328: {"mhpmevent8", nil, nil, nil},
	0x329: {"mhpmevent9", nil, nil, nil},
	0x32a: {"mhpmevent10", nil, nil, nil},
	0x32b: {"mhpmevent11", nil, nil, nil},
	0x32c: {"mhpmevent12", nil, nil, nil},
	0x32d: {"mhpmevent13", nil, nil, nil},
	0x32e: {"mhpmevent14", nil, nil, nil},
	0x32f: {"mhpmevent15", nil, nil, nil},
	0x330: {"mhpmevent16", nil, nil, nil},
	0x331: {"mhpmevent17", nil, nil, nil},
	0x332: {"mhpmevent18", nil, nil, nil},
	0x333: {"mhpmevent19", nil, nil, nil},
	0x334: {"mhpmevent20", nil, nil, nil},
	0x335: {"mhpmevent21", nil, nil, nil},
	0x336: {"mhpmevent22", nil, nil, nil},
	0x337: {"mhpmevent23", nil, nil, nil},
	0x338: {"mhpmevent24", nil, nil, nil},
	0x339: {"mhpmevent25", nil, nil, nil},
	0x33a: {"mhpmevent26", nil, nil, nil},
	0x33b: {"mhpmevent27", nil, nil, nil},
	0x33c: {"mhpmevent28", nil, nil, nil},
	0x33d: {"mhpmevent29", nil, nil, nil},
	0x33e: {"mhpmevent30", nil, nil, nil},
	0x33f: {"mhpmevent31", nil, nil, nil},
	0x340: {"mscratch", wrMSCRATCH, rdMSCRATCH, nil},
	0x341: {"mepc", wrMEPC, rdMEPC, nil},
	0x342: {"mcause", nil, rdMCAUSE, nil},
	0x343: {"mtval", wrMTVAL, rdMTVAL, nil},
	0x344: {"mip", wrMIP, rdMIP, nil},
	0x380: {"mbase", nil, nil, nil},
	0x381: {"mbound", nil, nil, nil},
	0x382: {"mibase", nil, nil, nil},
	0x383: {"mibound", nil, nil, nil},
	0x384: {"mdbase", nil, nil, nil},
	0x385: {"mdbound", nil, nil, nil},
	0x3a0: {"pmpcfg0", wrIgnore, nil, nil},
	0x3a1: {"pmpcfg1", wrIgnore, nil, nil},
	0x3a2: {"pmpcfg2", wrIgnore, nil, nil},
	0x3a3: {"pmpcfg3", wrIgnore, nil, nil},
	0x3b0: {"pmpaddr0", wrIgnore, nil, nil},
	0x3b1: {"pmpaddr1", wrIgnore, nil, nil},
	0x3b2: {"pmpaddr2", wrIgnore, nil, nil},
	0x3b3: {"pmpaddr3", wrIgnore, nil, nil},
	0x3b4: {"pmpaddr4", wrIgnore, nil, nil},
	0x3b5: {"pmpaddr5", wrIgnore, nil, nil},
	0x3b6: {"pmpaddr6", wrIgnore, nil, nil},
	0x3b7: {"pmpaddr7", wrIgnore, nil, nil},
	0x3b8: {"pmpaddr8", wrIgnore, nil, nil},
	0x3b9: {"pmpaddr9", wrIgnore, nil, nil},
	0x3ba: {"pmpaddr10", wrIgnore, nil, nil},
	0x3bb: {"pmpaddr11", wrIgnore, nil, nil},
	0x3bc: {"pmpaddr12", wrIgnore, nil, nil},
	0x3bd: {"pmpaddr13", wrIgnore, nil, nil},
	0x3be: {"pmpaddr14", wrIgnore, nil, nil},
	0x3bf: {"pmpaddr15", wrIgnore, nil, nil},
	// Machine CSRs 0xb00 - 0xb7f (read/write)
	0xb00: {"mcycle", nil, rdMCYCLE, nil},
	0xb02: {"minstret", nil, rdMINSTRET, nil},
	0xb03: {"mhpmcounter3", nil, nil, nil},
	0xb04: {"mhpmcounter4", nil, nil, nil},
	0xb05: {"mhpmcounter5", nil, nil, nil},
	0xb06: {"mhpmcounter6", nil, nil, nil},
	0xb07: {"mhpmcounter7", nil, nil, nil},
	0xb08: {"mhpmcounter8", nil, nil, nil},
	0xb09: {"mhpmcounter9", nil, nil, nil},
	0xb0a: {"mhpmcounter10", nil, nil, nil},
	0xb0b: {"mhpmcounter11", nil, nil, nil},
	0xb0c: {"mhpmcounter12", nil, nil, nil},
	0xb0d: {"mhpmcounter13", nil, nil, nil},
	0xb0e: {"mhpmcounter14", nil, nil, nil},
	0xb0f: {"mhpmcounter15", nil, nil, nil},
	0xb10: {"mhpmcounter16", nil, nil, nil},
	0xb11: {"mhpmcounter17", nil, nil, nil},
	0xb12: {"mhpmcounter18", nil, nil, nil},
	0xb13: {"mhpmcounter19", nil, nil, nil},
	0xb14: {"mhpmcounter20", nil, nil, nil},
	0xb15: {"mhpmcounter21", nil, nil, nil},
	0xb16: {"mhpmcounter22", nil, nil, nil},
	0xb17: {"mhpmcounter23", nil, nil, nil},
	0xb18: {"mhpmcounter24", nil, nil, nil},
	0xb19: {"mhpmcounter25", nil, nil, nil},
	0xb1a: {"mhpmcounter26", nil, nil, nil},
	0xb1b: {"mhpmcounter27", nil, nil, nil},
	0xb1c: {"mhpmcounter28", nil, nil, nil},
	0xb1d: {"mhpmcounter29", nil, nil, nil},
	0xb1e: {"mhpmcounter30", nil, nil, nil},
	0xb1f: {"mhpmcounter31", nil, nil, nil},
	// Machine CSRs 0xb80 - 0xbbf (read/write)
	0xb80: {"mcycleh", nil, rdMCYCLEH, nil},
	0xb82: {"minstreth", nil, rdMINSTRETH, nil},
	0xb83: {"mhpmcounter3h", nil, nil, nil},
	0xb84: {"mhpmcounter4h", nil, nil, nil},
	0xb85: {"mhpmcounter5h", nil, nil, nil},
	0xb86: {"mhpmcounter6h", nil, nil, nil},
	0xb87: {"mhpmcounter7h", nil, nil, nil},
	0xb88: {"mhpmcounter8h", nil, nil, nil},
	0xb89: {"mhpmcounter9h", nil, nil, nil},
	0xb8a: {"mhpmcounter10h", nil, nil, nil},
	0xb8b: {"mhpmcounter11h", nil, nil, nil},
	0xb8c: {"mhpmcounter12h", nil, nil, nil},
	0xb8d: {"mhpmcounter13h", nil, nil, nil},
	0xb8e: {"mhpmcounter14h", nil, nil, nil},
	0xb8f: {"mhpmcounter15h", nil, nil, nil},
	0xb90: {"mhpmcounter16h", nil, nil, nil},
	0xb91: {"mhpmcounter17h", nil, nil, nil},
	0xb92: {"mhpmcounter18h", nil, nil, nil},
	0xb93: {"mhpmcounter19h", nil, nil, nil},
	0xb94: {"mhpmcounter20h", nil, nil, nil},
	0xb95: {"mhpmcounter21h", nil, nil, nil},
	0xb96: {"mhpmcounter22h", nil, nil, nil},
	0xb97: {"mhpmcounter23h", nil, nil, nil},
	0xb98: {"mhpmcounter24h", nil, nil, nil},
	0xb99: {"mhpmcounter25h", nil, nil, nil},
	0xb9a: {"mhpmcounter26h", nil, nil, nil},
	0xb9b: {"mhpmcounter27h", nil, nil, nil},
	0xb9c: {"mhpmcounter28h", nil, nil, nil},
	0xb9d: {"mhpmcounter29h", nil, nil, nil},
	0xb9e: {"mhpmcounter30h", nil, nil, nil},
	0xb9f: {"mhpmcounter31h", nil, nil, nil},
	// Machine Debug CSRs 0x7a0 - 0x7af (read/write)
	0x7a0: {"tselect", wrIgnore, rdZero, nil},
	0x7a1: {"tdata1", wrIgnore, rdZero, nil},
	0x7a2: {"tdata2", wrIgnore, rdZero, nil},
	0x7a3: {"tdata3", wrIgnore, rdZero, nil},
	// Machine Debug Mode Only CSRs 0x7b0 - 0x7bf (read/write)
	0x7b0: {"dcsr", nil, nil, nil},
	0x7b1: {"dpc", nil, nil, nil},
	0x7b2: {"dscratch", nil, nil, nil},
	// Hypervisor CSRs 0x200 - 0x2ff (read/write)
	0x200: {"hstatus", nil, nil, nil},
	0x202: {"hedeleg", nil, nil, nil},
	0x203: {"hideleg", nil, nil, nil},
	0x204: {"hie", nil, nil, nil},
	0x205: {"htvec", nil, nil, nil},
	0x240: {"hscratch", nil, nil, nil},
	0x241: {"hepc", nil, nil, nil},
	0x242: {"hcause", nil, nil, nil},
	0x243: {"hbadaddr", nil, nil, nil},
	0x244: {"hip", nil, nil, nil},
}

//-----------------------------------------------------------------------------

// canAccess returns true if the register can be accessed in the current mode.
func (s *State) canAccess(reg uint) bool {
	mode := Mode((reg >> 8) & 3)
	return s.GetMode() >= mode
}

// canWr returns true if the register can be written.
func canWr(reg uint) bool {
	rw := (reg >> 10) & 3
	return rw != 3
}

//-----------------------------------------------------------------------------

// State stores the CSR state for the CPU.
type State struct {
	mode   Mode // current privilege mode
	xlen   uint // cpu register length 32/64/128
	mxlen  uint // machine register length
	uxlen  uint // user register length
	sxlen  uint // supervisor register length
	ialign uint // instruction alignment 16/32
	vm     VM   // cached virtual memory mode from SATP
	ppn    uint // cached physical page number from SATP
	// combined u/s/m CSRs
	mstatus mStatus // u/s/m status
	mie     uint    // u/s/m interrupt enable register
	mip     uint    // u/s/m interrupt pending register
	// machine CSRs
	mcause   uint   // machine cause register
	mepc     uint   // machine exception program counter
	mscratch uint   // machine scratch register
	mtvec    uint   // machine trap vector base address register
	mtval    uint   // machine trap value register
	misa     uint   // machine isa register
	medeleg  uint   // machine exception delegation register
	mideleg  uint   // machine interrupt delegation register
	mcycle   uint64 // machine clock cycles
	minstret uint64 // number of retired instructions
	// Supervisor CSRs
	scause   uint // supervisor cause register
	sepc     uint // supervisor exception program counter
	sscratch uint // supervisor scratch register
	stval    uint // supervisor trap value register
	stvec    uint // supervisor trap vector base address register
	sedeleg  uint // supervisor exception delegation register
	sideleg  uint // supervisor interrupt delegation register
	satp     uint // supervisor address translation and protection
	// User CSRs
	ucause   uint // user cause register
	uepc     uint // user exception program counter
	uscratch uint // user scratch register
	utval    uint // user trap value register
	utvec    uint // user trap vector base address register
	fcsr     uint // floating point control and status register
}

// NewState returns a CSR state object.
func NewState(xlen, ext uint) *State {
	s := &State{
		mode:  ModeM, // start in machine mode
		xlen:  xlen,
		mxlen: xlen,
		uxlen: xlen,
		sxlen: xlen,
	}
	initMISA(s, ext)
	s.mstatus.init(s.mxlen)
	return s
}

// Reset resets the state of the CSR sub-system.
func (s *State) Reset() {
	s.setMode(ModeM)
	wrSATP(s, 0)
	// etc..
}

// Rd reads from a CSR.
func (s *State) Rd(reg uint) (uint64, error) {
	if !s.canAccess(reg) {
		return 0, &Error{reg, ErrPrivilege}
	}
	if x, ok := lookup[reg]; ok {
		if x.rd == nil {
			return 0, &Error{reg, ErrNoRead}
		}
		return uint64(x.rd(s)), nil
	}
	return 0, &Error{reg, ErrTodo}
}

// Wr writes to a CSR.
func (s *State) Wr(reg uint, val uint64) error {
	if !canWr(reg) {
		return &Error{reg, ErrReadOnly}
	}
	if !s.canAccess(reg) {
		return &Error{reg, ErrPrivilege}
	}
	if x, ok := lookup[reg]; ok {
		if x.wr == nil {
			return &Error{reg, ErrNoWrite}
		}
		x.wr(s, uint(val))
		return nil
	}
	return &Error{reg, ErrTodo}
}

//-----------------------------------------------------------------------------

// Name returns the name of a given CSR.
func Name(reg uint) string {
	if x, ok := lookup[reg]; ok {
		return x.name
	}
	return fmt.Sprintf("0x%03x", reg)
}

// getMode returns the mode bits from a register address.
func getMode(reg uint) uint {
	return (reg >> 8) & 3
}

type regStrings struct {
	num    string // 3-nybble register number
	name   string // register name
	access string // current access mode
	val    string // raw register value
	field  string // bit field decodes within value
}

func (s *State) regDisplay(reg uint) (*regStrings, error) {
	r, ok := lookup[reg]
	if !ok || r.rd == nil {
		return nil, errors.New("no read")
	}
	// access string
	mode := [4]string{"u", "s", "h", "m"}[getMode(reg)]
	var rw string
	if s.canAccess(reg) {
		if canWr(reg) {
			rw = "rw"
		} else {
			rw = "r_"
		}
	} else {
		rw = ".."
	}
	accessStr := mode + rw
	// value string
	valStr := "0"
	val := r.rd(s)
	if val != 0 {
		rlen := []uint{s.uxlen, s.sxlen, 64, s.mxlen}[getMode(reg)]
		fmtStr := fmt.Sprintf("%%0%dx", rlen>>2)
		valStr = fmt.Sprintf(fmtStr, r.rd(s))
	}
	// field string
	fieldStr := ""
	if r.display != nil {
		fieldStr = r.display(s)
	}
	return &regStrings{
		num:    fmt.Sprintf("%03x", reg),
		name:   r.name,
		access: accessStr,
		val:    valStr,
		field:  fieldStr,
	}, nil
}

// Display displays the CSR state.
func (s *State) Display() string {
	x := [][]string{}
	x = append(x, []string{"mode", fmt.Sprintf("%s", s.GetMode()), ""})
	// read all registers
	for reg := uint(0); reg < 4096; reg++ {
		d, err := s.regDisplay(reg)
		if err != nil {
			continue
		}
		regStr := fmt.Sprintf("%s %s %s", d.num, d.access, d.name)
		x = append(x, []string{regStr, d.val, d.field})
	}
	// return the table string
	return cli.TableString(x, []int{0, 0, 0}, 1)
}

//-----------------------------------------------------------------------------

// MRET returns from a machine-mode exception.
func (s *State) MRET() uint {
	// restore previous MIE
	s.mstatusWrMIE(s.mstatusRdMPIE())
	// switch to target mode
	s.setMode(Mode(s.mstatusRdMPP()))
	// MPIE=1
	s.mstatusWrMPIE(1)
	// MPP=U (or M if no user mode)
	s.mstatusWrMPP(uint(s.getMinMode()))
	// jump to exception address
	return rdMEPC(s)
}

// SRET returns from a supervisor-mode exception.
func (s *State) SRET() uint {
	// restore previous SIE
	s.mstatusWrSIE(s.mstatusRdSPIE())
	// switch to target mode
	s.setMode(Mode(s.mstatusRdSPP()))
	// SPIE=1
	s.mstatusWrSPIE(1)
	// SPP=U (or M if no user mode)
	s.mstatusWrSPP(uint(s.getMinMode()))
	// jump to exception address
	return rdSEPC(s)
}

// URET returns from a user-mode exception.
func (s *State) URET() uint {
	// restore previous UIE
	s.mstatusWrUIE(s.mstatusRdUPIE())
	// switch to target mode
	s.setMode(ModeU)
	// UPIE=1
	s.mstatusWrUPIE(1)
	// jump to exception address
	return rdUEPC(s)
}

//-----------------------------------------------------------------------------

// ECALL performs an environment call exception.
func (s *State) ECALL(epc uint64, val uint) uint64 {
	switch s.GetMode() {
	case ModeU:
		return s.Exception(epc, uint(ExEnvCallFromUserMode), val, false)
	case ModeS:
		return s.Exception(epc, uint(ExEnvCallFromSupervisorMode), val, false)
	case ModeM:
		return s.Exception(epc, uint(ExEnvCallFromMachineMode), val, false)
	}
	return 0
}

// Exception performs a cpu exception.
func (s *State) Exception(epc uint64, code, val uint, isInterrupt bool) uint64 {
	// what's the next cpu mode?
	nextMode := s.getNextMode(code, isInterrupt)
	// update interrupt enable and interrupt enable stack
	s.updateMSTATUS(nextMode)
	// set the cause register
	s.setCause(code, isInterrupt, nextMode)
	// set the exception program counter
	s.setEPC(epc, nextMode)
	// set the trap value
	s.setTrapValue(val, nextMode)
	// get exception base address and mode
	base, mode := s.getTrapVector(nextMode)
	// handle direct or vectored exception
	var pc uint64
	if (mode == 0) || !isInterrupt {
		pc = uint64(base)
	} else {
		pc = uint64(base + (4 * code))
	}
	// update the mode
	switch nextMode {
	case ModeS:
		s.mstatusWrSPP(uint(s.GetMode()))
	case ModeM:
		s.mstatusWrMPP(uint(s.GetMode()))
	}
	s.setMode(nextMode)
	return pc
}

//-----------------------------------------------------------------------------
