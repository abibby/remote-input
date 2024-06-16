package main

import (
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
)

func main() {
	// conn, err := dbus.ConnectSystemBus()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "Failed to connect to SystemBus bus:", err)
	// 	os.Exit(1)
	// }
	// defer conn.Close()

	// var s string
	// err = conn.Object("org.bluez", "/").Call("org.freedesktop.DBus.Introspectable.Introspect", 0).Store(&s)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "Failed to introspect bluez", err)
	// 	os.Exit(1)
	// }

	// fmt.Println(s)
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	if err = conn.AddMatchSignal(
		// dbus.WithMatchObjectPath("/org/freedesktop/DBus"),
		// dbus.WithMatchInterface("org.bluez"),
		dbus.WithMatchSender("org.bluez"),
	); err != nil {
		panic(err)
	}

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)
	for v := range c {
		fmt.Printf("%#v\n", v)
	}
}
