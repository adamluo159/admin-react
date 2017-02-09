package machine

import (
	"fmt"
	"net/http"

	"github.com/adamluo159/admin-react/server/db"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//机器信息
type Machine struct {
	Key      string `bson:"key" json:"key"`
	Hostname string `json:"hostname" bson:"hostname"`
	IP       string
	OutIP    string `json:"outIP" bson:"outIP"`
	Ct       string `json:"type" bson:"type"`
	Edit     bool   `json:"edit" bson:"edit"`
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
	m, err := getM(&c)
	if err != nil {
		return err
	}
	fmt.Println("get save info:", m)
	ret := MachineRsp{
		Result: "OK",
	}
	query := bson.M{"hostname": m.Hostname}
	err = cl.Update(query, &m)
	if err != nil {
		ret.Result = "FALSE"
	} else {
		ret.Item = *m
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
	e.GET("/machine", GetMachines)
	e.POST("/machine/add", AddMachine)
	e.POST("/machine/save", SaveMachine)
	e.POST("/machine/del", DelMachine)
}
