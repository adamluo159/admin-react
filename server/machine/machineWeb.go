package machine

import (
	"encoding/json"
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

//
type MachineAllRsp struct {
	Result string
	Items  []comInterface.Machine
}

func MachineRspFunc(Items *[]comInterface.Machine) {
	cl.Find(nil).All(Items)

	var host string
	for index := 0; index < len(*Items); index++ {
		host = (*Items)[index].Hostname
		(*Items)[index].Online, (*Items)[index].CodeVersion = mhMgr.as.CheckOnlineMachine(host)
	}
}

//获取机器信息
func GetMachines(c echo.Context) error {
	rsp := InitMachine{}
	MachineRspFunc(&rsp.Items)
	return c.JSON(http.StatusOK, rsp)
}

//添加机器信息
func AddMachine(c echo.Context) error {
	m := comInterface.Machine{
		IP:       c.FormValue("IP"),
		Hostname: c.FormValue("hostname"),
		OutIP:    c.FormValue("outIP"),
	}
	ret := MachineAllRsp{Result: "OK"}
	err := cl.Insert(m)
	if err != nil {
		ret.Result = err.Error()
	} else {
		MachineRspFunc(&ret.Items)
	}
	return c.JSON(http.StatusOK, ret)
}

//保存
func SaveMachine(c echo.Context) error {
	m := comInterface.Machine{}
	rsp := MachineAllRsp{Result: "OK"}

	jdata := c.FormValue("Item")
	err := json.Unmarshal([]byte(jdata), &m)
	if err != nil {
		rsp.Result = err.Error()
		return c.JSON(http.StatusOK, rsp)
	}

	oldhost := c.FormValue("Oldhost")
	if oldhost != m.Hostname {
		online, _ := mhMgr.as.CheckOnlineMachine(oldhost)
		if online {
			rsp.Result = "已连接的机器不能修改主机名"
			return c.JSON(http.StatusOK, rsp)
		}
		del := bson.M{"hostname": oldhost}
		err = cl.Remove(del)
		if err != nil {
			log.Println(err.Error())
		}
		err = cl.Insert(m)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		query := bson.M{"hostname": m.Hostname}
		err = cl.Update(query, &m)
		if err != nil {
			log.Println("SaveMachine, update:", err.Error())
		}
	}

	MachineRspFunc(&rsp.Items)
	return c.JSON(http.StatusOK, rsp)
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
	MachineRspFunc(&ret.Items)

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
	rsp := MachineRsp{Result: "OK"}

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
	err = AccountDBLua(dir, accountDBM.IP, accountDBM.OutIP, masterM.IP)
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

func SvnUpdate(c echo.Context) error {
	hostName := c.FormValue("HostName")
	rsp := MachineAllRsp{Result: "OK"}
	suc := mhMgr.as.UpdateSvn(hostName)
	if !suc {
		rsp.Result = "Fail"
	}
	MachineRspFunc(&rsp.Items)
	return c.JSON(http.StatusOK, rsp)
}

func SvnUpdateAll(c echo.Context) error {
	rsp := MachineAllRsp{Result: "OK"}
	suc := mhMgr.as.UpdateSvnAll()
	if !suc {
		rsp.Result = "Fail"
	}
	MachineRspFunc(&rsp.Items)
	return c.JSON(http.StatusOK, rsp)
}
