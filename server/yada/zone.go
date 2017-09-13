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
		GetAllZoneInfo() *[]Zone
		GetZoneInfoByZid(zid int) *Zone
		AddZone(zone *Zone) error
		SaveZone(oldZid int, oldZoneName string, newZone *Zone) error
		DelZone(zid int) error
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

	err = z.cl.Update(query, &newZone)
	return err

	//newRelation := &comInterface.RelationZone{
	//	Zid:           m.Item.Zid,
	//	ZoneDBHost:    m.Item.ZoneDBHost,
	//	ZoneHost:      m.Item.ZoneHost,
	//	ZonelogdbHost: m.Item.ZonelogdbHost,
	//}

	//oldRelation := zMgr.GetZoneRelation(m.OldZid)
	//zMgr.machineMgr.UpdateZone(oldRelation, newRelation)

}

//删除
func (z *zoneMgr) DelZone(zid int) error {
	//dzone := Zone{}
	//err := z.cl.Find(bson.M{"zid": zid}).One(&dzone)
	//if err != nil {
	//	return err
	//}

	//r := comInterface.RelationZone{
	//	ZoneDBHost:    dzone.ZoneDBHost,
	//	ZoneHost:      dzone.ZoneHost,
	//	ZonelogdbHost: dzone.ZonelogdbHost,
	//	Zid:           dzone.Zid,
	//}
	//log.Println("delete zone :", r)
	//err = DelZoneConfig(r.Zid, r.ZoneHost)
	//if err != nil {
	//	ret.Result = "FALSE"
	//	return c.JSON(http.StatusOK, ret)
	//}
	//zMgr.machineMgr.OpZoneRelation(&r, comInterface.RelationDel)
	//err = cl.Remove(query)
	//if err != nil {
	//	ret.Result = "FALSE"
	//}
	//ret.Item = dzone
	//return c.JSON(http.StatusOK, ret)
	return z.cl.Remove(bson.M{"zid": zid})
}

//func (z *zoneMgr) UpdateZonelogdb(c echo.Context) error {
//	zReq := ZoneReq{}
//	err := c.Bind(&zReq)
//	ret := ZoneRsp{}
//	if err != nil {
//		log.Println(err.Error())
//		ret.Result = "FALSE"
//		return c.JSON(http.StatusOK, ret)
//	}
//	m := zMgr.machineMgr.GetMachineByName(zReq.Host)
//	if m == nil {
//		ret.Result = "FAlse"
//		return c.JSON(http.StatusOK, ret)
//	}
//	logdb := "zonelog" + strconv.Itoa(zReq.Zid)
//	s, _ := utils.ExeShellArgs2("sh", "update_zonelogdb", logdb, m.IP)
//	if logdb != s {
//		ret.Result = fmt.Sprintf("update zonelogdb fail,%s", s)
//	} else {
//		ret.Result = "OK"
//	}
//	return c.JSON(http.StatusOK, ret)
//}

//func (z *zoneMgr) SynMachine(c echo.Context) error {
//	ret := ZoneRsp{
//		Result: "更新失败",
//	}
//
//	zid, err := strconv.Atoi(c.QueryParam("zid"))
//	if err != nil {
//		ret.Result = "zid" + err.Error()
//		return c.JSON(http.StatusOK, ret)
//	}
//
//	hostname := c.QueryParam("hostname")
//	if hostname == "" {
//		ret.Result = "hostname none"
//		return c.JSON(http.StatusOK, ret)
//	}
//
//	WriteZoneConfigLua(zid, &ret, hostname)
//	ncode := zMgr.aserver.UpdateZone(hostname)
//	if ncode == protocol.NotifyDoSuc {
//		ret.Result = "更新成功"
//	}
//
//	return c.JSON(http.StatusOK, ret)
//}

//func (z *zoneMgr) Zonelist(c echo.Context) error {
//	var zones []Zone
//	zlst := ZoneYunYingLst{}
//	err := cl.Find(nil).All(&zones)
//	if err != nil {
//	}
//	for k := range zones {
//		z := ZoneYunYing{
//			Zid:      zones[k].Zid,
//			ZoneName: zones[k].ZoneName,
//			Channels: zones[k].Channels,
//		}
//		m := zMgr.machineMgr.GetMachineByName(zones[k].ZonelogdbHost)
//		if m == nil {
//			log.Println("zonelist cannt find machine info machineName:", zones[k].ZonelogdbHost)
//			continue
//		}
//		z.ZonelogdbIP = m.IP
//		zlst.Zlist = append(zlst.Zlist, z)
//	}
//	return c.JSON(http.StatusOK, zlst)
//}

