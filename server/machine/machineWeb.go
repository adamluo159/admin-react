package machine

import (
	"log"
	"net/http"
	"os"

	"github.com/adamluo159/admin-react/server/comInterface"
	"github.com/adamluo159/gameAgent/utils"
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

func MachineRspFunc(rsp *InitMachine) {
	cl.Find(nil).All(&(*rsp).Items)

	var host string
	for index := 0; index < len(rsp.Items); index++ {
		host = (*rsp).Items[index].Hostname
		(*rsp).Items[index].Online = mhMgr.as.CheckOnlineMachine(host)
	}
}

//获取机器信息
func GetMachines(c echo.Context) error {
	rsp := InitMachine{}
	MachineRspFunc(&rsp)
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
	ret := InitMachine{}
	MachineRspFunc(&ret)

	err := c.Bind(&m)
	if err != nil {
		log.Println("save machine:", err.Error())
		return c.JSON(http.StatusOK, ret)
	}
	log.Println("get save info:", m)
	if m.Oldhost != m.Item.Hostname {
		if mhMgr.as.CheckOnlineMachine(m.Oldhost) {
			log.Println("已连接的机器不能修改主机名")
			return c.JSON(http.StatusOK, ret)
		}
		del := bson.M{"hostname": m.Oldhost}
		err = cl.Remove(del)
		if err != nil {
			log.Println(err.Error())
		}
		err = cl.Insert(m.Item)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		query := bson.M{"hostname": m.Item.Hostname}
		err = cl.Update(query, &m.Item)
		if err != nil {
			log.Println("SaveMachine, update:", err.Error())
		}
	}

	MachineRspFunc(&ret)
	return c.JSON(http.StatusOK, ret)
}

//删除
func DelMachine(c echo.Context) error {
	m, err := getM(&c)
	if err != nil {
		return err
	}
	query := bson.M{"hostname": m.Hostname}
	cl.Remove(query)

	ret := InitMachine{}
	MachineRspFunc(&ret)

	return c.JSON(http.StatusOK, ret)
}

func getM(c *echo.Context) (*comInterface.Machine, error) {
	m := comInterface.Machine{}
	err := (*c).Bind(&m)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &m, err
}

//生成登陆服、master/masterLog等上层服务器配置
func CommonConfig(c echo.Context) error {
	rsp := MachineRsp{}

	dir := os.Getenv("HOME") + comInterface.ConfDir + "/commonConfig/"
	os.Mkdir(dir, os.ModePerm)
	errLogM := mhMgr.GetMachineByName("errLog")
	accountDBM := mhMgr.GetMachineByName("accountDB")
	loginWebM := mhMgr.GetMachineByName("loginWeb")
	masterM := mhMgr.GetMachineByName("master")
	if errLogM == nil {
		rsp.Result = "cannt find errLog machine"
		return c.JSON(http.StatusOK, rsp)
	}
	if accountDBM == nil {
		rsp.Result = "cannt find accountDB machine"
		return c.JSON(http.StatusOK, rsp)
	}
	if loginWebM == nil {
		rsp.Result = "cannt find loginWeb machine"
		return c.JSON(http.StatusOK, rsp)
	}
	if masterM == nil {
		rsp.Result = "cannt find master machine"
		return c.JSON(http.StatusOK, rsp)
	}

	err := LoginLua(1, dir, loginWebM.IP, masterM.IP, accountDBM.IP)
	if err != nil {
		rsp.Result = err.Error()
		return c.JSON(http.StatusOK, rsp)
	}
	err = MasterLua(dir, masterM.IP)
	if err != nil {
		rsp.Result = err.Error()
		return c.JSON(http.StatusOK, rsp)
	}
	err = AccountDBLua(dir, accountDBM.IP, masterM.IP)
	if err != nil {
		rsp.Result = err.Error()
		return c.JSON(http.StatusOK, rsp)
	}
	err = MasterLogLua(dir, masterM.IP, errLogM.IP)
	if err != nil {
		rsp.Result = err.Error()
		return c.JSON(http.StatusOK, rsp)
	}

	commitstr := os.Getenv("HOME") + comInterface.ConfDir + "gitCommit"
	_, exeErr := utils.ExeShell("sh", commitstr, "updata common Config")
	if exeErr != nil {
		rsp.Result = exeErr.Error()
	}

	return c.JSON(http.StatusOK, rsp)
}
