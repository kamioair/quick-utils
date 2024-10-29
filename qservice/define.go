package qservice

import "github.com/liaozhibinair/quick-utils/qdefine"

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

func NewSetting(module string, host string, version string, onReqHandler ReqHandler) Setting {
	return Setting{
		Module:       module,
		Host:         Host{Addr: host},
		Version:      version,
		OnReqHandler: onReqHandler,
	}
}

type ReqHandler func(route string, ctx qdefine.Context) (any, error)
