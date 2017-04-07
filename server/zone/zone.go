package zone

import (
	"log"

	"github.com/adamluo159/admin-react/server/db"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
)

//区服模块注册
func Register(e *echo.Echo) {
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
	//e.POST("/zone/del", DelZone)
}
