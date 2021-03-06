package main

import (
	"fmt"
	"math/rand"
)

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

var fontset = []byte{
	// Standard sprites for showing (some) text.
	0xf0, 0x90, 0x90, 0x90, 0xf0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xf0, 0x10, 0xf0, 0x80, 0xf0, // 2
	0xf0, 0x10, 0xf0, 0x10, 0xf0, // 3
	0x90, 0x90, 0xf0, 0x10, 0x10, // 4
	0xf0, 0x80, 0xf0, 0x10, 0xf0, // 5
	0xf0, 0x80, 0xf0, 0x90, 0xf0, // 6
	0xf0, 0x10, 0x20, 0x40, 0x40, // 7
	0xf0, 0x90, 0xf0, 0x90, 0xf0, // 8
	0xf0, 0x90, 0xf0, 0x10, 0xf0, // 9
	0xf0, 0x90, 0xf0, 0x90, 0x90, // A
	0xe0, 0x90, 0xe0, 0x90, 0xe0, // B
	0xf0, 0x80, 0x80, 0x80, 0xf0, // C
	0xe0, 0x90, 0x90, 0x90, 0xe0, // D
	0xf0, 0x80, 0xf0, 0x80, 0xf0, // E
	0xf0, 0x80, 0xf0, 0x80, 0x80, // F
}

//Chip8Engine manages Chip8 CPU
type Chip8Engine struct {
	memory       [0xFFF]byte
	register     [16]byte
	iRegister    uint16
	stack        [16]uint16
	delayTimer   byte
	soundTimer   byte
	pc           uint16
	screen       [64][32]bool
	stackPointer byte
	sdlWindow    SDLWindow
}

func (c *Chip8Engine) Init(rom []byte) int {
	c.sdlWindow.Init()
	for i, fontByte := range fontset {
		c.memory[i] = fontByte
	}
	for i := 0; i < len(rom); i++ {
		c.memory[i+512] = rom[i]
	}
	c.pc = 512
	c.stackPointer = 0
	return len(rom)
}

