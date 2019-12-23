//-----------------------------------------------------------------------------
/*

RISC-V ISA Definition

*/
//-----------------------------------------------------------------------------

package rv

import "github.com/deadsy/riscv/csr"

//-----------------------------------------------------------------------------

// daFunc is an instruction disassembly function
type daFunc func(name string, pc, ins uint) string

// emuFunc is an emulation function
type emuFunc func(m *RV, ins uint) error

// insDefn is the definition of an instruction
type insDefn struct {
	defn string  // instruction definition string (from the standard)
	da   daFunc  // disassembly function
	emu  emuFunc // emulation function
}

// ISAModule is a set/module of RISC-V instructions.
type ISAModule struct {
	ext  uint      // ISA extension bits per CSR misa
	ilen int       // instruction length
	defn []insDefn // instruction definitions
}

//-----------------------------------------------------------------------------
// RV32 instructions

// ISArv32i integer instructions.
var ISArv32i = ISAModule{
	ext:  csr.IsaExtI,
	ilen: 32,
	defn: []insDefn{
		{"imm[31:12] rd 0110111 LUI", daTypeUa, emu_LUI},                             // U
		{"imm[31:12] rd 0010111 AUIPC", daTypeUa, emu_AUIPC},                         // U
		{"imm[20|10:1|11|19:12] rd 1101111 JAL", daTypeJa, emu_JAL},                  // J
		{"imm[11:0] rs1 000 rd 1100111 JALR", daTypeIe, emu_JALR},                    // I
		{"imm[12|10:5] rs2 rs1 000 imm[4:1|11] 1100011 BEQ", daTypeBa, emu_BEQ},      // B
		{"imm[12|10:5] rs2 rs1 001 imm[4:1|11] 1100011 BNE", daTypeBa, emu_BNE},      // B
		{"imm[12|10:5] rs2 rs1 100 imm[4:1|11] 1100011 BLT", daTypeBa, emu_BLT},      // B
		{"imm[12|10:5] rs2 rs1 101 imm[4:1|11] 1100011 BGE", daTypeBa, emu_BGE},      // B
		{"imm[12|10:5] rs2 rs1 110 imm[4:1|11] 1100011 BLTU", daTypeBa, emu_BLTU},    // B
		{"imm[12|10:5] rs2 rs1 111 imm[4:1|11] 1100011 BGEU", daTypeBa, emu_BGEU},    // B
		{"imm[11:0] rs1 000 rd 0000011 LB", daTypeIc, emu_LB},                        // I
		{"imm[11:0] rs1 001 rd 0000011 LH", daTypeIc, emu_LH},                        // I
		{"imm[11:0] rs1 010 rd 0000011 LW", daTypeIc, emu_LW},                        // I
		{"imm[11:0] rs1 100 rd 0000011 LBU", daTypeIc, emu_LBU},                      // I
		{"imm[11:0] rs1 101 rd 0000011 LHU", daTypeIc, emu_LHU},                      // I
		{"imm[11:5] rs2 rs1 000 imm[4:0] 0100011 SB", daTypeSa, emu_SB},              // S
		{"imm[11:5] rs2 rs1 001 imm[4:0] 0100011 SH", daTypeSa, emu_SH},              // S
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100011 SW", daTypeSa, emu_SW},              // S
		{"imm[11:0] rs1 000 rd 0010011 ADDI", daTypeIb, emu_ADDI},                    // I
		{"imm[11:0] rs1 010 rd 0010011 SLTI", daTypeIa, emu_SLTI},                    // I
		{"imm[11:0] rs1 011 rd 0010011 SLTIU", daTypeIa, emu_SLTIU},                  // I
		{"imm[11:0] rs1 100 rd 0010011 XORI", daTypeIf, emu_XORI},                    // I
		{"imm[11:0] rs1 110 rd 0010011 ORI", daTypeIa, emu_ORI},                      // I
		{"imm[11:0] rs1 111 rd 0010011 ANDI", daTypeIa, emu_ANDI},                    // I
		{"000000 shamt6 rs1 001 rd 0010011 SLLI", daTypeId, emu_SLLI},                // I
		{"000000 shamt6 rs1 101 rd 0010011 SRLI", daTypeId, emu_SRLI},                // I
		{"010000 shamt6 rs1 101 rd 0010011 SRAI", daTypeId, emu_SRAI},                // I
		{"0000000 rs2 rs1 000 rd 0110011 ADD", daTypeRa, emu_ADD},                    // R
		{"0100000 rs2 rs1 000 rd 0110011 SUB", daTypeRa, emu_SUB},                    // R
		{"0000000 rs2 rs1 001 rd 0110011 SLL", daTypeRa, emu_SLL},                    // R
		{"0000000 rs2 rs1 010 rd 0110011 SLT", daTypeRa, emu_SLT},                    // R
		{"0000000 rs2 rs1 011 rd 0110011 SLTU", daTypeRa, emu_SLTU},                  // R
		{"0000000 rs2 rs1 100 rd 0110011 XOR", daTypeRa, emu_XOR},                    // R
		{"0000000 rs2 rs1 101 rd 0110011 SRL", daTypeRa, emu_SRL},                    // R
		{"0100000 rs2 rs1 101 rd 0110011 SRA", daTypeRa, emu_SRA},                    // R
		{"0000000 rs2 rs1 110 rd 0110011 OR", daTypeRa, emu_OR},                      // R
		{"0000000 rs2 rs1 111 rd 0110011 AND", daTypeRa, emu_AND},                    // R
		{"0000 pred succ 00000 000 00000 0001111 FENCE", daTypeIi, emu_FENCE},        // I
		{"0000 0000 0000 00000 001 00000 0001111 FENCE.I", daTypeIi, emu_FENCE_I},    // I
		{"0000000 00000 00000 000 00000 1110011 ECALL", daTypeIi, emu_ECALL},         // I
		{"0000000 00001 00000 000 00000 1110011 EBREAK", daTypeIi, emu_EBREAK},       // I
		{"0000000 00010 00000 000 00000 1110011 URET", daTypeIi, emu_URET},           // I
		{"0001000 00010 00000 000 00000 1110011 SRET", daTypeIi, emu_SRET},           // I
		{"0011000 00010 00000 000 00000 1110011 MRET", daTypeIi, emu_MRET},           // I
		{"0001000 00101 00000 000 00000 1110011 WFI", daTypeIi, emu_WFI},             // I
		{"0001001 rs2 rs1 000 00000 1110011 SFENCE.VMA", daTypeIk, emu_SFENCE_VMA},   // I
		{"0010001 rs2 rs1 000 00000 1110011 HFENCE.BVMA", daTypeIk, emu_HFENCE_BVMA}, // I
		{"1010001 rs2 rs1 000 00000 1110011 HFENCE.GVMA", daTypeIk, emu_HFENCE_GVMA}, // I
		{"csr rs1 001 rd 1110011 CSRRW", daTypeIh, emu_CSRRW},                        // I
		{"csr rs1 010 rd 1110011 CSRRS", daTypeIh, emu_CSRRS},                        // I
		{"csr rs1 011 rd 1110011 CSRRC", daTypeIh, emu_CSRRC},                        // I
		{"csr zimm 101 rd 1110011 CSRRWI", daTypeIj, emu_CSRRWI},                     // I
		{"csr zimm 110 rd 1110011 CSRRSI", daTypeIj, emu_CSRRSI},                     // I
		{"csr zimm 111 rd 1110011 CSRRCI", daTypeIj, emu_CSRRCI},                     // I
	},
}

