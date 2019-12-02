//-----------------------------------------------------------------------------
/*

RISC-V 32-bit CPU Emulation

*/
//-----------------------------------------------------------------------------

package rv

import (
	"math"

	"github.com/deadsy/riscv/csr"
	"github.com/deadsy/riscv/mem"
)

//-----------------------------------------------------------------------------

func emu32_Illegal(m *RV32, ins uint) {
	m.ex.N = ExIllegal
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
	x1 := int32(m.X[rs1])
	x2 := int32(m.X[rs2])
	if x1 < x2 {
		m.PC = uint32(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu32_BGE(m *RV32, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	x1 := int32(m.X[rs1])
	x2 := int32(m.X[rs2])
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
	m.wrX(rd, uint32(int8(val)))
	m.PC += 4
}

func emu32_LH(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd16(adr)
	m.checkMemory(adr, ex)
	m.wrX(rd, uint32(int16(val)))
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
	imm, rs1, rd := decodeIa(ins)
	var result uint32
	if int32(m.X[rs1]) < int32(imm) {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu32_SLTIU(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	var result uint32
	if m.X[rs1] < uint32(imm) {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
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
		m.ex.N = ExIllegal
		return
	}
	m.wrX(rd, m.X[rs1]<<shamt)
	m.PC += 4
}

func emu32_SRLI(m *RV32, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if shamt > 31 {
		m.ex.N = ExIllegal
		return
	}
	m.wrX(rd, m.X[rs1]>>shamt)
	m.PC += 4
}

func emu32_SRAI(m *RV32, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if shamt > 31 {
		m.ex.N = ExIllegal
		return
	}
	m.wrX(rd, uint32(int32(m.X[rs1])>>shamt))
	m.PC += 4
}

func emu32_ADD(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]+m.X[rs2])
	m.PC += 4
}

func emu32_SUB(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]-m.X[rs2])
	m.PC += 4
}

func emu32_SLL(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]<<(m.X[rs2]&31))
	m.PC += 4
}

func emu32_SLT(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var result uint32
	if int32(m.X[rs1]) < int32(m.X[rs2]) {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu32_SLTU(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var result uint32
	if m.X[rs1] < m.X[rs2] {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu32_XOR(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]^m.X[rs2])
	m.PC += 4
}

func emu32_SRL(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	shamt := m.X[rs2] & 31
	m.wrX(rd, m.X[rs1]>>shamt)
	m.PC += 4
}

func emu32_SRA(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	shamt := m.X[rs2] & 31
	m.wrX(rd, uint32(int32(m.X[rs1])>>shamt))
	m.PC += 4
}

func emu32_OR(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]|m.X[rs2])
	m.PC += 4
}

func emu32_AND(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]&m.X[rs2])
	m.PC += 4
}

func emu32_FENCE(m *RV32, ins uint) {
	// no-op for a sw emulator
	m.PC += 4
}

