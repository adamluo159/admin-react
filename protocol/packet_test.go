package protocol

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"testing"

	"bytes"
)

func TestPacket(t *testing.T) {
	a := C2sToken{
		Host:  "aa",
		Token: "111",
	}

	s, _ := json.Marshal(a)
	d := Packet(CmdToken, s)
	log.Println("packet:", string(s), len(s))

	buff := bytes.NewBuffer(d)
	ulength := uint32(0)

	binary.Read(buff, binary.BigEndian, &ulength)
	log.Println("ret:", string(d), ulength, len(d))

}

func TestUnPacket(t *testing.T) {
	a := C2sToken{
		Host:  "aa",
		Token: "111",
	}

	s, _ := json.Marshal(a)
	d := Packet(CmdToken, s)
	log.Println("Unpacket :", string(s), len(s))

	buff := bytes.NewBuffer(d)

	length := 0
	cmd, data := UnPacket(&length, buff)
	log.Println("unPack ret:", cmd, string(data), len(data))
}
