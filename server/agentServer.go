package yada

import (
	"log"
	"net"
	"strconv"

	"github.com/adamluo159/admin-react/protocol"
)

type (
	ZoneStates struct {
		Host     string `json:"host" bson:"host"`
		Online   bool   `json:"online" bson:"online"`
		ZoneName string `json:"zoneName" bson:"zoneName"`
	}
	Aserver interface {
		Listen()
		StartZone(host string, zid int) string
		StopZone(host string, zid int) string
		CheckOnlineMachine(mName string) (bool, string)
		UpdateZone(host string) string
		StartAllZone() string
		StopAllZone() string
		OnlineZones() []ZoneStates
		UpdateSvn(host string) string
		UpdateSvnAll() string
	}

	// TCP server
	aserver struct {
		clients      map[string]*Client
		address      string // Address to open connection: localhost:9999
		allOperating bool
	}
)

func NewAS(address string) Aserver {
	log.Println("Creating server with address", address)
	return &aserver{
		address: address,
		clients: make(map[string]*Client),
	}
}

func (s *aserver) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn:    &conn,
			gserver: s,
			opCmdDoing: map[uint32]bool{
				protocol.CmdToken:         false,
				protocol.CmdStartZone:     false,
				protocol.CmdStopZone:      false,
				protocol.CmdUpdateHost:    false,
				protocol.CmdStartHostZone: false,
				protocol.CmdStopHostZone:  false,
				protocol.CmdZoneState:     false,
				protocol.CmdUpdateSvn:     false,
			},
			curZoneNotify: make(map[string]bool),
		}
		log.Println("agent server accpet socket, ip:", conn.RemoteAddr().String())
		go client.OnMessage()
	}
}

func (s *aserver) ClientDisconnect(host string) {
	delete(s.clients, host)
}

func (s *aserver) StartZone(host string, zid int) string {
	log.Println(" agentSserver startzone", host, " zid:", zid)
	r := protocol.C2sNotifyDone{
		Do: protocol.NotifyDoFail,
	}
	p := protocol.S2cNotifyDo{
		Name: "zone" + strconv.Itoa(zid),
	}

	c, ok := s.clients[host]
	if !ok {
		return "起zone服失败,找不到在线机器"
	}

	if do, ok := c.opCmdDoing[protocol.CmdStartZone]; ok && do {
		return "zone服正在启动中, 请勿重复启动"
	}

	c.opCmdDoing[protocol.CmdStartZone] = true
	if err := protocol.SendJsonWaitCB(c.conn, protocol.CmdStartZone, p, &r); err != nil {
		r.Result = err.Error()
	}
	c.opCmdDoing[protocol.CmdStartZone] = false
	c.curZoneNotify[p.Name] = false

	return r.Result

}

func (s *aserver) StopZone(host string, zid int) string {
	log.Println(" agentServer stopzone", host, " zid:", zid)

	r := protocol.C2sNotifyDone{
		Do: protocol.NotifyDoFail,
	}

	p := protocol.S2cNotifyDo{
		Name: "zone" + strconv.Itoa(zid),
	}

	c, ok := s.clients[host]
	if !ok {
		return "停zone服失败,找不到在线机器"
	}

	if do, ok := c.opCmdDoing[protocol.CmdStopZone]; ok && do {
		return "zone服正在关服中, 请勿重复操作"
	}

	c.opCmdDoing[protocol.CmdStopZone] = true
	if err := protocol.SendJsonWaitCB(c.conn, protocol.CmdStopZone, p, &r); err != nil {
		r.Result = err.Error()
	}
	c.opCmdDoing[protocol.CmdStopZone] = false
	return r.Result
}

