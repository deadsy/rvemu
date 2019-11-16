//-----------------------------------------------------------------------------
/*

RISC-V 32-bit Emulator

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
)

//-----------------------------------------------------------------------------
// default emulation

func emu32_Illegal(m *RV32, ins uint) {
	// Trying to run an rv64 instruction on an rv32.
	m.flag |= flagIllegal
}

//-----------------------------------------------------------------------------
// rv32i

func emu32_LUI(m *RV32, ins uint) {
	imm, rd := decodeU(ins)
	m.wrX(rd, uint32(imm<<12))
	m.PC += 4
}

func emu32_AUIPC(m *RV32, ins uint) {
	imm, rd := decodeU(ins)
	m.wrX(rd, uint32(int(m.PC)+(imm<<12)))
	m.PC += 4
}

func emu32_JAL(m *RV32, ins uint) {
	imm, rd := decodeJ(ins)
	m.wrX(rd, m.PC+4)
	m.PC = uint32(int(m.PC) + int(imm))
}

func emu32_JALR(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.PC+4)
	m.PC = uint32((int(m.X[rs1]) + imm) & ^1)
}

func emu32_BEQ(m *RV32, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.X[rs1] == m.X[rs2] {
		m.PC = uint32(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu32_BNE(m *RV32, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.X[rs1] != m.X[rs2] {
		m.PC = uint32(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu32_BLT(m *RV32, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	x1 := bitSex(int(m.X[rs1]), 31)
	x2 := bitSex(int(m.X[rs2]), 31)
	if x1 < x2 {
		m.PC = uint32(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu32_BGE(m *RV32, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	x1 := bitSex(int(m.X[rs1]), 31)
	x2 := bitSex(int(m.X[rs2]), 31)
	if x1 >= x2 {
		m.PC = uint32(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu32_BLTU(m *RV32, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.X[rs1] < m.X[rs2] {
		m.PC = uint32(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu32_BGEU(m *RV32, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.X[rs1] >= m.X[rs2] {
		m.PC = uint32(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu32_LB(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd8(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, uint32(bitSex(int(val), 7)))
	m.PC += 4
}

func emu32_LH(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd16(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, uint32(bitSex(int(val), 15)))
	m.PC += 4
}

func emu32_LW(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd32(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, val)
	m.PC += 4
}

func emu32_LBU(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd8(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, uint32(val))
	m.PC += 4
}

func emu32_LHU(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd16(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, uint32(val))
	m.PC += 4
}

func emu32_SB(m *RV32, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.X[rs1]) + imm)
	ex := m.Mem.Wr8(adr, uint8(m.X[rs2]))
	m.checkMemory(adr, ex)
	m.PC += 4
}

func emu32_SH(m *RV32, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.X[rs1]) + imm)
	ex := m.Mem.Wr16(adr, uint16(m.X[rs2]))
	m.checkMemory(adr, ex)
	m.PC += 4
}

func emu32_SW(m *RV32, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.X[rs1]) + imm)
	ex := m.Mem.Wr32(adr, m.X[rs2])
	m.checkMemory(adr, ex)
	m.PC += 4
}

func emu32_ADDI(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, uint32(int(m.X[rs1])+imm))
	m.PC += 4
}

func emu32_SLTI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SLTIU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_XORI(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.X[rs1]^uint32(imm))
	m.PC += 4
}

func emu32_ORI(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.X[rs1]|uint32(imm))
	m.PC += 4
}

func emu32_ANDI(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.X[rs1]&uint32(imm))
	m.PC += 4
}

func emu32_SLLI(m *RV32, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if shamt > 31 {
		m.flag |= flagIllegal
		return
	}
	m.wrX(rd, m.X[rs1]<<shamt)
	m.PC += 4
}

func emu32_SRLI(m *RV32, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if shamt > 31 {
		m.flag |= flagIllegal
		return
	}
	m.wrX(rd, m.X[rs1]>>shamt)
	m.PC += 4
}

func emu32_SRAI(m *RV32, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if shamt > 31 {
		m.flag |= flagIllegal
		return
	}
	m.wrX(rd, uint32(int(m.X[rs1])>>shamt))
	m.PC += 4
}

func emu32_ADD(m *RV32, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]+m.X[rs2])
	m.PC += 4
}

func emu32_SUB(m *RV32, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]-m.X[rs2])
	m.PC += 4
}

func emu32_SLL(m *RV32, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]<<(m.X[rs2]&31))
	m.PC += 4
}

func emu32_SLT(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SLTU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_XOR(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SRL(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SRA(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_OR(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_AND(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FENCE(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FENCE_I(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_ECALL(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_EBREAK(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_CSRRW(m *RV32, ins uint) {
	csr, rs1, rd := decodeIb(ins)
	t := m.rdCSR(csr)
	m.wrCSR(csr, m.X[rs1])
	m.wrX(rd, t)
	m.PC += 4
}

func emu32_CSRRS(m *RV32, ins uint) {
	csr, rs1, rd := decodeIb(ins)
	t := m.rdCSR(csr)
	if rs1 != 0 {
		m.wrCSR(csr, t|m.X[rs1])
	}
	m.wrX(rd, t)
	m.PC += 4
}

func emu32_CSRRC(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_CSRRWI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_CSRRSI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_CSRRCI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv32m

func emu32_MUL(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_MULH(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_MULHSU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_MULHU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_DIV(m *RV32, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	result := -1
	if m.X[rs2] != 0 {
		result = int(m.X[rs1]) / int(m.X[rs2])
	}
	m.wrX(rd, uint32(result))
	m.PC += 4
}

func emu32_DIVU(m *RV32, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	result := uint32((1 << 32) - 1)
	if m.X[rs2] != 0 {
		result = m.X[rs1] / m.X[rs2]
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu32_REM(m *RV32, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	result := int(m.X[rs1])
	if m.X[rs2] != 0 {
		result %= int(m.X[rs2])
	}
	m.wrX(rd, uint32(result))
	m.PC += 4
}

func emu32_REMU(m *RV32, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	result := m.X[rs1]
	if m.X[rs2] != 0 {
		result %= m.X[rs2]
	}
	m.wrX(rd, result)
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv32a

func emu32_LR_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SC_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_AMOSWAP_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_AMOADD_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_AMOXOR_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_AMOAND_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_AMOOR_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_AMOMIN_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_AMOMAX_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_AMOMINU_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_AMOMAXU_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv32f

func emu32_FLW(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSW(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMADD_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMSUB_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FNMSUB_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FNMADD_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FADD_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSUB_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMUL_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FDIV_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSQRT_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSGNJ_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSGNJN_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSGNJX_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMIN_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMAX_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCVT_W_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCVT_WU_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMV_X_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FEQ_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FLT_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FLE_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCLASS_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCVT_S_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCVT_S_WU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMV_W_X(m *RV32, ins uint) {
	_, rs1, rd := decodeR(ins)
	m.F[rd] = u32Upper | uint64(m.X[rs1])
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv32d

func emu32_FLD(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSD(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMADD_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMSUB_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FNMSUB_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FNMADD_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FADD_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSUB_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMUL_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FDIV_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSQRT_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSGNJ_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSGNJN_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FSGNJX_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMIN_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FMAX_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCVT_S_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCVT_D_S(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FEQ_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FLT_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FLE_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCLASS_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCVT_W_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCVT_WU_D(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCVT_D_W(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_FCVT_D_WU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv32c

func emu32_C_ILLEGAL(m *RV32, ins uint) {
	m.flag |= flagIllegal
}

func emu32_C_ADDI4SPN(m *RV32, ins uint) {
	uimm, rd := decodeCIW(ins)
	m.X[rd] = m.X[regSp] + uint32(uimm)
	m.PC += 2
}

func emu32_C_FLD(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_LW(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_FLW(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_FSD(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_SW(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_FSW(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_NOP(m *RV32, ins uint) {
	m.PC += 2
}

func emu32_C_ADDI(m *RV32, ins uint) {
	imm, rd := decodeCIa(ins)
	if rd != 0 {
		m.X[rd] = uint32(int(m.X[rd]) + imm)
	}
	m.PC += 2
}

func emu32_C_JAL(m *RV32, ins uint) {
	imm := decodeCJb(ins)
	m.X[regRa] = m.PC + 2
	m.PC = uint32(int(m.PC) + imm)
}

func emu32_C_LI(m *RV32, ins uint) {
	imm, rd := decodeCIa(ins)
	m.wrX(rd, uint32(imm))
	m.PC += 2
}

func emu32_C_ADDI16SP(m *RV32, ins uint) {
	imm := decodeCIb(ins)
	m.X[regSp] = uint32(int(m.X[regSp]) + imm)
	m.PC += 2
}

func emu32_C_LUI(m *RV32, ins uint) {
	imm, rd := decodeCIf(ins)
	if imm == 0 {
		m.flag |= flagIllegal
		return
	}
	if rd != 0 && rd != 2 {
		m.X[rd] = uint32(imm << 12)
	}
	m.PC += 2
}

func emu32_C_SRLI(m *RV32, ins uint) {
	shamt, rd := decodeCIc(ins)
	if shamt > 31 {
		m.flag |= flagIllegal
		return
	}
	m.X[rd] = m.X[rd] >> shamt
	m.PC += 2
}

func emu32_C_SRAI(m *RV32, ins uint) {
	shamt, rd := decodeCIc(ins)
	if shamt > 31 {
		m.flag |= flagIllegal
		return
	}
	m.X[rd] = uint32(int(m.X[rd]) >> shamt)
	m.PC += 2
}

func emu32_C_ANDI(m *RV32, ins uint) {
	imm, rd := decodeCIe(ins)
	m.X[rd] = uint32(int(m.X[rd]) & imm)
	m.PC += 2
}

func emu32_C_SUB(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_XOR(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_OR(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_AND(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_J(m *RV32, ins uint) {
	imm := decodeCJb(ins)
	m.PC = uint32(int(m.PC) + imm)
}

func emu32_C_BEQZ(m *RV32, ins uint) {
	imm, rs := decodeCB(ins)
	if m.X[rs] == 0 {
		m.PC = uint32(int(m.PC) + imm)
	} else {
		m.PC += 2
	}
}

func emu32_C_BNEZ(m *RV32, ins uint) {
	imm, rs := decodeCB(ins)
	if m.X[rs] != 0 {
		m.PC = uint32(int(m.PC) + imm)
	} else {
		m.PC += 2
	}
}

func emu32_C_SLLI(m *RV32, ins uint) {
	shamt, rd := decodeCId(ins)
	if shamt > 31 {
		m.flag |= flagIllegal
		return
	}
	if rd != 0 && shamt != 0 {
		m.X[rd] = m.X[rd] << shamt
	}
	m.PC += 2
}

func emu32_C_SLLI64(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_FLDSP(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_LWSP(m *RV32, ins uint) {
	uimm, rd := decodeCSSa(ins)
	if rd == 0 {
		m.flag |= flagIllegal
		return
	}
	adr := uint(m.X[regSp]) + uimm
	val, ex := m.Mem.Rd32(adr)
	m.checkMemory(adr, ex)
	m.X[rd] = val
	m.PC += 2
}

func emu32_C_FLWSP(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_JR(m *RV32, ins uint) {
	rs1 := decodeCJa(ins)
	if rs1 == 0 {
		m.flag |= flagIllegal
		return
	}
	m.PC = m.X[rs1]
}

func emu32_C_MV(m *RV32, ins uint) {
	rd, rs := decodeCR(ins)
	if rs != 0 {
		m.wrX(rd, m.X[rs])
	}
	m.PC += 2
}

func emu32_C_EBREAK(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_JALR(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_ADD(m *RV32, ins uint) {
	rd, rs := decodeCR(ins)
	m.wrX(rd, m.X[rd]+m.X[rs])
	m.PC += 2
}

func emu32_C_FSDSP(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_SWSP(m *RV32, ins uint) {
	uimm, rs2 := decodeCSSb(ins)
	adr := uint(m.X[regSp]) + uimm
	ex := m.Mem.Wr32(adr, m.X[rs2])
	m.checkMemory(adr, ex)
	m.PC += 2
}

func emu32_C_FSWSP(m *RV32, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------

// Reset the RV32 CPU.
func (m *RV32) Reset() {
	m.PC = 0
	m.insCount = 0
	m.lastPC = 0
}

// Run the RV32 CPU for a single instruction.
func (m *RV32) Run() error {

	// read the next instruction
	ins, ex := m.Mem.RdIns(uint(m.PC))
	if ex != 0 {
		return fmt.Errorf("memory exception %s", m.mx)
	}

	// lookup and emulate the instruction
	im := m.isa.lookup(ins)
	im.defn.emu32(m, ins)
	m.insCount++

	// check exception flags
	if m.flag != 0 {
		if m.flag&flagIllegal != 0 {
			return fmt.Errorf("illegal instruction at PC %08x", m.PC)
		}
		if m.flag&flagMemory != 0 {
			return fmt.Errorf("memory exception %s", m.mx)
		}
		if m.flag&flagExit != 0 {
			return fmt.Errorf("exit at PC %08x, status %08x (%d instructions)", m.PC, m.X[1], m.insCount)
		}
		if m.flag&flagTodo != 0 {
			return fmt.Errorf("unimplemented instruction at PC %08x", m.PC)
		}
		panic("unknown flag")
	}

	// stuck PC detection
	if m.PC == m.lastPC {
		return fmt.Errorf("PC is stuck at %08x (%d instructions)", m.PC, m.insCount)
	} else {
		m.lastPC = m.PC
	}

	return nil
}

//-----------------------------------------------------------------------------
