package qservice

import (
	"github.com/liaozhibinair/quick-utils/qconfig"
	"github.com/liaozhibinair/quick-utils/qdefine"
	"os"
)

type Setting struct {
	Module       string     // 模块服务名称
	Version      string     // 模块服务版本
	Host         Host       // 主服务配置
	OnReqHandler ReqHandler // 请求回调
	ExitProcess  string     // 如果监听的进程不存在，则立即退出
}

type Host struct {
	Addr string
	UId  string
	Pwd  string
}

func NewSetting(module string, version string, onReqHandler ReqHandler) Setting {
	return Setting{
		Module: module,
		Host: Host{
			Addr: qconfig.Get(module, "mqtt.addr", "ws://127.0.0.1:5002/ws"),
			UId:  qconfig.Get(module, "mqtt.username", ""),
			Pwd:  qconfig.Get(module, "mqtt.password", ""),
		},
		Version:      version,
		OnReqHandler: onReqHandler,
	}
}

type ReqHandler func(route string, ctx qdefine.Context) (any, error)

func GetArgs(defArgs ...string) []string {
	args := make([]string, len(defArgs))
	for i := 0; i < len(defArgs); i++ {
		if len(os.Args) > i+1 {
			args[i] = os.Args[i+1]
		} else {
			args[i] = defArgs[i]
		}
	}
	return args
}
