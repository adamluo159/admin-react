package machine

import (
	"encoding/json"
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

type MachineRsp struct {
	Result string
	Item   Machine
}

var (
	cl *mgo.Collection
)

//获取机器信息
func GetMachines(c echo.Context) error {
	var a []Machine
	err := cl.Find(nil).All(&a)
	fmt.Println("machine request:", err, &a)
	b, err := json.Marshal(&a)
	return c.Blob(http.StatusOK, "application/json", b)
}

//添加机器信息
func AddMachine(c echo.Context) error {
	m := Machine{}
	err := c.Bind(&m)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//s := []string{"{'hostname': 'host101'}"}
	//i := mgo.Index{
	//	Key:    s,
	//	Unique: true,
	//}
	//err = cl.EnsureIndex(i)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return err
	//}
	err = cl.Insert(m)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("aaaa", m)
	ret := MachineRsp{}
	ret.Item = m
	ret.Result = "OK"
	return c.JSON(http.StatusOK, ret)
}

//保存
func SaveMachine(c echo.Context) error {
	fmt.Println("recv save machine")
	m, err := getM(&c)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("get save info:", m)
	query := bson.M{"hostname": m.Hostname}
	err = cl.Update(query, &m)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	ret := MachineRsp{}
	ret.Item = *m
	ret.Result = "OK"
	err = c.JSON(http.StatusOK, ret)
	fmt.Println("wwwww", err.Error())
	return err
}

//删除
func DelMachine(c echo.Context) error {
	m, err := getM(&c)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	query := bson.M{"hostname": m.Hostname}
	err = cl.Remove(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
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
