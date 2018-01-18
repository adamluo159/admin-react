package yada

import (
	"bytes"
	"encoding/json"
	"log"
	"net"

	"github.com/adamluo159/admin-react/protocol"
	"github.com/adamluo159/admin-react/utils"
)

type (
	Client struct {
		conn          *net.Conn
		host          string
		curServices   map[string]bool
		codeVersion   string
		gserver       *aserver
		opCmdDoing    map[uint32]bool
		curZoneNotify map[string]bool
	}
)

func (c *Client) OnMessage() {
	// 消息缓冲
	msgbuf := bytes.NewBuffer(make([]byte, 0, 1024))
	// 数据缓冲
	databuf := make([]byte, 1024)
	// 消息长度
	length := 0
	msgMap := map[uint32]func(data []byte){
		protocol.CmdToken:         c.TokenCheck,
		protocol.CmdStartZone:     c.CallBackHandle,
		protocol.CmdStopZone:      c.CallBackHandle,
		protocol.CmdUpdateHost:    c.CallBackHandle,
		protocol.CmdStartHostZone: c.CallBackHandle,
		protocol.CmdStopHostZone:  c.CallBackHandle,
		protocol.CmdZoneState:     c.ZoneState,
		protocol.CmdUpdateSvn:     c.CallBackHandle,
		protocol.CmdProcessErr:    c.CheckProcessErr,
	}

	for {
		// 读取数据
		n, err := (*c.conn).Read(databuf)
		if err != nil {
			log.Printf("Read error: %s host:%s\n", err, c.host)
			c.gserver.ClientDisconnect(c.host)
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
			mfunc := msgMap[cmd]
			if mfunc == nil {
				log.Printf("cannt find msg handle Client cmd: %d data: %s\n", cmd, string(data))
			} else {
				log.Printf("Client cmd: %d data: %s\n", cmd, string(data))
				mfunc(data)
			}
		}
	}
}

func (c *Client) TokenCheck(data []byte) {
	p := protocol.C2sToken{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Println("TokenCheck, uncode error: ", string(data))
		return
	}
	if utils.Md5Check(p.Token, "cgyx2017") == false {
		(*c.conn).Close()
		return
	}

	c.host = p.Host
	c.codeVersion = p.CodeVersion
	c.gserver.clients[p.Host] = c
	c.curServices = p.Mservice

	protocol.Send(c.conn, protocol.CmdToken, "OK")
}

func (c *Client) CallBackHandle(data []byte) {
	r := protocol.C2sNotifyDone{}
	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Println("CallBackHandle, uncode error: ", string(data))
		return
	}
	protocol.NotifyWait(r.Req, r)
}

func (c *Client) ZoneState(data []byte) {
	p := protocol.C2sZoneState{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Println("ZoneState, uncode error: ", string(data))
		return
	}
	c.curServices[p.Zone] = p.Open
	log.Println("notify zone state", p, c.curServices[p.Zone])
}

func (c *Client) CheckProcessErr(data []byte) {
	p := protocol.C2sDx{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Println("ZoneState, uncode error: ", string(data))
		return
	}
	arg := ""
	for _, z := range p.Zones {
		v, ok := c.curZoneNotify[z]
		if !ok || (ok && !v) {
			c.curZoneNotify[z] = true
			arg += z + ","
		}
	}
	if arg != "" {
		s, serr := utils.ExeShell("php", "shells/aliyunDX/demo/sendSms.php", arg)
		log.Println("exshell rsp:", s, "err:", serr)
	}

	log.Println("CheckProcessErr", arg)
}
