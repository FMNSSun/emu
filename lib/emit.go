package lib

import "strings"
import "io"
import "errors"
import "strconv"

func Put_Imm12(memory io.Writer, opc, a, b, dst uint8, offs uint16) (int, error) {
	buf := []byte{0x00, 0x00, 0x00, 0x00}
	buf[0] = opc
	buf[1] = (a << 4) | (b << 0)
	buf[2] = uint8(dst<<4) | uint8((offs>>8)&0x0F)
	buf[3] = uint8(offs & 0xFF)

	return memory.Write(buf)
}

func Put_Imm16(memory io.Writer, opc uint8, dst uint8, v uint16) (int, error) {
	vl := uint8(v & 0x00FF)
	vh := uint8((v >> 8) & 0x00FF)

	buf := []byte{0x00, 0x00, 0x00, 0x00}
	buf[0] = opc
	buf[1] = (dst & 0x0F)
	buf[2] = vl
	buf[3] = vh

	return memory.Write(buf)
}

func Put_Regs(memory io.Writer, opc uint8, src uint8, dst uint8) (int, error) {
	buf := []byte{0x00, 0x00}
	buf[0] = opc
	buf[1] = (src << 4) | (dst << 0)

	return memory.Write(buf)
}

func Put_Lines(memory io.Writer, lines []string) (int, error) {
	for i, line := range lines {
		_, err := Put_Line(memory, line)

		if err != nil {
			return i, err
		}
	}

	return -1, nil
}

func Put_Line(memory io.Writer, line string) (int, error) {
	parts := splitLine(line)

	if len(parts) < 1 {
		return -1, errors.New("Syntax error: Empty line!")
	}

	opc, ok := OPC_Table_Name2Opc[parts[0]]

	if !ok {
		return -2, errors.New("Error: Unknown instruction!")
	}

	iflag := (opc >> 6) & 0x03

	if iflag == IF_REGS {
		if len(parts) != 3 {
			if len(parts) == 1 {
				return Put_Regs(memory, opc, 0, 0)
			} else {
				return -5, errors.New("Syntax error: Invalid number of arguments!")
			}
		}

		// need two registers
		reg_src, ok := REG_Table_Name2Reg[parts[1]]

		if !ok {
			return -3, errors.New("Error: Unknown register!")
		}

		reg_dst, ok := REG_Table_Name2Reg[parts[2]]

		if !ok {
			return -4, errors.New("Error: Unknown register!")
		}

		return Put_Regs(memory, opc, reg_src, reg_dst)
	} else if iflag == IF_IMM16 {
		// need a register and a imm16 constant
	} else if iflag == IF_IMM12 {
		if len(parts) != 5 {
			return -6, errors.New("Syntax error: Invalid number of arguments!")
		}

		reg_a, ok := REG_Table_Name2Reg[parts[1]]

		if !ok {
			return -7, errors.New("Error: Unknown register!")
		}

		reg_b, ok := REG_Table_Name2Reg[parts[2]]

		if !ok {
			return -8, errors.New("Error: Unknown register!")
		}

		reg_dst, ok := REG_Table_Name2Reg[parts[3]]

		if !ok {
			return -9, errors.New("Error: Unknown register!")
		}

		u, err := strconv.ParseUint(parts[4], 16, 12)

		if err != nil {
			return -10, errors.New("Syntax error: Invalid literal for imm12!")
		}

		return Put_Imm12(memory, opc, reg_a, reg_b, reg_dst, uint16(u))
	}

	return -2, nil
}

func splitLine(line string) []string {
	return strings.FieldsFunc(line, func(c rune) bool {
		return c == ' ' || c == '\t' || c == ','
	})
}
