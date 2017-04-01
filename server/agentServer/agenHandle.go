package agentServer

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type AgentMsg struct {
	Cmd  string
	Host string
	Data string
}

// Read client data from channel
func (c *Client) OnMessage() {
	buffer := make([]byte, 1024)
	for {
		reader := bufio.NewReader(c.conn)
		len, err := reader.Read(buffer)
		if err != nil {
			log.Println("msg error:", err.Error())
			return
		}
		dataLength := binary.LittleEndian.Uint32(buffer)
		if dataLength <= 0 || dataLength > 1020 {
			continue
		}
		a := AgentMsg{}
		json.Unmarshal(buffer[4:dataLength+4], &a)
		msgfunc := msgMap[a.Cmd]
		if msgfunc == nil {
			log.Println("cannt recv agent msg, msg: ", a, dataLength, len)
		} else {
			msgfunc(c, &a)
			log.Println("recv agent msg, msg: ", a, dataLength, len)
		}
	}
}

// Send bytes to client
func (c *Client) SendBytes(cmd string, jdata string) error {
	a := AgentMsg{
		Cmd:  cmd,
		Data: jdata,
	}
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	lenData := (uint32)(len(data))
	s := make([]byte, 4)
	binary.LittleEndian.PutUint32(s, lenData)
	buff := bytes.NewBuffer(s)
	buff.Write(data)

	empty := make([]byte, 1024-buff.Len())
	buff.Write(empty)

	_, serr := c.conn.Write(buff.Bytes())
	//log.Println("send msg:", len(buff.Bytes()), lenData, string(buff.Bytes()))
	return serr
}

func (c *Client) SendBytesCmd(cmd string) error {
	a := AgentMsg{
		Cmd: cmd,
	}
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	lenData := (uint32)(len(data))
	s := make([]byte, 4)
	binary.LittleEndian.PutUint32(s, lenData)
	buff := bytes.NewBuffer(s)
	buff.Write(data)

	empty := make([]byte, 1024-buff.Len())
	buff.Write(empty)

	_, serr := c.conn.Write(buff.Bytes())
	log.Println("send msg:", len(buff.Bytes()), lenData, string(buff.Bytes()))
	return serr
}

func TokenCheck(c *Client, a *AgentMsg) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte("cgyx2017"))
	cipherStr := md5Ctx.Sum(nil)
	token := hex.EncodeToString(cipherStr)
	if token != (*a).Data {
		log.Println("token cannt be checked, ip and adress:", c.conn.RemoteAddr().String())
		return
	}

	c.host = a.Host
	gserver.clients[a.Host] = c
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
		log.Println("cannt find client hostname:", host, gserver.clients)
		return
	}
	err := c.SendBytesCmd("update")
	if err != nil {
		log.Println(host + "  update: " + err.Error())
	}
}
