package zone

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"strconv"

	"errors"

	"fmt"

	"github.com/adamluo159/admin-react/server/machine"
	"github.com/adamluo159/gameAgent/agentServer"
	"github.com/labstack/echo"
)

//区服信息
type Zone struct {
	Zid           int      `bson:"zid" json:"zid"`
	ZoneName      string   `json:"zoneName" bson:"zoneName"`
	ZoneHost      string   `json:"zoneHost" bson:"zoneHost"`
	ZoneDBHost    string   `json:"zoneDBHost" bson:"zoneDBHost"`
	ZonelogdbHost string   `json:"zonelogdbHost" bson:"zonelogdbHost"`
	Channels      []string `json:"channels" bson:"channels"`
	Whitelst      bool     `json:"whitelst" bson:"whitelst"`
}

type ZoneReq struct {
	Zid  int
	Host string
}

type SaveZoneReq struct {
	OldZoneName string
	OldZid      int
	Item        Zone
}

//回复信息
type ZoneRsp struct {
	Result string
	Item   Zone
	Items  []Zone
}

var (
	cl              *mgo.Collection
	GlobalDB        MysqlLua
	Str2IntChannels map[string]int
	LogicMap        map[int][]int
)

//获取区服信息
func GetZones(c echo.Context) error {
	rsp := ZoneRsp{}
	cl.Find(nil).All(&rsp.Items)
	if len(rsp.Items) > 0 {
		rsp.Result = "OK"
	} else {
		rsp.Result = "NOT ITEM"
	}
	return c.JSON(http.StatusOK, rsp)
}

//添加区服信息
func AddZone(c echo.Context) error {
	m, err := getM(&c)
	if err != nil {
		return err
	}
	ret := ZoneRsp{
		Result: "OK",
	}
	err = cl.Insert(m)
	if err != nil {
		ret.Result = "FALSE"
	} else {
		ret.Item = *m
		//新增的zone用到的机器加入到各自的用途中
		r := machine.RelationZone{
			Zid:           m.Zid,
			ZoneHost:      m.ZoneHost,
			ZoneDBHost:    m.ZoneDBHost,
			ZonelogdbHost: m.ZonelogdbHost,
		}
		machine.OpZoneRelation(&r, machine.RelationAdd)
	}
	return c.JSON(http.StatusOK, ret)
}

//保存
func SaveZone(c echo.Context) error {
	m := SaveZoneReq{}
	ret := ZoneRsp{
		Result: "OK",
	}
	err := c.Bind(&m)
	if err != nil {
		log.Println(err.Error())
		ret.Result = "FALSE"
		return c.JSON(http.StatusOK, ret)
	}
	log.Println("get save info:", m)

	oldRelation := GetZoneRelation(m.OldZid)

	query := bson.M{"zid": m.OldZid, "zoneName": m.OldZoneName}
	err = cl.Update(query, &m.Item)
	if err != nil {
		log.Println("SaveZone, update:", err.Error())
		ret.Result = "FALSE"
	} else {
		ret.Item = m.Item
		newRelation := machine.RelationZone{
			Zid:           m.Item.Zid,
			ZoneDBHost:    m.Item.ZoneDBHost,
			ZoneHost:      m.Item.ZoneHost,
			ZonelogdbHost: m.Item.ZonelogdbHost,
		}
		machine.UpdateZone(oldRelation, &newRelation)
	}

	return c.JSON(http.StatusOK, ret)
}

//删除
func DelZone(c echo.Context) error {
	m, err := getM(&c)
	if err != nil {
		return err
	}
	ret := ZoneRsp{}
	dzone := Zone{}
	ret.Result = "OK"
	query := bson.M{"zid": m.Zid}
	err = cl.Find(query).One(&dzone)
	if err != nil {
		log.Println("delete zone ok:", dzone, err.Error(), m.Zid)
		return c.JSON(http.StatusOK, ret)
	}
	r := machine.RelationZone{
		ZoneDBHost:    dzone.ZoneDBHost,
		ZoneHost:      dzone.ZoneHost,
		ZonelogdbHost: dzone.ZonelogdbHost,
		Zid:           dzone.Zid,
	}
	log.Println("delete zone :", r)
	err = DelZoneConfig(r.Zid, r.ZoneHost)
	if err != nil {
		ret.Result = "FALSE"
		return c.JSON(http.StatusOK, ret)
	}
	machine.OpZoneRelation(&r, machine.RelationDel)
	err = cl.Remove(query)
	if err != nil {
		ret.Result = "FALSE"
	}
	return c.JSON(http.StatusOK, ret)
}

