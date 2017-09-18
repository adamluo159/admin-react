package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"log"
	"regexp"
	"time"

	"errors"
	"fmt"
	"net"
)

const (
	HSize int = 4

	//agent反馈的状态
	NotifyDoFail int = 1
	NotifyDoSuc  int = 2
	NotifyDoing  int = 3

	Tzone      int = 1
	TzoneDB    int = 2
	TzonelogDB int = 3
	TSvn       int = 4

	Regzone      string = "^zone[0-9]*[1-9][0-9]*$"
	RegzoneDB    string = "^zonedb[0-9]*[1-9][0-9]*$"
	RegzonelogDB string = "^zonelogdb[0-9]*[1-9][0-9]*$"
)

var (
	gReqIndex    int                      = 1 //同步消息序号
	indexChanMap map[int]chan interface{} = make(map[int]chan interface{})
	TserviceReg                           = map[int]string{
		Tzone: Regzone,
	}
)

const (
	CmdNone          uint32 = 0
	CmdToken         uint32 = 1 //token验证
	CmdStartZone     uint32 = 2 //区服启动
	CmdStopZone      uint32 = 3 //区服停
	CmdStartHostZone uint32 = 4 //启动该机器上的所有区服
	CmdUpdateHost    uint32 = 5 //机器配置更新
	CmdStopHostZone  uint32 = 6 //关闭该机器上所有区服
	CmdZoneState     uint32 = 7 //区服状态(true:开启，false:关闭)
	CmdNewZone       uint32 = 8 //新增区服
	CmdUpdateSvn     uint32 = 9 //svn更新
)

type C2sToken struct {
	Token       string
	Host        string
	CodeVersion string
	Mservice    map[string]bool
}

type S2cNotifyDo struct {
	Name string
	Req  int
}
type C2sNotifyDone struct {
	Req    int
	Do     int
	Result string
}
type C2sZoneState struct {
	Zone string
	Open bool
}

func GetReqIndex() int {
	gReqIndex++
	return gReqIndex
}
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func NotifyWait(req int, v interface{}) {
	log.Println("notify wait , req:", req, indexChanMap, v)
	if c, ok := indexChanMap[req]; ok {
		c <- v
		delete(indexChanMap, req)
	}
}

func WaitCallBack(req int, reply interface{}, waitSec time.Duration) error {
	log.Println("wait call back msg, req:", req, indexChanMap)
	ch := make(chan interface{})
	indexChanMap[req] = ch
	t := time.NewTimer(time.Second * waitSec)
	select {
	case r := <-ch:
		DeepCopy(reply, r)
		log.Println("wait callback get reply:", reply)
	case <-t.C:
		delete(indexChanMap, req)
		log.Println("wait cb overtime in 30 second, req:", req)
		return errors.New(fmt.Sprintf("wait cb overtime in 30 second, req:%d", req))
	}
	return nil
}

func Packet(cmd uint32, data []byte) []byte {
	dataLen := len(data)
	HeadLen := HSize * 2
	packetlen := dataLen + HeadLen

	var buf = make([]byte, packetlen)
	binary.BigEndian.PutUint32(buf[:HSize], uint32(dataLen+HSize))
	binary.BigEndian.PutUint32(buf[HSize:HeadLen], cmd)
	copy(buf[HeadLen:], data)
	return buf
}

func UnPacket(length *int, msgbuf *bytes.Buffer) (uint32, []byte) {
	cmd := uint32(CmdNone)
	ulength := uint32(0)

	// 消息头
	if *length == 0 && msgbuf.Len() >= 4 {
		binary.Read(msgbuf, binary.BigEndian, &ulength)
		*length = int(ulength)
		// 检查超长消息
		if *length > 1024 {
			log.Printf("Message too length: %d\n", length)
			return 0, nil
		}
	}
	// 消息体
	if *length <= 0 || msgbuf.Len() < *length {
		//log.Printf("Message, len:%d, bufLen:%d", *length, msgbuf.Len())
		return 0, nil
	}

	binary.Read(msgbuf, binary.BigEndian, &cmd)
	*length = *length - HSize

	data := msgbuf.Next(*length)
	*length = 0

	return cmd, data
}

func Send(conn *net.Conn, cmd uint32, s string) error {
	if conn == nil {
		return errors.New(fmt.Sprintf("send msg cmd:%d, conn pointer is nil", cmd, conn))
	}

	b := Packet(cmd, []byte(s))
	_, err := (*conn).Write(b)
	if err != nil {
		return err
	}
	log.Println("send msg:", len(b), cmd, s)
	return nil
}

func SendJson(conn *net.Conn, cmd uint32, v interface{}) error {
	if conn == nil {
		return errors.New(fmt.Sprintf("sendjson send msg cmd:%d, conn pointer is nil", cmd, conn))
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	jsonBytes := Packet(cmd, data)
	_, jerr := (*conn).Write(jsonBytes)
	if jerr != nil {
		return jerr
	}
	log.Println("send msg:", cmd, len(jsonBytes), string(jsonBytes))
	return nil
}

func SendJsonWaitCB(conn *net.Conn, cmd uint32, v S2cNotifyDo, reply interface{}) error {
	if conn == nil {
		return errors.New(fmt.Sprintf("sendjson send msg cmd:%d, conn pointer is nil", cmd, conn))
	}
	v.Req = gReqIndex
	gReqIndex++

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	jsonBytes := Packet(cmd, data)
	_, jerr := (*conn).Write(jsonBytes)
	if jerr != nil {
		return jerr
	}

	ch := make(chan interface{})
	indexChanMap[v.Req] = ch

	t := time.NewTimer(time.Second * 30)
	select {
	case r := <-ch:
		DeepCopy(reply, r)
		log.Println("wait callback get reply:", reply)
	case <-t.C:
		delete(indexChanMap, v.Req)
		log.Println("wait cb overtime in 30 second, req:", v.Req)
		return errors.New(fmt.Sprintf("wait cb overtime in 30 second, req:%d", v.Req))
	}
	return nil
}

func MatchType(a string, f string) bool {
	reg := regexp.MustCompile(f)
	s := reg.FindAllString(a, -1)
	if len(s) > 0 {

		return true
	}
	return false
}
