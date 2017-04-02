package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func readOpCodes(romFileName string) []byte {
	var rom []byte
	file, err := os.Open(romFileName)
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
	var romFileName string
	flag.StringVar(&romFileName, "rom", "games/PONG2", "ROM path")
	flag.Parse()
	rom := readOpCodes(romFileName)
	engine := Chip8Engine{}
	engine.Init(rom)

	for {
		engine.RunCycle()
		timer := time.NewTimer(time.Second / 120) //60Hz
		<-timer.C
	}
}
