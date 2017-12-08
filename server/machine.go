package yada

import (
	"log"
	"strconv"
	"strings"

	"github.com/adamluo159/admin-react/utils"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	RelationZone struct {
		Zid           int
		ZoneHost      string
		ZoneDBHost    string
		ZonelogdbHost string
		DatalogdbHost string
		ZonedbBakHost string
	}

	//机器信息
	Machine struct {
		Hostname     string `json:"hostname" bson:"hostname"`
		IP           string
		OutIP        string   `json:"outIP" bson:"outIP"`
		Applications []string `json:"applications" bson:"applications"`
		Online       bool
		CodeVersion  string `json:"codeVersion" bson:"codeVersion"`
	}

	MachineMgr interface {
		//删除机器配置
		DelMachine(hostname string) error
		//保存机器配置
		SaveMachine(oldhost string, m *Machine) error
		//新增机器配置
		AddMachine(m *Machine) error
		//获取所有机器配置
		GetAllMachines() []Machine
		//获取所有机器名
		GetAllMachineName() []string

		//查找机器配置
		GetMachineByName(name string) *Machine
		//更新机器svn版本号
		UpdateSvnVersion(host string, ver string) error

		//更新用途关系
		UpdateZone(old *RelationZone, new *RelationZone)
		//删除游戏服配置
		DelZoneConf(zr *RelationZone) error
		//机器用途关系设置
		OpZoneRelation(r *RelationZone, op int)

		//写游戏服配置文件
		ZoneLua(zone *Zone, Dir string) error
		GateLua(zone *Zone, Dir string, arrayClientPorts *[]int) error
		CenterLua(zone *Zone, Dir string) error
		LogLua(zone *Zone, Dir string) error
		LogicLua(zone *Zone, Dir string) error
		CharDBLua(zone *Zone, Dir string) error
		LoginLua() error
		AccountDBLua() error
		MasterLogLua() error
		MasterLua() error
	}

	machineMgr struct {
		cl   *mgo.Collection
		conf Conf
	}
)

const (
	RelationDel int = 1
	RelationAdd int = 2
)

//机器模块注册
func NewMachineMgr(session *mgo.Session, mconf Conf) MachineMgr {
	mcl := session.DB("gameAdmin").C("machine")
	if mcl == nil {
		log.Fatal("cannt find Collection about machine")
	}
	if err := mcl.EnsureIndex(mgo.Index{Key: []string{"hostname"}, Unique: true}); err != nil {
		log.Fatalf("mongodb ensureindex err:%s", err.Error())
	}
	return &machineMgr{
		cl:   mcl,
		conf: mconf,
	}
}

func (m *machineMgr) UpdateMachineApps(host string, name string, op int) {

	if name == "" {
		log.Println("update Machine apps, name", name, host, op)
		return
	}

	updateM := make(bson.M)
	s := bson.M{"applications": name}

	switch op {
	case RelationDel:
		updateM["$pull"] = s
	case RelationAdd:
		updateM["$push"] = s
	default:
		log.Println("UpdateMachineApps op wrong ", op)
	}

	h := bson.M{"hostname": host}

	if err := m.cl.Update(h, updateM); err != nil {
		log.Println(" UpdateMachineApps err, ", err.Error())
	}
}

func (m *machineMgr) UpdateSvnVersion(host string, ver string) error {
	h := bson.M{"hostname": host}
	setv := bson.M{"$set": bson.M{"codeVersion": ver}}
	return m.cl.Update(h, setv)

}

func (m *machineMgr) GetMachineByName(name string) *Machine {
	d := Machine{}
	err := m.cl.Find(bson.M{"hostname": name}).One(&d)
	if err != nil {
		log.Println(" GetMachineByName name: ", name, err.Error())
		return nil
	}
	return &d
}

func (m *machineMgr) GetAllMachines() []Machine {
	ms := []Machine{}
	err := m.cl.Find(nil).All(&ms)
	if err != nil {
		log.Println(" GetAllMachines", err.Error())
		return nil
	}
	return ms
}

func (m *machineMgr) GetAllMachineName() []string {
	names := make([]string, 0)
	ms := m.GetAllMachines()
	if ms != nil {
		for _, v := range ms {
			if strings.Contains(v.Hostname, "cghost") {
				names = append(names, v.Hostname)
			}
		}
	}

	return names
}

func (m *machineMgr) UpdateZone(old *RelationZone, new *RelationZone) {
	if old == nil || new == nil {
		log.Println("machine Relation UpdateZone old or new is nil", old, new)
		return
	}
	log.Println("update:", *old, *new)
	m.OpZoneRelation(old, RelationDel)
	m.OpZoneRelation(new, RelationAdd)
}

func (m *machineMgr) OpZoneRelation(r *RelationZone, op int) {
	z := m.GetMachineByName((*r).ZoneHost)
	if z != nil {
		name := "zone" + strconv.Itoa((*r).Zid)
		m.UpdateMachineApps(z.Hostname, name, op)
	}
	db := m.GetMachineByName((*r).ZoneDBHost)
	if db != nil {
		name := "zonedb" + strconv.Itoa((*r).Zid)
		m.UpdateMachineApps(db.Hostname, name, op)
	}
	logdb := m.GetMachineByName((*r).ZonelogdbHost)
	if logdb != nil {
		name := "zonelogdb" + strconv.Itoa((*r).Zid)
		m.UpdateMachineApps(logdb.Hostname, name, op)
	}
	datalogdb := m.GetMachineByName((*r).DatalogdbHost)
	if datalogdb != nil {
		name := "datalogdb" + strconv.Itoa((*r).Zid)
		m.UpdateMachineApps(datalogdb.Hostname, name, op)
	}
	dbbak := m.GetMachineByName((*r).ZonedbBakHost)
	if dbbak != nil {
		name := "zonedbBak" + strconv.Itoa((*r).Zid)
		m.UpdateMachineApps(dbbak.Hostname, name, op)
	}
}

//添加机器信息
func (m *machineMgr) AddMachine(machine *Machine) error {
	return m.cl.Insert(*machine)
}

//保存
func (m *machineMgr) SaveMachine(oldhost string, machine *Machine) error {
	if oldhost != machine.Hostname {
		if err := m.cl.Remove(bson.M{"hostname": oldhost}); err != nil {
			return err
		}
		if err := m.cl.Insert(machine); err != nil {
			return err
		}

	} else {
		if err := m.cl.Update(bson.M{"hostname": machine.Hostname}, machine); err != nil {
			return err
		}
	}
	return nil
}

//删除
func (m *machineMgr) DelMachine(hostname string) error {
	return m.cl.Remove(bson.M{"hostname": hostname})
}

func (m *machineMgr) DelZoneConf(zr *RelationZone) error {
	dir := zr.ZoneHost + "/zone" + strconv.Itoa(zr.Zid)
	if _, err := utils.ExeShell("sh", m.conf.GitCommit, dir); err != nil {
		return err
	}
	m.OpZoneRelation(zr, RelationDel)

	return nil
}
