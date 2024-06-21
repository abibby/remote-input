package main

import (
	"fmt"
	"log"
	"math"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/vigem"
	"github.com/abibby/remote-input/windows"
)

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

var joysticks *Joysticks

func handleJoystick(events []common.InputEvent, index uint16) error {
	keyboardEvents := []common.InputEvent{}

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
