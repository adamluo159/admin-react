package zone

import (
	"fmt"
	"net/http"
	"os"

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

var (
	cl       *mgo.Collection
	GlobalDB MysqlLua
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
	os.Mkdir(machine.GameConfigDir, os.ModeDir)
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

	zonequery := bson.M{"zoneHost": zone.ZoneHost}
	zoneCount, zerr := cl.Find(zonequery).Count()
	if zerr != nil {
		return zerr
	}
	zonem, err := machine.GetMachineByName(zone.ZoneHost)
	if err != nil {
		return err
	}

	dir := machine.GameConfigDir + "zone" + strconv.Itoa(zone.Zid)
	os.Mkdir(dir, os.ModeDir)
	curDir := dir + "/"
	gerr := GateLua(&zone, zonem, zoneCount, curDir)
	if gerr != nil {
		ret.Result = gerr.Error()
		return c.JSON(http.StatusOK, ret)
	}
	cerr := CenterLua(&zone, zonem, zoneCount, curDir)
	if cerr != nil {
		ret.Result = cerr.Error()
		return c.JSON(http.StatusOK, ret)
	}
	lerr := LogLua(&zone, zonem, zoneCount, curDir)
	if lerr != nil {
		ret.Result = lerr.Error()
		return c.JSON(http.StatusOK, ret)
	}
	logicerr := LogicLua(&zone, zonem, zoneCount, curDir)
	if logicerr != nil {
		ret.Result = logicerr.Error()
		return c.JSON(http.StatusOK, ret)
	}

	return c.JSON(http.StatusOK, ret)
}

func GateLua(zone *Zone, zonem *machine.Machine, zoneCount int, Dir string) error {
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
	trans := struct2lua.ToLuaConfig(Dir, "Gate", gateLua)
	if trans == false {
		fmt.Println("gate cannt wirte lua file")
	}
	return nil
}

func CenterLua(zone *Zone, zonem *machine.Machine, zoneCount int, Dir string) error {
	centerLua := Center{
		ID:   zone.Zid,
		Zid:  zone.Zid,
		IP:   zonem.IP,
		Port: machine.CenterPort + zoneCount,
		OnlineNumberCheckTime: 60 * 5,
		SingleServerLoad:      4000,
		ConnectServers:        make(map[string]interface{}),
	}

	centerLua.ConnectServers["CharDB"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.CharDBPort + zoneCount,
	}
	centerLua.ConnectServers["Gate"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.GatePort + zoneCount,
	}
	centerLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.LogPort + zoneCount,
	}

	trans := struct2lua.ToLuaConfig(Dir, "Center", centerLua)
	if trans == false {
		fmt.Println("center cannt wirte lua file")
	}
	return nil
}

func CharDBLua(zone *Zone, zonem *machine.Machine, zoneCount int, Dir string) error {
	zoneDBquery := bson.M{"zoneDBHost": zone.ZoneDBHost}
	zoneDBCount, zdberr := cl.Find(zoneDBquery).Count()
	if zdberr != nil {
		return zdberr
	}
	zonedbm, dberr := machine.GetMachineByName(zone.ZoneDBHost)
	if dberr != nil {
		return dberr
	}

	mysqldbName := "zone" + strconv.Itoa(zone.Zid)
	charDBLua := CharDB{
		ID:   zone.Zid,
		Zid:  zone.Zid,
		IP:   zonem.IP,
		Port: machine.CharDBPort + zoneCount,
		Mysql: MysqlLua{
			IP:             zonedbm.IP,
			Port:           machine.MysqlPort,
			UserName:       machine.UserName,
			PassWord:       machine.PassWord,
			FlushFrequency: 300,
			DataBase:       mysqldbName,
		},
		Redis: RedisLua{
			IP:       zonedbm.IP,
			Port:     machine.RedisPort + zoneDBCount,
			Password: machine.PassWord,
		},
	}
	trans := struct2lua.ToLuaConfig(Dir, "CharDB", charDBLua)
	if trans == false {
		fmt.Println("chardb cannt wirte lua file")
	}
	return nil
}

func LogicLua(zone *Zone, zonem *machine.Machine, zoneCount int, Dir string) error {
	logicLua := Logic{
		ID:             1,
		Zid:            zone.Zid,
		IP:             zonem.IP,
		Port:           machine.LogicPort + zoneCount*3 + 1,
		ConnectServers: make(map[string]interface{}),
	}
	logicLua.ConnectServers["CharDB"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.CharDBPort + zoneCount,
	}
	logicLua.ConnectServers["Gate"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.LogicPort + zoneCount,
	}
	logicLua.ConnectServers["Center"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.CenterPort + zoneCount,
	}
	logicLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.LogPort + zoneCount,
	}

	trans := struct2lua.ToLuaConfig(Dir, "Logic", logicLua)
	if trans == false {
		fmt.Println("logic cannt wirte lua file")
	}

	return nil
}

func LogLua(zone *Zone, zonem *machine.Machine, zoneCount int, Dir string) error {
	logm, logerr := machine.GetMachineByName(zone.ZonelogdbHost)
	if logerr != nil {
		return logerr
	}
	logLua := Log{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.LogPort + zoneCount,
		ZoneLogMysql: MysqlLua{
			IP:             logm.IP,
			Port:           machine.MysqlPort,
			UserName:       machine.UserName,
			PassWord:       machine.PassWord,
			FlushFrequency: 300,
			DataBase:       "zonelog" + strconv.Itoa(zone.Zid),
		},
		GlobalLogMysql: GlobalDB,
	}
	trans := struct2lua.ToLuaConfig(Dir, "Log", logLua)
	if trans == false {
		fmt.Println("log cannt wirte lua file")
	}

	return nil
}
