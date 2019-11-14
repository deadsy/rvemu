//-----------------------------------------------------------------------------
/*

RISC-V 32-bit Emulator

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
	"math"
)

//-----------------------------------------------------------------------------
// default emulation

func emu32_None(m *RV32, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------
// rv32i

func emu32_LUI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_AUIPC(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_JAL(m *RV32, ins uint) {
	imm, rd := decodeJ(ins)
	m.wrX(rd, m.PC+4)
	m.PC = uint32(int(m.PC) + int(imm))
}

func emu32_JALR(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_BEQ(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_BNE(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_BLT(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_BGE(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_BLTU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_BGEU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_LB(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_LH(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_LW(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_LBU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_LHU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SB(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SH(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SW(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_ADDI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SLTI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SLTIU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_XORI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_ORI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_ANDI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SLLI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SRLI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SRAI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_ADD(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SUB(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_SLL(m *RV32, ins uint) {
	m.flag |= flagTodo
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
	m.flag |= flagTodo
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
	m.flag |= flagTodo
}

func emu32_DIVU(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_REM(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_REMU(m *RV32, ins uint) {
	m.flag |= flagTodo
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
	m.F[rd] = math.Float32frombits(m.X[rs1])
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
	m.flag |= flagTodo
}

func emu32_C_ADDI4SPN(m *RV32, ins uint) {
	m.flag |= flagTodo
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
	m.flag |= flagTodo
}

func emu32_C_ADDI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_JAL(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_LI(m *RV32, ins uint) {
	imm, rd := decodeCIa(ins)
	m.wrX(rd, uint32(imm))
	m.PC += 2
}

func emu32_C_ADDI16SP(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_LUI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_SRLI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_SRAI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_ANDI(m *RV32, ins uint) {
	m.flag |= flagTodo
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
	m.flag |= flagTodo
}

func emu32_C_BEQZ(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_BNEZ(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_SLLI(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_SLLI64(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_FLDSP(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_LWSP(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_FLWSP(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_JR(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_MV(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_EBREAK(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_JALR(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_ADD(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_FSDSP(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_SWSP(m *RV32, ins uint) {
	m.flag |= flagTodo
}

func emu32_C_FSWSP(m *RV32, ins uint) {
	m.flag |= flagTodo
}

//-----------------------------------------------------------------------------

// Reset the RV32 CPU.
func (m *RV32) Reset() {
}

// Run the RV32 CPU for a single instruction.
func (m *RV32) Run() error {

	// normal instructions
	ins, _ := m.Mem.RdIns(uint(m.PC))
	im := m.isa.lookup(ins)
	im.defn.emu32(m, ins)

	if m.flag != 0 {
		if m.flag&flagIllegal != 0 {
			return fmt.Errorf("illegal instruction at %08x", m.PC)
		}
		if m.flag&flagExit != 0 {
			return fmt.Errorf("exit at %08x, status %08x", m.PC, m.X[1])
		}
		if m.flag&flagTodo != 0 {
			return fmt.Errorf("unimplemented instruction at %08x", m.PC)
		}
	}

	return nil
}

//-----------------------------------------------------------------------------
