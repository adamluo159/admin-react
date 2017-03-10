package agentServer

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"strconv"
	"time"
)

type AgentMsg struct {
	Cmd  string
	Host string
	Data string
}

func NewClient(c *Client) {
	log.Println("new Client", c.conn.RemoteAddr().String())

	tmpCid++
	c.tmpCid = tmpCid
	c.Server.tmpClients[tmpCid] = c

	c.SendBytesCmd("connected")
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

func TokenCheck(c *Client, a *AgentMsg) {
	log.Println("recv agen token msg, token:", (*a).Data, " host:", (*a).Host)
	md5Ctx := md5.New()
	md5Ctx.Write([]byte("cgyx2017"))
	cipherStr := md5Ctx.Sum(nil)
	token := hex.EncodeToString(cipherStr)
	if token != (*a).Data {
		log.Println("token cannt be checked, ip and adress:", c.conn.RemoteAddr().String())
		return
	}

	c.host = a.Host
	c.Server.clients[a.Host] = c
	delete(c.Server.tmpClients, c.tmpCid)
	c.tmpCid = 0
	c.SendBytes("checked", "OK")
}

func Ping(c *Client, a *AgentMsg) {
	c.pingTime = time.Now()
}

func StartZone(host string, zid int) bool {
	log.Println(" recv web cmd startzone", host, " zid:", zid)
	c := gserver.clients[host]
	if c == nil {
		return false
	}
	zone := "zone" + strconv.Itoa(zid)
	err := c.SendBytes("start", zone)
	if err != nil {
		log.Println(host + "  startzone: " + err.Error())
	}
	return true
}

func StopZone(host string, zid int) bool {
	log.Println(" recv web cmd stopzone", host, " zid:", zid)
	c := gserver.clients[host]
	if c == nil {
		return false
	}
	zone := "zone" + strconv.Itoa(zid)
	err := c.SendBytes("stop", zone)
	if err != nil {
		log.Println(host + "  stopzone: " + err.Error())
	}
	return true
}

func Update(host string) {
	log.Println(" recv web cmd update", host)
	c := gserver.clients[host]
	if c == nil {
		log.Println("cannt find client hostname:", host, gserver.clients, gserver.tmpClients)
		return
	}
	err := c.SendBytesCmd("update")
	if err != nil {
		log.Println(host + "  update: " + err.Error())
	}
}
