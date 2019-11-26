//-----------------------------------------------------------------------------
/*

RISC-V 64-bit CPU Emulation

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
	"strings"

	"github.com/deadsy/riscv/mem"
)

//-----------------------------------------------------------------------------

func emu64_Illegal(m *RV64, ins uint) {
	m.flag |= flagIllegal
}

//-----------------------------------------------------------------------------
// rv32i

func emu64_LUI(m *RV64, ins uint) {
	imm, rd := decodeU(ins)
	m.wrX(rd, uint64(imm<<12))
	m.PC += 4
}

func emu64_AUIPC(m *RV64, ins uint) {
	imm, rd := decodeU(ins)
	m.wrX(rd, uint64(int(m.PC)+(imm<<12)))
	m.PC += 4
}

func emu64_JAL(m *RV64, ins uint) {
	imm, rd := decodeJ(ins)
	m.wrX(rd, m.PC+4)
	m.PC = uint64(int(m.PC) + int(imm))
}

func emu64_JALR(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.PC+4)
	m.PC = uint64((int(m.X[rs1]) + imm) & ^1)
}

func emu64_BEQ(m *RV64, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.X[rs1] == m.X[rs2] {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu64_BNE(m *RV64, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.X[rs1] != m.X[rs2] {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu64_BLT(m *RV64, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	x1 := int64(m.X[rs1])
	x2 := int64(m.X[rs2])
	if x1 < x2 {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu64_BGE(m *RV64, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	x1 := int64(m.X[rs1])
	x2 := int64(m.X[rs2])
	if x1 >= x2 {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu64_BLTU(m *RV64, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.X[rs1] < m.X[rs2] {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu64_BGEU(m *RV64, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.X[rs1] >= m.X[rs2] {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu64_LB(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd8(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, uint64(int8(val)))
	m.PC += 4
}

func emu64_LH(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd16(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, uint64(int16(val)))
	m.PC += 4
}

func emu64_LW(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd32(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, uint64(int(val)))
	m.PC += 4
}

func emu64_LBU(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd8(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, uint64(val))
	m.PC += 4
}

func emu64_LHU(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd16(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, uint64(val))
	m.PC += 4
}

func emu64_SB(m *RV64, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.X[rs1]) + imm)
	ex := m.Mem.Wr8(adr, uint8(m.X[rs2]))
	m.checkMemory(adr, ex)
	m.PC += 4
}

func emu64_SH(m *RV64, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.X[rs1]) + imm)
	ex := m.Mem.Wr16(adr, uint16(m.X[rs2]))
	m.checkMemory(adr, ex)
	m.PC += 4
}

func emu64_SW(m *RV64, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.X[rs1]) + imm)
	ex := m.Mem.Wr32(adr, uint32(m.X[rs2]))
	m.checkMemory(adr, ex)
	m.PC += 4
}

func emu64_ADDI(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, uint64(int(m.X[rs1])+imm))
	m.PC += 4
}

func emu64_SLTI(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_SLTIU(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_XORI(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.X[rs1]^uint64(imm))
	m.PC += 4
}

func emu64_ORI(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.X[rs1]|uint64(imm))
	m.PC += 4
}

func emu64_ANDI(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.X[rs1]&uint64(imm))
	m.PC += 4
}

func emu64_ADD(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]+m.X[rs2])
	m.PC += 4
}

func emu64_SUB(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]-m.X[rs2])
	m.PC += 4
}

func emu64_SLL(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]<<(m.X[rs2]&63))
	m.PC += 4
}

func emu64_SLT(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_SLTU(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_XOR(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]^m.X[rs2])
	m.PC += 4
}

func emu64_SRL(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_SRA(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_OR(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]|m.X[rs2])
	m.PC += 4
}

func emu64_AND(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]&m.X[rs2])
	m.PC += 4
}

func emu64_FENCE(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FENCE_I(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_ECALL(m *RV64, ins uint) {
	s := scLookup(int(m.X[regA7]))
	if s == nil {
		m.flag |= flagSyscall
		return
	}
	s.sc64(m, s)
	m.PC += 4
}

func emu64_EBREAK(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_CSRRW(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_CSRRS(m *RV64, ins uint) {
	csr, rs1, rd := decodeIb(ins)
	t := m.rdCSR(csr)
	if rs1 != 0 {
		m.wrCSR(csr, t|m.X[rs1])
	}
	m.wrX(rd, t)
	m.PC += 4
}

func emu64_CSRRC(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_CSRRWI(m *RV64, ins uint) {
	csr, zimm, rd := decodeIb(ins)
	if rd != 0 {
		m.X[rd] = m.rdCSR(csr)
	}
	if zimm != 0 {
		m.wrCSR(csr, uint64(zimm))
	}
	m.PC += 4
}

func emu64_CSRRSI(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_CSRRCI(m *RV64, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv32i privileged

func emu64_URET(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_SRET(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_MRET(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_WFI(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_SFENCE_VMA(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_HFENCE_BVMA(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_HFENCE_GVMA(m *RV64, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv32m

func emu64_MUL(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]*m.X[rs2])
	m.PC += 4
}

func emu64_MULH(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_MULHSU(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_MULHU(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_DIV(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	result := -1
	if m.X[rs2] != 0 {
		result = int(m.X[rs1]) / int(m.X[rs2])
	}
	m.wrX(rd, uint64(result))
	m.PC += 4
}

func emu64_DIVU(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	result := uint64((1 << 64) - 1)
	if m.X[rs2] != 0 {
		result = m.X[rs1] / m.X[rs2]
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu64_REM(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	result := int(m.X[rs1])
	if m.X[rs2] != 0 {
		result %= int(m.X[rs2])
	}
	m.wrX(rd, uint64(result))
	m.PC += 4
}

func emu64_REMU(m *RV64, ins uint) {
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

func emu64_LR_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_SC_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOSWAP_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOADD_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOXOR_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOAND_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOOR_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOMIN_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOMAX_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOMINU_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOMAXU_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv32f

func emu64_FLW(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSW(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMADD_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMSUB_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FNMSUB_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FNMADD_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FADD_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSUB_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMUL_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FDIV_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSQRT_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSGNJ_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSGNJN_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSGNJX_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMIN_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMAX_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_W_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_WU_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMV_X_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FEQ_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FLT_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FLE_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCLASS_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_S_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_S_WU(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMV_W_X(m *RV64, ins uint) {
	_, rs1, rd := decodeR(ins)
	m.F[rd] = u32Upper | (m.X[rs1] & u32Lower)
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv32d

func emu64_FLD(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSD(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMADD_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMSUB_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FNMSUB_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FNMADD_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FADD_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSUB_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMUL_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FDIV_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSQRT_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSGNJ_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSGNJN_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FSGNJX_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMIN_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMAX_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_S_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_D_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FEQ_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FLT_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FLE_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCLASS_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_W_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_WU_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_D_W(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_D_WU(m *RV64, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv32c

func emu64_C_ILLEGAL(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_C_ADDI4SPN(m *RV64, ins uint) {
	uimm, rd := decodeCIW(ins)
	m.X[rd] = m.X[regSp] + uint64(uimm)
	m.PC += 2
}

func emu64_C_FLD(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_C_LW(m *RV64, ins uint) {
	uimm, rs1, rd := decodeCS(ins)
	adr := uint(m.X[rs1]) + uimm
	val, ex := m.Mem.Rd32(adr)
	m.checkMemory(adr, ex)
	m.X[rd] = uint64(int(val))
	m.PC += 2
}

func emu64_C_FSD(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_C_SW(m *RV64, ins uint) {
	uimm, rs1, rs2 := decodeCS(ins)
	adr := uint(m.X[rs1]) + uimm
	ex := m.Mem.Wr32(adr, uint32(m.X[rs2]))
	m.checkMemory(adr, ex)
	m.PC += 2
}

func emu64_C_FSW(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_C_NOP(m *RV64, ins uint) {
	m.PC += 2
}

func emu64_C_ADDI(m *RV64, ins uint) {
	imm, rd := decodeCIa(ins)
	if rd != 0 {
		m.X[rd] = uint64(int(m.X[rd]) + imm)
	}
	m.PC += 2
}

func emu64_C_LI(m *RV64, ins uint) {
	imm, rd := decodeCIa(ins)
	m.wrX(rd, uint64(imm))
	m.PC += 2
}

func emu64_C_ADDI16SP(m *RV64, ins uint) {
	imm := decodeCIb(ins)
	m.X[regSp] = uint64(int(m.X[regSp]) + imm)
	m.PC += 2
}

func emu64_C_LUI(m *RV64, ins uint) {
	imm, rd := decodeCIf(ins)
	if imm == 0 {
		m.flag |= flagIllegal
		return
	}
	if rd != 0 && rd != 2 {
		m.X[rd] = uint64(imm << 12)
	}
	m.PC += 2
}

func emu64_C_SRLI(m *RV64, ins uint) {
	shamt, rd := decodeCIc(ins)
	m.X[rd] = m.X[rd] << shamt
	m.PC += 2
}

func emu64_C_SRAI(m *RV64, ins uint) {
	shamt, rd := decodeCIc(ins)
	m.X[rd] = uint64(int(m.X[rd]) >> shamt)
	m.PC += 2
}

func emu64_C_ANDI(m *RV64, ins uint) {
	imm, rd := decodeCIe(ins)
	m.X[rd] = uint64(int(m.X[rd]) & imm)
	m.PC += 2
}

func emu64_C_SUB(m *RV64, ins uint) {
	rd, rs := decodeCRa(ins)
	m.X[rd] -= m.X[rs]
	m.PC += 2
}

func emu64_C_XOR(m *RV64, ins uint) {
	rd, rs := decodeCRa(ins)
	m.X[rd] ^= m.X[rs]
	m.PC += 2
}

func emu64_C_OR(m *RV64, ins uint) {
	rd, rs := decodeCRa(ins)
	m.X[rd] |= m.X[rs]
	m.PC += 2
}

func emu64_C_AND(m *RV64, ins uint) {
	rd, rs := decodeCRa(ins)
	m.X[rd] &= m.X[rs]
	m.PC += 2
}

func emu64_C_J(m *RV64, ins uint) {
	imm := decodeCJ(ins)
	m.PC = uint64(int(m.PC) + imm)
}

func emu64_C_BEQZ(m *RV64, ins uint) {
	imm, rs := decodeCB(ins)
	if m.X[rs] == 0 {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 2
	}
}

func emu64_C_BNEZ(m *RV64, ins uint) {
	imm, rs := decodeCB(ins)
	if m.X[rs] != 0 {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 2
	}
}

func emu64_C_SLLI(m *RV64, ins uint) {
	shamt, rd := decodeCId(ins)
	if rd != 0 && shamt != 0 {
		m.X[rd] = m.X[rd] << shamt
	}
	m.PC += 2
}

func emu64_C_SLLI64(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_C_FLDSP(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_C_LWSP(m *RV64, ins uint) {
	uimm, rd := decodeCSSa(ins)
	if rd == 0 {
		m.flag |= flagIllegal
		return
	}
	adr := uint(m.X[regSp]) + uimm
	val, ex := m.Mem.Rd32(adr)
	m.checkMemory(adr, ex)
	m.X[rd] = uint64(int(val))
	m.PC += 2
}

func emu64_C_JR(m *RV64, ins uint) {
	rs1, _ := decodeCR(ins)
	if rs1 == 0 {
		m.flag |= flagIllegal
		return
	}
	m.PC = m.X[rs1]
}

func emu64_C_MV(m *RV64, ins uint) {
	rd, rs := decodeCR(ins)
	if rs != 0 {
		m.wrX(rd, m.X[rs])
	}
	m.PC += 2
}

func emu64_C_EBREAK(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_C_JALR(m *RV64, ins uint) {
	rs1, _ := decodeCR(ins)
	if rs1 == 0 {
		m.flag |= flagIllegal
		return
	}
	t := m.PC + 2
	m.PC = m.X[rs1]
	m.X[regRa] = t
}

func emu64_C_ADD(m *RV64, ins uint) {
	rd, rs := decodeCR(ins)
	m.wrX(rd, m.X[rd]+m.X[rs])
	m.PC += 2
}

func emu64_C_FSDSP(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_C_SWSP(m *RV64, ins uint) {
	uimm, rs2 := decodeCSSb(ins)
	adr := uint(m.X[regSp]) + uimm
	ex := m.Mem.Wr32(adr, uint32(m.X[rs2]))
	m.checkMemory(adr, ex)
	m.PC += 2
}

//-----------------------------------------------------------------------------
// rv64i

func emu64_LWU(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_LD(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd64(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, val)
	m.PC += 4
}

func emu64_SD(m *RV64, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.X[rs1]) + imm)
	ex := m.Mem.Wr64(adr, m.X[rs2])
	m.checkMemory(adr, ex)
	m.PC += 4
}

func emu64_SLLI(m *RV64, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	m.wrX(rd, m.X[rs1]<<shamt)
	m.PC += 4
}

func emu64_SRLI(m *RV64, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	m.wrX(rd, m.X[rs1]>>shamt)
	m.PC += 4
}

func emu64_SRAI(m *RV64, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	m.wrX(rd, uint64(int(m.X[rs1])>>shamt))
	m.PC += 4
}

func emu64_ADDIW(m *RV64, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, uint64(int32(int(m.X[rs1])+imm)))
	m.PC += 4
}

func emu64_SLLIW(m *RV64, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if shamt&32 != 0 {
		m.flag |= flagIllegal
		return
	}
	m.wrX(rd, uint64(int32(m.X[rs1])<<shamt))
	m.PC += 4
}

func emu64_SRLIW(m *RV64, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if shamt&32 != 0 {
		m.flag |= flagIllegal
		return
	}
	m.wrX(rd, uint64(int32(uint32(m.X[rs1])>>shamt)))
	m.PC += 4
}

func emu64_SRAIW(m *RV64, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if shamt&32 != 0 {
		m.flag |= flagIllegal
		return
	}
	m.wrX(rd, uint64(int32(m.X[rs1])>>shamt))
	m.PC += 4
}

func emu64_ADDW(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, uint64(int32(m.X[rs1]+m.X[rs2])))
	m.PC += 4
}

func emu64_SUBW(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, uint64(int32(m.X[rs1]-m.X[rs2])))
	m.PC += 4
}

func emu64_SLLW(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, uint64(int32(m.X[rs1]<<m.X[rs2]&31)))
	m.PC += 4
}

func emu64_SRLW(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_SRAW(m *RV64, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv64m

func emu64_MULW(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_DIVW(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_DIVUW(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, uint64(int32(uint32(m.X[rs1])/uint32(m.X[rs2]))))
	m.PC += 4
}

func emu64_REMW(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_REMUW(m *RV64, ins uint) {
	rs2, rs1, rd := decodeR(ins)
	m.wrX(rd, uint64(int32(uint32(m.X[rs1])%uint32(m.X[rs2]))))
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv64a

func emu64_LR_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_SC_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOSWAP_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOADD_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOXOR_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOAND_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOOR_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOMIN_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOMAX_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOMINU_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_AMOMAXU_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv64f

func emu64_FCVT_L_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_LU_S(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_S_L(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_S_LU(m *RV64, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv64d

func emu64_FCVT_L_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_LU_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMV_X_D(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_D_L(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FCVT_D_LU(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_FMV_D_X(m *RV64, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv64c

func emu64_C_ADDIW(m *RV64, ins uint) {
	imm, rd := decodeCIa(ins)
	if rd != 0 {
		m.X[rd] = uint64(int32(int(m.X[rd]) + imm))
	} else {
		m.flag |= flagIllegal
	}
	m.PC += 2
}

func emu64_C_LDSP(m *RV64, ins uint) {
	uimm, rd := decodeCIg(ins)
	adr := uint(m.X[regSp]) + uimm
	val, ex := m.Mem.Rd64(adr)
	m.checkMemory(adr, ex)
	if rd != 0 {
		m.X[rd] = val
	} else {
		m.flag |= flagIllegal
	}
	m.PC += 2
}

func emu64_C_SDSP(m *RV64, ins uint) {
	uimm, rs2 := decodeCSSc(ins)
	adr := uint(m.X[regSp]) + uimm
	ex := m.Mem.Wr64(adr, m.X[rs2])
	m.checkMemory(adr, ex)
	m.PC += 2
}

func emu64_C_LD(m *RV64, ins uint) {
	uimm, rs1, rd := decodeCSa(ins)
	adr := uint(m.X[rs1]) + uimm
	val, ex := m.Mem.Rd64(adr)
	m.checkMemory(adr, ex)
	m.X[rd] = val
	m.PC += 2
}

func emu64_C_SD(m *RV64, ins uint) {
	uimm, rs1, rs2 := decodeCSa(ins)
	adr := uint(m.X[rs1]) + uimm
	ex := m.Mem.Wr64(adr, m.X[rs2])
	m.checkMemory(adr, ex)
	m.PC += 2
}

func emu64_C_SUBW(m *RV64, ins uint) {
	m.flag |= flagTodo
}

func emu64_C_ADDW(m *RV64, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// private methods

// wrX sets a register value (no writes to zero).
func (m *RV64) wrX(i uint, val uint64) {
	if i != 0 {
		m.X[i] = val
	}
}

// checkMemory records a memory exception.
func (m *RV64) checkMemory(adr uint, ex mem.Exception) {
	if ex == 0 {
		return
	}
	m.flag |= flagMemory
	m.mx = memoryException{uint(m.PC), adr, ex}
}

//-----------------------------------------------------------------------------

// RV64 is a 64-bit RISC-V CPU.
type RV64 struct {
	Mem      *mem.Memory     // memory of the target system
	X        [32]uint64      // registers
	F        [32]uint64      // float registers
	PC       uint64          // program counter
	insCount uint            // number of instructions run
	lastPC   uint64          // stuck PC detection
	flag     emuFlags        // event flags
	mx       memoryException // memory exceptions
	isa      *ISA            // ISA implemented for the CPU
}

// NewRV64 returns a 64-bit RISC-V CPU.
func NewRV64(isa *ISA, mem *mem.Memory) *RV64 {
	return &RV64{
		Mem: mem,
		isa: isa,
	}
}

// IRegs returns a display string for the integer registers.
func (m *RV64) IRegs() string {
	nregs := 32
	s := make([]string, nregs+1)
	for i := 0; i < nregs; i++ {
		x := fmt.Sprintf("x%d", i)
		r := "0"
		if m.X[i] != 0 {
			r = fmt.Sprintf("%016x", m.X[i])
		}
		s[i] = fmt.Sprintf("%-4s %-4s %s", x, abiXName[i], r)
	}
	s[nregs] = fmt.Sprintf("%-9s %016x", "pc", m.PC)
	return strings.Join(s, "\n")
}

// Exit sets a status code and exits the emulation
func (m *RV64) Exit(status uint64) {
	m.X[1] = status
	m.flag |= flagExit
}

// Disassemble the instruction at the address.
func (m *RV64) Disassemble(adr uint) *Disassembly {
	return m.isa.Disassemble(m.Mem, adr)
}

// Reset the RV64 CPU.
func (m *RV64) Reset() {
	m.PC = uint64(m.Mem.Entry)
	m.X[regSp] = uint64(uint(1<<32) - 16)
	m.insCount = 0
	m.lastPC = 0
	m.flag = 0
}

// Run the RV64 CPU for a single instruction.
func (m *RV64) Run() error {

	// read the next instruction
	ins, ex := m.Mem.RdIns(uint(m.PC))
	if ex != 0 {
		m.checkMemory(uint(m.PC), ex)
		return fmt.Errorf("memory exception %s", m.mx)
	}

	// lookup and emulate the instruction
	im := m.isa.lookup(ins)
	if im != nil {
		im.defn.emu64(m, ins)
		m.insCount++
	} else {
		m.flag |= flagIllegal
	}

	// check exception flags
	if m.flag != 0 {
		if m.flag&flagIllegal != 0 {
			return fmt.Errorf("illegal instruction at PC %016x", m.PC)
		}
		if m.flag&flagMemory != 0 {
			return fmt.Errorf("memory exception %s", m.mx)
		}
		if m.flag&flagExit != 0 {
			return fmt.Errorf("exit at PC %016x, status %016x (%d instructions)", m.PC, m.X[1], m.insCount)
		}
		if m.flag&flagSyscall != 0 {
			return fmt.Errorf("unrecognized system call at PC %016x, %d", m.PC, m.X[regA7])
		}
		if m.flag&flagBreak != 0 {
			m.flag &= ^flagBreak
			return fmt.Errorf("breakpoint at PC %016x", m.PC)
		}
		if m.flag&flagTodo != 0 {
			return fmt.Errorf("unimplemented instruction at PC %016x", m.PC)
		}
		panic("unknown flag")
	}

	// stuck PC detection
	if m.PC == m.lastPC {
		return fmt.Errorf("PC is stuck at %016x (%d instructions)", m.PC, m.insCount)
	}
	m.lastPC = m.PC

	return nil
}

//-----------------------------------------------------------------------------
