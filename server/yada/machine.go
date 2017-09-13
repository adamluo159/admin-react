package yada

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/adamluo159/gameAgent/utils"
	"github.com/adamluo159/struct2lua"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	RelationZone struct {
		Zid           int
		ZoneHost      string
		ZoneDBHost    string
		ZonelogdbHost string
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
		DelMachine(hostname string) error
		SaveMachine(oldhost string, m *Machine) error
		AddMachine(m *Machine) error
		GetAllMachines() []Machine
		UpdateZone(old *RelationZone, new *RelationZone)
		GetMachineByName(name string) *Machine

		DelZoneConf(zr *RelationZone) error

		OpZoneRelation(r *RelationZone, op int)

		ZoneLua(zone *Zone, Dir string) error
		GateLua(zone *Zone, Dir string) error
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

func (m *machineMgr) UpdateMachineApplications(host string, apps []string) {
	err := m.cl.Update(bson.M{"hostname": host}, bson.M{"$set": bson.M{"applications": apps}})
	if err != nil {
		log.Println("UpdateMachineApplications update err, ", err.Error())
		return
	}
}

func (m *machineMgr) SliceString(A *[]string, name string, op int) {
	index := -1
	for i := range *A {
		if name == (*A)[i] {
			index = i
			break
		}
	}
	if (index == -1 && RelationDel == op) || (index >= 0 && RelationAdd == op) {
		return
	}

	switch op {
	case RelationDel:
		(*A) = append((*A)[:index], (*A)[index+1:]...)
	case RelationAdd:
		(*A) = append((*A), name)
	default:
		log.Println("SliceString op wrong ", op)
	}
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
		m.SliceString(&z.Applications, name, op)
		m.UpdateMachineApplications(z.Hostname, z.Applications)
	}
	db := m.GetMachineByName((*r).ZoneDBHost)
	if db != nil {
		name := "zonedb" + strconv.Itoa((*r).Zid)
		m.SliceString(&db.Applications, name, op)
		m.UpdateMachineApplications(db.Hostname, db.Applications)
	}
	logdb := m.GetMachineByName((*r).ZonelogdbHost)
	if logdb != nil {
		name := "zonelogdb" + strconv.Itoa((*r).Zid)
		m.SliceString(&logdb.Applications, name, op)
		m.UpdateMachineApplications(logdb.Hostname, logdb.Applications)
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
	dir := zr.ZoneHost + "/" + "zone" + strconv.Itoa(zr.Zid)
	if _, err := utils.ExeShell("sh", m.conf.GitCommit, dir); err != nil {
		return err
	}
	m.OpZoneRelation(zr, RelationDel)

	return nil
}

func (m *machineMgr) ZoneLua(zone *Zone, Dir string) error {
	zonem := m.GetMachineByName(zone.ZoneHost)
	if zonem == nil {
		return errors.New(fmt.Sprintf("zone machine info err"))
	}
	masterm := m.GetMachineByName("master")
	if zonem == nil {
		return errors.New(fmt.Sprintf("master machine info err"))
	}

	s := make([]int, len(zone.Channels))
	for k, v := range zone.Channels {
		s[k] = m.conf.Channels[v]
	}

	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(longForm, zone.OpenTime, loc)

	zoneLua := ZoneConf{
		ID:             zone.Zid,
		Zid:            zone.Zid,
		ServerIP:       zonem.IP,
		ServerPort:     ZonePort + zone.PortNumber,
		ClientIP:       zonem.OutIP,
		ClientPort:     ZoneClientPort + zone.PortNumber,
		ChannelIds:     s,
		Open:           zone.Whitelst,
		Name:           zone.ZoneName,
		OpenTime:       theTime.Unix(),
		ConnectServers: make(map[string]interface{}),
	}
	zoneLua.ConnectServers["Master"] = Connect{
		ID:   1,
		IP:   masterm.IP,
		Port: MasterPort + MasterCount,
	}
	zoneLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: LogPort + zone.PortNumber,
	}
	srvHead := ServerConfigHead{
		NET_TIMEOUT_MSEC:  30000,
		NET_MAX_CONNETION: NetMaxConnection,
		StartService:      []map[string]int{{"nType": ZoneServer}},
		LOG_INDEX:         "zone",
		LOG_MAXLINE:       LogMaxLine,
		OpenGM:            1,
	}
	trans := struct2lua.ToLuaConfig(Dir, "Zone", zoneLua, srvHead, 0)
	if trans == false {
		log.Println("gate cannt wirte lua file")
	}
	return nil
}

