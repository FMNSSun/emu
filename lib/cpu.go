package lib


type RunContext struct {
	PC_Init   uint32
	Registers []uint32
	Memory    []uint8
	PC_End    uint32
	debugCh	chan uint32
}

func Run(rc *RunContext) {
	pc := rc.PC_Init
	regs := rc.Registers
	mem := rc.Memory
	memsz := uint32(len(mem))
	hiaddr := memsz - 4

	for {

		if pc > hiaddr {
			// hiaddr is the highest address with still 4 bytes left. 
			// can't do this!
			// TODO: handle this properly
			break
		}

		loc := pc

		// Fetch first two bytes of instruction
		b1 := uint8(mem[pc+0])
		b2 := uint8(mem[pc+1])

		// Decode it
		imm16 := uint16(0x00)

		iflag := (b1 >> 6) & 0x03
		opc := b1
		src := (b2 >> 4) & 0xFF
		dst := (b2 >> 0) & 0x0F
		target := uint8(0x00)
		imm12 := uint16(0x00)
		imm12_s := int16(0x00)

		pc += 2

		// If IMM16 fetch additional bytes
		if iflag == IF_IMM16 {
			b3 := uint8(mem[pc+0])
			b4 := uint8(mem[pc+1])

			imm16 = (uint16(b4) << 8) | (uint16(b3) << 0)

			pc += 2
		} else if iflag == IF_IMM12 {
			b3 := uint8(mem[pc+0])
			b4 := uint8(mem[pc+1])

			target = (b3 >> 4) & 0x03

			imm12 = (uint16((b3 & 0x03)) << 8) | (uint16(b4) << 0)

			if imm12 & 0x800 == 0x800 { // sign bit check?
				imm12_s = - int16(imm12 & 0x7FF)
			} else {
				imm12_s = int16(imm12 & 0x7FF)
			}
		}

		switch opc {
		case OPC_NOP:
		case OPC_HLT:
			// TODO: handle this properly
			goto exit_hlt
		case OPC_LDL:
			i_ldl_16(regs, mem, dst, imm16)
		case OPC_LDH:
			i_ldh_16(regs, mem, dst, imm16)
		case OPC_XOR:
			i_xor(regs, mem, src, dst)
		case OPC_ADD:
			i_add(regs, mem, src, dst)
		case OPC_BEQ:
			pc = i_beq_12(regs, mem, src, dst, target, pc, imm12_s)
		default:
			// invalid opcode.
			// TODO: handle this properly
		}

		if rc.debugCh != nil {
			rc.debugCh <- loc
			_ = <- rc.debugCh
		}

		continue

	exit_hlt:
		break
	}

	rc.PC_End = pc
	
	if rc.debugCh != nil {
		close(rc.debugCh)
	}
}

func i_ldl_16(regs []uint32, mem []byte, dst uint8, v uint16) {
	regs[dst] = (regs[dst] & MASK_H16) | uint32(v)
}

func i_ldh_16(regs []uint32, mem []byte, dst uint8, v uint16) {
	regs[dst] = (uint32(v) << 16) | (regs[dst] & MASK_L16)
}

func i_xor(regs []uint32, mem []byte, src, dst uint8) {
	regs[dst] ^= regs[src]
}

func i_add(regs []uint32, mem []byte, src, dst uint8) {
	regs[dst] += regs[src]
}

func i_beq_12(regs []uint32, mem []byte, a, b, dst uint8, curPc uint32, imm12_s int16) uint32 {
	if regs[a] == regs[b] {
		return uint32(int64(regs[dst]) + int64(imm12_s))
	} else {
		return curPc
	}
}
