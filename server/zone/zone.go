package zone

import (
	"fmt"
	"log"

	"errors"

	"github.com/adamluo159/admin-react/server/comInterface"
	"github.com/adamluo159/admin-react/server/db"
	"github.com/adamluo159/gameAgent/utils"
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
	cl = db.Session.DB("gameAdmin").C("zone")
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

	clPort = db.Session.DB("gameAdmin").C("portCount")
	if cl == nil {
		log.Printf("cannt find Collection about portCount")
		panic(0)
	}
	w := mgo.Index{
		Key:    []string{"host"},
		Unique: true,
	}
	werr := clPort.EnsureIndex(w)
	if err != nil {
		log.Printf("mongodb ensureindex err:%s", werr.Error())
		panic(0)
	}

	iName := mgo.Index{
		Key:    []string{"zoneName"},
		Unique: true,
	}
	err = cl.EnsureIndex(iName)
	if err != nil {
		panic(fmt.Sprintf("mongodb ensureindex err:%s", err.Error()))
	}

	err = LoadZoneConfig()
	if err != nil {
		panic(err.Error())
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
	e.GET("/zone/zonelist", Zonelist)

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

func LoadZoneConfig() error {
	utils.LoadConfigJson()

	Str2IntChannels = make(map[string]int)
	jsonerr := utils.GetConfigMap("channels", &Str2IntChannels)
	if jsonerr != nil {
		return errors.New(fmt.Sprintf("zone register load channels json file fail %v\n", jsonerr))
	}
	LogicMap = make(map[int][]int)
	logic1 := []int{}
	logic2 := []int{}
	jsonerr = utils.GetConfigArray("logic1_maps", &logic1)
	if jsonerr != nil {
		return errors.New(fmt.Sprintf("zone register load logic1 maps json file fail %v\n", jsonerr))
	}
	jsonerr = utils.GetConfigArray("logic2_maps", &logic2)
	if jsonerr != nil {
		return errors.New(fmt.Sprintf("zone register load logic2 maps json file fail %v\n", jsonerr))
	}
	LogicMap[1] = logic1
	LogicMap[2] = logic2
	return nil
}