func (m *machineMgr) GateLua(zone *Zone, Dir string) error {
	zonem := m.GetMachineByName(zone.ZoneHost)
	if zonem == nil {
		return errors.New(fmt.Sprintf("gateLua zone machine info err"))
	}

	gateLua := GateConf{
		Zid:            zone.Zid,
		ServerIP:       zonem.IP,
		ClientIP:       zonem.OutIP,
		ConnectServers: make(map[string]interface{}),
	}
	gateLua.ConnectServers["Zone"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: ZonePort + zone.PortNumber,
	}
	gateLua.ConnectServers["Log"] = Connect{
		ID:   1,
		IP:   zonem.IP,
		Port: LogPort + zone.PortNumber,
	}
	srvHead := ServerConfigHead{
		NET_TIMEOUT_MSEC:  30000,
		NET_MAX_CONNETION: NetMaxConnection,
		StartService:      []map[string]int{{"nType": GateServer}},
		LOG_INDEX:         "zone",
		LOG_MAXLINE:       LogMaxLine,
		OpenGM:            1,
	}

	for i := 1; i <= m.conf.GateCount; i++ {
		gateLua.ID = i
		gateLua.ClientPort = ClientPort + zone.PortNumber*10 + i - 1
		gateLua.ServerPort = GatePort + zone.PortNumber*10 + i - 1
		srvHead.LOG_INDEX = "gate" + strconv.Itoa(i)
		trans := struct2lua.ToLuaConfig(Dir, "Gate", gateLua, srvHead, i)
		if trans == false {
			log.Printf("gate cannt wirte lua file, gateid:%d\n", i)
		}
	}

	return nil
}

func (m *machineMgr) CenterLua(zone *Zone, Dir string) error {
	zonem := m.GetMachineByName(zone.ZoneHost)
	if zonem == nil {
		return errors.New(fmt.Sprintf("centerLua zone machine info err"))
	}

	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(longForm, zone.OpenTime, loc)
	centerLua := CenterConf{
		ID:               1,
		Zid:              zone.Zid,
		IP:               zonem.IP,
		Port:             CenterPort + zone.PortNumber,
		SingleServerLoad: 7000,
		ConnectServers:   make(map[string]interface{}),
		OpenTime:         theTime.Unix(),
	}
	centerLua.ConnectServers["CharDB"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: CharDBPort + zone.PortNumber,
	}
	gateArray := make([]Connect, m.conf.GateCount)
	for i := 1; i <= m.conf.GateCount; i++ {
		gateArray[i-1] = Connect{
			ID:   i,
			IP:   zonem.IP,
			Port: GatePort + zone.PortNumber*10 + i - 1,
		}
	}
	centerLua.ConnectServers["Gate"] = gateArray
	centerLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: LogPort + zone.PortNumber,
	}
	centerLua.ConnectServers["WebOp"] = Connect{
		ID:   1,
		IP:   zonem.IP,
		Port: LogPort + zone.PortNumber,
	}
	srvHead := ServerConfigHead{
		NET_TIMEOUT_MSEC:  30000,
		NET_MAX_CONNETION: NetMaxConnection,
		StartService:      []map[string]int{{"nType": CenterServer}},
		LOG_INDEX:         "center" + strconv.Itoa(zone.Zid),
		LOG_MAXLINE:       LogMaxLine,
		OpenGM:            1,
	}

	trans := struct2lua.ToLuaConfig(Dir, "Center", centerLua, srvHead, 0)
	if trans == false {
		log.Println("center cannt wirte lua file")
	}
	return nil
}

