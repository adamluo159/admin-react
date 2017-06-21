package machine

import (
	"fmt"
	"log"
	"strconv"

	"github.com/adamluo159/admin-react/server/comInterface"
	"github.com/adamluo159/admin-react/server/db"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	cl    *mgo.Collection
	mhMgr MachineMgr
)

type MachineMgr struct {
	as comInterface.Aserver
}

//机器模块注册
func Register(e *echo.Echo) *MachineMgr {
	cl = db.Session.DB("gameAdmin").C("machine")
	if cl == nil {
		fmt.Printf("cannt find Collection about machine")
		panic(0)
	}
	i := mgo.Index{
		Key:    []string{"hostname"},
		Unique: true,
	}
	err := cl.EnsureIndex(i)
	if err != nil {
		fmt.Printf("mongodb ensureindex err:%s", err.Error())
		panic(0)
	}
	e.GET("/machine", GetMachines)
	e.POST("/machine/add", AddMachine)
	e.POST("/machine/save", SaveMachine)
	e.POST("/machine/del", DelMachine)
	e.GET("machine/common", CommonConfig)
	return &mhMgr
}

func (m *MachineMgr) InitMgr(as comInterface.Aserver) {
	m.as = as
}

func UpdateMachineApplications(host string, apps []string) {
	err := cl.Update(bson.M{"hostname": host}, bson.M{"$set": bson.M{"applications": apps}})
	if err != nil {
		log.Println("UpdateMachineApplications update err, ", err.Error())
		return
	}
}

func SliceString(A *[]string, name string, op int) {
	index := -1
	for i := range *A {
		if name == (*A)[i] {
			index = i
			break
		}
	}
	if (index == -1 && comInterface.RelationDel == op) ||
		(index >= 0 && comInterface.RelationAdd == op) {
		return
	}

	switch op {
	case comInterface.RelationDel:
		(*A) = append((*A)[:index], (*A)[index+1:]...)
	case comInterface.RelationAdd:
		(*A) = append((*A), name)
	default:
		log.Println("SliceString op wrong ", op)
	}
}

func (m *MachineMgr) GetMachineByName(name string) *comInterface.Machine {
	d := comInterface.Machine{}
	err := cl.Find(bson.M{"hostname": name}).One(&d)
	if err != nil {
		log.Println(" GetMachineByName name: ", name, err.Error())
		return nil
	}
	return &d
}

func (m *MachineMgr) GetAllMachines() []comInterface.Machine {
	var ms []comInterface.Machine
	err := cl.Find(nil).All(&ms)
	if err != nil {
		log.Println(" GetAllMachines", err.Error())
		return nil
	}
	return ms
}

func (m *MachineMgr) UpdateZone(old *comInterface.RelationZone, new *comInterface.RelationZone) {
	if old == nil || new == nil {
		log.Println("machine Relation UpdateZone old or new is nil", old, new)
		return
	}
	log.Println("update:", *old, *new)
	m.OpZoneRelation(old, comInterface.RelationDel)
	m.OpZoneRelation(new, comInterface.RelationAdd)
}

func (m *MachineMgr) OpZoneRelation(r *comInterface.RelationZone, op int) {
	z := m.GetMachineByName((*r).ZoneHost)
	if z != nil {
		name := "zone" + strconv.Itoa((*r).Zid)
		SliceString(&z.Applications, name, op)
		UpdateMachineApplications(z.Hostname, z.Applications)
	}
	db := m.GetMachineByName((*r).ZoneDBHost)
	if db != nil {
		name := "zonedb" + strconv.Itoa((*r).Zid)
		SliceString(&db.Applications, name, op)
		UpdateMachineApplications(db.Hostname, db.Applications)
	}
	logdb := m.GetMachineByName((*r).ZonelogdbHost)
	if logdb != nil {
		name := "zonelogdb" + strconv.Itoa((*r).Zid)
		SliceString(&logdb.Applications, name, op)
		UpdateMachineApplications(logdb.Hostname, logdb.Applications)
	}
}
