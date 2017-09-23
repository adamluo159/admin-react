package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/adamluo159/admin-react/protocol"
	"github.com/adamluo159/admin-react/utils"
)

//生成agent实例
func New(cfPath string) {

	h, err := os.Hostname()
	if err != nil {
		log.Fatal("get host err", err)
	}
	LoadConfig(cfPath)
	conf.RemoteConfDir += h

	a := agent{
		srvs:     make(map[string]*ServiceInfo),
		msgMap:   make(map[uint32]func([]byte)),
		SvnVer:   SvnInfo(),
		hostName: h,
	}

	UpdateGameConf()
	a.LoadServices()

	go a.RegularlyCheckProcess()
	go a.GetOnlinePlayers()

	a.Connect()
}

func (a *agent) Connect() {
	for {
		conn, err := net.Dial("tcp", conf.ConAddress)
		if err == nil {
			log.Println("connect to agent server sucess:", conf.ConAddress)
			a.conn = &conn
			a.OnMessage()
		} else {

			log.Println("connect to agent server fail :", conf.ConAddress)
		}

		time.Sleep(5 * time.Second)
	}

}

func (a *agent) GetOnlinePlayers() {
	for {
		rt := OnlineReport{
			Method:    "Backend.onlineData",
			TimeStamp: time.Now().Unix(),
			Param:     []ZoneOnline{},
		}
		n := 0
		zo := ZoneOnline{}
		sigkey := conf.HttpKey + "method=Backend.onlineDataparam"
		for _, v := range a.srvs {
			count := 0
			if v.ClientPorts != nil {
				for _, p := range *v.ClientPorts {
					port := strconv.Itoa(p)
					count += OnlinePlayers(port)
				}
				zo.Zid = v.Zid
				zo.OnlineNum = count
				sigkey += "=" + strconv.Itoa(n) + "=onlineNum=" + strconv.Itoa(count) + "zid=" + strconv.Itoa(v.Zid)
				log.Println(v.Sname, "online players : ", count)
				rt.Param = append(rt.Param, zo)
			}
		}

		timeStamp := time.Now().Unix()
		sigkey += "timestamp=" + strconv.FormatInt(timeStamp, 10)

		if jdata, err := json.Marshal(&rt); err == nil {
			webStr := fmt.Sprintf(conf.HttpStr, string(jdata), utils.CreateMd5(sigkey))
			if resp, err := http.Get(webStr); err != nil {
				log.Println("http get err", err)
			} else {
				body, err := ioutil.ReadAll(resp.Body)
				log.Println("http result:", string(body), err)
			}
		}

		time.Sleep(3 * time.Minute)
	}
}

//加载本地游戏配置
func (a *agent) LoadServices() {
	hostDir := conf.LocalConfDir + a.hostName

	dir, err := ioutil.ReadDir(hostDir)
	if err != nil {
		log.Println("LoadServices, read dir err, ", err.Error())
	}
	for index := 0; index < len(dir); index++ {
		a.InitSrv(dir[index].Name())

	}
}

func (a *agent) OnMessage() {

	a.msgMap[protocol.CmdToken] = a.S2cCheckRsp
	a.msgMap[protocol.CmdStartZone] = a.S2cStartZone
	a.msgMap[protocol.CmdStopZone] = a.S2cStopZone
	a.msgMap[protocol.CmdUpdateHost] = a.S2cUpdateZoneConfig
	a.msgMap[protocol.CmdStartHostZone] = a.S2cStartHostZones
	a.msgMap[protocol.CmdStopHostZone] = a.S2cStopHostZones
	a.msgMap[protocol.CmdUpdateSvn] = a.S2cUpdateSvn

	// 消息缓冲
	msgbuf := bytes.NewBuffer(make([]byte, 0, 1024))
	// 数据缓冲
	databuf := make([]byte, 1024)
	// 消息长度
	length := 0

	a.C2sCheckReq()
	conn := a.conn

	for {
		// 读取数据
		n, err := (*conn).Read(databuf)
		if err == io.EOF {
			log.Printf("Client exit: %s\n", (*conn).RemoteAddr())
		}
		if err != nil {
			log.Printf("Read error: %s\n", err)
			return
		}
		// 数据添加到消息缓冲
		n, err = msgbuf.Write(databuf[:n])
		if err != nil {
			log.Printf("Buffer write error: %s\n", err)
			return
		}
		// 消息分割循环
		for {
			cmd, data := protocol.UnPacket(&length, msgbuf)
			if cmd <= 0 {
				break
			}
			mfunc := a.msgMap[cmd]
			if mfunc == nil {
				log.Printf("cannt find msg handle server cmd: %d data: %s\n", cmd, string(data))
			} else {
				mfunc(data)
				log.Printf("server cmd: %d data: %s\n", cmd, string(data))
			}
		}
	}
}

