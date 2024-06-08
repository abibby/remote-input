package main

// Various MouseInput dwFlags.
//
// Refer to the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-mouseinput
const (
	MouseEventFAbsolute       uint32 = 0x8000
	MouseEventFHWheel         uint32 = 0x01000
	MouseEventFMove           uint32 = 0x0001
	MouseEventFMoveNoCoalesce uint32 = 0x2000
	MouseEventFLeftDown       uint32 = 0x0002
	MouseEventFLeftUp         uint32 = 0x0004
	MouseEventFRightDown      uint32 = 0x0008
	MouseEventFRightUp        uint32 = 0x0010
	MouseEventFMiddleDown     uint32 = 0x0020
	MouseEventFMiddleUp       uint32 = 0x0040
	MouseEventFVirtualDesk    uint32 = 0x4000
	MouseEventFWheel          uint32 = 0x0800
	MouseEventFXDown          uint32 = 0x0080
	MouseEventFXUp            uint32 = 0x0100
)
