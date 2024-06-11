package vigem

import "golang.org/x/sys/windows"

var (
	client = windows.NewLazyDLL("ViGEmClient.dll")

	procAlloc                            = client.NewProc("vigem_alloc")
	procFree                             = client.NewProc("vigem_free")
	procConnect                          = client.NewProc("vigem_connect")
	procDisconnect                       = client.NewProc("vigem_disconnect")
	procTargetAdd                        = client.NewProc("vigem_target_add")
	procTargetFree                       = client.NewProc("vigem_target_free")
	procTargetRemove                     = client.NewProc("vigem_target_remove")
	procTargetX360Alloc                  = client.NewProc("vigem_target_x360_alloc")
	procTargetX360RegisterNotification   = client.NewProc("vigem_target_x360_register_notification")
	procTargetX360UnregisterNotification = client.NewProc("vigem_target_x360_unregister_notification")
	procTargetX360Update                 = client.NewProc("vigem_target_x360_update")
)

var windowsERROR_SUCCESS = windows.ERROR_SUCCESS

var windowsNewCallback = windows.NewCallback