func emu32_FENCE_I(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_ECALL(m *RV32, ins uint) {
	if m.ecall == nil {
		m.ex.N = ExEcall
		return
	}
	m.ecall.Call32(m)
	m.PC += 4
}

func emu32_EBREAK(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_CSRRW(m *RV32, ins uint) {
	csr, rs1, rd := decodeIb(ins)
	var t uint32
	if rd != 0 {
		t = m.rdCSR(csr)
	}
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
	csr, rs1, rd := decodeIb(ins)
	t := m.rdCSR(csr)
	if rs1 != 0 {
		m.wrCSR(csr, t & ^m.X[rs1])
	}
	m.wrX(rd, t)
	m.PC += 4
}

func emu32_CSRRWI(m *RV32, ins uint) {
	csr, zimm, rd := decodeIb(ins)
	if rd != 0 {
		m.X[rd] = m.rdCSR(csr)
	}
	m.wrCSR(csr, uint32(zimm))
	m.PC += 4
}

func emu32_CSRRSI(m *RV32, ins uint) {
	csr, zimm, rd := decodeIb(ins)
	t := m.rdCSR(csr)
	m.wrCSR(csr, t|uint32(zimm))
	m.wrX(rd, t)
	m.PC += 4
}

func emu32_CSRRCI(m *RV32, ins uint) {
	csr, zimm, rd := decodeIb(ins)
	t := m.rdCSR(csr)
	m.wrCSR(csr, t & ^uint32(zimm))
	m.wrX(rd, t)
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv32i privileged

func emu32_URET(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_SRET(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_MRET(m *RV32, ins uint) {
	pc, ex := m.CSR.MRET()
	m.checkCSR(csr.MSTATUS, ex)
	m.PC = uint32(pc)
}

func emu32_WFI(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_SFENCE_VMA(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_HFENCE_BVMA(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_HFENCE_GVMA(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

//-----------------------------------------------------------------------------
// rv32m

func emu32_MUL(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.X[rs1]*m.X[rs2])
	m.PC += 4
}

func emu32_MULH(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	a := int64(int32(m.X[rs1]))
	b := int64(int32(m.X[rs2]))
	c := (a * b) >> 32
	m.wrX(rd, uint32(c))
	m.PC += 4
}

func emu32_MULHSU(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	a := int64(int32(m.X[rs1]))
	b := int64(m.X[rs2])
	c := (a * b) >> 32
	m.wrX(rd, uint32(c))
	m.PC += 4
}

func emu32_MULHU(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	a := uint64(m.X[rs1])
	b := uint64(m.X[rs2])
	c := (a * b) >> 32
	m.wrX(rd, uint32(c))
	m.PC += 4
}

func emu32_DIV(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	result := int32(-1)
	a := int32(m.X[rs1])
	b := int32(m.X[rs2])
	if b != 0 {
		result = a / b
	}
	m.wrX(rd, uint32(result))
	m.PC += 4
}

func emu32_DIVU(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	result := uint32((1 << 32) - 1)
	if m.X[rs2] != 0 {
		result = m.X[rs1] / m.X[rs2]
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu32_REM(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	result := int32(m.X[rs1])
	b := int32(m.X[rs2])
	if b != 0 {
		result %= b
	}
	m.wrX(rd, uint32(result))
	m.PC += 4
}

func emu32_REMU(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
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
	m.ex.N = ExTodo
}

func emu32_SC_W(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_AMOSWAP_W(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_AMOADD_W(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_AMOXOR_W(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_AMOAND_W(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_AMOOR_W(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_AMOMIN_W(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_AMOMAX_W(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_AMOMINU_W(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_AMOMAXU_W(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

//-----------------------------------------------------------------------------
// rv32f

func emu32_FLW(m *RV32, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.X[rs1]) + imm)
	val, ex := m.Mem.Rd32(adr)
	m.checkMemory(adr, ex)
	m.F[rd] = uint64(val)
	m.PC += 4
}

func emu32_FSW(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FMADD_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FMSUB_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FNMSUB_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FNMADD_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FADD_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FSUB_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FMUL_S(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	f1 := math.Float32frombits(uint32(m.F[rs1]))
	f2 := math.Float32frombits(uint32(m.F[rs2]))
	m.F[rd] = uint64(math.Float32bits(f1 * f2))
	m.PC += 4
}

func emu32_FDIV_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FSQRT_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FSGNJ_S(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	sign := m.F[rs2] & mask31
	m.F[rd] = sign | (m.F[rs1] & mask30to0)
	m.PC += 4
}

func emu32_FSGNJN_S(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	sign := ^m.F[rs2] & mask31
	m.F[rd] = sign | (m.F[rs1] & mask30to0)
	m.PC += 4
}

func emu32_FSGNJX_S(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	sign := (m.F[rs1] ^ m.F[rs2]) & mask31
	m.F[rd] = sign | (m.F[rs1] & mask30to0)
	m.PC += 4
}

func emu32_FMIN_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FMAX_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FCVT_W_S(m *RV32, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	f := math.Float32frombits(uint32(m.F[rs1]))
	x, err := convertF32toI32(f, rm, m.CSR)
	if err != nil {
		m.ex.N = ExIllegal
		return
	}
	m.wrX(rd, uint32(x))
	m.PC += 4
}

func emu32_FCVT_WU_S(m *RV32, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	f := math.Float32frombits(uint32(m.F[rs1]))
	x, err := convertF32toU32(f, rm, m.CSR)
	if err != nil {
		m.ex.N = ExIllegal
		return
	}
	m.wrX(rd, x)
	m.PC += 4
}

func emu32_FMV_X_W(m *RV32, ins uint) {
	_, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint32(int32(m.F[rs1])))
	m.PC += 4
}

func emu32_FEQ_S(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var result uint32
	if uint32(m.F[rs1]) == uint32(m.F[rs2]) {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu32_FLT_S(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var result uint32
	if uint32(m.F[rs1]) < uint32(m.F[rs2]) {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu32_FLE_S(m *RV32, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var result uint32
	if uint32(m.F[rs1]) <= uint32(m.F[rs2]) {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu32_FCLASS_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FCVT_S_W(m *RV32, ins uint) {
	_, rs1, _, rd := decodeR(ins)
	f1 := float32(int32(m.X[rs1]))
	m.F[rd] = uint64(math.Float32bits(f1))
	m.PC += 4
}

func emu32_FCVT_S_WU(m *RV32, ins uint) {
	_, rs1, _, rd := decodeR(ins)
	f1 := float32(m.X[rs1])
	m.F[rd] = uint64(math.Float32bits(f1))
	m.PC += 4
}

func emu32_FMV_W_X(m *RV32, ins uint) {
	_, rs1, _, rd := decodeR(ins)
	m.F[rd] = uint64(m.X[rs1])
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv32d

func emu32_FLD(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FSD(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FMADD_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FMSUB_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FNMSUB_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FNMADD_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FADD_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FSUB_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FMUL_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FDIV_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FSQRT_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FSGNJ_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FSGNJN_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FSGNJX_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FMIN_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FMAX_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FCVT_S_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FCVT_D_S(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FEQ_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FLT_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FLE_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FCLASS_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FCVT_W_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FCVT_WU_D(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FCVT_D_W(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_FCVT_D_WU(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

//-----------------------------------------------------------------------------
// rv32c

func emu32_C_ILLEGAL(m *RV32, ins uint) {
	m.ex.N = ExIllegal
}

func emu32_C_ADDI4SPN(m *RV32, ins uint) {
	uimm, rd := decodeCIW(ins)
	m.X[rd] = m.X[RegSp] + uint32(uimm)
	m.PC += 2
}

func emu32_C_FLD(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_C_LW(m *RV32, ins uint) {
	uimm, rs1, rd := decodeCS(ins)
	adr := uint(m.X[rs1]) + uimm
	val, ex := m.Mem.Rd32(adr)
	m.checkMemory(adr, ex)
	m.X[rd] = val
	m.PC += 2
}

func emu32_C_FLW(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_C_FSD(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_C_SW(m *RV32, ins uint) {
	uimm, rs1, rs2 := decodeCS(ins)
	adr := uint(m.X[rs1]) + uimm
	ex := m.Mem.Wr32(adr, m.X[rs2])
	m.checkMemory(adr, ex)
	m.PC += 2
}

func emu32_C_FSW(m *RV32, ins uint) {
	m.ex.N = ExTodo
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
	imm := decodeCJ(ins)
	m.X[RegRa] = m.PC + 2
	m.PC = uint32(int(m.PC) + imm)
}

func emu32_C_LI(m *RV32, ins uint) {
	imm, rd := decodeCIa(ins)
	m.wrX(rd, uint32(imm))
	m.PC += 2
}

func emu32_C_ADDI16SP(m *RV32, ins uint) {
	imm := decodeCIb(ins)
	m.X[RegSp] = uint32(int(m.X[RegSp]) + imm)
	m.PC += 2
}

func emu32_C_LUI(m *RV32, ins uint) {
	imm, rd := decodeCIf(ins)
	if imm == 0 {
		m.ex.N = ExIllegal
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
		m.ex.N = ExIllegal
		return
	}
	m.X[rd] = m.X[rd] >> shamt
	m.PC += 2
}

func emu32_C_SRAI(m *RV32, ins uint) {
	shamt, rd := decodeCIc(ins)
	if shamt > 31 {
		m.ex.N = ExIllegal
		return
	}
	m.X[rd] = uint32(int32(m.X[rd]) >> shamt)
	m.PC += 2
}

func emu32_C_ANDI(m *RV32, ins uint) {
	imm, rd := decodeCIe(ins)
	m.X[rd] = uint32(int(m.X[rd]) & imm)
	m.PC += 2
}

func emu32_C_SUB(m *RV32, ins uint) {
	rd, rs := decodeCRa(ins)
	m.X[rd] -= m.X[rs]
	m.PC += 2
}

func emu32_C_XOR(m *RV32, ins uint) {
	rd, rs := decodeCRa(ins)
	m.X[rd] ^= m.X[rs]
	m.PC += 2
}

func emu32_C_OR(m *RV32, ins uint) {
	rd, rs := decodeCRa(ins)
	m.X[rd] |= m.X[rs]
	m.PC += 2
}

func emu32_C_AND(m *RV32, ins uint) {
	rd, rs := decodeCRa(ins)
	m.X[rd] &= m.X[rs]
	m.PC += 2
}

func emu32_C_J(m *RV32, ins uint) {
	imm := decodeCJ(ins)
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
		m.ex.N = ExIllegal
		return
	}
	if rd != 0 && shamt != 0 {
		m.X[rd] = m.X[rd] << shamt
	}
	m.PC += 2
}

func emu32_C_SLLI64(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_C_FLDSP(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_C_LWSP(m *RV32, ins uint) {
	uimm, rd := decodeCSSa(ins)
	if rd == 0 {
		m.ex.N = ExIllegal
		return
	}
	adr := uint(m.X[RegSp]) + uimm
	val, ex := m.Mem.Rd32(adr)
	m.checkMemory(adr, ex)
	m.X[rd] = val
	m.PC += 2
}

func emu32_C_FLWSP(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_C_JR(m *RV32, ins uint) {
	rs1, _ := decodeCR(ins)
	if rs1 == 0 {
		m.ex.N = ExIllegal
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
	m.ex.N = ExTodo
}

func emu32_C_JALR(m *RV32, ins uint) {
	rs1, _ := decodeCR(ins)
	if rs1 == 0 {
		m.ex.N = ExIllegal
		return
	}
	t := m.PC + 2
	m.PC = m.X[rs1]
	m.X[RegRa] = t
}

func emu32_C_ADD(m *RV32, ins uint) {
	rd, rs := decodeCR(ins)
	m.wrX(rd, m.X[rd]+m.X[rs])
	m.PC += 2
}

func emu32_C_FSDSP(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

func emu32_C_SWSP(m *RV32, ins uint) {
	uimm, rs2 := decodeCSSb(ins)
	adr := uint(m.X[RegSp]) + uimm
	ex := m.Mem.Wr32(adr, m.X[rs2])
	m.checkMemory(adr, ex)
	m.PC += 2
}

func emu32_C_FSWSP(m *RV32, ins uint) {
	m.ex.N = ExTodo
}

//-----------------------------------------------------------------------------
// private methods

// wrX sets a register value (no writes to zero).
func (m *RV32) wrX(i uint, val uint32) {
	if i != 0 {
		m.X[i] = val
	}
}

// checkMemory records a memory exception.
func (m *RV32) checkMemory(adr uint, ex mem.Exception) {
	if ex == 0 {
		return
	}
	m.ex.N = ExMemory
	m.ex.mem = memoryException{adr, ex}
}

// rdCSR reads a CSR.
func (m *RV32) rdCSR(reg uint) uint32 {
	val, ex := m.CSR.Rd(reg)
	m.checkCSR(reg, ex)
	return uint32(val)
}

// wrCSR writes a CSR.
func (m *RV32) wrCSR(reg uint, val uint32) {
	ex := m.CSR.Wr(reg, uint(val))
	m.checkCSR(reg, ex)
}

// checkCSR records a memory exception.
func (m *RV32) checkCSR(reg uint, ex csr.Exception) {
	if ex == 0 {
		return
	}
	m.ex.N = ExCSR
	m.ex.csr = csrException{reg, ex}
}

//-----------------------------------------------------------------------------

// RV32 is a 32-bit RISC-V CPU.
type RV32 struct {
	Mem      *mem.Memory // memory of the target system
	X        [32]uint32  // integer registers
	F        [32]uint64  // float registers
	PC       uint32      // program counter
	CSR      *csr.State  // CSR state
	insCount uint        // number of instructions run
	lastPC   uint32      // stuck PC detection
	ex       Exception   // emulation exceptions
	isa      *ISA        // ISA implemented for the CPU
	ecall    Ecall       // ecall interface
}

// NewRV32 returns a 32-bit RISC-V CPU.
func NewRV32(isa *ISA, mem *mem.Memory, ecall Ecall) *RV32 {
	m := RV32{
		Mem:   mem,
		CSR:   csr.NewState(32),
		isa:   isa,
		ecall: ecall,
	}
	m.Reset()
	return &m
}

// Run the RV32 CPU for a single instruction.
func (m *RV32) Run() error {

	// set the pc for the exception (if there is one)
	m.ex.pc = uint(m.PC)

	// read the next instruction
	ins, ex := m.Mem.RdIns(uint(m.PC))
	if ex != 0 {
		m.checkMemory(uint(m.PC), ex)
		return &m.ex
	}

	// lookup and emulate the instruction
	im := m.isa.lookup(ins)
	if im != nil {
		im.defn.emu32(m, ins)
		m.insCount++
	} else {
		m.ex.N = ExIllegal
	}

	// check exception flags
	if m.ex.N != 0 {
		return &m.ex
	}

	// stuck PC detection
	if m.PC == m.lastPC {
		m.ex.N = ExStuck
		return &m.ex
	}
	m.lastPC = m.PC

	return nil
}

//-----------------------------------------------------------------------------

// IntRegs returns a display string for the integer registers.
func (m *RV32) IntRegs() string {
	reg := make([]uint, 32)
	for i := range reg {
		reg[i] = uint(m.X[i])
	}
	return intRegString(reg, uint(m.PC), 32)
}

// FloatRegs returns a display string for the float registers.
func (m *RV32) FloatRegs() string {
	return floatRegString(m.F[:])
}

// Disassemble the instruction at the address.
func (m *RV32) Disassemble(adr uint) *Disassembly {
	return m.isa.Disassemble(m.Mem, adr)
}

//-----------------------------------------------------------------------------

// Reset the RV32 CPU.
func (m *RV32) Reset() {
	m.PC = uint32(m.Mem.Entry)
	m.X[RegSp] = uint32(uint(1<<32) - 16)
	m.insCount = 0
	m.lastPC = 0
	m.ex = Exception{alen: 32}
}

// Exit sets a status code and exits the emulation
func (m *RV32) Exit(status uint32) {
	m.X[RegA0] = status
	m.ex.N = ExExit
}

// SetBreak sets the break flag.
func (m *RV32) SetBreak() {
	m.ex.N = ExBreak
}

//-----------------------------------------------------------------------------
