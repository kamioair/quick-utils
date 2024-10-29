package qservice

import (
	"github.com/liaozhibinair/quick-utils/qdefine"
	"github.com/liaozhibinair/quick-utils/qlauncher"
	easyCon "github.com/qiu-tec/easy-con.golang"
	"strings"
	"time"
)

type MicroService struct {
	adapter easyCon.IAdapter
	setting Setting
}

func New(setting Setting) *MicroService {
	serv := &MicroService{
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
	//if serv.setting.OnReqHandler != nil {
	//	return serv.setting.OnReqHandler(pack)
	//}
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