// ISArv32m integer multiplication/division instructions.
var ISArv32m = ISAModule{
	ext:  csr.IsaExtM,
	ilen: 32,
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0110011 MUL", daTypeRa, emu_MUL},       // R
		{"0000001 rs2 rs1 001 rd 0110011 MULH", daTypeRa, emu_MULH},     // R
		{"0000001 rs2 rs1 010 rd 0110011 MULHSU", daTypeRa, emu_MULHSU}, // R
		{"0000001 rs2 rs1 011 rd 0110011 MULHU", daTypeRa, emu_MULHU},   // R
		{"0000001 rs2 rs1 100 rd 0110011 DIV", daTypeRa, emu_DIV},       // R
		{"0000001 rs2 rs1 101 rd 0110011 DIVU", daTypeRa, emu_DIVU},     // R
		{"0000001 rs2 rs1 110 rd 0110011 REM", daTypeRa, emu_REM},       // R
		{"0000001 rs2 rs1 111 rd 0110011 REMU", daTypeRa, emu_REMU},     // R
	},
}

// ISArv32a atomic operation instructions.
var ISArv32a = ISAModule{
	ext:  csr.IsaExtA,
	ilen: 32,
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 010 rd 0101111 LR.W", daTypeRb, emu_LR_W},         // R
		{"00011 aq rl rs2 rs1 010 rd 0101111 SC.W", daTypeRb, emu_SC_W},           // R
		{"00001 aq rl rs2 rs1 010 rd 0101111 AMOSWAP.W", daTypeRb, emu_AMOSWAP_W}, // R
		{"00000 aq rl rs2 rs1 010 rd 0101111 AMOADD.W", daTypeRb, emu_AMOADD_W},   // R
		{"00100 aq rl rs2 rs1 010 rd 0101111 AMOXOR.W", daTypeRb, emu_AMOXOR_W},   // R
		{"01100 aq rl rs2 rs1 010 rd 0101111 AMOAND.W", daTypeRb, emu_AMOAND_W},   // R
		{"01000 aq rl rs2 rs1 010 rd 0101111 AMOOR.W", daTypeRb, emu_AMOOR_W},     // R
		{"10000 aq rl rs2 rs1 010 rd 0101111 AMOMIN.W", daTypeRb, emu_AMOMIN_W},   // R
		{"10100 aq rl rs2 rs1 010 rd 0101111 AMOMAX.W", daTypeRb, emu_AMOMAX_W},   // R
		{"11000 aq rl rs2 rs1 010 rd 0101111 AMOMINU.W", daTypeRb, emu_AMOMINU_W}, // R
		{"11100 aq rl rs2 rs1 010 rd 0101111 AMOMAXU.W", daTypeRb, emu_AMOMAXU_W}, // R
	},
}

