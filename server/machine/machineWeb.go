package machine

import (
	"fmt"
	"net/http"

	"github.com/adamluo159/admin-react/server/comInterface"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

type InitMachine struct {
	Items []comInterface.Machine
}

type SaveMachineReq struct {
	Oldhost string
	Item    comInterface.Machine
}

//回复信息
type MachineRsp struct {
	Result string
	Item   comInterface.Machine
}

//获取机器信息
func GetMachines(c echo.Context) error {
	rsp := InitMachine{}
	cl.Find(nil).All(&rsp.Items)

	var host string
	for index := 0; index < len(rsp.Items); index++ {
		host = rsp.Items[index].Hostname
		rsp.Items[index].Online = mhMgr.as.CheckOnlineMachine(host)
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
		fmt.Println("save machine:", err.Error())
		ret.Result = "FALSE"
		return c.JSON(http.StatusOK, ret)
	}
	fmt.Println("get save info:", m)
	if m.Oldhost != m.Item.Hostname {
		if mhMgr.as.CheckOnlineMachine(m.Oldhost) {
			ret.Result = "已连接的机器不能修改主机名"
			return c.JSON(http.StatusOK, ret)
		}
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

func getM(c *echo.Context) (*comInterface.Machine, error) {
	m := comInterface.Machine{}
	err := (*c).Bind(&m)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &m, err
}