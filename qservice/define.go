package qservice

import (
	"github.com/liaozhibinair/quick-utils/qconfig"
	"github.com/liaozhibinair/quick-utils/qdefine"
	"github.com/liaozhibinair/quick-utils/qio"
	"os"
)

type Setting struct {
	Module          string        // 模块服务名称
	Version         string        // 模块服务版本
	Broker          BrokerConfig  // 主服务配置
	OnReqHandler    ReqHandler    // 请求回调
	OnNoticeHandler NoticeHandler // 通知回调
	ExitProcess     string        // 如果监听的进程不存在，则立即退出
}

type BrokerConfig struct {
	Addr    string
	UId     string
	Pwd     string
	LogMode string
	TimeOut int
	Retry   int
}

type ReqHandler func(route string, ctx qdefine.Context) (any, error)

type NoticeHandler func(route string, ctx qdefine.Context)

func NewSetting(defModule string, version string, onReqHandler ReqHandler, onNoticeHandler NoticeHandler) Setting {
	// 修改系统路径为当前目录
	cd, err := qio.GetCurrentDirectory()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(qio.GetDirectory(cd))
	if err != nil {
		panic(err)
	}

	module := defModule

	// 根据默认值和启动器传值，获取模块名称
	args := getArgs(defModule)
	if args[0] != "" {
		module = args[0]
	}
	setting := Setting{
		Module: module,
		Broker: BrokerConfig{
			Addr:    qconfig.Get(module, "mqtt.addr", "ws://127.0.0.1:5002/ws"),
			UId:     qconfig.Get(module, "mqtt.username", ""),
			Pwd:     qconfig.Get(module, "mqtt.password", ""),
			LogMode: qconfig.Get(module, "mqtt.logMode", "NONE"),
			TimeOut: qconfig.Get(module, "mqtt.timeOut", 3000),
			Retry:   qconfig.Get(module, "mqtt.retry", 3),
		},
		Version:         version,
		OnReqHandler:    onReqHandler,
		OnNoticeHandler: onNoticeHandler,
	}
	return setting
}

func getArgs(defArgs ...string) []string {
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
