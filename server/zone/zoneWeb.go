package zone

import (
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"strconv"

	"fmt"

	"github.com/adamluo159/admin-react/server/comInterface"
	"github.com/adamluo159/gameAgent/protocol"
	"github.com/adamluo159/gameAgent/utils"
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
		r := comInterface.RelationZone{
			Zid:           m.Zid,
			ZoneHost:      m.ZoneHost,
			ZoneDBHost:    m.ZoneDBHost,
			ZonelogdbHost: m.ZonelogdbHost,
		}
		zMgr.machineMgr.OpZoneRelation(&r, comInterface.RelationAdd)
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

	oldRelation := zMgr.GetZoneRelation(m.OldZid)
	query := bson.M{"zid": m.OldZid, "zoneName": m.OldZoneName}
	err = cl.Update(query, &m.Item)
	if err != nil {
		log.Println("SaveZone, update:", err.Error())
		ret.Result = "FALSE"
	} else {
		ret.Item = m.Item
		newRelation := &comInterface.RelationZone{
			Zid:           m.Item.Zid,
			ZoneDBHost:    m.Item.ZoneDBHost,
			ZoneHost:      m.Item.ZoneHost,
			ZonelogdbHost: m.Item.ZonelogdbHost,
		}
		zMgr.machineMgr.UpdateZone(oldRelation, newRelation)
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
	r := comInterface.RelationZone{
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
	zMgr.machineMgr.OpZoneRelation(&r, comInterface.RelationDel)
	err = cl.Remove(query)
	if err != nil {
		ret.Result = "FALSE"
	}
	return c.JSON(http.StatusOK, ret)
}

func UpdateZonelogdb(c echo.Context) error {
	zReq := ZoneReq{}
	err := c.Bind(&zReq)
	ret := ZoneRsp{}
	if err != nil {
		log.Println(err.Error())
		ret.Result = "FALSE"
		return c.JSON(http.StatusOK, ret)
	}
	m := zMgr.machineMgr.GetMachineByName(zReq.Host)
	if m == nil {
		ret.Result = "FAlse"
		return c.JSON(http.StatusOK, ret)
	}
	logdb := "zonelog" + strconv.Itoa(zReq.Zid)
	s, _ := utils.ExeShellArgs2("sh", "update_zonelogdb", logdb, m.IP)
	if logdb != s {
		ret.Result = fmt.Sprintf("update zonelogdb fail,%s", s)
	} else {
		ret.Result = "OK"
	}
	return c.JSON(http.StatusOK, ret)
}

func SynMachine(c echo.Context) error {
	ret := ZoneRsp{
		Result: "更新失败",
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
	ncode := zMgr.aserver.UpdateZone(hostname)
	if ncode == protocol.NotifyDoSuc {
		ret.Result = "更新成功"
	}

	return c.JSON(http.StatusOK, ret)
}

func StartZone(c echo.Context) error {
	ret := ZoneRsp{
		Result: "启服失败",
	}
	m := ZoneReq{}
	err := c.Bind(&m)
	if err != nil {
		ret.Result = "启服失败, parse post info fail"
		return c.JSON(http.StatusOK, ret)
	}
	zone := Zone{}
	query := bson.M{"zid": m.Zid}
	cl.Find(query).One(&zone)
	if zone.ZoneHost != m.Host {
		log.Printf(ret.Result, "send zid cannt match zonehost, zid:%d zonehost:%s", m.Zid, m.Host)
		return c.JSON(http.StatusOK, ret)
	}
	s := zMgr.aserver.StartZone(m.Host, m.Zid)
	log.Println("start result", s)
	switch s {
	case protocol.NotifyDoFail:
		ret.Result = "启服失败"
	case protocol.NotifyDoSuc:
		ret.Result = "启服成功"
	case protocol.NotifyDoing:
		ret.Result = "正在启服中，请勿重复启服"
	}
	log.Println("start ", ret)
	return c.JSON(http.StatusOK, ret)
}

func StopZone(c echo.Context) error {
	ret := ZoneRsp{
		Result: "关服失败",
	}
	m := ZoneReq{}
	err := c.Bind(&m)
	if err != nil {
		ret.Result = "关服失败, parse post info fail"
		return c.JSON(http.StatusOK, ret)
	}
	zone := Zone{}
	query := bson.M{"zid": m.Zid}
	cl.Find(query).One(&zone)
	if zone.ZoneHost != m.Host {
		log.Printf(ret.Result, "stopzone send zid cannt match zonehost, zid:%d zonehost:%s", m.Zid, m.Host)
		return c.JSON(http.StatusOK, ret)
	}
	s := zMgr.aserver.StopZone(m.Host, m.Zid)
	switch s {
	case protocol.NotifyDoFail:
		ret.Result = "关服失败"
	case protocol.NotifyDoSuc:
		ret.Result = "关服成功"
	case protocol.NotifyDoing:
		ret.Result = "正在关服中，请勿重复关服"
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
