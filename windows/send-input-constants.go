package windows

type VirtualKey uint16

const (
	VK_LBUTTON  = VirtualKey(0x01) //Left mouse button
	VK_RBUTTON  = VirtualKey(0x02) //Right mouse button
	VK_CANCEL   = VirtualKey(0x03) //Control-break processing
	VK_MBUTTON  = VirtualKey(0x04) //Middle mouse button
	VK_XBUTTON1 = VirtualKey(0x05) //X1 mouse button
	VK_XBUTTON2 = VirtualKey(0x06) //X2 mouse button
	// -	0x07	Reserved
	VK_BACK = VirtualKey(0x08) //BACKSPACE key
	VK_TAB  = VirtualKey(0x09) //TAB key
	// -	0x0A-0B	Reserved
	VK_CLEAR  = VirtualKey(0x0C) //CLEAR key
	VK_RETURN = VirtualKey(0x0D) //ENTER key
	// -	0x0E-0F	Unassigned
	VK_SHIFT      = VirtualKey(0x10) //SHIFT key
	VK_CONTROL    = VirtualKey(0x11) //CTRL key
	VK_MENU       = VirtualKey(0x12) //ALT key
	VK_PAUSE      = VirtualKey(0x13) //PAUSE key
	VK_CAPITAL    = VirtualKey(0x14) //CAPS LOCK key
	VK_KANA       = VirtualKey(0x15) //IME Kana mode
	VK_HANGUL     = VirtualKey(0x15) //IME Hangul mode
	VK_IME_ON     = VirtualKey(0x16) //IME On
	VK_JUNJA      = VirtualKey(0x17) //IME Junja mode
	VK_FINAL      = VirtualKey(0x18) //IME final mode
	VK_HANJA      = VirtualKey(0x19) //IME Hanja mode
	VK_KANJI      = VirtualKey(0x19) //IME Kanji mode
	VK_IME_OFF    = VirtualKey(0x1A) //IME Off
	VK_ESCAPE     = VirtualKey(0x1B) //ESC key
	VK_CONVERT    = VirtualKey(0x1C) //IME convert
	VK_NONCONVERT = VirtualKey(0x1D) //IME nonconvert
	VK_ACCEPT     = VirtualKey(0x1E) //IME accept
	VK_MODECHANGE = VirtualKey(0x1F) //IME mode change request
	VK_SPACE      = VirtualKey(0x20) //SPACEBAR
	VK_PRIOR      = VirtualKey(0x21) //PAGE UP key
	VK_NEXT       = VirtualKey(0x22) //PAGE DOWN key
	VK_END        = VirtualKey(0x23) //END key
	VK_HOME       = VirtualKey(0x24) //HOME key
	VK_LEFT       = VirtualKey(0x25) //LEFT ARROW key
	VK_UP         = VirtualKey(0x26) //UP ARROW key
	VK_RIGHT      = VirtualKey(0x27) //RIGHT ARROW key
	VK_DOWN       = VirtualKey(0x28) //DOWN ARROW key
	VK_SELECT     = VirtualKey(0x29) //SELECT key
	VK_PRINT      = VirtualKey(0x2A) //PRINT key
	VK_EXECUTE    = VirtualKey(0x2B) //EXECUTE key
	VK_SNAPSHOT   = VirtualKey(0x2C) //PRINT SCREEN key
	VK_INSERT     = VirtualKey(0x2D) //INS key
	VK_DELETE     = VirtualKey(0x2E) //DEL key
	VK_HELP       = VirtualKey(0x2F) //HELP key
	VK_0          = VirtualKey(0x30) // 0 key
	VK_1          = VirtualKey(0x31) // 1 key
	VK_2          = VirtualKey(0x32) // 2 key
	VK_3          = VirtualKey(0x33) // 3 key
	VK_4          = VirtualKey(0x34) // 4 key
	VK_5          = VirtualKey(0x35) // 5 key
	VK_6          = VirtualKey(0x36) // 6 key
	VK_7          = VirtualKey(0x37) // 7 key
	VK_8          = VirtualKey(0x38) // 8 key
	VK_9          = VirtualKey(0x39) // 9 key
	// -	0x3A-40	Undefined
	VK_A    = VirtualKey(0x41) // A key
	VK_B    = VirtualKey(0x42) // B key
	VK_C    = VirtualKey(0x43) // C key
	VK_D    = VirtualKey(0x44) // D key
	VK_E    = VirtualKey(0x45) // E key
	VK_F    = VirtualKey(0x46) // F key
	VK_G    = VirtualKey(0x47) // G key
	VK_H    = VirtualKey(0x48) // H key
	VK_I    = VirtualKey(0x49) // I key
	VK_J    = VirtualKey(0x4A) // J key
	VK_K    = VirtualKey(0x4B) // K key
	VK_L    = VirtualKey(0x4C) // L key
	VK_M    = VirtualKey(0x4D) // M key
	VK_N    = VirtualKey(0x4E) // N key
	VK_O    = VirtualKey(0x4F) // O key
	VK_P    = VirtualKey(0x50) // P key
	VK_Q    = VirtualKey(0x51) // Q key
	VK_R    = VirtualKey(0x52) // R key
	VK_S    = VirtualKey(0x53) // S key
	VK_T    = VirtualKey(0x54) // T key
	VK_U    = VirtualKey(0x55) // U key
	VK_V    = VirtualKey(0x56) // V key
	VK_W    = VirtualKey(0x57) // W key
	VK_X    = VirtualKey(0x58) // X key
	VK_Y    = VirtualKey(0x59) // Y key
	VK_Z    = VirtualKey(0x5A) // Z key
	VK_LWIN = VirtualKey(0x5B) //Left Windows key
	VK_RWIN = VirtualKey(0x5C) //Right Windows key
	VK_APPS = VirtualKey(0x5D) //Applications key
	// -	0x5E	Reserved
	VK_SLEEP     = VirtualKey(0x5F) //Computer Sleep key
	VK_NUMPAD0   = VirtualKey(0x60) //Numeric keypad 0 key
	VK_NUMPAD1   = VirtualKey(0x61) //Numeric keypad 1 key
	VK_NUMPAD2   = VirtualKey(0x62) //Numeric keypad 2 key
	VK_NUMPAD3   = VirtualKey(0x63) //Numeric keypad 3 key
	VK_NUMPAD4   = VirtualKey(0x64) //Numeric keypad 4 key
	VK_NUMPAD5   = VirtualKey(0x65) //Numeric keypad 5 key
	VK_NUMPAD6   = VirtualKey(0x66) //Numeric keypad 6 key
	VK_NUMPAD7   = VirtualKey(0x67) //Numeric keypad 7 key
	VK_NUMPAD8   = VirtualKey(0x68) //Numeric keypad 8 key
	VK_NUMPAD9   = VirtualKey(0x69) //Numeric keypad 9 key
	VK_MULTIPLY  = VirtualKey(0x6A) //Multiply key
	VK_ADD       = VirtualKey(0x6B) //Add key
	VK_SEPARATOR = VirtualKey(0x6C) //Separator key
	VK_SUBTRACT  = VirtualKey(0x6D) //Subtract key
	VK_DECIMAL   = VirtualKey(0x6E) //Decimal key
	VK_DIVIDE    = VirtualKey(0x6F) //Divide key
	VK_F1        = VirtualKey(0x70) //F1 key
	VK_F2        = VirtualKey(0x71) //F2 key
	VK_F3        = VirtualKey(0x72) //F3 key
	VK_F4        = VirtualKey(0x73) //F4 key
	VK_F5        = VirtualKey(0x74) //F5 key
	VK_F6        = VirtualKey(0x75) //F6 key
	VK_F7        = VirtualKey(0x76) //F7 key
	VK_F8        = VirtualKey(0x77) //F8 key
	VK_F9        = VirtualKey(0x78) //F9 key
	VK_F10       = VirtualKey(0x79) //F10 key
	VK_F11       = VirtualKey(0x7A) //F11 key
	VK_F12       = VirtualKey(0x7B) //F12 key
	VK_F13       = VirtualKey(0x7C) //F13 key
	VK_F14       = VirtualKey(0x7D) //F14 key
	VK_F15       = VirtualKey(0x7E) //F15 key
	VK_F16       = VirtualKey(0x7F) //F16 key
	VK_F17       = VirtualKey(0x80) //F17 key
	VK_F18       = VirtualKey(0x81) //F18 key
	VK_F19       = VirtualKey(0x82) //F19 key
	VK_F20       = VirtualKey(0x83) //F20 key
	VK_F21       = VirtualKey(0x84) //F21 key
	VK_F22       = VirtualKey(0x85) //F22 key
	VK_F23       = VirtualKey(0x86) //F23 key
	VK_F24       = VirtualKey(0x87) //F24 key
	// -	0x88-8F	Reserved
	VK_NUMLOCK = VirtualKey(0x90) //NUM LOCK key
	VK_SCROLL  = VirtualKey(0x91) //SCROLL LOCK key
	// -	0x92-96	OEM specific
	// -	0x97-9F	Unassigned
	VK_LSHIFT              = VirtualKey(0xA0) //Left SHIFT key
	VK_RSHIFT              = VirtualKey(0xA1) //Right SHIFT key
	VK_LCONTROL            = VirtualKey(0xA2) //Left CONTROL key
	VK_RCONTROL            = VirtualKey(0xA3) //Right CONTROL key
	VK_LMENU               = VirtualKey(0xA4) //Left ALT key
	VK_RMENU               = VirtualKey(0xA5) //Right ALT key
	VK_BROWSER_BACK        = VirtualKey(0xA6) //Browser Back key
	VK_BROWSER_FORWARD     = VirtualKey(0xA7) //Browser Forward key
	VK_BROWSER_REFRESH     = VirtualKey(0xA8) //Browser Refresh key
	VK_BROWSER_STOP        = VirtualKey(0xA9) //Browser Stop key
	VK_BROWSER_SEARCH      = VirtualKey(0xAA) //Browser Search key
	VK_BROWSER_FAVORITES   = VirtualKey(0xAB) //Browser Favorites key
	VK_BROWSER_HOME        = VirtualKey(0xAC) //Browser Start and Home key
	VK_VOLUME_MUTE         = VirtualKey(0xAD) //Volume Mute key
	VK_VOLUME_DOWN         = VirtualKey(0xAE) //Volume Down key
	VK_VOLUME_UP           = VirtualKey(0xAF) //Volume Up key
	VK_MEDIA_NEXT_TRACK    = VirtualKey(0xB0) //Next Track key
	VK_MEDIA_PREV_TRACK    = VirtualKey(0xB1) //Previous Track key
	VK_MEDIA_STOP          = VirtualKey(0xB2) //Stop Media key
	VK_MEDIA_PLAY_PAUSE    = VirtualKey(0xB3) //Play/Pause Media key
	VK_LAUNCH_MAIL         = VirtualKey(0xB4) //Start Mail key
	VK_LAUNCH_MEDIA_SELECT = VirtualKey(0xB5) //Select Media key
	VK_LAUNCH_APP1         = VirtualKey(0xB6) //Start Application 1 key
	VK_LAUNCH_APP2         = VirtualKey(0xB7) //Start Application 2 key
	// -	0xB8-B9	Reserved
	VK_OEM_1      = VirtualKey(0xBA) //Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the ;: key
	VK_OEM_PLUS   = VirtualKey(0xBB) //For any country/region, the + key
	VK_OEM_COMMA  = VirtualKey(0xBC) //For any country/region, the , key
	VK_OEM_MINUS  = VirtualKey(0xBD) //For any country/region, the - key
	VK_OEM_PERIOD = VirtualKey(0xBE) //For any country/region, the . key
	VK_OEM_2      = VirtualKey(0xBF) //Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the /? key
	VK_OEM_3      = VirtualKey(0xC0) //Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the `~ key
	// -	0xC1-DA	Reserved
	VK_OEM_4 = VirtualKey(0xDB) //Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the [{ key
	VK_OEM_5 = VirtualKey(0xDC) //Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the \\| key
	VK_OEM_6 = VirtualKey(0xDD) //Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the ]} key
	VK_OEM_7 = VirtualKey(0xDE) //Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the '" key
	VK_OEM_8 = VirtualKey(0xDF) //Used for miscellaneous characters; it can vary by keyboard.
	// -	0xE0	Reserved
	// -	0xE1	OEM specific
	VK_OEM_102 = VirtualKey(0xE2) //The <> keys on the US standard keyboard, or the \\| key on the non-US 102-key keyboard
	// -	0xE3-E4	OEM specific
	VK_PROCESSKEY = VirtualKey(0xE5) //IME PROCESS key
	// -	0xE6	OEM specific
	VK_PACKET = VirtualKey(0xE7) //Used to pass Unicode characters as if they were keystrokes. The VK_PACKET= VirtualKet(key)//is the low word of a 32-bit Virtual Key value used for non-keyboard input methods. For more information, see Remark in KEYBDINPUT, SendInput, WM_KEYDOWN, and WM_KEYUP
	// -	0xE8	Unassigned
	// -	0xE9-F5	OEM specific
	VK_ATTN      = VirtualKey(0xF6) //Attn key
	VK_CRSEL     = VirtualKey(0xF7) //CrSel key
	VK_EXSEL     = VirtualKey(0xF8) //ExSel key
	VK_EREOF     = VirtualKey(0xF9) //Erase EOF key
	VK_PLAY      = VirtualKey(0xFA) //Play key
	VK_ZOOM      = VirtualKey(0xFB) //Zoom key
	VK_NONAME    = VirtualKey(0xFC) //Reserved
	VK_PA1       = VirtualKey(0xFD) //PA1 key
	VK_OEM_CLEAR = VirtualKey(0xFE) //Clear key
)

type KeyEventFlag int32

const (
	KEYEVENTF_KEYPRESS = KeyEventFlag(0x0000)
	KEYEVENTF_KEYUP    = KeyEventFlag(0x0002)
)