// ISArv32f 32-bit floating point instructions.
var ISArv32f = ISAModule{
	ext:  csr.IsaExtF,
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 010 rd 0000111 FLW", daTypeIg, emu_FLW},                // I
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100111 FSW", daTypeSb, emu_FSW},      // S
		{"rs3 00 rs2 rs1 rm rd 1000011 FMADD.S", daTypeR4a, emu_FMADD_S},       // R4
		{"rs3 00 rs2 rs1 rm rd 1000111 FMSUB.S", daTypeR4a, emu_FMSUB_S},       // R4
		{"rs3 00 rs2 rs1 rm rd 1001011 FNMSUB.S", daTypeR4a, emu_FNMSUB_S},     // R4
		{"rs3 00 rs2 rs1 rm rd 1001111 FNMADD.S", daTypeR4a, emu_FNMADD_S},     // R4
		{"0000000 rs2 rs1 rm rd 1010011 FADD.S", daTypeRc, emu_FADD_S},         // R
		{"0000100 rs2 rs1 rm rd 1010011 FSUB.S", daTypeRc, emu_FSUB_S},         // R
		{"0001000 rs2 rs1 rm rd 1010011 FMUL.S", daTypeRc, emu_FMUL_S},         // R
		{"0001100 rs2 rs1 rm rd 1010011 FDIV.S", daTypeRc, emu_FDIV_S},         // R
		{"0101100 00000 rs1 rm rd 1010011 FSQRT.S", daTypeRh, emu_FSQRT_S},     // R
		{"0010000 rs2 rs1 000 rd 1010011 FSGNJ.S", daTypeRc, emu_FSGNJ_S},      // R
		{"0010000 rs2 rs1 001 rd 1010011 FSGNJN.S", daTypeRc, emu_FSGNJN_S},    // R
		{"0010000 rs2 rs1 010 rd 1010011 FSGNJX.S", daTypeRc, emu_FSGNJX_S},    // R
		{"0010100 rs2 rs1 000 rd 1010011 FMIN.S", daTypeRc, emu_FMIN_S},        // R
		{"0010100 rs2 rs1 001 rd 1010011 FMAX.S", daTypeRc, emu_FMAX_S},        // R
		{"1100000 00000 rs1 rm rd 1010011 FCVT.W.S", daTypeRk, emu_FCVT_W_S},   // R
		{"1100000 00001 rs1 rm rd 1010011 FCVT.WU.S", daTypeRk, emu_FCVT_WU_S}, // R
		{"1110000 00000 rs1 000 rd 1010011 FMV.X.W", daTypeRd, emu_FMV_X_W},    // R
		{"1010000 rs2 rs1 010 rd 1010011 FEQ.S", daTypeRf, emu_FEQ_S},          // R
		{"1010000 rs2 rs1 001 rd 1010011 FLT.S", daTypeRf, emu_FLT_S},          // R
		{"1010000 rs2 rs1 000 rd 1010011 FLE.S", daTypeRf, emu_FLE_S},          // R
		{"1110000 00000 rs1 001 rd 1010011 FCLASS.S", daTypeRd, emu_FCLASS_S},  // R
		{"1101000 00000 rs1 rm rd 1010011 FCVT.S.W", daTypeRj, emu_FCVT_S_W},   // R
		{"1101000 00001 rs1 rm rd 1010011 FCVT.S.WU", daTypeRj, emu_FCVT_S_WU}, // R
		{"1111000 00000 rs1 000 rd 1010011 FMV.W.X", daTypeRe, emu_FMV_W_X},    // R
	},
}