//func (z *zoneMgr) StartZone(c echo.Context) error {
//	ret := ZoneRsp{
//		Result: "启服失败",
//	}
//	m := ZoneReq{}
//	err := c.Bind(&m)
//	if err != nil {
//		ret.Result = "启服失败, parse post info fail"
//		return c.JSON(http.StatusOK, ret)
//	}
//	zone := Zone{}
//	query := bson.M{"zid": m.Zid}
//	cl.Find(query).One(&zone)
//	if zone.ZoneHost != m.Host {
//		log.Printf(ret.Result, "send zid cannt match zonehost, zid:%d zonehost:%s", m.Zid, m.Host)
//		return c.JSON(http.StatusOK, ret)
//	}
//	s := zMgr.aserver.StartZone(m.Host, m.Zid)
//	log.Println("start result", s)
//	switch s {
//	case protocol.NotifyDoFail:
//		ret.Result = "启服失败"
//	case protocol.NotifyDoSuc:
//		ret.Result = "启服成功"
//	case protocol.NotifyDoing:
//		ret.Result = "正在启服中，请勿重复启服"
//	}
//	log.Println("start ", ret)
//	ret.Zstates = zMgr.aserver.OnlineZones()
//	return c.JSON(http.StatusOK, ret)
//}

//func (z *zoneMgr) StopZone(c echo.Context) error {
//	ret := ZoneRsp{
//		Result: "关服失败",
//	}
//	m := ZoneReq{}
//	err := c.Bind(&m)
//	if err != nil {
//		ret.Result = "关服失败, parse post info fail"
//		return c.JSON(http.StatusOK, ret)
//	}
//	zone := Zone{}
//	query := bson.M{"zid": m.Zid}
//	cl.Find(query).One(&zone)
//	if zone.ZoneHost != m.Host {
//		log.Printf(ret.Result, "stopzone send zid cannt match zonehost, zid:%d zonehost:%s", m.Zid, m.Host)
//		return c.JSON(http.StatusOK, ret)
//	}
//	s := zMgr.aserver.StopZone(m.Host, m.Zid)
//	switch s {
//	case protocol.NotifyDoFail:
//		ret.Result = "关服失败"
//	case protocol.NotifyDoSuc:
//		ret.Result = "关服成功"
//	case protocol.NotifyDoing:
//		ret.Result = "正在关服中，请勿重复关服"
//	}
//	ret.Zstates = zMgr.aserver.OnlineZones()
//	return c.JSON(http.StatusOK, ret)
//}

//func (z *zoneMgr) StartAllZone(c echo.Context) error {
//	ret := ZoneRsp{}
//	s := zMgr.aserver.StartAllZone()
//	switch s {
//	case protocol.NotifyDoFail:
//		ret.Result = "全服启动失败"
//	case protocol.NotifyDoSuc:
//		ret.Result = "全服启动成功"
//	case protocol.NotifyDoing:
//		ret.Result = "正在启服中，请勿重复启服"
//	}
//	ret.Zstates = zMgr.aserver.OnlineZones()
//	log.Println("startaaaaa:", ret)
//	return c.JSON(http.StatusOK, ret)
//}

//func (z *zoneMgr) StopAllZone(c echo.Context) error {
//	ret := ZoneRsp{}
//	s := zMgr.aserver.StopAllZone()
//	switch s {
//	case protocol.NotifyDoFail:
//		ret.Result = "全服关闭失败"
//	case protocol.NotifyDoSuc:
//		ret.Result = "全服关闭成功"
//	case protocol.NotifyDoing:
//		ret.Result = "正在关服中，请勿重复关服"
//	}
//	ret.Zstates = zMgr.aserver.OnlineZones()
//	log.Println("stopaaaaa:", ret)
//	return c.JSON(http.StatusOK, ret)
//}

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
