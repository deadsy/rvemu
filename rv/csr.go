//-----------------------------------------------------------------------------
/*

RISC-V Control and Status Register Definitions

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
)

//-----------------------------------------------------------------------------
// Privilege Levels

const privU = 0 // user
const privS = 1 // supervisor
const privM = 3 // machine

//-----------------------------------------------------------------------------
// Known CSRs

// User CSRs 0x000 - 0x0ff (read/write)
const csrUSTATUS = 0x0
const csrFFLAGS = 0x1
const csrFRM = 0x2
const csrFCSR = 0x3
const csrUIE = 0x4
const csrUTVEC = 0x5
const csrUSCRATCH = 0x40
const csrUEPC = 0x41
const csrUCAUSE = 0x42
const csrUTVAL = 0x43
const csrUIP = 0x44

// User CSRs 0xc00 - 0xc7f (read only)
const csrCYCLE = 0xc00
const csrTIME = 0xc01
const csrINSTRET = 0xc02
const csrHPMCOUNTER3 = 0xc03
const csrHPMCOUNTER4 = 0xc04
const csrHPMCOUNTER5 = 0xc05
const csrHPMCOUNTER6 = 0xc06
const csrHPMCOUNTER7 = 0xc07
const csrHPMCOUNTER8 = 0xc08
const csrHPMCOUNTER9 = 0xc09
const csrHPMCOUNTER10 = 0xc0a
const csrHPMCOUNTER11 = 0xc0b
const csrHPMCOUNTER12 = 0xc0c
const csrHPMCOUNTER13 = 0xc0d
const csrHPMCOUNTER14 = 0xc0e
const csrHPMCOUNTER15 = 0xc0f
const csrHPMCOUNTER16 = 0xc10
const csrHPMCOUNTER17 = 0xc11
const csrHPMCOUNTER18 = 0xc12
const csrHPMCOUNTER19 = 0xc13
const csrHPMCOUNTER20 = 0xc14
const csrHPMCOUNTER21 = 0xc15
const csrHPMCOUNTER22 = 0xc16
const csrHPMCOUNTER23 = 0xc17
const csrHPMCOUNTER24 = 0xc18
const csrHPMCOUNTER25 = 0xc19
const csrHPMCOUNTER26 = 0xc1a
const csrHPMCOUNTER27 = 0xc1b
const csrHPMCOUNTER28 = 0xc1c
const csrHPMCOUNTER29 = 0xc1d
const csrHPMCOUNTER30 = 0xc1e
const csrHPMCOUNTER31 = 0xc1f

// User CSRs 0xc80 - 0xcbf (read only)
const csrCYCLEH = 0xc80
const csrTIMEH = 0xc81
const csrINSTRETH = 0xc82
const csrHPMCOUNTER3H = 0xc83
const csrHPMCOUNTER4H = 0xc84
const csrHPMCOUNTER5H = 0xc85
const csrHPMCOUNTER6H = 0xc86
const csrHPMCOUNTER7H = 0xc87
const csrHPMCOUNTER8H = 0xc88
const csrHPMCOUNTER9H = 0xc89
const csrHPMCOUNTER10H = 0xc8a
const csrHPMCOUNTER11H = 0xc8b
const csrHPMCOUNTER12H = 0xc8c
const csrHPMCOUNTER13H = 0xc8d
const csrHPMCOUNTER14H = 0xc8e
const csrHPMCOUNTER15H = 0xc8f
const csrHPMCOUNTER16H = 0xc90
const csrHPMCOUNTER17H = 0xc91
const csrHPMCOUNTER18H = 0xc92
const csrHPMCOUNTER19H = 0xc93
const csrHPMCOUNTER20H = 0xc94
const csrHPMCOUNTER21H = 0xc95
const csrHPMCOUNTER22H = 0xc96
const csrHPMCOUNTER23H = 0xc97
const csrHPMCOUNTER24H = 0xc98
const csrHPMCOUNTER25H = 0xc99
const csrHPMCOUNTER26H = 0xc9a
const csrHPMCOUNTER27H = 0xc9b
const csrHPMCOUNTER28H = 0xc9c
const csrHPMCOUNTER29H = 0xc9d
const csrHPMCOUNTER30H = 0xc9e
const csrHPMCOUNTER31H = 0xc9f

// Supervisor CSRs 0x100 - 0x1ff (read/write)
const csrSSTATUS = 0x100
const csrSEDELEG = 0x102
const csrSIDELEG = 0x103
const csrSIE = 0x104
const csrSTVEC = 0x105
const csrSCOUNTEREN = 0x106
const csrSSCRATCH = 0x140
const csrSEPC = 0x141
const csrSCAUSE = 0x142
const csrSTVAL = 0x143
const csrSIP = 0x144
const csrSATP = 0x180

// Machine CSRs 0xf00 - 0xf7f (read only)
const csrMVENDORID = 0xf11
const csrMARCHID = 0xf12
const csrMIMPID = 0xf13
const csrMHARTID = 0xf14

// Machine CSRs 0x300 - 0x3ff (read/write)
const csrMSTATUS = 0x300
const csrMISA = 0x301
const csrMEDELEG = 0x302
const csrMIDELEG = 0x303
const csrMIE = 0x304
const csrMTVEC = 0x305
const csrMCOUNTEREN = 0x306
const csrMUCOUNTEREN = 0x320
const csrMSCOUNTEREN = 0x321
const csrMHCOUNTEREN = 0x322
const csrMHPMEVENT3 = 0x323
const csrMHPMEVENT4 = 0x324
const csrMHPMEVENT5 = 0x325
const csrMHPMEVENT6 = 0x326
const csrMHPMEVENT7 = 0x327
const csrMHPMEVENT8 = 0x328
const csrMHPMEVENT9 = 0x329
const csrMHPMEVENT10 = 0x32a
const csrMHPMEVENT11 = 0x32b
const csrMHPMEVENT12 = 0x32c
const csrMHPMEVENT13 = 0x32d
const csrMHPMEVENT14 = 0x32e
const csrMHPMEVENT15 = 0x32f
const csrMHPMEVENT16 = 0x330
const csrMHPMEVENT17 = 0x331
const csrMHPMEVENT18 = 0x332
const csrMHPMEVENT19 = 0x333
const csrMHPMEVENT20 = 0x334
const csrMHPMEVENT21 = 0x335
const csrMHPMEVENT22 = 0x336
const csrMHPMEVENT23 = 0x337
const csrMHPMEVENT24 = 0x338
const csrMHPMEVENT25 = 0x339
const csrMHPMEVENT26 = 0x33a
const csrMHPMEVENT27 = 0x33b
const csrMHPMEVENT28 = 0x33c
const csrMHPMEVENT29 = 0x33d
const csrMHPMEVENT30 = 0x33e
const csrMHPMEVENT31 = 0x33f
const csrMSCRATCH = 0x340
const csrMEPC = 0x341
const csrMCAUSE = 0x342
const csrMTVAL = 0x343
const csrMIP = 0x344
const csrMBASE = 0x380
const csrMBOUND = 0x381
const csrMIBASE = 0x382
const csrMIBOUND = 0x383
const csrMDBASE = 0x384
const csrMDBOUND = 0x385
const csrPMPCFG0 = 0x3a0
const csrPMPCFG1 = 0x3a1
const csrPMPCFG2 = 0x3a2
const csrPMPCFG3 = 0x3a3
const csrPMPADDR0 = 0x3b0
const csrPMPADDR1 = 0x3b1
const csrPMPADDR2 = 0x3b2
const csrPMPADDR3 = 0x3b3
const csrPMPADDR4 = 0x3b4
const csrPMPADDR5 = 0x3b5
const csrPMPADDR6 = 0x3b6
const csrPMPADDR7 = 0x3b7
const csrPMPADDR8 = 0x3b8
const csrPMPADDR9 = 0x3b9
const csrPMPADDR10 = 0x3ba
const csrPMPADDR11 = 0x3bb
const csrPMPADDR12 = 0x3bc
const csrPMPADDR13 = 0x3bd
const csrPMPADDR14 = 0x3be
const csrPMPADDR15 = 0x3bf

// Machine CSRs 0xb00 - 0xb7f (read/write)
const csrMCYCLE = 0xb00
const csrMINSTRET = 0xb02
const csrMHPMCOUNTER3 = 0xb03
const csrMHPMCOUNTER4 = 0xb04
const csrMHPMCOUNTER5 = 0xb05
const csrMHPMCOUNTER6 = 0xb06
const csrMHPMCOUNTER7 = 0xb07
const csrMHPMCOUNTER8 = 0xb08
const csrMHPMCOUNTER9 = 0xb09
const csrMHPMCOUNTER10 = 0xb0a
const csrMHPMCOUNTER11 = 0xb0b
const csrMHPMCOUNTER12 = 0xb0c
const csrMHPMCOUNTER13 = 0xb0d
const csrMHPMCOUNTER14 = 0xb0e
const csrMHPMCOUNTER15 = 0xb0f
const csrMHPMCOUNTER16 = 0xb10
const csrMHPMCOUNTER17 = 0xb11
const csrMHPMCOUNTER18 = 0xb12
const csrMHPMCOUNTER19 = 0xb13
const csrMHPMCOUNTER20 = 0xb14
const csrMHPMCOUNTER21 = 0xb15
const csrMHPMCOUNTER22 = 0xb16
const csrMHPMCOUNTER23 = 0xb17
const csrMHPMCOUNTER24 = 0xb18
const csrMHPMCOUNTER25 = 0xb19
const csrMHPMCOUNTER26 = 0xb1a
const csrMHPMCOUNTER27 = 0xb1b
const csrMHPMCOUNTER28 = 0xb1c
const csrMHPMCOUNTER29 = 0xb1d
const csrMHPMCOUNTER30 = 0xb1e
const csrMHPMCOUNTER31 = 0xb1f

// Machine CSRs 0xb80 - 0xbbf (read/write)
const csrMCYCLEH = 0xb80
const csrMINSTRETH = 0xb82
const csrMHPMCOUNTER3H = 0xb83
const csrMHPMCOUNTER4H = 0xb84
const csrMHPMCOUNTER5H = 0xb85
const csrMHPMCOUNTER6H = 0xb86
const csrMHPMCOUNTER7H = 0xb87
const csrMHPMCOUNTER8H = 0xb88
const csrMHPMCOUNTER9H = 0xb89
const csrMHPMCOUNTER10H = 0xb8a
const csrMHPMCOUNTER11H = 0xb8b
const csrMHPMCOUNTER12H = 0xb8c
const csrMHPMCOUNTER13H = 0xb8d
const csrMHPMCOUNTER14H = 0xb8e
const csrMHPMCOUNTER15H = 0xb8f
const csrMHPMCOUNTER16H = 0xb90
const csrMHPMCOUNTER17H = 0xb91
const csrMHPMCOUNTER18H = 0xb92
const csrMHPMCOUNTER19H = 0xb93
const csrMHPMCOUNTER20H = 0xb94
const csrMHPMCOUNTER21H = 0xb95
const csrMHPMCOUNTER22H = 0xb96
const csrMHPMCOUNTER23H = 0xb97
const csrMHPMCOUNTER24H = 0xb98
const csrMHPMCOUNTER25H = 0xb99
const csrMHPMCOUNTER26H = 0xb9a
const csrMHPMCOUNTER27H = 0xb9b
const csrMHPMCOUNTER28H = 0xb9c
const csrMHPMCOUNTER29H = 0xb9d
const csrMHPMCOUNTER30H = 0xb9e
const csrMHPMCOUNTER31H = 0xb9f

// Machine Debug CSRs 0x7a0 - 0x7af (read/write)
const csrTSELECT = 0x7a0
const csrTDATA1 = 0x7a1
const csrTDATA2 = 0x7a2
const csrTDATA3 = 0x7a3

// Machine Debug Mode Only CSRs 0x7b0 - 0x7bf (read/write)
const csrDCSR = 0x7b0
const csrDPC = 0x7b1
const csrDSCRATCH = 0x7b2

// Hypervisor CSRs 0x200 - 0x2ff (read/write)
const csrHSTATUS = 0x200
const csrHEDELEG = 0x202
const csrHIDELEG = 0x203
const csrHIE = 0x204
const csrHTVEC = 0x205
const csrHSCRATCH = 0x240
const csrHEPC = 0x241
const csrHCAUSE = 0x242
const csrHBADADDR = 0x243
const csrHIP = 0x244

//-----------------------------------------------------------------------------

type csrWr func(csr, val uint)
type csrRd func(csr uint) uint

type csrDefn struct {
	name string // name of CSR
	wr   csrWr  // write function for CSR
	rd   csrRd  // read function for CSR
}

var csrLUT = map[uint]csrDefn{
	csrMTVEC:     {"mtvec", nil, nil},
	csrPMPADDR0:  {"pmpaddr0", nil, nil},
	csrPMPADDR1:  {"pmpaddr1", nil, nil},
	csrPMPADDR2:  {"pmpaddr2", nil, nil},
	csrPMPADDR3:  {"pmpaddr3", nil, nil},
	csrPMPADDR4:  {"pmpaddr4", nil, nil},
	csrPMPADDR5:  {"pmpaddr5", nil, nil},
	csrPMPADDR6:  {"pmpaddr6", nil, nil},
	csrPMPADDR7:  {"pmpaddr7", nil, nil},
	csrPMPADDR8:  {"pmpaddr8", nil, nil},
	csrPMPADDR9:  {"pmpaddr9", nil, nil},
	csrPMPADDR10: {"pmpaddr10", nil, nil},
	csrPMPADDR11: {"pmpaddr11", nil, nil},
	csrPMPADDR12: {"pmpaddr12", nil, nil},
	csrPMPADDR13: {"pmpaddr13", nil, nil},
	csrPMPADDR14: {"pmpaddr14", nil, nil},
	csrPMPADDR15: {"pmpaddr15", nil, nil},
	csrMCAUSE:    {"mcause", nil, nil},
	csrMHARTID:   {"mhartid", nil, nil},
	csrSATP:      {"satp", nil, nil},
}

//-----------------------------------------------------------------------------

/*

var csrNameTable = map[uint]string{
	csrUSTATUS:        "ustatus",
	csrUIE:            "uie",
	csrUTVEC:          "utvec",
	csrUSCRATCH:       "uscratch",
	csrUEPC:           "uepc",
	csrUCAUSE:         "ucause",
	csrUTVAL:          "utval",
	csrUIP:            "uip",
	csrFFLAGS:         "fflags",
	csrFRM:            "frm",
	csrFCSR:           "fcsr",
	csrCYCLE:          "cycle",
	csrTIME:           "time",
	csrINSTRET:        "instret",
	csrHPMCOUNTER3:    "hpmcounter3",
	csrHPMCOUNTER4:    "hpmcounter4",
	csrHPMCOUNTER5:    "hpmcounter5",
	csrHPMCOUNTER6:    "hpmcounter6",
	csrHPMCOUNTER7:    "hpmcounter7",
	csrHPMCOUNTER8:    "hpmcounter8",
	csrHPMCOUNTER9:    "hpmcounter9",
	csrHPMCOUNTER10:   "hpmcounter10",
	csrHPMCOUNTER11:   "hpmcounter11",
	csrHPMCOUNTER12:   "hpmcounter12",
	csrHPMCOUNTER13:   "hpmcounter13",
	csrHPMCOUNTER14:   "hpmcounter14",
	csrHPMCOUNTER15:   "hpmcounter15",
	csrHPMCOUNTER16:   "hpmcounter16",
	csrHPMCOUNTER17:   "hpmcounter17",
	csrHPMCOUNTER18:   "hpmcounter18",
	csrHPMCOUNTER19:   "hpmcounter19",
	csrHPMCOUNTER20:   "hpmcounter20",
	csrHPMCOUNTER21:   "hpmcounter21",
	csrHPMCOUNTER22:   "hpmcounter22",
	csrHPMCOUNTER23:   "hpmcounter23",
	csrHPMCOUNTER24:   "hpmcounter24",
	csrHPMCOUNTER25:   "hpmcounter25",
	csrHPMCOUNTER26:   "hpmcounter26",
	csrHPMCOUNTER27:   "hpmcounter27",
	csrHPMCOUNTER28:   "hpmcounter28",
	csrHPMCOUNTER29:   "hpmcounter29",
	csrHPMCOUNTER30:   "hpmcounter30",
	csrHPMCOUNTER31:   "hpmcounter31",
	csrCYCLEH:         "cycleh",
	csrTIMEH:          "timeh",
	csrINSTRETH:       "instreth",
	csrHPMCOUNTER3H:   "hpmcounter3h",
	csrHPMCOUNTER4H:   "hpmcounter4h",
	csrHPMCOUNTER5H:   "hpmcounter5h",
	csrHPMCOUNTER6H:   "hpmcounter6h",
	csrHPMCOUNTER7H:   "hpmcounter7h",
	csrHPMCOUNTER8H:   "hpmcounter8h",
	csrHPMCOUNTER9H:   "hpmcounter9h",
	csrHPMCOUNTER10H:  "hpmcounter10h",
	csrHPMCOUNTER11H:  "hpmcounter11h",
	csrHPMCOUNTER12H:  "hpmcounter12h",
	csrHPMCOUNTER13H:  "hpmcounter13h",
	csrHPMCOUNTER14H:  "hpmcounter14h",
	csrHPMCOUNTER15H:  "hpmcounter15h",
	csrHPMCOUNTER16H:  "hpmcounter16h",
	csrHPMCOUNTER17H:  "hpmcounter17h",
	csrHPMCOUNTER18H:  "hpmcounter18h",
	csrHPMCOUNTER19H:  "hpmcounter19h",
	csrHPMCOUNTER20H:  "hpmcounter20h",
	csrHPMCOUNTER21H:  "hpmcounter21h",
	csrHPMCOUNTER22H:  "hpmcounter22h",
	csrHPMCOUNTER23H:  "hpmcounter23h",
	csrHPMCOUNTER24H:  "hpmcounter24h",
	csrHPMCOUNTER25H:  "hpmcounter25h",
	csrHPMCOUNTER26H:  "hpmcounter26h",
	csrHPMCOUNTER27H:  "hpmcounter27h",
	csrHPMCOUNTER28H:  "hpmcounter28h",
	csrHPMCOUNTER29H:  "hpmcounter29h",
	csrHPMCOUNTER30H:  "hpmcounter30h",
	csrHPMCOUNTER31H:  "hpmcounter31h",
	csrSSTATUS:        "sstatus",
	csrSEDELEG:        "sedeleg",
	csrSIDELEG:        "sideleg",
	csrSIE:            "sie",
	csrSTVEC:          "stvec",
	csrSCOUNTEREN:     "scounteren",
	csrSSCRATCH:       "sscratch",
	csrSEPC:           "sepc",
	csrSCAUSE:         "scause",
	csrSTVAL:          "stval",
	csrSIP:            "sip",
	csrMVENDORID:      "mvendorid",
	csrMARCHID:        "marchid",
	csrMIMPID:         "mimpid",
	csrMSTATUS:        "mstatus",
	csrMISA:           "misa",
	csrMEDELEG:        "medeleg",
	csrMIDELEG:        "mideleg",
	csrMIE:            "mie",
	csrMCOUNTEREN:     "mcounteren",
	csrMSCRATCH:       "mscratch",
	csrMEPC:           "mepc",
	csrMTVAL:          "mtval",
	csrMIP:            "mip",
	csrPMPCFG0:        "pmpcfg0",
	csrPMPCFG1:        "pmpcfg1",
	csrPMPCFG2:        "pmpcfg2",
	csrPMPCFG3:        "pmpcfg3",
	csrMCYCLE:         "mcycle",
	csrMINSTRET:       "minstret",
	csrMHPMCOUNTER3:   "mhpmcounter3",
	csrMHPMCOUNTER4:   "mhpmcounter4",
	csrMHPMCOUNTER5:   "mhpmcounter5",
	csrMHPMCOUNTER6:   "mhpmcounter6",
	csrMHPMCOUNTER7:   "mhpmcounter7",
	csrMHPMCOUNTER8:   "mhpmcounter8",
	csrMHPMCOUNTER9:   "mhpmcounter9",
	csrMHPMCOUNTER10:  "mhpmcounter10",
	csrMHPMCOUNTER11:  "mhpmcounter11",
	csrMHPMCOUNTER12:  "mhpmcounter12",
	csrMHPMCOUNTER13:  "mhpmcounter13",
	csrMHPMCOUNTER14:  "mhpmcounter14",
	csrMHPMCOUNTER15:  "mhpmcounter15",
	csrMHPMCOUNTER16:  "mhpmcounter16",
	csrMHPMCOUNTER17:  "mhpmcounter17",
	csrMHPMCOUNTER18:  "mhpmcounter18",
	csrMHPMCOUNTER19:  "mhpmcounter19",
	csrMHPMCOUNTER20:  "mhpmcounter20",
	csrMHPMCOUNTER21:  "mhpmcounter21",
	csrMHPMCOUNTER22:  "mhpmcounter22",
	csrMHPMCOUNTER23:  "mhpmcounter23",
	csrMHPMCOUNTER24:  "mhpmcounter24",
	csrMHPMCOUNTER25:  "mhpmcounter25",
	csrMHPMCOUNTER26:  "mhpmcounter26",
	csrMHPMCOUNTER27:  "mhpmcounter27",
	csrMHPMCOUNTER28:  "mhpmcounter28",
	csrMHPMCOUNTER29:  "mhpmcounter29",
	csrMHPMCOUNTER30:  "mhpmcounter30",
	csrMHPMCOUNTER31:  "mhpmcounter31",
	csrMCYCLEH:        "mcycleh",
	csrMINSTRETH:      "minstreth",
	csrMHPMCOUNTER3H:  "mhpmcounter3h",
	csrMHPMCOUNTER4H:  "mhpmcounter4h",
	csrMHPMCOUNTER5H:  "mhpmcounter5h",
	csrMHPMCOUNTER6H:  "mhpmcounter6h",
	csrMHPMCOUNTER7H:  "mhpmcounter7h",
	csrMHPMCOUNTER8H:  "mhpmcounter8h",
	csrMHPMCOUNTER9H:  "mhpmcounter9h",
	csrMHPMCOUNTER10H: "mhpmcounter10h",
	csrMHPMCOUNTER11H: "mhpmcounter11h",
	csrMHPMCOUNTER12H: "mhpmcounter12h",
	csrMHPMCOUNTER13H: "mhpmcounter13h",
	csrMHPMCOUNTER14H: "mhpmcounter14h",
	csrMHPMCOUNTER15H: "mhpmcounter15h",
	csrMHPMCOUNTER16H: "mhpmcounter16h",
	csrMHPMCOUNTER17H: "mhpmcounter17h",
	csrMHPMCOUNTER18H: "mhpmcounter18h",
	csrMHPMCOUNTER19H: "mhpmcounter19h",
	csrMHPMCOUNTER20H: "mhpmcounter20h",
	csrMHPMCOUNTER21H: "mhpmcounter21h",
	csrMHPMCOUNTER22H: "mhpmcounter22h",
	csrMHPMCOUNTER23H: "mhpmcounter23h",
	csrMHPMCOUNTER24H: "mhpmcounter24h",
	csrMHPMCOUNTER25H: "mhpmcounter25h",
	csrMHPMCOUNTER26H: "mhpmcounter26h",
	csrMHPMCOUNTER27H: "mhpmcounter27h",
	csrMHPMCOUNTER28H: "mhpmcounter28h",
	csrMHPMCOUNTER29H: "mhpmcounter29h",
	csrMHPMCOUNTER30H: "mhpmcounter30h",
	csrMHPMCOUNTER31H: "mhpmcounter31h",
	csrMHPMEVENT3:     "mhpmevent3",
	csrMHPMEVENT4:     "mhpmevent4",
	csrMHPMEVENT5:     "mhpmevent5",
	csrMHPMEVENT6:     "mhpmevent6",
	csrMHPMEVENT7:     "mhpmevent7",
	csrMHPMEVENT8:     "mhpmevent8",
	csrMHPMEVENT9:     "mhpmevent9",
	csrMHPMEVENT10:    "mhpmevent10",
	csrMHPMEVENT11:    "mhpmevent11",
	csrMHPMEVENT12:    "mhpmevent12",
	csrMHPMEVENT13:    "mhpmevent13",
	csrMHPMEVENT14:    "mhpmevent14",
	csrMHPMEVENT15:    "mhpmevent15",
	csrMHPMEVENT16:    "mhpmevent16",
	csrMHPMEVENT17:    "mhpmevent17",
	csrMHPMEVENT18:    "mhpmevent18",
	csrMHPMEVENT19:    "mhpmevent19",
	csrMHPMEVENT20:    "mhpmevent20",
	csrMHPMEVENT21:    "mhpmevent21",
	csrMHPMEVENT22:    "mhpmevent22",
	csrMHPMEVENT23:    "mhpmevent23",
	csrMHPMEVENT24:    "mhpmevent24",
	csrMHPMEVENT25:    "mhpmevent25",
	csrMHPMEVENT26:    "mhpmevent26",
	csrMHPMEVENT27:    "mhpmevent27",
	csrMHPMEVENT28:    "mhpmevent28",
	csrMHPMEVENT29:    "mhpmevent29",
	csrMHPMEVENT30:    "mhpmevent30",
	csrMHPMEVENT31:    "mhpmevent31",
	csrTSELECT:        "tselect",
	csrTDATA1:         "tdata1",
	csrTDATA2:         "tdata2",
	csrTDATA3:         "tdata3",
	csrDCSR:           "dcsr",
	csrDPC:            "dpc",
	csrDSCRATCH:       "dscratch",
	csrHSTATUS:        "hstatus",
	csrHEDELEG:        "hedeleg",
	csrHIDELEG:        "hideleg",
	csrHIE:            "hie",
	csrHTVEC:          "htvec",
	csrHSCRATCH:       "hscratch",
	csrHEPC:           "hepc",
	csrHCAUSE:         "hcause",
	csrHBADADDR:       "hbadaddr",
	csrHIP:            "hip",
	csrMBASE:          "mbase",
	csrMBOUND:         "mbound",
	csrMIBASE:         "mibase",
	csrMIBOUND:        "mibound",
	csrMDBASE:         "mdbase",
	csrMDBOUND:        "mdbound",
	csrMUCOUNTEREN:    "mucounteren",
	csrMSCOUNTEREN:    "mscounteren",
	csrMHCOUNTEREN:    "mhcounteren",
}

*/

func csrName(csr uint) string {
	if cd, ok := csrLUT[csr]; ok {
		return cd.name
	}
	return fmt.Sprintf("0x%03x", csr)
}

//-----------------------------------------------------------------------------

func (m *RV64) rdCSR(csr uint) uint64 {
	return 0
}

func (m *RV64) wrCSR(csr uint, val uint64) {
}

func (m *RV32) rdCSR(csr uint) uint32 {
	return 0
}

func (m *RV32) wrCSR(csr uint, val uint32) {
}

//-----------------------------------------------------------------------------
