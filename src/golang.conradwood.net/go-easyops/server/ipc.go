package server

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/utils"
	ad "golang.yacloud.eu/apis/autodeployer2"
	"golang.yacloud.eu/unixipc"
	"strconv"
	"sync"
)

var (
	enable_ipc  = flag.Bool("ge_enable_ipc", true, "enable the internal ipc between code and autodeployer")
	ipc_fd_env  = cmdline.ENV("GE_AUTODEPLOYER_IPC_FD", "if set it is assumed to be a filedescriptor over which an IPC can be initiated with the autodeployer")
	ipc_lock    sync.Mutex
	ipc_started = false
	srv         *unixipc.IPCServer
)

func ipc_enabled() bool {
	if *enable_ipc {
		return true
	}
	return false
}

func start_ipc() {
	if !ipc_enabled() {
		return
	}
	ipc_lock.Lock()
	defer ipc_lock.Unlock()
	if ipc_started {
		return
	}
	ipc_started = true
	if ipc_fd_env.Value() == "" {
		//fmt.Printf("[go-easyops] no ipc fd\n")
		return
	}
	fd, err := strconv.Atoi(ipc_fd_env.Value())
	if err != nil {
		panic(fmt.Sprintf("GE_AUTODEPLOYER_IPC_FD invalid value: %s", err))
	}
	srv, err = unixipc.NewConnectedServer(fd)
	if err != nil {
		panic(fmt.Sprintf("failed to start autodeployer IPC: %s", err))
	}
}
func ipc_send_startup(sd *serverDef) error {
	if !ipc_enabled() || !ipc_started {
		return nil
	}
	proto_payload := &ad.INTRPCStartup{
		ServiceName: sd.name,
		Port:        uint32(sd.port),
		Healthz:     health,
	}
	payload, err := utils.MarshalBytes(proto_payload)
	if err != nil {
		return err
	}
	_, err = srv.Send("startup", payload)
	if err != nil {
		return err
	}
	return nil
}
func ipc_send_health(sd *serverDef, h common.Health) error {
	if !ipc_enabled() || !ipc_started {
		return nil
	}
	proto_payload := &ad.INTRPCHealthz{
		ServiceName: sd.name,
		Port:        uint32(sd.port),
		Healthz:     h,
	}
	payload, err := utils.MarshalBytes(proto_payload)
	if err != nil {
		return err
	}
	_, err = srv.Send("healthz", payload)
	if err != nil {
		return err
	}
	return nil
}
