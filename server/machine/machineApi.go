package machine

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//机器信息
type Machine struct {
	Hostname     string `json:"hostname" bson:"hostname"`
	IP           string
	OutIP        string `json:"outIP" bson:"outIP"`
	Applications string `json:"applications"`
}

type InitMachine struct {
	Items []Machine
}

type SaveMachineReq struct {
	Oldhost string
	Item    Machine
}

//回复信息
type MachineRsp struct {
	Result string
	Item   Machine
}

var (
	cl *mgo.Collection
)

//获取机器信息
func GetMachines(c echo.Context) error {
	rsp := InitMachine{}
	cl.Find(nil).All(&rsp.Items)
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
