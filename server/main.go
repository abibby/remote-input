package main

import (
	"log"
	"net/http"
	"time"

	"github.com/abibby/remote-input/windows"
)

func greet(w http.ResponseWriter, r *http.Request) {
	go func() {
		// time.Sleep(time.Second)
		err := windows.SendInput(windows.VK_A, windows.KEYEVENTF_KEYPRESS)
		if err != nil {
			log.Print(err)
		}
		time.Sleep(time.Second)
		err = windows.SendInput(windows.VK_A, windows.KEYEVENTF_KEYUP)
		if err != nil {
			log.Print(err)
		}
	}()
}

func main() {
	http.HandleFunc("/", greet)
	http.ListenAndServe(":8080", nil)
}