func (a *agent) ReadGameJson(name string) (*[]int, int) {
	fpath := conf.LocalConfDir + a.hostName + "/" + name + "/gameJsonConf"
	cp := struct {
		ClientPorts []int
		Zid         int
	}{}
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, 0
	}

	datajson := []byte(data)
	err = json.Unmarshal(datajson, &cp)
	if err != nil {
		log.Fatal(err)
	}
	return &cp.ClientPorts, cp.Zid
}

//目前只有zone级服务初始化,后面添加登陆、充值等
func (a *agent) InitSrv(name string) {
	if _, ok := a.srvs[name]; ok {
		return
	}
	cp, zid := a.ReadGameJson(name)

	run := CheckProcess(name)
	a.srvs[name] = &ServiceInfo{
		Started:        run,
		RegularlyCheck: run,
		Sname:          name,
		ClientPorts:    cp,
		Zid:            zid,
	}
}

//定时检查已启动的进程是否现在存在
func (a *agent) RegularlyCheckProcess() {
	for {
		for k, v := range a.srvs {
			if v.RegularlyCheck && CheckProcess(k) == false {
				log.Println("check process error ", k)
			}
		}

		time.Sleep(time.Minute * 30)
	}
}

//同步本机信息(机器名、机器上服务以及起停状态、svn代码版本号)
func (a *agent) C2sCheckReq() {
	p := protocol.C2sToken{
		Mservice: make(map[string]bool),
	}

	p.Host = a.hostName
	p.Token = utils.CreateMd5("cgyx2017")
	p.CodeVersion = a.SvnVer

	for k, v := range a.srvs {
		p.Mservice[k] = v.Started
	}

	protocol.SendJson(a.conn, protocol.CmdToken, &p)
}

//同步回复
func (a *agent) S2cCheckRsp(data []byte) {
	r := string(data)
	if r != "OK" {
		log.Fatal("register agentserver callback not ok")
	}
}

//更新zone配置信息
func (a *agent) S2cUpdateZoneConfig(data []byte) {
	p := protocol.S2cNotifyDo{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Println(" Stop Zone uncode json err, zone:", err.Error())
		return
	}
	r := protocol.C2sNotifyDone{
		Req:    p.Req,
		Do:     protocol.NotifyDoSuc,
		Result: "更新zone配置成功",
	}

	log.Println("update zoneConfig, Name:", p.Name, "req:", p.Req)
	if err := UpdateGameConf(); err != nil {
		r.Do = protocol.NotifyDoFail
		r.Result = fmt.Sprintf("更新zone配置失败,%v", err)
	}

	if _, ok := a.srvs[p.Name]; !ok {
		a.InitSrv(p.Name)
	}

	protocol.SendJson(a.conn, protocol.CmdUpdateHost, r)
}

//启动游戏服
func (a *agent) S2cStartZone(data []byte) {
	p := protocol.S2cNotifyDo{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Println(" StartZone uncode json err, zone:", err.Error())
		return
	}

	zone := p.Name
	r := protocol.C2sNotifyDone{
		Req:    p.Req,
		Do:     protocol.NotifyDoFail,
		Result: "启zone服成功",
	}
	run := StartZone(zone)
	a.srvs[zone].Started = run

	if run {
		a.srvs[zone].RegularlyCheck = true
		a.C2sZoneState(zone)
	} else {
		r.Do = protocol.NotifyDoFail
		r.Result = "启zone服失败"
	}

	protocol.SendJson(a.conn, protocol.CmdStartZone, r)
}

