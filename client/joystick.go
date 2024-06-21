package main

import (
	"errors"

	"github.com/abibby/remote-input/vigem"
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