func (s *aserver) UpdateZone(host string) string {
	log.Println(" agentServer update host info", host)
	r := protocol.C2sNotifyDone{
		Do:     protocol.NotifyDoFail,
		Result: "OK",
	}

	p := protocol.S2cNotifyDo{}
	c, ok := s.clients[host]
	if !ok {
		return "更新机器配置失败,找不到在线机器"
	}
	if do, ok := c.opCmdDoing[protocol.CmdUpdateHost]; ok && do {
		return "正在更新机器配置, 请勿重复操作"
	}

	c.opCmdDoing[protocol.CmdUpdateHost] = true
	if err := protocol.SendJsonWaitCB(c.conn, protocol.CmdUpdateHost, p, &r); err != nil {
		r.Result = err.Error()
	}
	c.opCmdDoing[protocol.CmdUpdateHost] = false

	return r.Result
}

func (s *aserver) StartAllZone() string {
	if s.allOperating {
		return "正在全zone服启动中，请勿重复启动"
	}
	r := protocol.C2sNotifyDone{
		Do: protocol.NotifyDoFail,
	}

	retStr := ""
	p := protocol.S2cNotifyDo{}
	s.allOperating = true

	for _, v := range s.clients {
		for kz, _ := range v.curZoneNotify {
			v.curZoneNotify[kz] = false
		}

		if err := protocol.SendJsonWaitCB(v.conn, protocol.CmdStartHostZone, p, &r); err != nil {
			r.Result = err.Error()
		}

		if r.Do != protocol.NotifyDoSuc {
			retStr += "机器:" + v.host + "  " + r.Result + "\n"
		}
	}
	s.allOperating = false
	return retStr
}

func (s *aserver) StopAllZone() string {
	if s.allOperating {
		return "正在停全zone服中，请勿重复操作"
	}
	r := protocol.C2sNotifyDone{
		Do: protocol.NotifyDoFail,
	}

	retStr := ""
	p := protocol.S2cNotifyDo{}
	s.allOperating = true
	for _, v := range s.clients {
		if err := protocol.SendJsonWaitCB(v.conn, protocol.CmdStopHostZone, p, &r); err != nil {
			r.Result = err.Error()
		}
		if r.Do != protocol.NotifyDoSuc {
			retStr += "机器:" + v.host + "  " + r.Result + "\n"
		}
	}
	s.allOperating = false

	return retStr
}

func (s *aserver) OnlineZones() []ZoneStates {
	sz := make([]ZoneStates, 0)
	for k, v := range s.clients {
		for sk, sv := range v.curServices {
			state := ZoneStates{
				Host:     k,
				Online:   sv,
				ZoneName: sk,
			}
			sz = append(sz, state)
		}
	}
	return sz
}

func (a *aserver) CheckOnlineMachine(mName string) (bool, string) {
	if v, ok := a.clients[mName]; ok {
		return true, v.codeVersion
	}
	return false, ""
}

func (s *aserver) UpdateSvn(host string) string {
	c := s.clients[host]
	if c == nil {
		return " Update, cannt find host client:"
	}

	if do, ok := c.opCmdDoing[protocol.CmdUpdateSvn]; ok && do {
		return "正在更新机器svn文件, 请勿重复操作"
	}

	p := protocol.S2cNotifyDo{}
	r := protocol.C2sNotifyDone{}

	c.opCmdDoing[protocol.CmdUpdateSvn] = true
	if err := protocol.SendJsonWaitCB(c.conn, protocol.CmdUpdateSvn, p, &r); err != nil {
		r.Result = err.Error()
	}
	c.opCmdDoing[protocol.CmdUpdateSvn] = false

	if r.Do == protocol.NotifyDoSuc {
		c.codeVersion = r.Result
		r.Result = "OK"
	}

	return r.Result
}

func (s *aserver) UpdateSvnAll() string {
	r := protocol.C2sNotifyDone{}
	p := protocol.S2cNotifyDo{}

	if s.allOperating {
		return "正在全机器更新svn文件中，请勿重复操作"
	}

	retStr := ""
	s.allOperating = true
	for _, v := range s.clients {
		if err := protocol.SendJsonWaitCB(v.conn, protocol.CmdUpdateSvn, p, &r); err != nil {
			r.Result = err.Error()
		}
		if r.Do != protocol.NotifyDoSuc {
			retStr += "机器:" + v.host + "  " + r.Result + "\n"
		} else {
			v.codeVersion = r.Result
		}
	}
	s.allOperating = false
	return retStr
}