//关闭游戏服
func (a *agent) S2cStopZone(data []byte) {
	p := protocol.S2cNotifyDo{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Println(" Stop Zone uncode json err, zone:", err.Error())
		return
	}
	zone := p.Name
	r := protocol.C2sNotifyDone{
		Req:    p.Req,
		Do:     protocol.NotifyDoSuc,
		Result: "关zone服成功",
	}
	log.Println("recv stop msg, Name:", zone, "req:", p.Req)
	if StopZone(zone) {
		a.srvs[zone].Started = false
		a.srvs[zone].RegularlyCheck = false
		a.C2sZoneState(zone)
	} else {
		r.Do = protocol.NotifyDoFail
		r.Result = "关zone服失败"
	}
	protocol.SendJson(a.conn, protocol.CmdStopZone, r)
}

//启动机器上所有的游戏服
func (a *agent) S2cStartHostZones(data []byte) {
	p := protocol.S2cNotifyDo{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Println(" Start hostZones uncode json err, zone:", err.Error())
		return
	}
	r := protocol.C2sNotifyDone{
		Req:    p.Req,
		Do:     protocol.NotifyDoSuc,
		Result: "启动机器所有zone服成功",
	}

	for k, v := range a.srvs {
		run := StartZone(k)
		a.srvs[k].Started = run
		if run {
			a.srvs[k].RegularlyCheck = true
			a.C2sZoneState(k)
		} else {

			r.Do = protocol.NotifyDoFail
			if r.Result == "OK" {
				r.Result = "启动zone服失败,失败的服有:" + v.Sname
			} else {
				r.Result += "," + v.Sname
			}
		}
	}
	protocol.SendJson(a.conn, protocol.CmdStartHostZone, r)
}

//关闭机器上所有的游戏服
func (a *agent) S2cStopHostZones(data []byte) {
	p := protocol.S2cNotifyDo{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Println(" Stop hostZones uncode json err, zone:", err.Error())
		return
	}
	r := protocol.C2sNotifyDone{
		Req:    p.Req,
		Do:     protocol.NotifyDoSuc,
		Result: "关闭机器所有zone服成功",
	}
	for k, v := range a.srvs {
		stop := StopZone(k)
		if stop {
			a.srvs[k].Started = false
			a.srvs[k].RegularlyCheck = false
			a.C2sZoneState(k)
		} else {
			r.Do = protocol.NotifyDoFail
			if r.Result == "OK" {
				r.Result = "关zone服失败，失败的服有:" + v.Sname
			} else {
				r.Result += "," + v.Sname
			}
		}
	}
	protocol.SendJson(a.conn, protocol.CmdStopHostZone, r)
}

//agent同步游戏状态
func (a *agent) C2sZoneState(zone string) {
	p := protocol.C2sZoneState{
		Zone: zone,
		Open: a.srvs[zone].Started,
	}
	err := protocol.SendJson(a.conn, protocol.CmdZoneState, p)
	if err != nil {
		log.Println("sysn zone state err, ", err.Error())
	}
}

//更新svn
func (a *agent) S2cUpdateSvn(data []byte) {
	p := protocol.S2cNotifyDo{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Println(" Stop Zone uncode json err, zone:", err.Error())
		return
	}
	r := protocol.C2sNotifyDone{
		Req: p.Req,
		Do:  protocol.NotifyDoSuc,
	}

	log.Println("update zoneConfig, Name:", p.Name, "req:", p.Req)
	if err := SvnUp(); err != nil {
		r.Result = "更新svn失败"
		r.Do = protocol.NotifyDoFail
		protocol.SendJson(a.conn, protocol.CmdUpdateSvn, r)
	}

	a.SvnVer = SvnInfo()
	r.Result = a.SvnVer
	protocol.SendJson(a.conn, protocol.CmdUpdateSvn, r)
}
