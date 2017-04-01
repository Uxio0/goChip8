package main

import (
	"fmt"
	"log"
	"os"
)

type uint4 byte

type Chip8Engine struct {
	memory       [0xFFF]byte
	register     [16]byte
	iRegister    uint16
	stack        [16]uint16
	delayTimer   uint4
	soundTimer   uint4
	pc           uint16
	screen       [64][32]bool
	stackPointer uint4
}

func (c *Chip8Engine) storeRom(rom []byte) int {
	for i := 0; i < len(rom); i++ {
		c.memory[i+512] = rom[i]
	}
	c.pc = 512
	c.stackPointer = 0
	return len(rom)
}

func (c *Chip8Engine) getOpCode(i uint16) uint16 {
	return (uint16(c.memory[i]) << 8) + uint16(c.memory[i+1])
}

func (c *Chip8Engine) currentInstruction() uint16 {
	return c.getOpCode(c.pc)
}

func (c *Chip8Engine) showAllOpCodes() uint16 {
	for i := 512; c.memory[i] != 0; i += 2 {
		opcode := (uint16(c.memory[i]) << 8) + uint16(c.memory[i+1])
		selectOpCode(i, opcode)
	}
	return 0
}

func (c *Chip8Engine) runCycle() {
	var opcode uint16 = c.currentInstruction()
	fmt.Printf("%03x\tOpcode\t%04x\t", c.pc, opcode)
	c.pc += 2
	switch {
	case opcode == 0x00E0:
		fmt.Println("disp_clear()")
	case opcode == 0x00EE:
		fmt.Println("Return")
		c.stackPointer -= 1
		c.pc = c.stack[c.stackPointer]
	case opcode>>12 == 1:
		address := 0x0FFF & opcode
		fmt.Printf("Goto %x\n", address)
		c.pc = address
	case opcode>>12 == 2:
		address := 0x0FFF & opcode
		fmt.Println("Calls subroutine on %x\n", address)
		c.stack[c.stackPointer] = c.pc
		c.stackPointer++
		c.pc = address
	case opcode>>12 == 3:
		register := 0x0F00 & opcode >> 8
		value := byte(0x00FF & opcode)
		fmt.Printf("if(V%x==%x)\n", register, value)
		if c.register[register] == value {
			c.pc += 2
		}
	case opcode>>12 == 4:
		register := 0x0F00 & opcode >> 8
		value := byte(0x00FF & opcode)
		fmt.Printf("if(V%x!=%x)\n", register, value)
		if c.register[register] != value {
			c.pc += 2
		}
	case opcode>>12 == 5 && opcode&0x000F == 0:
		register1 := 0x0F00 & opcode >> 8
		register2 := 0x00F0 & opcode >> 4
		fmt.Printf("if(V%x==V%x)\n", register1, register2)
		if c.register[register1] == c.register[register2] {
			c.pc += 2
		}
	case opcode>>12 == 6:
		register := 0x0F00 & opcode >> 8
		value := byte(0x00FF & opcode)
		fmt.Printf("V%x=%x\n", register, value)
		c.register[register] = value
	case opcode>>12 == 7:
		register := 0x0F00 & opcode >> 8
		value := byte(0x00FF & opcode)
		fmt.Printf("V%x+=%x\n", register, value)
		c.register[register] += value
	case opcode>>12 == 8:
		register1 := 0x0F00 & opcode >> 8
		register2 := 0x00F0 & opcode >> 4
		switch {
		case opcode&0x000F == 0:
			fmt.Printf("V%x=V%x\n", register1, register2)
			c.register[register1] = c.register[register2]
		case opcode&0x000F == 1:
			fmt.Printf("V%x=V%x|V%x\n", register1, register1, register2)
			c.register[register1] |= c.register[register2]
			c.register[0xF] = 0
		case opcode&0x000F == 2:
			fmt.Printf("V%x=V%x&V%x\n", register1, register1, register2)
			c.register[register1] &= c.register[register2]
			c.register[0xF] = 0
		case opcode&0x000F == 3:
			fmt.Printf("V%x=V%x^V%x\n", register1, register1, register2)
			c.register[register1] ^= c.register[register2]
			c.register[0xF] = 0
		case opcode&0x000F == 4:
			fmt.Printf("V%x=V%x+V%x\n", register1, register1, register2)
			if (uint16(c.register[register1]) + uint16(c.register[register2])) > 0xFF {
				c.register[0xF] = 1
			} else {
				c.register[0xF] = 0
			}
			c.register[register1] += c.register[register2]
		case opcode&0x000F == 5:
			fmt.Printf("V%x=V%x-V%x\n", register1, register1, register2)
			if c.register[register1] > c.register[register2] {
				c.register[0xF] = 1
			} else {
				c.register[0xF] = 0
			}
			c.register[register1] -= c.register[register2]
		case opcode&0x000F == 6:
			fmt.Printf("V%x >> 1\n", register1)
			c.register[0xF] = c.register[register1] & 0x01
			c.register[register1] = c.register[register1] >> 1
		case opcode&0x000F == 7:
			fmt.Printf("V%x=V%x-V%x\n", register1, register2, register1)
			if c.register[register2] > c.register[register1] {
				c.register[0xF] = 1
			} else {
				c.register[0xF] = 0
			}
			c.register[register1] = c.register[register2] - c.register[register1]
		case opcode&0x000F == 0xE:
			fmt.Printf("V%x << 1\n", register1)
			c.register[0xF] = c.register[register1] >> 7
			c.register[register1] = c.register[register1] << 1

		}
	case opcode>>12 == 9 && opcode&0x000F == 0:
		register1 := 0x0F00 & opcode >> 8
		register2 := 0x00F0 & opcode >> 4
		fmt.Printf("if(V%x!=V%x)\n", register1, register2)
		if c.register[register1] != c.register[register2] {
			c.pc += 2
		}
	case opcode>>12 == 0xA:
		fmt.Println("I = NNN")
	case opcode>>12 == 0xB:
		fmt.Println("PC=V0+NNN")
	case opcode>>12 == 0xC:
		fmt.Println("Vx=rand()&NN")
	case opcode>>12 == 0xD:
		fmt.Println("draw(Vx,Vy,N)")
	case opcode&0xF0FF == 0xE09E:
		fmt.Println("if(key()==Vx)")
	case opcode&0xF0FF == 0xE0A1:
		fmt.Println("if(key()!=Vx)")
	case opcode&0xF0FF == 0xF007:
		fmt.Println("Vx = get_delay()")
	case opcode&0xF0FF == 0xF00A:
		fmt.Println("Vx = get_key()")
	case opcode&0xF0FF == 0xF015:
		fmt.Println("delay_timer(Vx)")
	case opcode&0xF0FF == 0xF018:
		fmt.Println("sound_timer(Vx)")
	case opcode&0xF0FF == 0xF01E:
		fmt.Println("I +=Vx")
	case opcode&0xF0FF == 0xF029:
		fmt.Println("I=sprite_addr[Vx]")
	case opcode&0xF0FF == 0xF033:
		fmt.Println("set_BCD(Vx)")
	case opcode&0xF0FF == 0xF055:
		fmt.Println("reg_dump(Vx,&I)")
	case opcode&0xF0FF == 0xF065:
		fmt.Println("reg_load(Vx,&I)")
	}
}

