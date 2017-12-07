package yada

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//区服信息
type (
	Zone struct {
		Zid           int      `bson:"zid" json:"zid"`
		ZoneName      string   `json:"zoneName" bson:"zoneName"`
		ZoneHost      string   `json:"zoneHost" bson:"zoneHost"`
		ZoneDBHost    string   `json:"zoneDBHost" bson:"zoneDBHost"`
		ZonelogdbHost string   `json:"zonelogdbHost" bson:"zonelogdbHost"`
		DatalogdbHost string   `json:"datalogdbHost" bson:"datalogdbHost"`
		Channels      []string `json:"channels" bson:"channels"`
		Whitelst      bool     `json:"whitelst" bson:"whitelst"`
		PortNumber    int      `json:"portNumber" bson:"PortNumber"`
		OpenTime      string   `json:"opentime" bson:"opentime"`
	}

	PortCount struct {
		Host  string `bson:"host" json:"host"`
		Count int    `bson:count json:count`
	}

	ZoneYunYing struct {
		Zid         int      `bson:"zid" json:"zid"`
		ZoneName    string   `json:"zoneName" bson:"zoneName"`
		ZonelogdbIP string   `json:"zonelogdbIP" bson:"zonelogdbIP"`
		Channels    []string `json:"channels" bson:"channels"`
	}
	ZoneYunYingLst struct {
		Zlist []ZoneYunYing
	}

	ZoneMgr interface {
		//获取所有服的信息
		GetAllZoneInfo() *[]Zone

		//通过zid找服信息
		GetZoneInfoByZid(zid int) *Zone

		//新增服信息
		AddZone(zone *Zone) error

		//保存更新服信息
		SaveZone(oldZid int, oldZoneName string, newZone *Zone) error

		//删除服信息
		DelZone(zid int) error

		//获取服相关机器信息
		GetZoneRelation(zid int) *RelationZone
	}

	zoneMgr struct {
		cl     *mgo.Collection
		clPort *mgo.Collection
	}
)

func NewZoneMgr(session *mgo.Session) ZoneMgr {
	z := zoneMgr{
		cl:     session.DB("gameAdmin").C("zone"),
		clPort: session.DB("gameAdmin").C("Port"),
	}
	if z.cl == nil || z.clPort == nil {
		log.Fatalf(" zone collecttion nil, cl:%v, clPort:%v", z.cl, z.clPort)
	}

	err := z.cl.EnsureIndex(mgo.Index{
		Key:    []string{"zid", "zoneName"},
		Unique: true,
	})

	errPort := z.clPort.EnsureIndex(mgo.Index{
		Key:    []string{"host"},
		Unique: true,
	})
	if err != nil || errPort != nil {
		log.Fatalf(" zone collecttion ensureIndex err, cl:%v, clPort:%v", err, errPort)
	}

	return &z
}

func (z *zoneMgr) GetZoneRelation(zid int) *RelationZone {
	dz := Zone{}
	err := z.cl.Find(bson.M{"zid": zid}).One(&dz)
	if err != nil {
		log.Println("GetZoneRelation err", err.Error())
		return nil
	}
	return &RelationZone{
		Zid:           zid,
		ZoneHost:      dz.ZoneHost,
		ZoneDBHost:    dz.ZoneDBHost,
		ZonelogdbHost: dz.ZonelogdbHost,
		DatalogdbHost: dz.DatalogdbHost,
	}
}

//获取所有区服的配置信息
func (z *zoneMgr) GetAllZoneInfo() *[]Zone {
	zones := []Zone{}
	z.cl.Find(nil).All(&zones)
	return &zones
}

//获取所有区服的配置信息
func (z *zoneMgr) GetZoneInfoByZid(zid int) *Zone {
	zone := Zone{}
	if err := z.cl.Find(bson.M{"zid": zid}).One(&zone); err != nil {
		log.Println("getzoneinfo by zid err, ", err)
		return nil
	}
	return &zone
}

func (z *zoneMgr) GetZoneInfoByName(zoneName string) *Zone {
	zone := Zone{}
	err := z.cl.Find(bson.M{"zoneName": zoneName}).One(&zone)
	if err != nil {
		log.Println("getzoneinfo by name err, ", err)
		return nil
	}
	return &zone
}

//添加区服信息
func (z *zoneMgr) AddZone(zone *Zone) error {
	n, err := z.ReqPortCount(zone.ZoneHost, zone.Zid)
	if err != nil {
		return err
	}
	zone.PortNumber = n
	err = z.cl.Insert(zone)
	return err
}

//保存
func (z *zoneMgr) SaveZone(oldZid int, oldZoneName string, newZone *Zone) error {
	query := bson.M{"zid": oldZid, "zoneName": oldZoneName}
	oldZone := Zone{}
	err := z.cl.Find(query).One(&oldZone)

	if err != nil {
		return err
	}
	newZone.PortNumber = oldZone.PortNumber
	if newZone.ZoneHost != oldZone.ZoneHost || newZone.PortNumber <= 0 {
		n, rerr := z.ReqPortCount(newZone.ZoneHost, newZone.Zid)
		if rerr != nil {
			return rerr
		}
		newZone.PortNumber = n
	}

	return z.cl.Update(query, &newZone)
}

//删除
func (z *zoneMgr) DelZone(zid int) error {
	return z.cl.Remove(bson.M{"zid": zid})
}

func (z *zoneMgr) ReqPortCount(host string, zid int) (int, error) {
	c := bson.M{"host": host}
	item := PortCount{}

	err := z.clPort.Find(c).One(&item)
	if err != nil {
		item.Host = host
		item.Count = 1
		err = z.clPort.Insert(item)
	} else {
		item.Count++
		err = z.clPort.Update(c, &item)
	}
	if err != nil {
		log.Println(" ReqPortCount bbbb, ", err.Error(), host, zid)
		return 0, err
	}
	return item.Count, err
}
