package handlers

import "tinygo.org/x/bluetooth"

func ParseAddress(s string) (bluetooth.Address, error) {
	mac, err := bluetooth.ParseMAC(s)
	if err != nil {
		return bluetooth.Address{}, err
	}
	return bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: mac}}, nil
}
