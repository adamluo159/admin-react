package yada

import (
	"log"
	"net"
	"strconv"

	"github.com/adamluo159/admin-react/server/comInterface"
	"github.com/adamluo159/gameAgent/protocol"
)

type (
	ZoneStates struct {
		Host     string
		Online   bool
		ZoneName string
	}
	Aserver interface {
		Listen()
		StartZone(host string, zid int) int
		StopZone(host string, zid int) int
		CheckOnlineMachine(mName string) (bool, string)
		UpdateZone(host string) int
		StartAllZone() int
		StopAllZone() int
		OnlineZones() []ZoneStates
		AddNewZone(host string, zone string, zid int)
		UpdateSvn(host string) bool
		UpdateSvnAll() bool
	}

	// TCP server
	aserver struct {
		clients             map[string]*Client
		address             string // Address to open connection: localhost:9999
		mhMgr               comInterface.MachineMgr
		zonelogDBserviceMap map[string][]string
		allOperating        bool
	}
)

func NewAS(address string) Aserver {
	log.Println("Creating server with address", address)
	return &aserver{
		address:             address,
		clients:             make(map[string]*Client),
		zonelogDBserviceMap: make(map[string][]string),
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
		}
		go client.OnMessage()
	}
}

func (s *aserver) ClientDisconnect(host string) {
	delete(s.clients, host)
}

func (s *aserver) StartZone(host string, zid int) int {
	log.Println(" agentSserver startzone", host, " zid:", zid)
	r := protocol.C2sNotifyDone{
		Do: protocol.NotifyDoFail,
	}
	p := protocol.S2cNotifyDo{
		Name: "zone" + strconv.Itoa(zid),
	}

	c, ok := s.clients[host]
	if !ok {
		return r.Do
	}
	protocol.SendJsonWaitCB(c.conn, protocol.CmdStartZone, p, &r)
	return r.Do

}

func (s *aserver) StopZone(host string, zid int) int {
	log.Println(" agentServer stopzone", host, " zid:", zid)

	r := protocol.C2sNotifyDone{
		Do: protocol.NotifyDoFail,
	}

	p := protocol.S2cNotifyDo{
		Name: "zone" + strconv.Itoa(zid),
	}

	c, ok := s.clients[host]
	if !ok {
		return r.Do
	}
	protocol.SendJsonWaitCB(c.conn, protocol.CmdStopZone, p, &r)
	return r.Do
}

func (s *aserver) UpdateZone(host string) int {
	log.Println(" agentServer update host info", host)
	r := protocol.C2sNotifyDone{
		Do: protocol.NotifyDoFail,
	}

	p := protocol.S2cNotifyDo{}
	c, ok := s.clients[host]
	if !ok {
		return r.Do
	}
	protocol.SendJsonWaitCB(c.conn, protocol.CmdUpdateHost, p, &r)
	return r.Do
}

func (s *aserver) StartAllZone() int {
	if s.allOperating {
		return protocol.NotifyDoing
	}
	r := protocol.C2sNotifyDone{
		Do: protocol.NotifyDoFail,
	}

	p := protocol.S2cNotifyDo{}
	s.allOperating = true
	for _, v := range s.clients {
		protocol.SendJsonWaitCB(v.conn, protocol.CmdStartHostZone, p, &r)
	}
	s.allOperating = false
	return r.Do
}

func (s *aserver) StopAllZone() int {
	if s.allOperating {
		return protocol.NotifyDoing
	}
	r := protocol.C2sNotifyDone{
		Do: protocol.NotifyDoFail,
	}

	p := protocol.S2cNotifyDo{}
	s.allOperating = true
	for _, v := range s.clients {
		protocol.SendJsonWaitCB(v.conn, protocol.CmdStopHostZone, p, &r)
	}

	s.allOperating = false

	return r.Do
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
	log.Println("wwwwwwwwwwww", sz)
	return sz
}

func (s *aserver) CheckOnlineMachine(mName string) (bool, string) {
	if v, ok := (*s).clients[mName]; ok {
		return true, v.codeVersion
	}
	return false, ""
}

func (s *aserver) AddNewZone(host string, zone string, zid int) {
	c := s.clients[host]
	if c == nil {
		log.Println(" AddNewZone, cannt find host client:", host, zone)
		return
	}
	p := protocol.S2cNotifyDo{
		Name: "zone" + strconv.Itoa(zid),
	}
	r := protocol.C2sNotifyDone{}
	protocol.SendJsonWaitCB(c.conn, protocol.CmdNewZone, p, &r)
}

func (s *aserver) UpdateSvn(host string) bool {
	c := s.clients[host]
	if c == nil {
		log.Println(" Update, cannt find host client:", host)
		return false
	}
	p := protocol.S2cNotifyDo{}
	r := protocol.C2sNotifyDone{}
	protocol.SendJsonWaitCB(c.conn, protocol.CmdNewZone, p, &r)
	return r.Do == protocol.NotifyDoSuc
}

func (s *aserver) UpdateSvnAll() bool {
	r := protocol.C2sNotifyDone{}
	p := protocol.S2cNotifyDo{}
	suc := true
	for _, v := range s.clients {
		protocol.SendJsonWaitCB(v.conn, protocol.CmdUpdateSvn, p, &r)
		if r.Do == protocol.NotifyDoFail {
			suc = false
		}
	}
	return suc
}
