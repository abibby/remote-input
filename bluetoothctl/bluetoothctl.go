package bluetoothctl

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Device struct {
	Type    string
	Address string
	Name    string
}

func run(args ...string) (io.ReadCloser, io.ReadCloser, error) {
	fmt.Printf("bluetoothctl %s\n", strings.Join(args, " "))
	cmd := exec.Command("bluetoothctl", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdout.Close()
		return nil, nil, err
	}
	err = cmd.Start()
	if err != nil {
		stdout.Close()
		stderr.Close()
		return nil, nil, err
	}
	return stdout, stderr, nil
}

func Scan() (chan *Device, chan error) {
	devices := make(chan *Device, 2)
	errs := make(chan error)
	go func() {
		defer func() {
			close(devices)
			close(errs)
		}()

		stdout, stderr, err := run("scan", "on")
		if err != nil {
			errs <- err
			return
		}
		defer func() {
			stdout.Close()
			stderr.Close()
		}()
		io.Copy(os.Stdout, stderr)
		// s := bufio.NewScanner(stderr)
		// for s.Scan() {
		// 	line := s.Text()
		// 	devices <- &Device{
		// 		Name: line,
		// 	}
		// }
	}()
	return devices, errs
}