// ISArv32d 64-bit floating point instructions.
var ISArv32d = ISAModule{
	ext:  csr.IsaExtD,
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 011 rd 0000111 FLD", daTypeIg, emu_FLD},                // I
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100111 FSD", daTypeSb, emu_FSD},      // S
		{"rs3 01 rs2 rs1 rm rd 1000011 FMADD.D", daTypeR4a, emu_FMADD_D},       // R4
		{"rs3 01 rs2 rs1 rm rd 1000111 FMSUB.D", daTypeR4a, emu_FMSUB_D},       // R4
		{"rs3 01 rs2 rs1 rm rd 1001011 FNMSUB.D", daTypeR4a, emu_FNMSUB_D},     // R4
		{"rs3 01 rs2 rs1 rm rd 1001111 FNMADD.D", daTypeR4a, emu_FNMADD_D},     // R4
		{"0000001 rs2 rs1 rm rd 1010011 FADD.D", daTypeRc, emu_FADD_D},         // R
		{"0000101 rs2 rs1 rm rd 1010011 FSUB.D", daTypeRc, emu_FSUB_D},         // R
		{"0001001 rs2 rs1 rm rd 1010011 FMUL.D", daTypeRc, emu_FMUL_D},         // R
		{"0001101 rs2 rs1 rm rd 1010011 FDIV.D", daTypeRc, emu_FDIV_D},         // R
		{"0101101 00000 rs1 rm rd 1010011 FSQRT.D", daTypeRh, emu_FSQRT_D},     // R
		{"0010001 rs2 rs1 000 rd 1010011 FSGNJ.D", daTypeRc, emu_FSGNJ_D},      // R
		{"0010001 rs2 rs1 001 rd 1010011 FSGNJN.D", daTypeRc, emu_FSGNJN_D},    // R
		{"0010001 rs2 rs1 010 rd 1010011 FSGNJX.D", daTypeRc, emu_FSGNJX_D},    // R
		{"0010101 rs2 rs1 000 rd 1010011 FMIN.D", daTypeRc, emu_FMIN_D},        // R
		{"0010101 rs2 rs1 001 rd 1010011 FMAX.D", daTypeRc, emu_FMAX_D},        // R
		{"0100000 00001 rs1 rm rd 1010011 FCVT.S.D", daTypeRi, emu_FCVT_S_D},   // R
		{"0100001 00000 rs1 rm rd 1010011 FCVT.D.S", daTypeRi, emu_FCVT_D_S},   // R
		{"1010001 rs2 rs1 010 rd 1010011 FEQ.D", daTypeRf, emu_FEQ_D},          // R
		{"1010001 rs2 rs1 001 rd 1010011 FLT.D", daTypeRf, emu_FLT_D},          // R
		{"1010001 rs2 rs1 000 rd 1010011 FLE.D", daTypeRf, emu_FLE_D},          // R
		{"1110001 00000 rs1 001 rd 1010011 FCLASS.D", daTypeRd, emu_FCLASS_D},  // R
		{"1100001 00000 rs1 rm rd 1010011 FCVT.W.D", daTypeRk, emu_FCVT_W_D},   // R
		{"1100001 00001 rs1 rm rd 1010011 FCVT.WU.D", daTypeRk, emu_FCVT_WU_D}, // R
		{"1101001 00000 rs1 rm rd 1010011 FCVT.D.W", daTypeRj, emu_FCVT_D_W},   // R
		{"1101001 00001 rs1 rm rd 1010011 FCVT.D.WU", daTypeRj, emu_FCVT_D_WU}, // R
	},
}

