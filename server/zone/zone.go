package zone

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"strconv"

	"github.com/adamluo159/admin-react/server/agentServer"
	"github.com/adamluo159/admin-react/server/db"
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
	Zid      int    `bson:"zid" json:"zid"`
	ZoneHost string `json:"zoneHost" bson:"zoneHost"`
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
		fmt.Println(err.Error())
		ret.Result = "FALSE"
		return c.JSON(http.StatusOK, ret)
	}
	fmt.Println("get save info:", m)
	query := bson.M{"zid": m.OldZid, "zoneName": m.OldZoneName}
	err = cl.Update(query, &m.Item)
	if err != nil {
		fmt.Println("SaveZone, update:", err.Error())
		ret.Result = "FALSE"
	} else {
		ret.Item = m.Item
	}

	return c.JSON(http.StatusOK, ret)
}

//删除
//func DelZone(c echo.Context) error {
//	m, err := getM(&c)
//	if err != nil {
//		return err
//	}
//	ret := ZoneRsp{}
//	ret.Result = "OK"
//	query := bson.M{"hostname": m.Hostname}
//	err = cl.Remove(query)
//	if err != nil {
//		ret.Result = "FALSE"
//	}
//	return c.JSON(http.StatusOK, ret)
//}

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
	if hostname != "" {
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
	if zone.ZoneHost != m.ZoneHost {
		fmt.Sprintf(ret.Result, "send zid cannt match zonehost, zid:%d zonehost:%s", m.Zid, m.ZoneHost)
		return c.JSON(http.StatusOK, ret)
	}
	suc := agentServer.StartZone(m.ZoneHost, m.Zid)
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
	if zone.ZoneHost != m.ZoneHost {
		fmt.Sprintf(ret.Result, "stopzone send zid cannt match zonehost, zid:%d zonehost:%s", m.Zid, m.ZoneHost)
		return c.JSON(http.StatusOK, ret)
	}
	suc := agentServer.StopZone(m.ZoneHost, m.Zid)
	if suc == false {
		return c.JSON(http.StatusOK, ret)
	}
	return c.JSON(http.StatusOK, ret)
}

func getM(c *echo.Context) (*Zone, error) {
	m := Zone{}
	err := (*c).Bind(&m)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &m, err
}

//区服模块注册
func Register(e *echo.Echo) {
	Str2IntChannels = make(map[string]int)
	Str2IntChannels["ios"] = 1
	Str2IntChannels["yyb"] = 2
	Str2IntChannels["xiaomi"] = 3

	cl = db.Session.DB("zone").C("zone")
	if cl == nil {
		fmt.Printf("cannt find Collection about zone")
		panic(0)
	}
	i := mgo.Index{
		Key:    []string{"zid"},
		Unique: true,
	}
	err := cl.EnsureIndex(i)
	if err != nil {
		fmt.Printf("mongodb ensureindex err:%s", err.Error())
		panic(0)
	}
	LogicMap = make(map[int][]int)
	LogicMap[1] = []int{
		210106002,
		210109001,
		210109002,
		210109003,
		210109004,
		210109005,
		210109006,
		210109007,
		210109008,
		210109009,
		210109999,
		210181001,
		210181002,
		210181003,
		210181004,
		210182001,
		210182002,
		210182003,
		210182004,
	}
	LogicMap[2] = []int{
		210106001,
		210102101,
		210102102,
		210102201,
		210102202,
		210104001,
		210104002,
		210104005,
		210104099,
		210104006,
		210105001,
		210105002,
		210106003,
		210107001,
		210107002,
	}

	iName := mgo.Index{
		Key:    []string{"zoneName"},
		Unique: true,
	}
	err = cl.EnsureIndex(iName)
	if err != nil {
		fmt.Printf("mongodb ensureindex err:%s", err.Error())
		panic(0)
	}
	e.GET("/zone", GetZones)
	e.POST("/zone/add", AddZone)
	e.POST("/zone/save", SaveZone)
	e.GET("/zone/synMachine", SynMachine)
	e.POST("/zone/startZone", StartZone)
	e.POST("/zone/stopZone", StopZone)
	//e.POST("/zone/del", DelZone)
}