func UpdateZonelogdb(c echo.Context) error {
	z := ZoneReq{}
	err := c.Bind(&z)
	ret := ZoneRsp{}
	if err != nil {
		log.Println(err.Error())
		ret.Result = "FALSE"
		return c.JSON(http.StatusOK, ret)
	}
	m := machine.GetMachineByName(z.Host)
	if m == nil {
		ret.Result = "FAlse"
		return c.JSON(http.StatusOK, ret)
	}
	err = ExecUpdatelogdb(z.Zid, m.IP)
	if err != nil {
		ret.Result = err.Error()
	} else {
		ret.Result = "OK"
	}
	return c.JSON(http.StatusOK, ret)
}

func SynMachine(c echo.Context) error {
	ret := ZoneRsp{
		Result: "OK",
	}

	zid, err := strconv.Atoi(c.QueryParam("zid"))
	if err != nil {
		ret.Result = "zid" + err.Error()
		return c.JSON(http.StatusOK, ret)
	}

	hostname := c.QueryParam("hostname")
	if hostname == "" {
		ret.Result = "hostname none"
		return c.JSON(http.StatusOK, ret)
	}

	WriteZoneConfigLua(zid, &ret, hostname)
	agentServer.Update(hostname)
	return c.JSON(http.StatusOK, ret)
}

func StartZone(c echo.Context) error {
	ret := ZoneRsp{
		Result: "OK",
	}
	m := ZoneReq{}
	err := c.Bind(&m)
	if err != nil {
		ret.Result = "Start zone, parse post info fail"
		return c.JSON(http.StatusOK, ret)
	}
	zone := Zone{}
	query := bson.M{"zid": m.Zid}
	cl.Find(query).One(&zone)
	if zone.ZoneHost != m.Host {
		log.Printf(ret.Result, "send zid cannt match zonehost, zid:%d zonehost:%s", m.Zid, m.Host)
		return c.JSON(http.StatusOK, ret)
	}
	suc := agentServer.StartZone(m.Host, m.Zid)
	if suc == false {
		return c.JSON(http.StatusOK, ret)
	}
	return c.JSON(http.StatusOK, ret)
}

func StopZone(c echo.Context) error {
	ret := ZoneRsp{
		Result: "OK",
	}
	m := ZoneReq{}
	err := c.Bind(&m)
	if err != nil {
		ret.Result = "Stop zone, parse post info fail"
		return c.JSON(http.StatusOK, ret)
	}
	zone := Zone{}
	query := bson.M{"zid": m.Zid}
	cl.Find(query).One(&zone)
	if zone.ZoneHost != m.Host {
		log.Printf(ret.Result, "stopzone send zid cannt match zonehost, zid:%d zonehost:%s", m.Zid, m.Host)
		return c.JSON(http.StatusOK, ret)
	}
	suc := agentServer.StopZone(m.Host, m.Zid)
	if suc == false {
		return c.JSON(http.StatusOK, ret)
	}
	return c.JSON(http.StatusOK, ret)
}

func getM(c *echo.Context) (*Zone, error) {
	m := Zone{}
	err := (*c).Bind(&m)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &m, err
}

func ExecUpdatelogdb(zid int, arg string) error {
	logdb := "zonelog" + strconv.Itoa(zid)
	log.Println("begin execute updatezonelogdb shell.....", zid, "--", arg, logdb)
	// 执行系统命令
	// 第一个参数是命令名称
	// 后面参数可以有多个，命令参数
	cmd := exec.Command("sh", "update_zonelogdb", logdb, arg) //"GameConfig/gitCommit", "zoneo")
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		return err
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
		return err
	}
	cmd.Wait()
	if logdb != string(opBytes) {
		return errors.New(fmt.Sprintf("update zonelogdb fail,%s", string(opBytes)))
	}
	return nil
}
