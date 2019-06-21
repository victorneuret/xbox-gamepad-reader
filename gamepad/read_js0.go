package gamepad

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/exec"
	"unsafe"
)

const jsButtonEvent = 0x01 /* button pressed/released */
const jsAxisEvent = 0x02   /* joystick moved */

type JsEvent struct {
	Time       uint32
	Value      int16
	ActionType uint8
	Number     uint8
}

type MyGamepad struct {
	Button   [11]int16
	Joystick [8]float32
}

func printValue(gamepad *MyGamepad) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	for i := 0; i < 11; i++ {
		fmt.Printf("Button %d: ", i)
		if gamepad.Button[i] == 1 {
			fmt.Printf("1\n")
		} else {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")
	for i := 0; i < 8; i++ {
		fmt.Printf("Joysitck %d: ", i)
		fmt.Printf("%f\n", gamepad.Joystick[i])
	}
}

func fillStruct(js *JsEvent, gamepad *MyGamepad) {
	if js.ActionType == jsButtonEvent {
		gamepad.Button[js.Number] = js.Value
	} else if js.ActionType == jsAxisEvent {
		gamepad.Joystick[js.Number] = float32(float32(js.Value) / (65535.0 / 2.0))
	}
}

func ReadJoystickInput() {
	js := JsEvent{0, 0, 0, 0}
	gamepad := MyGamepad{}
	file, err := os.Open("/dev/input/js0")
	buffer := make([]byte, unsafe.Sizeof(js))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for true {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		if bytesRead == int(unsafe.Sizeof(js)) {
			binary.Read(bytes.NewBuffer(buffer[0:4]), binary.LittleEndian, &js.Time)
			binary.Read(bytes.NewBuffer(buffer[4:6]), binary.LittleEndian, &js.Value)
			binary.Read(bytes.NewBuffer(buffer[6:7]), binary.LittleEndian, &js.ActionType)
			binary.Read(bytes.NewBuffer(buffer[7:8]), binary.LittleEndian, &js.Number)
			fillStruct(&js, &gamepad)
		}
		printValue(&gamepad)
	}
}