// ISArv32c compressed instructions (subset of RV64C).
var ISArv32c = ISAModule{
	ext:  csr.IsaExtC,
	ilen: 16,
	defn: []insDefn{
		{"000 00000000 000 00 C.ILLEGAL", daTypeCIWa, emu_C_ILLEGAL},                     // CIW (Quadrant 0)
		{"000 nzuimm[5:4|9:6|2|3] rd0 00 C.ADDI4SPN", daTypeCIWb, emu_C_ADDI4SPN},        // CIW
		{"010 uimm[5:3] rs10 uimm[2|6] rd0 00 C.LW", daTypeCSa, emu_C_LW},                // CL
		{"110 uimm[5:3] rs10 uimm[2|6] rs20 00 C.SW", daTypeCSa, emu_C_SW},               // CS
		{"000 nzimm[5] 00000 nzimm[4:0] 01 C.NOP", daNop, emu_C_NOP},                     // CI (Quadrant 1)
		{"000 nzimm[5] rs1/rd!=0 nzimm[4:0] 01 C.ADDI", daTypeCIc, emu_C_ADDI},           // CI
		{"010 imm[5] rd!=0 imm[4:0] 01 C.LI", daTypeCIa, emu_C_LI},                       // CI
		{"011 nzimm[9] 00010 nzimm[4|6|8:7|5] 01 C.ADDI16SP", daTypeCIb, emu_C_ADDI16SP}, // CI
		{"011 nzimm[17] rd!={0,2} nzimm[16:12] 01 C.LUI", daTypeCIg, emu_C_LUI},          // CI
		{"100 nzuimm[5] 00 rs10/rd0 nzuimm[4:0] 01 C.SRLI", daTypeCId, emu_C_SRLI},       // CI
		{"100 nzuimm[5] 01 rs10/rd0 nzuimm[4:0] 01 C.SRAI", daTypeCId, emu_C_SRAI},       // CI
		{"100 imm[5] 10 rs10/rd0 imm[4:0] 01 C.ANDI", daTypeCIf, emu_C_ANDI},             // CI
		{"100 0 11 rs10/rd0 00 rs20 01 C.SUB", daTypeCRc, emu_C_SUB},                     // CR
		{"100 0 11 rs10/rd0 01 rs20 01 C.XOR", daTypeCRc, emu_C_XOR},                     // CR
		{"100 0 11 rs10/rd0 10 rs20 01 C.OR", daTypeCRc, emu_C_OR},                       // CR
		{"100 0 11 rs10/rd0 11 rs20 01 C.AND", daTypeCRc, emu_C_AND},                     // CR
		{"101 imm[11|4|9:8|10|6|7|3:1|5] 01 C.J", daTypeCJb, emu_C_J},                    // CJ
		{"110 imm[8|4:3] rs10 imm[7:6|2:1|5] 01 C.BEQZ", daTypeCBa, emu_C_BEQZ},          // CB
		{"111 imm[8|4:3] rs10 imm[7:6|2:1|5] 01 C.BNEZ", daTypeCBa, emu_C_BNEZ},          // CB
		{"000 nzuimm[5] rs1/rd!=0 nzuimm[4:0] 10 C.SLLI", daTypeCIe, emu_C_SLLI},         // CI (Quadrant 2)
		{"000 0 rs1/rd!=0 00000 10 C.SLLI64", daNone, emu_C_SLLI64},                      // CI
		{"010 uimm[5] rd!=0 uimm[4:2|7:6] 10 C.LWSP", daTypeCSSa, emu_C_LWSP},            // CSS
		{"100 0 rs1!=0 00000 10 C.JR", daTypeCRd, emu_C_JR},                              // CR
		{"100 0 rd!=0 rs2!=0 10 C.MV", daTypeCRa, emu_C_MV},                              // CR
		{"100 1 00000 00000 10 C.EBREAK", daNone, emu_C_EBREAK},                          // CI
		{"100 1 rs1!=0 00000 10 C.JALR", daTypeCRe, emu_C_JALR},                          // CR
		{"100 1 rs1/rd!=0 rs2!=0 10 C.ADD", daTypeCRb, emu_C_ADD},                        // CR
		{"110 uimm[5:2|7:6] rs2 10 C.SWSP", daTypeCSSb, emu_C_SWSP},                      // CSS
	},
}