func (m *machineMgr) CharDBLua(zone *Zone, Dir string) error {
	//这个后面再加，现在假设所有服都用6379作缓存
	//zoneDBquery := bson.M{"zoneDBHost": zone.ZoneDBHost}
	//zoneDBCount, zdberr := cl.Find(zoneDBquery).Count()
	//if zdberr != nil {
	//	return zdberr
	//}

	zonem := m.GetMachineByName(zone.ZoneHost)
	if zonem == nil {
		return errors.New(fmt.Sprintf("chardbLua zone machine info err"))
	}
	zonedb := m.GetMachineByName(zone.ZoneDBHost)
	if zonem == nil {
		return errors.New(fmt.Sprintf("chardbLua db machine info err"))
	}

	mysqldbName := "cgzone" + strconv.Itoa(zone.Zid)
	charDBLua := CharDBConf{
		ID:   zone.Zid,
		Zid:  zone.Zid,
		IP:   zonem.IP,
		Port: CharDBPort + zone.PortNumber,
		Mysql: MysqlLua{
			IP:             zonedb.IP,
			Port:           MysqlPort,
			UserName:       UserName,
			Password:       PassWord,
			FlushFrequency: 300,
			DataBase:       mysqldbName,
		},
		Redis: RedisLua{
			IP:       zonedb.IP,
			Port:     RedisPort,
			Password: RedisPassWord,
		},
		ConnectServers: make(map[string]interface{}),
	}

	charDBLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: LogPort + zone.PortNumber,
	}
	srvHead := ServerConfigHead{
		NET_TIMEOUT_MSEC:  30000,
		NET_MAX_CONNETION: NetMaxConnection,
		StartService:      []map[string]int{{"nType": DbproxyServer}},
		LOG_INDEX:         "charDB" + strconv.Itoa(zone.Zid),
		LOG_MAXLINE:       LogMaxLine,
		OpenGM:            1,
	}

	trans := struct2lua.ToLuaConfig(Dir, "CharDB", charDBLua, srvHead, 0)
	if trans == false {
		log.Println("chardb cannt wirte lua file")
	}
	return nil
}

func (m *machineMgr) LogicLua(zone *Zone, Dir string) error {
	zonem := m.GetMachineByName(zone.ZoneHost)
	if zonem == nil {
		return errors.New(fmt.Sprintf("centerLua zone machine info err"))
	}

	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(longForm, zone.OpenTime, loc)
	logicLua := LogicConf{
		Zid:            zone.Zid,
		IP:             zonem.IP,
		ConnectServers: make(map[string]interface{}),
		LoadAllMapIds:  false,
		OpenTime:       theTime.Unix(),
	}
	logicLua.ConnectServers["CharDB"] = Connect{
		ID:   1,
		IP:   zonem.IP,
		Port: CharDBPort + zone.PortNumber,
	}
	gateArray := make([]Connect, m.conf.GateCount)
	for i := 1; i <= m.conf.GateCount; i++ {
		gateArray[i-1] = Connect{
			ID:   i,
			IP:   zonem.IP,
			Port: GatePort + zone.PortNumber*10 + i - 1,
		}
	}
	logicLua.ConnectServers["Gate"] = gateArray
	logicLua.ConnectServers["Center"] = Connect{
		ID:   1,
		IP:   zonem.IP,
		Port: CenterPort + zone.PortNumber,
	}
	logicLua.ConnectServers["Log"] = Connect{
		ID:   1,
		IP:   zonem.IP,
		Port: LogPort + zone.PortNumber,
	}

	srvHead := ServerConfigHead{
		NET_TIMEOUT_MSEC:  30000,
		NET_MAX_CONNETION: NetMaxConnection,
		StartService:      []map[string]int{{"nType": LogicServer}},
		LOG_MAXLINE:       LogMaxLine,
		OpenGM:            1,
	}

	for i := 1; i <= m.conf.LogicCount; i++ {
		logicLua.ID = i
		logicLua.Port = LogicPort + 10*zone.PortNumber + i - 1
		srvHead.LOG_INDEX = "logic" + strconv.Itoa(i)
		trans := struct2lua.ToLuaConfig(Dir, "Logic", logicLua, srvHead, i)
		if trans == false {
			log.Printf("logic:%d cannt wirte lua file\n", i)
		}
	}

	return nil
}

