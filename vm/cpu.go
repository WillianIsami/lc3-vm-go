package vm

import (
	"log"
	"math"
)

// Registers
const (
	R_R0 = iota
	R_R1
	R_R2
	R_R3
	R_R4
	R_R5
	R_R6
	R_R7
	R_PC    // Program Counter
	R_COND  // Condition (flag of status)
	R_COUNT // Register Counter
)

// Opcodes
const (
	OP_BR   uint16 = iota //  branch
	OP_ADD                //  add
	OP_SUB                //  sub
	OP_LD                 //  load
	OP_ST                 //  store
	OP_JSR                //  jump register
	OP_AND                //  bitwise and
	OP_LDR                //  load register
	OP_STR                //  store register
	OP_RTI                //  unused
	OP_NOT                //  bitwise not
	OP_LDI                //  load indirect
	OP_STI                //  store indirect
	OP_JMP                //  jump
	OP_RES                //  reserved (unused)
	OP_LEA                //  load effective address
	OP_TRAP               //  execute trap
)

// Cond flags
const (
	FL_POS = 1 << 0 // Posive
	FL_ZRO = 1 << 1 // Zero
	FL_NEG = 1 << 2 // Negative
)

// List of Trap codes
const (
	TRAP_GETC  uint16 = 0x20 // get character from keyboard
	TRAP_OUT   uint16 = 0x21 // output a character
	TRAP_PUTS  uint16 = 0x22 // output a word string
	TRAP_IN    uint16 = 0x23 // input a string
	TRAP_PUTSP uint16 = 0x24 // output a byte string
	TRAP_HALT  uint16 = 0x25 // halt the program
)

type CPU struct {
	registers        [R_COUNT]uint16
	memory           [math.MaxUint16 + 1]uint16 // 65536 locations (128kb)
	currInstruction  uint16
	currentOperation uint16
	isRunning        bool
	pc               uint16
	cond             uint16 // cond flags
}

func NewCPU(memorySize int) *CPU {
	cpu := &CPU{}
	cpu.pc = 0x3000
	return cpu
}

func (cpu *CPU) Reset() {
	cpu.registers = [R_COUNT]uint16{}
	cpu.currInstruction = 0
	cpu.currentOperation = 0
	cpu.isRunning = false
	cpu.pc = 0x3000
}

func (cpu *CPU) Execute() {
	cpu.registers[R_PC] = cpu.pc
	cpu.isRunning = true
	for cpu.isRunning {
		cpu.currInstruction = cpu.memory[cpu.registers[R_PC]]
		cpu.pc++

		cpu.currentOperation = cpu.currInstruction >> 12
		switch cpu.currentOperation {
		case OP_BR:
			cpu.branch()
		case OP_ADD:
			cpu.add()
		case OP_SUB:
			cpu.sub()
		case OP_LD:
			cpu.load()
		case OP_ST:
			cpu.store()
		case OP_JSR:
			cpu.jumpRegister()
		case OP_AND:
			cpu.bitwiseAnd()
		case OP_LDR:
			cpu.loadRegister()
		case OP_STR:
			cpu.storeRegister()
		case OP_NOT:
			cpu.bitwiseNot()
		case OP_LDI:
			cpu.loadIndirect()
		case OP_STI:
			cpu.storeIndirect()
		case OP_JMP:
			cpu.jump()
		case OP_LEA:
			cpu.loadEffectiveAddress()
		case OP_RES:
		case OP_RTI:
		case OP_TRAP:
			cpu.traps()
		default:
			log.Printf("operation code not implemented: 0x%04X - 0x%04x", cpu.currentOperation, cpu.currInstruction)
		}
	}
}

func (cpu *CPU) Stop() {
	cpu.isRunning = false
}

func sign_extend(x uint16, bit_count int) uint16 {
	if (x>>(bit_count-1))&1 != 0 {
		x |= (0xFFFF << bit_count)
	}
	return x
}

func (cpu *CPU) updateFlags(regIndex uint16) {
	if cpu.registers[regIndex] == 0 {
		cpu.cond = FL_ZRO
	} else if cpu.registers[regIndex]>>15 != 0 { // the most significant bit (leftmost) is 1, indicating a negative number
		cpu.cond = FL_NEG
	} else {
		cpu.cond = FL_POS
	}
}

// ============Instructions============
func (cpu *CPU) branch() {
	pcOffset := sign_extend(cpu.currInstruction&0x1ff, 9)
	condFlag := (cpu.currInstruction >> 9) & 0x7
	if (condFlag & cpu.registers[R_COND]) != 0 {
		cpu.registers[R_PC] += pcOffset
	}
}

func (cpu *CPU) add() {
	r0 := (cpu.currInstruction >> 9) & 0x7
	r1 := (cpu.currInstruction >> 6) & 0x7
	immFlag := (cpu.currInstruction >> 5) & 0x1
	if immFlag != 0 {
		imm5 := sign_extend(cpu.currInstruction&0x001F, 5)
		cpu.registers[r0] = cpu.registers[r1] + imm5
	} else {
		r2 := cpu.currInstruction & 0x7
		cpu.registers[r0] = cpu.registers[r1] + cpu.registers[r2]
	}
	cpu.updateFlags(r0)
}

func (cpu *CPU) sub() {

}

func (cpu *CPU) load() {

}

func (cpu *CPU) store() {

}

func (cpu *CPU) jumpRegister() {

}

func (cpu *CPU) bitwiseAnd() {

}

func (cpu *CPU) loadRegister() {

}

func (cpu *CPU) storeRegister() {

}

func (cpu *CPU) bitwiseNot() {

}

func (cpu *CPU) loadIndirect() {

}

func (cpu *CPU) storeIndirect() {

}

func (cpu *CPU) jump() {

}

func (cpu *CPU) loadEffectiveAddress() {

}

// ============Traps============
func (cpu *CPU) traps() {
	switch cpu.currInstruction & 0xFF {
	case TRAP_GETC:
		cpu.trapGetc()
	case TRAP_OUT:
		cpu.trapOut()
	case TRAP_PUTS:
		cpu.trapPuts()
	case TRAP_IN:
		cpu.trapIn()
	case TRAP_PUTSP:
		cpu.trapPutsp()
	case TRAP_HALT:
		cpu.trapHalt()
	default:
		log.Printf("trap code not implemented: 0x%04X", cpu.currInstruction)
	}
}

func (cpu *CPU) trapGetc() {

}

func (cpu *CPU) trapOut() {

}

func (cpu *CPU) trapPuts() {

}

func (cpu *CPU) trapIn() {

}

func (cpu *CPU) trapPutsp() {

}

func (cpu *CPU) trapHalt() {
	cpu.Stop()
}
