package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/config"
	"github.com/abibby/remote-input/vigem"
	"github.com/abibby/remote-input/windows"
)

type Joysticks struct {
	controllers map[uint16]*vigem.Xbox360Controller
	emulator    *vigem.Emulator
}

func NewJoysticks() (*Joysticks, error) {
	emu, err := vigem.NewEmulator(func(vibration vigem.Vibration) {})
	if err != nil {
		return nil, err
	}
	return &Joysticks{
		emulator:    emu,
		controllers: map[uint16]*vigem.Xbox360Controller{},
	}, nil
}

func (j *Joysticks) Connect(index uint16) (*vigem.Xbox360Controller, error) {
	var err error
	e, ok := j.controllers[index]
	if !ok {
		e, err = j.emulator.CreateXbox360Controller()
		if err != nil {
			return nil, err
		}
		err = e.Connect()
		if err != nil {
			return nil, err
		}
		j.controllers[index] = e
	}
	return e, nil
}
func (j *Joysticks) Disconnect(index uint16) error {
	e, ok := j.controllers[index]
	if !ok {
		return nil
	}
	err := e.Disconnect()
	if err != nil {
		return err
	}
	err = e.Close()
	if err != nil {
		return err
	}
	delete(j.controllers, index)
	return nil
}

func (j *Joysticks) Close() error {
	errs := []error{}

	for _, c := range j.controllers {
		err := c.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

var joysticks *Joysticks

func main() {
	log.Printf("started")

	serverIP := fmt.Sprintf("%s:%d", config.Host, config.Port)

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	joysticks, err = NewJoysticks()
	if err != nil {
		log.Printf("Joystick setup failed: %v", err)
	}
	defer joysticks.Close()

	log.Printf("connected to %s", serverIP)

	events := make([]common.InputEvent, 0, 8)

	b := make([]byte, 24)
	for {
		events = append(events, common.InputEvent{})
		e := &events[len(events)-1]

		_, err := conn.Read(b)
		if err == io.EOF {
			log.Print("disconnected")
			return
		} else if err != nil {
			log.Print(err)
			continue
		}

		err = e.UnmarshalBinary(b)
		if err != nil {
			log.Print(err)
			continue
		}

		if e.EventType == common.EV_SYN {
			err = handleEvent(events)
			if err != nil {
				log.Print(err)
			}
			events = events[:0]
		}
	}

}

func handleEvent(events []common.InputEvent) error {
	syn := events[len(events)-1]
	rest := events[:len(events)-1]
	if syn.Value == common.DeviceTypeKeyboard {
		return handleKeyboard(rest)
	}
	if syn.Value == common.DeviceTypeMouse {
		return handleMouse(rest)
	}
	if syn.Value == common.DeviceTypeJoystick {
		return handleJoystick(rest, syn.Code)
	}

	return nil
}

func i16(v int32) int16 {
	if v >= math.MaxInt16-1 {
		return math.MaxInt16 - 1
	}
	if v <= math.MinInt16+1 {
		return math.MinInt16 + 1
	}
	return int16(v)
}
func ui8(v int32) uint8 {
	if v > math.MaxUint8 {
		return math.MaxUint8
	}
	if v < 0 {
		return 0
	}
	return uint8(v)
}

func handleJoystick(events []common.InputEvent, index uint16) error {
	keyboardEvents := []common.InputEvent{}
	// spew.Dump(events)
	controller, err := joysticks.Connect(index)
	if err != nil {
		return err
	}
	report := controller.State()
	// report := vigem.NewXbox360ControllerReport()
	for _, e := range events {
		switch e.EventType {
		case common.EV_ABS:
			switch e.Code {
			case 0:
				report.SetLeftThumbX(i16(e.Value))
			case 1:
				log.Print(e.Value)
				report.SetLeftThumbY(-i16(e.Value))
			case 2:
				report.SetLeftTrigger(ui8(e.Value / (1024 / math.MaxUint8)))
			case 3:
				report.SetRightThumbX(i16(e.Value))
			case 4:
				report.SetRightThumbY(-i16(e.Value))
			case 5:
				report.SetRightTrigger(ui8(e.Value / (1024 / math.MaxUint8)))
			case 16:
				switch e.Value {
				case -1:
					report.SetButton(vigem.Xbox360ControllerButtonLeft)
				case 0:
					report.ClearButton(vigem.Xbox360ControllerButtonLeft)
					report.ClearButton(vigem.Xbox360ControllerButtonRight)
				case 1:
					report.SetButton(vigem.Xbox360ControllerButtonRight)
				}
			case 17:
				switch e.Value {
				case -1:
					report.SetButton(vigem.Xbox360ControllerButtonUp)
				case 0:
					report.ClearButton(vigem.Xbox360ControllerButtonUp)
					report.ClearButton(vigem.Xbox360ControllerButtonDown)
				case 1:
					report.SetButton(vigem.Xbox360ControllerButtonDown)
				}
			default:
				log.Printf("%v %v", e.Code, e.Value)
			}
		case common.EV_KEY:
			if e.Code < uint16(len(keyMap)) {
				keyboardEvents = append(keyboardEvents, e)
			} else {
				log.Printf("%v: %x %x\n", e.EventType, e.Code, e.Value)
				buttonID := e.Code - common.JOYSTICK_BASE
				log.Printf("gamepad button %d\n", buttonID)
				if int(buttonID) > len(gamepadMap) || gamepadMap[buttonID] == -1 {
					log.Printf("no mapping for button %d", buttonID)
					continue
				}
				btn := gamepadMap[buttonID]
				switch e.Value {
				case 0:
					report.ClearButton(btn)
				case 1:
					report.SetButton(btn)
				}
			}
		}
	}

	err = handleKeyboard(keyboardEvents)
	if err != nil {
		return err
	}
	return controller.Send(report)
}
func handleMouse(events []common.InputEvent) error {
	var flags uint32
	var data int32
	var dx int32
	var dy int32
	for _, e := range events {
		switch e.EventType {
		case common.EV_REL:
			switch e.Code {
			case 0:
				dx = e.Value
				flags |= MouseEventFMove
			case 1:
				dy = e.Value
				flags |= MouseEventFMove
			case 11:
				data = e.Value
				flags |= MouseEventFWheel
			}
		case common.EV_KEY:
			switch e.Code {
			case common.MOUSE_LEFT:
				if e.Value == 1 {
					flags |= MouseEventFLeftDown
				} else if e.Value == 0 {
					flags |= MouseEventFLeftUp
				}
			case common.MOUSE_RIGHT:
				if e.Value == 1 {
					flags |= MouseEventFRightDown
				} else if e.Value == 0 {
					flags |= MouseEventFRightUp
				}
			}
		}
	}
	return windows.SendMouseInput(dx, dy, data, flags)
}
func handleKeyboard(events []common.InputEvent) error {
	for _, e := range events {
		if e.EventType != common.EV_KEY {
			continue
		}
		vKey := keyMap[e.Code]
		if vKey == 0 {
			return fmt.Errorf("no map for key code %d", e.Code)
		}

		var flag windows.KeyEventFlag
		if e.Value == 0 {
			flag = windows.KEYEVENTF_KEYUP
		} else if e.Value == 1 {
			flag = windows.KEYEVENTF_KEYPRESS
		} else if e.Value == 2 {
			flag = windows.KEYEVENTF_KEYPRESS
		}
		return windows.SendInput(vKey, flag)
	}
	return nil
}
