package zone

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"strconv"

	"github.com/adamluo159/admin-react/server/db"
	"github.com/adamluo159/admin-react/server/machine"
	"github.com/adamluo159/struct2lua"
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

type Connect struct {
	ID   int
	Port int
	IP   string
}

type MysqlLua struct {
	IP             string
	Port           int
	UserName       string
	PassWord       string
	FlushFrequency int
	DataBase       string
}

type RedisLua struct {
	IP       string
	Port     int
	Password string
}

type Gate struct {
	ID             int
	Zid            int
	ServerIP       string
	ServerPort     int
	ClientIP       string
	ClientPort     int
	ChannelIds     []string
	Open           bool
	Name           string
	ConnectServers map[string]interface{}
}

var (
	cl *mgo.Collection
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
	//e.POST("/zone/del", DelZone)
}

func SynMachine(c echo.Context) error {
	ret := ZoneRsp{
		Result: "OK",
	}
	zid, err := strconv.Atoi(c.QueryParam("zid"))
	if err != nil {
		ret.Result = err.Error()
		return c.JSON(http.StatusOK, ret)
	}

	zone := Zone{}
	query := bson.M{"zid": zid}
	cl.Find(query).One(&zone)

	//zoneDBquery := bson.M{"zoneDBHost": zone.ZoneDBHost}
	//zonelogquery := bson.M{"zonelogdbHost": zone.ZonelogdbHost}

	//zoneDBCount, zdberr := cl.Find(zoneDBquery).Count()
	//if zdberr != nil {
	//	ret.Result = zdberr.Error()
	//	return c.JSON(http.StatusOK, ret)
	//}
	//zonelogDBCount, zlogerr := cl.Find(zonelogquery).Count()
	//if zlogerr != nil {
	//	ret.Result = zlogerr.Error()
	//	return c.JSON(http.StatusOK, ret)
	//}

	//zonedbm, dberr := machine.GetMachineByName(zone.ZoneDBHost)
	//if dberr != nil {
	//	ret.Result = dberr.Error()
	//}

	//logm, logerr := machine.GetMachineByName(zone.ZonelogdbHost)
	//if logerr != nil {
	//	ret.Result = logerr.Error()
	//}

	gerr := GateLua(&zone)
	if gerr != nil {
		ret.Result = gerr.Error()
		return c.JSON(http.StatusOK, ret)
	}

	return c.JSON(http.StatusOK, ret)
}

func GateLua(zone *Zone) error {
	zonequery := bson.M{"zoneHost": zone.ZoneHost}
	zoneCount, zerr := cl.Find(zonequery).Count()
	if zerr != nil {
		return zerr
	}
	zonem, err := machine.GetMachineByName(zone.ZoneHost)
	if err != nil {
		return err
	}
	masterm, merr := machine.GetMachineByName("cghost2")
	if merr != nil {
		return merr
	}
	gateLua := Gate{
		ID:             zone.Zid,
		Zid:            zone.Zid,
		ServerIP:       zonem.IP,
		ServerPort:     machine.GatePort + zoneCount,
		ClientIP:       zonem.OutIP,
		ClientPort:     machine.ClientPort + zoneCount,
		ChannelIds:     zone.Channels,
		Open:           zone.Whitelst,
		Name:           zone.ZoneName,
		ConnectServers: make(map[string]interface{}),
	}
	gateLua.ConnectServers["CharDB"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.CharDBPort + zoneCount,
	}
	gateLua.ConnectServers["Master"] = Connect{
		ID:   1,
		IP:   masterm.IP,
		Port: machine.MasterPort,
	}
	gateLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.LogPort + zoneCount,
	}

	trans := struct2lua.ToLuaConfig("Gate", zone.Zid, gateLua)
	if trans == false {
		fmt.Println("gate cannt wirte lua file")
	}
	return nil
}
