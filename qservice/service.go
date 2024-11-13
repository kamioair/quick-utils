package qservice

import (
	"errors"
	"fmt"
	"github.com/kamioair/quick-utils/qdefine"
	"github.com/kamioair/quick-utils/qio"
	"github.com/kamioair/quick-utils/qlauncher"
	easyCon "github.com/qiu-tec/easy-con.golang"
	"os"
	"strconv"
	"strings"
	"time"
)

type MicroService struct {
	Module  string
	adapter easyCon.IAdapter
	setting *Setting
}

// NewService 创建服务
func NewService(setting *Setting) *MicroService {
	// 修改系统路径为当前目录
	err := os.Chdir(qio.GetCurrentDirectory())
	if err != nil {
		panic(err)
	}

	// 创建服务
	serv := &MicroService{
		Module:  setting.Module,
		setting: setting,
	}

	// 启动访问器
	serv.initAdapter()

	return serv
}

// Run 启动服务
func (serv *MicroService) Run() {
	qlauncher.Run(serv.onStart, serv.onStop)
}

// ResetClient 重置客户端
func (serv *MicroService) ResetClient(code string) {
	serv.setting.deviceCode = code
	module := serv.setting.Module
	sp := strings.Split(module, ".")
	if len(sp) >= 2 {
		module = sp[0] + "." + code
	} else {
		module = module + "." + code
	}
	module = strings.Trim(module, ".")
	serv.setting.Module = module
	serv.Module = module

	// 重新创建服务
	serv.initAdapter()
}

// SendRequest 发送请求
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

// SendNotice 发送通知
func (serv *MicroService) SendNotice(route string, content any) {
	err := serv.adapter.SendNotice(route, content)
	if err != nil {
		serv.SendLog("error", "Service Send Notice Error", err)
	}
}

// SendLog 发送日志
func (serv *MicroService) SendLog(logType qdefine.ELog, content string, err error) {
	switch logType {
	case qdefine.ELogError:
		serv.adapter.Err(content, err)
	case qdefine.ELogWarn:
		serv.adapter.Warn(content)
	case qdefine.ELogDebug:
		serv.adapter.Debug(content)
	default:
		serv.adapter.Debug(content)
	}
}

func (serv *MicroService) initAdapter() {
	// 先停止
	if serv.adapter != nil {
		serv.adapter.Stop()
		serv.adapter = nil
	}
	// 重新创建
	apiSetting := easyCon.NewSetting(serv.setting.Module, serv.setting.Broker.Addr, serv.onReq, serv.onStatusChanged)
	apiSetting.OnNotice = serv.onNotice
	apiSetting.UID = serv.setting.Broker.UId
	apiSetting.PWD = serv.setting.Broker.Pwd
	apiSetting.TimeOut = time.Duration(serv.setting.Broker.TimeOut) * time.Second
	apiSetting.ReTry = serv.setting.Broker.Retry
	apiSetting.LogMode = easyCon.ELogMode(serv.setting.Broker.LogMode)
	serv.adapter = easyCon.NewMqttAdapter(apiSetting)
}

func (serv *MicroService) onReq(pack easyCon.PackReq) (code easyCon.EResp, resp any) {
	defer errRecover(func(err string) {
		code = easyCon.ERespError
		resp = err
		// 记录日志
		writeErrLog("service.onReq", err)
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
	if serv.setting.onReqHandler != nil {
		ctx, err1 := newControlReq(pack)
		if err1 != nil {
			return easyCon.ERespError, err1.Error()
		}
		rs, err2 := serv.setting.onReqHandler(pack.Route, ctx)
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

func (serv *MicroService) onNotice(notice easyCon.PackNotice) {
	defer errRecover(func(err string) {
		// 记录日志
		writeErrLog("service.onNotice", err)
	})

	// 外置方法
	if serv.setting.onNoticeHandler != nil {
		ctx, err := newControlNotice(notice)
		if err != nil {
			panic(err)
		}
		serv.setting.onNoticeHandler(notice.Route, ctx)
	}
}

func (serv *MicroService) onStatusChanged(adapter easyCon.IAdapter, status easyCon.EStatus) {
	//if status == easyCon.EStatusLinkLost {
	//	adapter.Reset()
	//}
	if serv.setting.onStateHandler != nil {
		sn := qdefine.ECommState(status)
		serv.setting.onStateHandler(sn)
	}
}

func (serv *MicroService) onStart() {
	if serv.setting.onInitHandler != nil {
		serv.setting.onInitHandler()
	}
}

func (serv *MicroService) onStop() {
	if serv.adapter != nil {
		serv.adapter.Stop()
		serv.adapter = nil
	}
}
