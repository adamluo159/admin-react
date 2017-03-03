package agentServer

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
)

type AgentMsg struct {
	Cmd  string
	Host string
	Data string
}

func NewClient(c *Client) {
	log.Println("new Client", c.conn.RemoteAddr().String())

	a := AgentMsg{}
	a.Cmd = "connected"
	data, err := json.Marshal(a)
	if err != nil {
		log.Println("newclient error, ", err.Error())
	}
	tmpCid++
	c.tmpCid = tmpCid
	c.Server.tmpClients[tmpCid] = c

	c.SendBytes(data)
}

func DisConnect(c *Client) {
	if c.tmpCid > 0 {
		delete(c.Server.tmpClients, c.tmpCid)
		log.Println(" client DisConnect: tmpCid:", c.tmpCid)
	} else {
		delete(c.Server.clients, c.host)
		log.Println(" client DisConnect: host:", c.host)
	}
	c.Close()
}

func TokenCheck(c *Client, a *AgentMsg) error {
	log.Println("recv agen token msg, token:", (*a).Data, " host:", (*a).Host)
	md5Ctx := md5.New()
	md5Ctx.Write([]byte("cgyx2017"))
	cipherStr := md5Ctx.Sum(nil)
	token := hex.EncodeToString(cipherStr)
	if token != (*a).Data {
		log.Println("token cannt be checked, ip and adress:", c.conn.RemoteAddr().String())
		return nil
	}

	rsp := AgentMsg{}
	rsp.Cmd = "checked"
	rsp.Data = "OK"
	data, err := json.Marshal(rsp)
	if err != nil {
		return err
	}
	c.host = a.Host
	c.Server.clients[a.Host] = c
	delete(c.Server.tmpClients, c.tmpCid)
	c.tmpCid = 0
	return c.SendBytes(data)
}

func Start(host string, zone string) {
	c := gserver.clients[host]
	if c == nil {
		return
	}
	a := AgentMsg{
		Cmd:  "start",
		Data: zone,
	}
	data, err := json.Marshal(a)
	if err != nil {
		log.Println(host + "  start: " + err.Error())
		return
	}
	err = c.SendBytes(data)
	if err != nil {
		log.Println(host + "  start: " + err.Error())
	}
}

func Stop(host string, zone string) {
	c := gserver.clients[host]
	if c == nil {
		return
	}
	a := AgentMsg{
		Cmd:  "stop",
		Data: zone,
	}
	data, err := json.Marshal(a)
	if err != nil {
		log.Println(host + "  stop: " + err.Error())
		return
	}
	err = c.SendBytes(data)
	if err != nil {
		log.Println(host + "  stop: " + err.Error())
	}
}