func (m *machineMgr) LogLua(zone *Zone, Dir string) error {
	errCollect := m.GetMachineByName("errLog")
	if errCollect == nil {
		return errors.New("LogLua cannt find errCollect")
	}

	zonem := m.GetMachineByName(zone.ZoneHost)
	if zonem == nil {
		return errors.New(fmt.Sprintf("zone machine info err"))
	}

	logLua := LogConf{
		ID:             zone.Zid,
		IP:             zonem.IP,
		Port:           LogPort + zone.PortNumber,
		ConnectServers: make(map[string]interface{}),
	}
	logLua.ConnectServers["Collect"] = Connect{
		ID:   0,
		IP:   errCollect.IP,
		Port: ErrLogPort,
	}

	srvHead := ServerConfigHead{
		NET_TIMEOUT_MSEC:  30000,
		NET_MAX_CONNETION: NetMaxConnection,
		StartService:      []map[string]int{{"nType": LogServer}},
		LOG_INDEX:         "logserver",
		LOG_MAXLINE:       LogMaxLine,
		OpenGM:            1,
	}

	trans := struct2lua.ToLuaConfig(Dir, "Log", logLua, srvHead, 0)
	if trans == false {
		log.Println("log cannt wirte lua file")
	}

	zonelogdbm := m.GetMachineByName(zone.ZonelogdbHost)
	if zonelogdbm == nil {
		return errors.New(fmt.Sprintf("zone logdb conf info err"))
	}

	l := LogDBConf{
		DirName: "zonelog" + strconv.Itoa(zone.Zid),
		IP:      zonelogdbm.IP,
	}

	c, err := json.Marshal(l)
	if err != nil {
		return err
	}
	f, cerr := os.Create(Dir + "logdbconf")
	if cerr != nil {
		return cerr
	}
	f.Write(c)
	defer f.Close()

	return nil
}

func (m *machineMgr) LoginLua() error {
	loginWebM := m.GetMachineByName("loginWeb")
	if loginWebM == nil {
		return errors.New("Login Lua cannt find loginWeb")
	}

	accountDBM := m.GetMachineByName("accountDB")
	if accountDBM == nil {
		return errors.New("Login Lua cannt find accountdb")
	}

	masterM := m.GetMachineByName("master")
	if masterM == nil {
		return errors.New("Login Lua cannt find master")
	}

	loginM := m.GetMachineByName("login1")
	if loginM == nil {
		return errors.New("Login Lua cannt find login1")
	}

	loginLua := LoginConf{
		ID:             1,
		IP:             loginM.IP,
		Port:           LoginPort + 1,
		VesionStr:      "3.4.1",
		ConnectServers: make(map[string]interface{}),
	}

	loginLua.ConnectServers["LoginWeb"] = Connect{
		ID:   1,
		IP:   loginWebM.IP,
		Port: LoginWebPort,
	}
	loginLua.ConnectServers["Log"] = Connect{
		ID:   1,
		IP:   masterM.IP,
		Port: LogPort,
	}
	loginLua.ConnectServers["AccountDB"] = Connect{
		ID:   1,
		IP:   accountDBM.IP,
		Port: AccountDBPort,
	}
	loginLua.ConnectServers["Master"] = Connect{
		ID:   1,
		IP:   masterM.IP,
		Port: MasterPort,
	}
	srvHead := ServerConfigHead{
		NET_TIMEOUT_MSEC:  30000,
		NET_MAX_CONNETION: NetMaxConnection,
		StartService:      []map[string]int{{"nType": LoginServer}},
		LOG_INDEX:         "login1",
		LOG_MAXLINE:       LogMaxLine,
		OpenGM:            1,
	}

	trans := struct2lua.ToLuaConfig(m.conf.CommonConf, "Login", loginLua, srvHead, 1)
	if trans == false {
		log.Println("log cannt wirte lua file")
	}
	return nil
}