func readOpCodes() []byte {
	var rom []byte
	file, err := os.Open("PONG")
	if err != nil {
		log.Fatal(err)
	}
	b1 := make([]byte, 1)
	for {
		_, err := file.Read(b1)
		if err != nil {
			break
		}
		rom = append(rom, b1[0])
		//fmt.Printf("%d bytes: %04x\n", n1, x)
	}
	return rom
}

func selectOpCode(index int, opcode uint16) {
	fmt.Printf("%03x\tOpcode\t%04x\t", index, opcode)
	switch {
	case opcode == 0x00E0:
		fmt.Println("disp_clear()")
	case opcode == 0x00EE:
		fmt.Println("Return")
	case opcode>>12 == 1:
		address := 0x0FFF & opcode
		fmt.Printf("Goto %x\n", address)
	case opcode>>12 == 2:
		fmt.Println("Calls subroutine")
	case opcode>>12 == 3:
		fmt.Println("if(Vx==NN)")
	case opcode>>12 == 4:
		fmt.Println("if(Vx!=NN)")
	case opcode>>12 == 5 && opcode&0x000F == 0:
		fmt.Println("if(Vx==Vy)")
	case opcode>>12 == 6:
		fmt.Println("Set VX to NN")
	case opcode>>12 == 7:
		fmt.Println("Adds NN to VX")
	case opcode>>12 == 8 && opcode&0x000F == 0:
		fmt.Println("Vx=Vy")
	case opcode>>12 == 8 && opcode&0x000F == 1:
		fmt.Println("Vx=Vx|Vy")
	case opcode>>12 == 8 && opcode&0x000F == 2:
		fmt.Println("Vx=Vx&Vy")
	case opcode>>12 == 8 && opcode&0x000F == 3:
		fmt.Println("Vx=Vx^Vy")
	case opcode>>12 == 8 && opcode&0x000F == 4:
		fmt.Println("Vx += Vy")
	case opcode>>12 == 8 && opcode&0x000F == 5:
		fmt.Println("Vx -= Vy")
	case opcode>>12 == 8 && opcode&0x000F == 6:
		fmt.Println("Vx >> 1")
	case opcode>>12 == 8 && opcode&0x000F == 7:
		fmt.Println("Vx=Vy-Vx")
	case opcode>>12 == 8 && opcode&0x000F == 0xE:
		fmt.Println("Vx << 1")
	case opcode>>12 == 9 && opcode&0x000F == 0:
		fmt.Println("Vx << 1")
	case opcode>>12 == 0xA:
		fmt.Println("I = NNN")
	case opcode>>12 == 0xB:
		fmt.Println("PC=V0+NNN")
	case opcode>>12 == 0xC:
		fmt.Println("Vx=rand()&NN")
	case opcode>>12 == 0xD:
		fmt.Println("draw(Vx,Vy,N)")
	case opcode&0xF0FF == 0xE09E:
		fmt.Println("if(key()==Vx)")
	case opcode&0xF0FF == 0xE0A1:
		fmt.Println("if(key()!=Vx)")
	case opcode&0xF0FF == 0xF007:
		fmt.Println("Vx = get_delay()")
	case opcode&0xF0FF == 0xF00A:
		fmt.Println("Vx = get_key()")
	case opcode&0xF0FF == 0xF015:
		fmt.Println("delay_timer(Vx)")
	case opcode&0xF0FF == 0xF018:
		fmt.Println("sound_timer(Vx)")
	case opcode&0xF0FF == 0xF01E:
		fmt.Println("I +=Vx")
	case opcode&0xF0FF == 0xF029:
		fmt.Println("I=sprite_addr[Vx]")
	case opcode&0xF0FF == 0xF033:
		fmt.Println("set_BCD(Vx)")
	case opcode&0xF0FF == 0xF055:
		fmt.Println("reg_dump(Vx,&I)")
	case opcode&0xF0FF == 0xF065:
		fmt.Println("reg_load(Vx,&I)")
	}
}

func main() {
	// Memory begins at 0x200
	// 0xEA0 - 0xEFF Call stack, internal use and other variables
	// 0xF00 - 0xFFF to display refresh

	// Registers. 16 8bit registers from V0 to VF
	// VF -> Carry flag on addition, not borrow flag on substraction and pixel collision on drawing
	// Address register -> I, 16 bits and used with serveral opcodes that involve memory operations

	// Stack. 16 levels of nesting. Used when subroutines are called

	// Timers
	// Delay timer
	// Sound timer. When value is nonzero, a beeping sound is made

	// Input
	// Input is done with a hex keyboard of 16 keys

	// Graphics and sound
	// Resolution -> 64x32 pixels monochrome. Sprites are 8 pixels wide and 1-15 pixels in height
	// Sound -> Beeping sound when sound timer is not zero
	engine := Chip8Engine{}
	rom := readOpCodes()
	engine.storeRom(rom)

	//fmt.Printf("current-instruction %04x\n", engine.currentInstruction())
	//engine.showAllOpCodes()
	engine.runCycle()
}
