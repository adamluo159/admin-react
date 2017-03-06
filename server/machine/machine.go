package machine

import (
	"fmt"
	"net/http"

	"github.com/adamluo159/admin-react/server/db"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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

	NetTimeOut       int = 1000 * 3600
	NetMaxConnection int = 5000
	DbproxyServer    int = 1
	LoginServer      int = 2
	CenterServer     int = 3
	LogicServer      int = 4
	LogServer        int = 5
	MasterServer     int = 6
	GateServer       int = 7
)

type SRV map[string]int
type ServerConfigHead struct {
	NET_TIMEOUT_MSEC  int
	NET_MAX_CONNETION int
	StartService      []SRV
}

//机器信息
type Machine struct {
	Key      string `bson:"key" json:"key"`
	Hostname string `json:"hostname" bson:"hostname"`
	IP       string
	OutIP    string `json:"outIP" bson:"outIP"`
	Ct       string `json:"type" bson:"type"`
	Edit     bool   `json:"edit" bson:"edit"`
}

type SaveMachineReq struct {
	Oldhost string
	Item    Machine
}

//回复信息
type MachineRsp struct {
	Result string
	Item   Machine
	Items  []Machine
}

var (
	cl *mgo.Collection
)

//获取机器信息
func GetMachines(c echo.Context) error {
	rsp := MachineRsp{}
	cl.Find(nil).All(&rsp.Items)
	if len(rsp.Items) > 0 {
		rsp.Result = "OK"
	} else {
		rsp.Result = "NOT ITEM"
	}
	return c.JSON(http.StatusOK, rsp)
}

//添加机器信息
func AddMachine(c echo.Context) error {
	m, err := getM(&c)
	if err != nil {
		return err
	}
	ret := MachineRsp{
		Result: "OK",
	}
	err = cl.Insert(m)
	if err != nil {
		ret.Result = "FALSE"
	} else {
		ret.Item = *m
	}
	return c.JSON(http.StatusOK, ret)
}

//保存
func SaveMachine(c echo.Context) error {
	m := SaveMachineReq{}
	ret := MachineRsp{
		Result: "OK",
	}
	err := c.Bind(&m)
	if err != nil {
		fmt.Println(err.Error())
		ret.Result = "FALSE"
		return c.JSON(http.StatusOK, ret)
	}
	fmt.Println("get save info:", m)
	if m.Oldhost != m.Item.Hostname {
		del := bson.M{"hostname": m.Oldhost}
		err = cl.Remove(del)
		if err != nil {
			fmt.Println(err.Error())
		}
		err = cl.Insert(m.Item)
		if err != nil {
			fmt.Println(err.Error())
			ret.Result = "FALSE"
		} else {
			ret.Item = m.Item
		}
	} else {
		query := bson.M{"hostname": m.Item.Hostname}
		err = cl.Update(query, &m.Item)
		if err != nil {
			fmt.Println("SaveMachine, update:", err.Error())
			ret.Result = "FALSE"
		} else {
			ret.Item = m.Item
		}
	}

	return c.JSON(http.StatusOK, ret)
}

//删除
func DelMachine(c echo.Context) error {
	m, err := getM(&c)
	if err != nil {
		return err
	}
	ret := MachineRsp{}
	ret.Result = "OK"
	query := bson.M{"hostname": m.Hostname}
	err = cl.Remove(query)
	if err != nil {
		ret.Result = "FALSE"
	}
	return c.JSON(http.StatusOK, ret)
}

func getM(c *echo.Context) (*Machine, error) {
	m := Machine{}
	err := (*c).Bind(&m)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &m, err
}

//机器模块注册
func Register(e *echo.Echo) {
	cl = db.Session.DB("machine").C("machine")
	if cl == nil {
		fmt.Printf("cannt find Collection about machine")
		panic(0)
	}
	i := mgo.Index{
		Key:    []string{"key", "hostname"},
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

func GetMachineByName(name string) (*Machine, error) {
	m := Machine{}
	get := bson.M{"hostname": name}
	err := cl.Find(get).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
