package zone

import (
	"fmt"
	"net/http"

	"github.com/adamluo159/admin-react/server/db"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	//if m.OldZoneName != m.Item.ZoneName || m.OldZid != m.Item.Zid {
	//	del := bson.M{"zid": m.OldZid, "zoneName": m.OldZoneName}
	//	err = cl.Remove(del)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//	err = cl.Insert(m.Item)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		ret.Result = "FALSE"
	//	} else {
	//		ret.Item = m.Item
	//	}
	//} else {
	//	query := bson.M{"zid": m.Item.Zid, "zoneName": m.Item.ZoneName}
	//	err = cl.Update(query, &m.Item)
	//	if err != nil {
	//		fmt.Println("SaveZone, update:", err.Error())
	//		ret.Result = "FALSE"
	//	} else {
	//		ret.Item = m.Item
	//	}
	//}

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
	//e.POST("/zone/del", DelZone)
}