// ISArv32cOnly compressed instructions (not in RV64C).
var ISArv32cOnly = ISAModule{
	ext:  csr.IsaExtC,
	ilen: 16,
	defn: []insDefn{
		{"001 imm[11|4|9:8|10|6|7|3:1|5] 01 C.JAL", daTypeCJc, emu_C_JAL}, // CJ
	},
}

// ISArv32fc compressed 32-bit floating point instructions.
var ISArv32fc = ISAModule{
	ext:  csr.IsaExtC,
	ilen: 16,
	defn: []insDefn{
		{"011 uimm[5:3] rs10 uimm[2|6] rd0 00 C.FLW", daTypeCSc, emu_C_FLW},  // CL
		{"011 uimm[5] rd uimm[4:2|7:6] 10 C.FLWSP", daNone, emu_C_FLWSP},     // CSS
		{"111 uimm[5:3] rs10 uimm[2|6] rs20 00 C.FSW", daTypeCSc, emu_C_FSW}, // CS
		{"111 uimm[5:2|7:6] rs2 10 C.FSWSP", daNone, emu_C_FSWSP},            // CSS
	},
}

// ISArv32dc compressed 64-bit floating point instructions.
var ISArv32dc = ISAModule{
	ext:  csr.IsaExtC,
	ilen: 16,
	defn: []insDefn{
		{"001 uimm[5:3] rs10 uimm[7:6] rd0 00 C.FLD", daTypeCSc, emu_C_FLD},  // CL
		{"001 uimm[5] rd uimm[4:3|8:6] 10 C.FLDSP", daNone, emu_C_FLDSP},     // CSS
		{"101 uimm[5:3] rs10 uimm[7:6] rs20 00 C.FSD", daTypeCSc, emu_C_FSD}, // CS
		{"101 uimm[5:3|8:6] rs2 10 C.FSDSP", daNone, emu_C_FSDSP},            // CSS
	},
}

//-----------------------------------------------------------------------------
// RV64 instructions (+ RV32)

// ISArv64i Integer
var ISArv64i = ISAModule{
	ext:  csr.IsaExtI,
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 110 rd 0000011 LWU", daTypeIc, emu_LWU},          // I
		{"imm[11:0] rs1 011 rd 0000011 LD", daTypeIa, emu_LD},            // I
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100011 SD", daTypeSa, emu_SD},  // S
		{"000000 shamt6 rs1 001 rd 0010011 SLLI", daTypeId, emu_SLLI},    // I
		{"000000 shamt6 rs1 101 rd 0010011 SRLI", daTypeId, emu_SRLI},    // I
		{"010000 shamt6 rs1 101 rd 0010011 SRAI", daTypeId, emu_SRAI},    // I
		{"imm[11:0] rs1 000 rd 0011011 ADDIW", daTypeIa, emu_ADDIW},      // I
		{"0000000 shamt5 rs1 001 rd 0011011 SLLIW", daTypeId, emu_SLLIW}, // I
		{"0000000 shamt5 rs1 101 rd 0011011 SRLIW", daTypeId, emu_SRLIW}, // I
		{"0100000 shamt5 rs1 101 rd 0011011 SRAIW", daTypeId, emu_SRAIW}, // I
		{"0000000 rs2 rs1 000 rd 0111011 ADDW", daTypeRa, emu_ADDW},      // R
		{"0100000 rs2 rs1 000 rd 0111011 SUBW", daTypeRa, emu_SUBW},      // R
		{"0000000 rs2 rs1 001 rd 0111011 SLLW", daTypeRa, emu_SLLW},      // R
		{"0000000 rs2 rs1 101 rd 0111011 SRLW", daTypeRa, emu_SRLW},      // R
		{"0100000 rs2 rs1 101 rd 0111011 SRAW", daTypeRa, emu_SRAW},      // R
	},
}

