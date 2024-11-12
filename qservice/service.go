package qservice

import (
	"errors"
	"fmt"
	"github.com/liaozhibinair/quick-utils/qdefine"
	"github.com/liaozhibinair/quick-utils/qio"
	"github.com/liaozhibinair/quick-utils/qlauncher"
	easyCon "github.com/qiu-tec/easy-con.golang"
	"os"
	"strconv"
	"strings"
	"time"
)

type MicroService struct {
	Module  string
	adapter easyCon.IAdapter
	setting Setting
}

func NewService(setting Setting) *MicroService {
	// 修改系统路径为当前目录
	cd, err := qio.GetCurrentDirectory()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(qio.GetDirectory(cd))
	if err != nil {
		panic(err)
	}

	// 创建服务
	serv := &MicroService{
		Module:  setting.Module,
		setting: setting,
	}

	// 初始化Api适配器
	apiSetting := easyCon.NewSetting(setting.Module, setting.Broker.Addr, serv.onReq, serv.onStatusChanged)
	apiSetting.OnNotice = serv.OnNotice
	apiSetting.UID = setting.Broker.UId
	apiSetting.PWD = setting.Broker.Pwd
	apiSetting.TimeOut = time.Duration(setting.Broker.TimeOut) * time.Second
	apiSetting.ReTry = setting.Broker.Retry
	apiSetting.LogMode = easyCon.ELogMode(setting.Broker.LogMode)
	serv.adapter = easyCon.NewMqttAdapter(apiSetting)

	return serv
}

func (serv *MicroService) Run() {
	qlauncher.Run(serv.onStart, serv.onStop)
}

func (serv *MicroService) SendRequest(module, route string, params any) (qdefine.Context, error) {
	var resp easyCon.PackResp

	if strings.Contains(module, "/") {
		// 路由请求
		newParams := map[string]any{}
		newParams["Module"] = module
		newParams["Route"] = route
		newParams["Content"] = params
		resp = serv.adapter.Req("Route", "Request", newParams)
	} else {
		// 常规请求
		resp = serv.adapter.Req(module, route, params)
	}
	if resp.RespCode == easyCon.ERespSuccess {
		// 返回成功
		return newControlResp(resp)
	}
	// 返回异常
	if resp.RespCode == easyCon.ERespTimeout {
		return nil, errors.New(fmt.Sprintf("%v:%s", resp.RespCode, "request timeout"))
	}
	if resp.RespCode == easyCon.ERespRouteNotFind {
		return nil, errors.New(fmt.Sprintf("%v:%s", resp.RespCode, "request route not find"))
	}
	if resp.RespCode == easyCon.ERespForbidden {
		return nil, errors.New(fmt.Sprintf("%v:%s", resp.RespCode, "request forbidden"))
	}
	return nil, errors.New(fmt.Sprintf("%v:%s,%s", resp.RespCode, resp.Content, resp.Error))
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

func (serv *MicroService) onReq(pack easyCon.PackReq) (code easyCon.EResp, resp any) {
	defer errRecover(func(err string) {
		code = easyCon.ERespError
		resp = err
		// 记录日志

	})

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
		ctx, err1 := newControlReq(pack)
		if err1 != nil {
			return easyCon.ERespError, err1.Error()
		}
		rs, err2 := serv.setting.OnReqHandler(pack.Route, ctx)
		if err2 != nil {
			c, _ := strconv.Atoi(err2.Error())
			switch c {
			case int(easyCon.ERespBadReq):
				return easyCon.ERespBadReq, "request bad"
			case int(easyCon.ERespRouteNotFind):
				return easyCon.ERespRouteNotFind, "request route not find"
			case int(easyCon.ERespForbidden):
				return easyCon.ERespForbidden, "request forbidden"
			case int(easyCon.ERespTimeout):
				return easyCon.ERespTimeout, "request timeout"
			default:
				return easyCon.ERespError, err2.Error()
			}
		}
		// 执行成功，返回结果
		return easyCon.ERespSuccess, rs
	}
	return easyCon.ERespRouteNotFind, "Route Not Matched"
}

func (serv *MicroService) OnNotice(notice easyCon.PackNotice) {
	defer errRecover(func(err string) {
		// 记录日志
	})

	// 外置方法
	if serv.setting.OnNoticeHandler != nil {
		ctx, err := newControlNotice(notice)
		if err != nil {
			panic(err)
		}
		serv.setting.OnNoticeHandler(notice.Route, ctx)
	}
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
