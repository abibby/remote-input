package vigem

import (
	"errors"
	"log"
)

type lazyProc struct {
	Name string
}

func newProc(name string) *lazyProc {
	return &lazyProc{
		Name: name,
	}
}

func (p *lazyProc) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	log.Printf("%s.Call(%v)", p.Name, a)
	return 0, 0, windowsERROR_SUCCESS
}

var (
	procAlloc                            = newProc("vigem_alloc")
	procFree                             = newProc("vigem_free")
	procConnect                          = newProc("vigem_connect")
	procDisconnect                       = newProc("vigem_disconnect")
	procTargetAdd                        = newProc("vigem_target_add")
	procTargetFree                       = newProc("vigem_target_free")
	procTargetRemove                     = newProc("vigem_target_remove")
	procTargetX360Alloc                  = newProc("vigem_target_x360_alloc")
	procTargetX360RegisterNotification   = newProc("vigem_target_x360_register_notification")
	procTargetX360UnregisterNotification = newProc("vigem_target_x360_unregister_notification")
	procTargetX360Update                 = newProc("vigem_target_x360_update")
)

var windowsERROR_SUCCESS = errors.New("success")

func windowsNewCallback(fn interface{}) uintptr {
	log.Printf("windowsNewCallback(%v)", fn)
	return 0
}
