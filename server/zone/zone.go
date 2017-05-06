package zone

import (
	"log"

	"github.com/adamluo159/admin-react/server/comInterface"
	"github.com/adamluo159/admin-react/server/db"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ZoneMgr struct {
	machineMgr comInterface.MachineMgr
	aserver    comInterface.Aserver
}

var zMgr ZoneMgr

//区服模块注册
func Register(e *echo.Echo) *ZoneMgr {
	Str2IntChannels = make(map[string]int)
	Str2IntChannels["ios"] = 1
	Str2IntChannels["yyb"] = 2
	Str2IntChannels["xiaomi"] = 3

	cl = db.Session.DB("zone").C("zone")
	if cl == nil {
		log.Printf("cannt find Collection about zone")
		panic(0)
	}
	i := mgo.Index{
		Key:    []string{"zid"},
		Unique: true,
	}
	err := cl.EnsureIndex(i)
	if err != nil {
		log.Printf("mongodb ensureindex err:%s", err.Error())
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
		210189001,
	}

	iName := mgo.Index{
		Key:    []string{"zoneName"},
		Unique: true,
	}
	err = cl.EnsureIndex(iName)
	if err != nil {
		log.Printf("mongodb ensureindex err:%s", err.Error())
		panic(0)
	}
	e.GET("/zone", GetZones)
	e.POST("/zone/add", AddZone)
	e.POST("/zone/save", SaveZone)
	e.GET("/zone/synMachine", SynMachine)
	e.POST("/zone/startZone", StartZone)
	e.POST("/zone/stopZone", StopZone)
	e.POST("/zone/del", DelZone)
	e.POST("/zone/updateZonelogdb", UpdateZonelogdb)
	e.POST("/zone/startAllZone", StartAllZone)
	e.POST("/zone/stopAllZone", StopAllZone)

	return &zMgr
}

func (z *ZoneMgr) InitMgr(m comInterface.MachineMgr, as comInterface.Aserver) {
	z.aserver = as
	z.machineMgr = m
}

func (z *ZoneMgr) GetZoneRelation(zid int) *comInterface.RelationZone {
	dz := Zone{}
	err := cl.Find(bson.M{"zid": zid}).One(&dz)
	if err != nil {
		log.Println("GetZoneRelation err", err.Error())
		return nil
	}
	return &comInterface.RelationZone{
		Zid:           zid,
		ZoneHost:      dz.ZoneHost,
		ZoneDBHost:    dz.ZoneDBHost,
		ZonelogdbHost: dz.ZonelogdbHost,
	}
}