// ISArv64m Integer Multiplication and Division
var ISArv64m = ISAModule{
	ext:  csr.IsaExtM,
	ilen: 32,
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0111011 MULW", daTypeRa, emu_MULW},   // R
		{"0000001 rs2 rs1 100 rd 0111011 DIVW", daTypeRa, emu_DIVW},   // R
		{"0000001 rs2 rs1 101 rd 0111011 DIVUW", daTypeRa, emu_DIVUW}, // R
		{"0000001 rs2 rs1 110 rd 0111011 REMW", daTypeRa, emu_REMW},   // R
		{"0000001 rs2 rs1 111 rd 0111011 REMUW", daTypeRa, emu_REMUW}, // R
	},
}

// ISArv64a Atomics
var ISArv64a = ISAModule{
	ext:  csr.IsaExtA,
	ilen: 32,
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 011 rd 0101111 LR.D", daTypeRb, emu_LR_D},         // R
		{"00011 aq rl rs2 rs1 011 rd 0101111 SC.D", daTypeRb, emu_SC_D},           // R
		{"00001 aq rl rs2 rs1 011 rd 0101111 AMOSWAP.D", daTypeRb, emu_AMOSWAP_D}, // R
		{"00000 aq rl rs2 rs1 011 rd 0101111 AMOADD.D", daTypeRb, emu_AMOADD_D},   // R
		{"00100 aq rl rs2 rs1 011 rd 0101111 AMOXOR.D", daTypeRb, emu_AMOXOR_D},   // R
		{"01100 aq rl rs2 rs1 011 rd 0101111 AMOAND.D", daTypeRb, emu_AMOAND_D},   // R
		{"01000 aq rl rs2 rs1 011 rd 0101111 AMOOR.D", daTypeRb, emu_AMOOR_D},     // R
		{"10000 aq rl rs2 rs1 011 rd 0101111 AMOMIN.D", daTypeRb, emu_AMOMIN_D},   // R
		{"10100 aq rl rs2 rs1 011 rd 0101111 AMOMAX.D", daTypeRb, emu_AMOMAX_D},   // R
		{"11000 aq rl rs2 rs1 011 rd 0101111 AMOMINU.D", daTypeRb, emu_AMOMINU_D}, // R
		{"11100 aq rl rs2 rs1 011 rd 0101111 AMOMAXU.D", daTypeRb, emu_AMOMAXU_D}, // R
	},
}

// ISArv64f Single-Precision Floating-Point
var ISArv64f = ISAModule{
	ext:  csr.IsaExtF,
	ilen: 32,
	defn: []insDefn{
		{"1100000 00010 rs1 rm rd 1010011 FCVT.L.S", daTypeRk, emu_FCVT_L_S},   // R
		{"1100000 00011 rs1 rm rd 1010011 FCVT.LU.S", daTypeRk, emu_FCVT_LU_S}, // R
		{"1101000 00010 rs1 rm rd 1010011 FCVT.S.L", daTypeRj, emu_FCVT_S_L},   // R
		{"1101000 00011 rs1 rm rd 1010011 FCVT.S.LU", daTypeRj, emu_FCVT_S_LU}, // R
	},
}

// ISArv64d Double-Precision Floating-Point
var ISArv64d = ISAModule{
	ext:  csr.IsaExtD,
	ilen: 32,
	defn: []insDefn{
		{"1100001 00010 rs1 rm rd 1010011 FCVT.L.D", daTypeRk, emu_FCVT_L_D},   // R
		{"1100001 00011 rs1 rm rd 1010011 FCVT.LU.D", daTypeRk, emu_FCVT_LU_D}, // R
		{"1110001 00000 rs1 000 rd 1010011 FMV.X.D", daTypeRd, emu_FMV_X_D},    // R
		{"1101001 00010 rs1 rm rd 1010011 FCVT.D.L", daTypeRj, emu_FCVT_D_L},   // R
		{"1101001 00011 rs1 rm rd 1010011 FCVT.D.LU", daTypeRj, emu_FCVT_D_LU}, // R
		{"1111001 00000 rs1 000 rd 1010011 FMV.D.X", daTypeRe, emu_FMV_D_X},    // R
	},
}

