package qservice

import (
	"github.com/liaozhibinair/quick-utils/qconfig"
	"github.com/liaozhibinair/quick-utils/qdefine"
	"github.com/liaozhibinair/quick-utils/qio"
	"github.com/liaozhibinair/quick-utils/qlauncher"
	easyCon "github.com/qiu-tec/easy-con.golang"
	"os"
	"strings"
	"time"
)

type MicroService struct {
	Module  string
	adapter easyCon.IAdapter
	setting Setting
}

func New(module string, version string, onReqHandler ReqHandler, configContent []byte) *MicroService {
	// 修改系统路径为当前目录
	cd, err := qio.GetCurrentDirectory()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(qio.GetDirectory(cd))
	if err != nil {
		panic(err)
	}

	// 初始化配置
	qconfig.Init("./config/config.yaml", configContent)

	setting := Setting{
		Module: module,
		Host: Host{
			Addr: qconfig.Get(module, "mqtt.addr", "ws://127.0.0.1:5002/ws"),
			UId:  qconfig.Get(module, "mqtt.username", ""),
			Pwd:  qconfig.Get(module, "mqtt.password", ""),
		},
		Version:      version,
		OnReqHandler: onReqHandler,
	}
	serv := &MicroService{
		Module:  module,
		setting: setting,
	}

	// 初始化Api适配器
	apiSetting := easyCon.NewSetting(setting.Module, setting.Host.Addr, serv.onReq, serv.onStatusChanged)
	apiSetting.UID = setting.Host.UId
	apiSetting.PWD = setting.Host.Pwd
	serv.adapter = easyCon.NewMqttAdapter(apiSetting)

	return serv
}

func (serv *MicroService) Run() {
	qlauncher.Run(serv.onStart, serv.onStop)
}

func (serv *MicroService) SendRequest(module, route string, params any) (qdefine.Context, error) {
	_ = serv.adapter.Req(module, route, params)
	return nil, nil
}

func (serv *MicroService) SendNotice(route string, content any) {
	err := serv.adapter.SendNotice(route, content)
	if err != nil {
		serv.SendLog("error", "Service Send Notice Error", err)
	}
}

func (serv *MicroService) SendLog(logType string, content string, err error) {
	switch strings.ToLower(logType) {
	case "error", "err":
		serv.adapter.Err(content, err)
	case "warn":
		serv.adapter.Warn(content)
	default:
		serv.adapter.Debug(content)
	}
}

func (serv *MicroService) onReq(pack easyCon.PackReq) (easyCon.EResp, any) {
	switch pack.Route {
	case "Exit":
		serv.onStop()
		go func() {
			time.Sleep(time.Millisecond * 100)
			qlauncher.Exit()
		}()
		return easyCon.ERespSuccess, nil
	case "Reset":
		serv.adapter.Reset()
		return easyCon.ERespSuccess, nil
	}
	if serv.setting.OnReqHandler != nil {
		ctx, err1 := newControl(pack)
		if err1 != nil {
			return easyCon.ERespError, err1
		}
		rs, err2 := serv.setting.OnReqHandler(pack.Route, ctx)
		if err2 != nil {
			return easyCon.ERespError, err2
		}
		// 执行成功，返回结果
		return easyCon.ERespSuccess, rs
	}
	return easyCon.ERespRouteNotFind, "Route Not Matched"
}

func (serv *MicroService) onStatusChanged(adapter easyCon.IAdapter, status easyCon.EStatus) {
	//if serv.setting.OnStatusChangedHandler != nil {
	//	serv.setting.OnStatusChangedHandler(adapter, status)
	//}
}

func (serv *MicroService) onStart() {

}

func (serv *MicroService) onStop() {
	if serv.adapter != nil {
		serv.adapter.Stop()
		serv.adapter = nil
	}
}