func (c *Chip8Engine) printScreen() {
	fmt.Println("")
	for y := 0; y < 32; y++ {
		fmt.Printf("%02x ", y)
		for x := 0; x < 64; x++ {
			if c.screen[x][y] {
				fmt.Print("-")
			} else {
				fmt.Print("*")
			}
		}
		fmt.Println("")
	}
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

func (c *Chip8Engine) RunCycle() {
	opcode := c.currentInstruction()
	fmt.Printf("%03x\tOpcode\t%04x\t", c.pc, opcode)
	c.pc += 2
	if c.delayTimer > 0 {
		c.delayTimer--
	}
	if c.soundTimer > 0 {
		c.soundTimer--
	}
	switch {
	case opcode == 0x00E0:
		fmt.Println("disp_clear()")
		c.sdlWindow.Clear()
		for i := 0; i < len(c.screen); i++ {
			for j := 0; j < len(c.screen[i]); j++ {
				c.screen[i][j] = false
			}
		}
	case opcode == 0x00EE:
		fmt.Println("Return")
		c.stackPointer--
		c.pc = c.stack[c.stackPointer]
	case opcode>>12 == 1:
		address := 0x0FFF & opcode
		fmt.Printf("Goto %x\n", address)
		c.pc = address
	case opcode>>12 == 2:
		address := 0x0FFF & opcode
		fmt.Printf("Calls subroutine on %x\n", address)
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
		address := 0x0FFF & opcode
		fmt.Printf("I = %x\n", address)
		c.iRegister = address
	case opcode>>12 == 0xB:
		address := 0x0FFF & opcode
		fmt.Printf("PC=V0+%x\n", address)
		c.pc = uint16(c.register[0]) + address
	case opcode>>12 == 0xC:
		register := 0x0F00 & opcode >> 8
		value := byte(0x00FF & opcode)
		// TODO Seed
		fmt.Printf("V%x=rand() & %x\n", register, value)
		c.register[register] = uint8(rand.Intn(0xFF)) & value
	case opcode>>12 == 0xD:
		register1 := 0x0F00 & opcode >> 8
		register2 := 0x00F0 & opcode >> 4
		constant := byte(0x000F & opcode)
		fmt.Printf("draw(V%x,V%x,%x)\n", register1, register2, constant)
		x1 := c.register[register1]
		y1 := c.register[register2]
		fmt.Printf("draw(%x,%x,%x)\n", x1, y1, constant)
		for y := y1; y < y1+constant && y < 32; y++ {
			for x := byte(0); x < 8 && (x+x1) < 64; x++ {
				current := c.screen[x+x1][y]
				c.screen[x+x1][y] = c.screen[x+x1][y] != ((c.memory[c.iRegister+uint16(y-y1)] >> (7 - x) & 0x1) == 1)
				if current != c.screen[x+x1][y] {
					c.register[0xF] = 0x1
				}
			}
		}
		//c.printScreen()
		c.sdlWindow.Draw(c.screen)
	case opcode&0xF0FF == 0xE09E:
		register := 0x0F00 & opcode >> 8
		if c.sdlWindow.CheckKeyPress(c.register[register]) {
			c.pc += 2
		}
		fmt.Printf("if(key()==V%x\n)", register)
	case opcode&0xF0FF == 0xE0A1:
		register := 0x0F00 & opcode >> 8
		if !c.sdlWindow.CheckKeyPress(c.register[register]) {
			c.pc += 2
		}
		fmt.Printf("if(key()!=V%x\n)", register)
	case opcode&0xF0FF == 0xF007:
		register := 0x0F00 & opcode >> 8
		fmt.Printf("V%x = get_delay()\n", register)
		c.register[register] = byte(c.delayTimer)
	case opcode&0xF0FF == 0xF00A:
		register := 0x0F00 & opcode >> 8
		fmt.Printf("V%x = get_key()\n", register)
		pressed := c.sdlWindow.WaitUntilKeyPress()
		c.register[register] = pressed
	case opcode&0xF0FF == 0xF015:
		register := 0x0F00 & opcode >> 8
		fmt.Printf("delay_timer(V%x)\n", register)
		c.delayTimer = c.register[register]
	case opcode&0xF0FF == 0xF018:
		register := 0x0F00 & opcode >> 8
		fmt.Printf("sound_timer(V%x)\n", register)
		c.soundTimer = c.register[register]
	case opcode&0xF0FF == 0xF01E:
		register := 0x0F00 & opcode >> 8
		fmt.Printf("I +=V%x\n", register)
		c.iRegister += uint16(c.register[register])
		// This is an undocumented feature of the CHIP-8 and used by the "Spacefight 2091!" game.
		if c.iRegister+uint16(c.register[register]) > 0xFFF {
			c.register[0xF] = 1
		} else {
			c.register[0xF] = 0
		}
	case opcode&0xF0FF == 0xF029:
		register := 0x0F00 & opcode >> 8
		fmt.Printf("I=sprite_addr[V%x]\n", register)
		// 0 is in 0x0, 1 in 0x5, 2 in 0xA
		c.iRegister = uint16(c.register[register]) * 5
	case opcode&0xF0FF == 0xF033:
		register := 0x0F00 & opcode >> 8
		fmt.Printf("set_BCD(V%x)\n", register)
		bcd := c.register[register]
		c.memory[c.iRegister] = bcd / 100 % 10
		c.memory[c.iRegister+1] = bcd / 10 % 10
		c.memory[c.iRegister+2] = bcd / 1 % 10
	case opcode&0xF0FF == 0xF055:
		register := 0x0F00 & opcode >> 8
		fmt.Printf("reg_dump(V%x,&I)\n", register)
		for i := uint16(0); i <= register; i++ {
			c.memory[c.iRegister+i] = c.register[i]
		}
	case opcode&0xF0FF == 0xF065:
		register := 0x0F00 & opcode >> 8
		fmt.Printf("reg_load(V%x,&I)\n", register)
		for i := uint16(0); i <= register; i++ {
			c.register[i] = c.memory[c.iRegister+i]
		}
	}
}
