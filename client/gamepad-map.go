package main

import (
	"github.com/abibby/remote-input/vigem"
)

var gamepadMap = []int{
	vigem.Xbox360ControllerButtonA, // 0
	vigem.Xbox360ControllerButtonB, // 1
	-1,                             // 2
	vigem.Xbox360ControllerButtonX, // 3
	vigem.Xbox360ControllerButtonY, // 4
	-1,                             // 5
	vigem.Xbox360ControllerButtonLeftShoulder,  // 6
	vigem.Xbox360ControllerButtonRightShoulder, // 7
	-1,                                      // 8
	-1,                                      // 9
	vigem.Xbox360ControllerButtonGuide,      // 10?
	vigem.Xbox360ControllerButtonStart,      // 11
	-1,                                      // 12
	vigem.Xbox360ControllerButtonLeftThumb,  // 13
	vigem.Xbox360ControllerButtonRightThumb, // 14
}
