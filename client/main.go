package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/config"
	"github.com/kardianos/service"
)

type program struct {
	logger service.Logger
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	p.logger.Info("started")

	serverIP := fmt.Sprintf("%s:%d", config.Host, config.Port)

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		p.logger.Error(err)
		return
	}
	defer conn.Close()

	joysticks, err = NewJoysticks()
	if err != nil {
		p.logger.Warning("Joystick setup failed: %v", err)
	}
	defer joysticks.Close()

	p.logger.Infof("connected to %s", serverIP)

	events := make([]common.InputEvent, 0, 8)

	b := make([]byte, 24)
	for {
		events = append(events, common.InputEvent{})
		e := &events[len(events)-1]

		_, err := conn.Read(b)
		if err == io.EOF {
			p.logger.Info("disconnected")
			return
		} else if err != nil {
			p.logger.Warning(err)
			continue
		}

		err = e.UnmarshalBinary(b)
		if err != nil {
			p.logger.Warning(err)
			continue
		}

		if e.EventType == common.EV_SYN {
			err = handleEvent(events)
			if err != nil {
				p.logger.Warning(err)
			}
			events = events[:0]
		}
	}
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "remote-input-client",
		DisplayName: "Remote Input Client",
		Description: "Client to receive inputs from a host server",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			err = s.Install()
		case "uninstall":
			err = s.Uninstall()
		default:
			err = fmt.Errorf("unknown command %s, expected install or uninstall", os.Args[1])
		}
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	prg.logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		prg.logger.Error(err)
	}
}
