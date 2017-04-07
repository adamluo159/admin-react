package machine

import (
	"fmt"
	"log"

	"github.com/adamluo159/admin-react/server/db"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	cl *mgo.Collection
)

//机器信息
type Machine struct {
	Hostname     string `json:"hostname" bson:"hostname"`
	IP           string
	OutIP        string   `json:"outIP" bson:"outIP"`
	Applications []string `json:"applications"`
}

const (
	CharDBPort    int    = 7000
	GatePort      int    = 7100
	CenterPort    int    = 7200
	LogicPort     int    = 7300
	ClientPort    int    = 7400
	MasterPort    int    = 9500
	LoginPort     int    = 7550
	LogPort       int    = 7600
	AccountDBPort int    = 6500
	UserName      string = "root"
	PassWord      string = "cg2016"
	RedisPort     int    = 6379
	MysqlPort     int    = 3306

	NetTimeOut       int = 1000 * 30
	NetMaxConnection int = 5000
	DbproxyServer    int = 1
	LoginServer      int = 2
	CenterServer     int = 3
	LogicServer      int = 4
	LogServer        int = 5
	MasterServer     int = 6
	GateServer       int = 7

	MasterCount int = 1
	LogMaxLine  int = 10000
)

type SRV map[string]int
type ServerConfigHead struct {
	NET_TIMEOUT_MSEC  int
	NET_MAX_CONNETION int
	StartService      []SRV
	LOG_INDEX         string
	LOG_MAXLINE       int
}

//机器模块注册
func Register(e *echo.Echo) {
	cl = db.Session.DB("machine").C("machine")
	if cl == nil {
		fmt.Printf("cannt find Collection about machine")
		panic(0)
	}
	i := mgo.Index{
		Key:    []string{"hostname"},
		Unique: true,
	}
	err := cl.EnsureIndex(i)
	if err != nil {
		fmt.Printf("mongodb ensureindex err:%s", err.Error())
		panic(0)
	}

	e.GET("/machine", GetMachines)
	e.POST("/machine/add", AddMachine)
	e.POST("/machine/save", SaveMachine)
	e.POST("/machine/del", DelMachine)
}

func GetMachineByName(name string) *Machine {
	m := Machine{}
	err := cl.Find(bson.M{"hostname": name}).One(&m)
	if err != nil {
		log.Println(" GetMachineByName name: ", name, err.Error())
		return nil
	}
	return &m
}

func UpdateMachineApplications(host string, apps []string) {
	log.Println(" UpdateMachineApplications, ", apps)
	err := cl.Update(bson.M{"hostname": host}, bson.M{"$set": bson.M{"applications": apps}})
	if err != nil {
		log.Println("UpdateMachineApplications update err, ", err.Error())
		return
	}

}
