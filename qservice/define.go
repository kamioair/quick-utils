package qservice

import (
	"encoding/json"
	"github.com/kamioair/quick-utils/qconfig"
	"github.com/kamioair/quick-utils/qdefine"
	"github.com/kamioair/quick-utils/qio"
	"os"
	"strings"
)

// Setting 模块配置
type Setting struct {
	Module          string        // 模块服务名称
	Desc            string        // 模块服务描述
	Version         string        // 模块服务版本
	Broker          BrokerConfig  // 主服务配置
	onInitHandler   InitHandler   // 初始化回调
	onReqHandler    ReqHandler    // 请求回调
	onNoticeHandler NoticeHandler // 通知回调
	onStateHandler  StateHandler  // 状态回调
	deviceCode      string        // 设备码
}

// BrokerConfig 主服务配置
type BrokerConfig struct {
	Addr    string
	UId     string
	Pwd     string
	LogMode string
	TimeOut int
	Retry   int
}

type (
	// InitHandler 初始化回调
	InitHandler func()
	// ReqHandler 请求回调
	ReqHandler func(route string, ctx qdefine.Context) (any, error)
	// NoticeHandler 通知回调
	NoticeHandler func(route string, ctx qdefine.Context)
	// StateHandler 状态回调
	StateHandler func(state qdefine.ECommState)
)

// NewSetting 创建模块配置
func NewSetting(moduleName, moduleDesc, version string) *Setting {
	// 修改系统路径为当前目录
	err := os.Chdir(qio.GetCurrentDirectory())
	if err != nil {
		panic(err)
	}

	// 默认值
	module := moduleName
	broker := BrokerConfig{
		Addr:    qconfig.Get(module, "mqtt.addr", "ws://127.0.0.1:5002/ws"),
		UId:     qconfig.Get(module, "mqtt.username", ""),
		Pwd:     qconfig.Get(module, "mqtt.password", ""),
		LogMode: qconfig.Get(module, "mqtt.logMode", "CONSOLE"),
		TimeOut: qconfig.Get(module, "mqtt.timeOut", 3000),
		Retry:   qconfig.Get(module, "mqtt.retry", 3),
	}
	devCode := ""
	// 根据传参更新配置
	if len(os.Args) > 1 {
		args := args{}
		err = json.Unmarshal([]byte(os.Args[1]), &args)
		if err != nil {
			panic(err)
		}
		devCode = args.DeviceCode
		if args.Module != "" {
			module = args.Module
		}
		if args.MqConfig != nil {
			broker = *args.MqConfig
		}
		if devCode != "" {
			module = module + "." + devCode
		}
	}
	// 返回配置
	return &Setting{
		Module:     module,
		Desc:       moduleDesc,
		Version:    version,
		Broker:     broker,
		deviceCode: devCode,
	}
}

func (s *Setting) BindInitFunc(onInitHandler InitHandler) *Setting {
	s.onInitHandler = onInitHandler
	return s
}

func (s *Setting) BindReqFunc(onReqHandler ReqHandler) *Setting {
	s.onReqHandler = onReqHandler
	return s
}

func (s *Setting) BindNoticeFunc(onNoticeHandler NoticeHandler) *Setting {
	s.onNoticeHandler = onNoticeHandler
	return s
}

func (s *Setting) BindStateFunc(onStateHandler StateHandler) *Setting {
	s.onStateHandler = onStateHandler
	return s
}

func (s *Setting) SetDeviceCode(code string) *Setting {
	s.deviceCode = code
	sp := strings.Split(s.Module, ".")
	if len(sp) >= 2 {
		s.Module = sp[0] + "." + code
	} else {
		s.Module = s.Module + "." + code
	}
	s.Module = strings.Trim(s.Module, ".")
	return s
}

type args struct {
	Module     string
	DeviceCode string
	MqConfig   *BrokerConfig
}

func writeErrLog(tp string, err string) {

}
