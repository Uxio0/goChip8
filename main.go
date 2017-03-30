package main

import (
	"fmt"
	"log"
	"os"
)

type uint4 uint8

type Chip8Engine struct {
	rom        []uint16
	register   [16]uint8
	stack      [16]uint4
	delayTimer uint4
	soundTimer uint4
	pc         uint16
	screen     [64][32]bool
}

func (c *Chip8Engine) currentInstruction() uint16 {
	return c.rom[c.pc]
}

func (c *Chip8Engine) showAllOpCodes() uint8 {
	for i, opcode := range c.rom {
		selectOpCode(i, opcode)
	}
	return 0
}

func readOpCodes() []uint16 {
	var rom []uint16
	file, err := os.Open("PONG")
	if err != nil {
		log.Fatal(err)
	}
	b1 := make([]byte, 2)
	for {
		_, err := file.Read(b1)
		if err != nil {
			break
		}
		var x uint16 = (uint16(b1[0]) << 8) + uint16(b1[1])
		rom = append(rom, x)
		//fmt.Printf("%d bytes: %04x\n", n1, x)
	}
	return rom
}

func selectOpCode(index int, opcode uint16) {
	fmt.Printf("%03x\tOpcode\t%04x\t", index<<1+0x200, opcode)
	switch {
	case opcode == 0x00E0:
		fmt.Println("disp_clear()")
	case opcode == 0x00EE:
		fmt.Println("Return")
	case opcode>>12 == 1:
		fmt.Println("Goto")
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
	engine.rom = rom

	fmt.Printf("current-instruction %04x\n", engine.currentInstruction())
	engine.showAllOpCodes()
}