// ISArv64c Compressed
var ISArv64c = ISAModule{
	ext:  csr.IsaExtC,
	ilen: 16,
	defn: []insDefn{
		{"001 imm[5] rd!=0 imm[4:0] 01 C.ADDIW", daTypeCIc, emu_C_ADDIW},   // CI
		{"011 uimm[5] rd uimm[4:3|8:6] 10 C.LDSP", daTypeCIh, emu_C_LDSP},  // CI
		{"011 uimm[5:3] rs10 uimm[7:6] rd0 00 C.LD", daTypeCSb, emu_C_LD},  // CL
		{"100 1 11 rs10/rd0 00 rs20 01 C.SUBW", daTypeCRc, emu_C_SUBW},     // CR
		{"100 1 11 rs10/rd0 01 rs20 01 C.ADDW", daTypeCRc, emu_C_ADDW},     // CR
		{"111 uimm[5:3] rs10 uimm[7:6] rs20 00 C.SD", daTypeCSb, emu_C_SD}, // CS
		{"111 uimm[5:3|8:6] rs2 10 C.SDSP", daTypeCSSc, emu_C_SDSP},        // CSS
	},
}

//-----------------------------------------------------------------------------

// ISArv128c Compressed
var ISArv128c = ISAModule{
	ext:  csr.IsaExtC,
	ilen: 16,
	defn: []insDefn{
		// C.SQ
		// C.LQ
		// C.LQSP
	},
}

//-----------------------------------------------------------------------------
// pre-canned ISA module sets

// ISArv32g = RV32imafd
var ISArv32g = []ISAModule{
	ISArv32i, ISArv32m, ISArv32a, ISArv32f, ISArv32d,
}

// ISArv32gc = RV32imafdc
var ISArv32gc = []ISAModule{
	ISArv32i, ISArv32m, ISArv32a, ISArv32f, ISArv32d,
	ISArv32c, ISArv32cOnly, ISArv32fc, ISArv32dc,
}

// ISArv64g = RV64imafd
var ISArv64g = []ISAModule{
	ISArv32i, ISArv32m, ISArv32a, ISArv32f, ISArv32d,
	ISArv64i, ISArv64m, ISArv64a, ISArv64f, ISArv64d,
}

// ISArv64gc = RV64imafdc
var ISArv64gc = []ISAModule{
	ISArv32i, ISArv32m, ISArv32a, ISArv32f, ISArv32d,
	ISArv32c, ISArv32dc,
	ISArv64i, ISArv64m, ISArv64a, ISArv64f, ISArv64d,
	ISArv64c,
}

//-----------------------------------------------------------------------------

// insMeta is instruction meta-data determined at runtime
type insMeta struct {
	defn      *insDefn   // the instruction definition
	name      string     // instruction mneumonic
	n         int        // instruction bit length
	val, mask uint       // value and mask of fixed bits in the instruction
	dt        decodeType // decode type
}

// ISA is an instruction set
type ISA struct {
	ext   uint       // ISA extension bits per CSR misa
	ins16 []*insMeta // the set of 16-bit instructions in the ISA
	ins32 []*insMeta // the set of 32-bit instructions in the ISA
}

// NewISA creates an empty instruction set.
func NewISA() *ISA {
	return &ISA{
		ins16: make([]*insMeta, 0),
		ins32: make([]*insMeta, 0),
	}
}

// Add a ISA sub-module to the ISA.
func (isa *ISA) Add(module []ISAModule) error {
	for i := range module {
		isa.ext |= module[i].ext
		for j := range module[i].defn {
			im, err := parseDefn(&module[i].defn[j], module[i].ilen)
			if err != nil {
				return err
			}
			if im.n == 16 {
				isa.ins16 = append(isa.ins16, im)
			} else {
				isa.ins32 = append(isa.ins32, im)
			}
		}
	}
	return nil
}

// lookup returns the instruction meta information for an instruction.
func (isa *ISA) lookup(ins uint) *insMeta {
	if ins&3 == 3 {
		// 32-bit instruction
		for _, im := range isa.ins32 {
			if ins&im.mask == im.val {
				return im
			}
		}
	} else {
		// 16-bit instruction
		for _, im := range isa.ins16 {
			if ins&im.mask == im.val {
				return im
			}
		}
	}
	return nil
}

// GetExtensions returns the ISA extension bits.
func (isa *ISA) GetExtensions() uint {
	return isa.ext
}

//-----------------------------------------------------------------------------