func (m *machineMgr) MasterLua() error {
	masterM := m.GetMachineByName("master")
	if masterM == nil {
		return errors.New("master Lua cannt find master")
	}
	masterlua := MasterConf{
		ID:             1,
		IP:             masterM.IP,
		Port:           MasterPort,
		AllZoneOpen:    true,
		ConnectServers: make(map[string]interface{}),
	}
	masterlua.ConnectServers["Log"] = Connect{
		ID:   1,
		IP:   masterM.IP,
		Port: LogPort,
	}

	srvHead := ServerConfigHead{
		NET_TIMEOUT_MSEC:  30000,
		NET_MAX_CONNETION: NetMaxConnection,
		StartService:      []map[string]int{{"nType": MasterServer}},
		LOG_INDEX:         "master",
		LOG_MAXLINE:       LogMaxLine,
		OpenGM:            1,
	}

	trans := struct2lua.ToLuaConfig(m.conf.CommonConf, "Master", masterlua, srvHead, 0)
	if trans == false {
		return errors.New("master cannt wirte lua file")
	}
	return nil
}

func (m *machineMgr) AccountDBLua() error {
	accountDBM := m.GetMachineByName("accountDB")
	if accountDBM == nil {
		return errors.New("accountdb Lua cannt find accountdb")
	}
	masterM := m.GetMachineByName("master")
	if masterM == nil {
		return errors.New("accountdb Lua cannt find master")
	}
	accountDBlua := AccountDBConf{
		ID:   0,
		IP:   accountDBM.IP,
		Port: AccountDBPort,
		Zid:  0,
		Mysql: MysqlLua{
			IP:             accountDBM.IP,
			Port:           MysqlPort,
			UserName:       UserName,
			Password:       PassWord,
			FlushFrequency: 300,
			DataBase:       "",
		},
		Redis: RedisLua{
			IP:       accountDBM.OutIP,
			Port:     RedisAccountPort,
			Password: RedisAccountPassWord,
		},
		ConnectServers: make(map[string]interface{}),
	}
	accountDBlua.ConnectServers["Log"] = Connect{
		ID:   1,
		IP:   masterM.IP,
		Port: LogPort,
	}

	srvHead := ServerConfigHead{
		NET_TIMEOUT_MSEC:  30000,
		NET_MAX_CONNETION: NetMaxConnection,
		StartService:      []map[string]int{{"nType": DbproxyServer}},
		LOG_INDEX:         "accountdb",
		LOG_MAXLINE:       LogMaxLine,
		OpenGM:            1,
	}

	trans := struct2lua.ToLuaConfig(m.conf.CommonConf, "AccountDB", accountDBlua, srvHead, 0)
	if trans == false {
		return errors.New("accountdblua cannt wirte lua file")
	}
	return nil
}

func (m *machineMgr) MasterLogLua() error {
	masterM := m.GetMachineByName("master")
	if masterM == nil {
		return errors.New("masterlog Lua cannt find master")
	}

	errLogM := m.GetMachineByName("errLog")
	if errLogM == nil {
		return errors.New("masterlog Lua cannt find errLog")
	}

	loglua := LogConf{
		ID:             1,
		IP:             masterM.IP,
		Port:           LogPort,
		ConnectServers: make(map[string]interface{}),
	}
	loglua.ConnectServers["Collect"] = Connect{
		ID:   1,
		IP:   errLogM.IP,
		Port: ErrLogPort,
	}

	srvHead := ServerConfigHead{
		NET_TIMEOUT_MSEC:  30000,
		NET_MAX_CONNETION: NetMaxConnection,
		StartService:      []map[string]int{{"nType": LogServer}},
		LOG_INDEX:         "masterlog",
		LOG_MAXLINE:       LogMaxLine,
		OpenGM:            1,
	}

	trans := struct2lua.ToLuaConfig(m.conf.CommonConf, "Log", loglua, srvHead, 0)
	if trans == false {
		return errors.New("masterLog cannt wirte lua file")
	}
	return nil

}
