package handlers

import "tinygo.org/x/bluetooth"

func ParseAddress(s string) (bluetooth.Address, error) {
	uuid, err := bluetooth.ParseUUID(s)
	if err != nil {
		return bluetooth.Address{}, err
	}
	return bluetooth.Address{UUID: uuid}, nil
}
