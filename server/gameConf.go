package yada

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/adamluo159/struct2lua"
)

type (
	ZoneConf struct {
		ID             int
		Zid            int
		ServerIP       string
		ServerPort     int
		ClientIP       string
		ClientPort     int
		ChannelIds     []int
		WhiteLst       bool
		Name           string
		OpenTime       int64
		DataLogDB      string
		ZoneLogDB      string
		ZoneDBBak      string
		TokenTimeCheck bool
		ConnectServers map[string]interface{}
	}
	Connect struct {
		ID   int
		Port int
		IP   string
	}

	MysqlLua struct {
		IP             string
		Port           int
		UserName       string
		Password       string
		FlushFrequency int
		DataBase       string
	}

	RedisLua struct {
		IP       string
		Port     int
		Password string
	}

	GateConf struct {
		ID             int
		Zid            int
		ServerIP       string
		ServerPort     int
		ClientIP       string
		ClientPort     int
		ConnectServers map[string]interface{}
	}

	CenterConf struct {
		ID               int
		Zid              int
		IP               string
		Port             int
		SingleServerLoad int
		ConnectServers   map[string]interface{}
		OpenTime         int64
		CheckPing        int
		TokenTimeCheck   bool
	}

	CharDBConf struct {
		ID             int
		Zid            int
		IP             string
		Port           int
		Mysql          MysqlLua
		Redis          RedisLua
		ConnectServers map[string]interface{}
	}

	LogicConf struct {
		ID             int
		Zid            int
		IP             string
		Port           int
		ConnectServers map[string]interface{}
		LoadAllMapIds  bool
		OpenTime       int64
	}

	LogConf struct {
		ID             int
		IP             string
		Port           int
		Zid            int
		ConnectServers map[string]interface{}
	}

	LoginConf struct {
		ID             int
		IP             string
		Port           int
		VesionStr      string
		TokenTimeCheck bool
		ConnectServers map[string]interface{}
	}

	MasterConf struct {
		ID             int
		IP             string
		Port           int
		AllZoneOpen    bool
		ConnectServers map[string]interface{}
	}
	AccountDBConf struct {
		ID             int
		Zid            int
		IP             string
		Port           int
		Mysql          MysqlLua
		Redis          RedisLua
		ConnectServers map[string]interface{}
	}

	ServerConfigHead struct {
		NET_TIMEOUT_MSEC    int
		NET_MAX_CONNETION   int
		StartService        []map[string]int
		LOG_INDEX           string
		LOG_MAXLINE         int
		OpenGM              int
		LOG_PRIORITY        int
		CLIENT_TIMEOUT_MSEC int
	}
)

const (
	CharDBPort int = 5000
	CenterPort int = 5100
	LogPort    int = 5200
	ZonePort   int = 5300
	LogicPort  int = 5400 //logic1 5500起 logic2 5500起

	ZoneClientPort int = 7000
	GatePort       int = 7100 //gate1 7100起 gate2 7110起
	ClientPort     int = 7200

	AccountDBPort    int = 6500
	RedisPort        int = 6379
	MysqlPort        int = 3306
	OpWebPort        int = 1252
	LoginWebPort     int = 1236
	ErrLogPort       int = 1237
	RedisAccountPort int = 6380

	LoginPort        int = 9550
	NetMaxConnection int = 5000

	DbproxyServer int = 1
	LoginServer   int = 2
	CenterServer  int = 3
	LogicServer   int = 4
	LogServer     int = 5
	MasterServer  int = 6
	GateServer    int = 7
	ZoneServer    int = 8
)

const longForm = "2006-01-02 15:04:05"

func (m *machineMgr) ZoneLua(zone *Zone, Dir string) error {
	zonem := m.GetMachineByName(zone.ZoneHost)
	if zonem == nil {
		return errors.New(fmt.Sprintf("zone machine info err, %s", zone.ZoneHost))
	}
	masterm := m.GetMachineByName("master")
	if zonem == nil {
		return errors.New(fmt.Sprintf("master machine info err"))
	}
	zonelogdb := m.GetMachineByName(zone.ZonelogdbHost)
	if zonelogdb == nil {
		return errors.New(fmt.Sprintf("zonelogdb machine info err %s", zone.ZonelogdbHost))
	}

	datalogdb := m.GetMachineByName(zone.DatalogdbHost)
	if datalogdb == nil {
		return errors.New(fmt.Sprintf("datalogdb machine info err, %s", zone.DatalogdbHost))
	}

	zonedbBak := m.GetMachineByName(zone.ZonedbBakHost)
	if zonedbBak == nil {
		return errors.New(fmt.Sprintf("zonedbBak machine info err, %s", zone.ZonedbBakHost))
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
		WhiteLst:       zone.Whitelst,
		Name:           zone.ZoneName,
		OpenTime:       theTime.Unix(),
		DataLogDB:      datalogdb.IP,
		ZoneLogDB:      zonelogdb.IP,
		ZoneDBBak:      zonedbBak.IP,
		TokenTimeCheck: !m.conf.Debug,
		ConnectServers: make(map[string]interface{}),
	}
	zoneLua.ConnectServers["Master"] = Connect{
		ID:   1,
		IP:   masterm.IP,
		Port: 9501,
	}
	zoneLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: LogPort + zone.PortNumber,
	}
	m.sHead.StartService[0]["nType"] = ZoneServer
	m.sHead.LOG_INDEX = "zone"
	trans := struct2lua.ToLuaConfig(Dir, "Zone", zoneLua, m.sHead, 0)
	if trans == false {
		log.Println("gate cannt wirte lua file")
	}
	return nil
}

func (m *machineMgr) GateLua(zone *Zone, Dir string, arrayClientPorts *[]int) error {
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
	m.sHead.StartService[0]["nType"] = GateServer

	for i := 1; i <= m.conf.GateCount; i++ {
		gateLua.ID = i
		gateLua.ClientPort = ClientPort + zone.PortNumber*10 + i - 1
		gateLua.ServerPort = GatePort + zone.PortNumber*10 + i - 1
		m.sHead.LOG_INDEX = "gate" + strconv.Itoa(i)
		trans := struct2lua.ToLuaConfig(Dir, "Gate", gateLua, m.sHead, i)
		if trans == false {
			log.Printf("gate cannt wirte lua file, gateid:%d\n", i)
		}
		(*arrayClientPorts)[i-1] = gateLua.ClientPort
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
		CheckPing:        1,
		TokenTimeCheck:   !m.conf.Debug,
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
		IP:   m.conf.OpWebIP,
		Port: OpWebPort,
	}
	centerLua.ConnectServers["WebPay"] = Connect{
		ID:   1,
		IP:   m.conf.PayWebIP,
		Port: m.conf.PayWebPort,
	}
	m.sHead.StartService[0]["nType"] = CenterServer
	m.sHead.LOG_INDEX = "center"

	trans := struct2lua.ToLuaConfig(Dir, "Center", centerLua, m.sHead, 0)
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
			UserName:       m.conf.MysqlUsr,
			Password:       m.conf.MysqlPwd,
			FlushFrequency: 300,
			DataBase:       mysqldbName,
		},
		Redis: RedisLua{
			IP:       zonedb.IP,
			Port:     RedisPort,
			Password: m.conf.RedisCharPwd,
		},
		ConnectServers: make(map[string]interface{}),
	}

	charDBLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: LogPort + zone.PortNumber,
	}
	m.sHead.StartService[0]["nType"] = DbproxyServer
	m.sHead.LOG_INDEX = "charDB"

	trans := struct2lua.ToLuaConfig(Dir, "CharDB", charDBLua, m.sHead, 0)
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
	m.sHead.StartService[0]["nType"] = LogicServer
	for i := 1; i <= m.conf.LogicCount; i++ {
		logicLua.ID = i
		logicLua.Port = LogicPort + 10*zone.PortNumber + i - 1
		m.sHead.LOG_INDEX = "logic" + strconv.Itoa(i)
		trans := struct2lua.ToLuaConfig(Dir, "Logic", logicLua, m.sHead, i)
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
		ID:             1,
		IP:             zonem.IP,
		Port:           LogPort + zone.PortNumber,
		Zid:            zone.Zid,
		ConnectServers: make(map[string]interface{}),
	}
	logLua.ConnectServers["Collect"] = Connect{
		ID:   0,
		IP:   errCollect.IP,
		Port: ErrLogPort,
	}
	logLua.ConnectServers["DataLog"] = Connect{
		ID:   0,
		IP:   m.conf.DataLogIP,
		Port: m.conf.DataLogPort,
	}

	m.sHead.StartService[0]["nType"] = LogServer
	m.sHead.LOG_INDEX = "logServer"

	trans := struct2lua.ToLuaConfig(Dir, "Log", logLua, m.sHead, 0)
	if trans == false {
		return errors.New("log cannt wirte lua file")
	}

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
		TokenTimeCheck: !m.conf.Debug,
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
		Port: m.conf.MasterPort,
	}
	m.sHead.StartService[0]["nType"] = LoginServer
	m.sHead.LOG_INDEX = "login1"
	trans := struct2lua.ToLuaConfig(m.conf.CommonConf, "Login", loginLua, m.sHead, 1)
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
		Port:           m.conf.MasterPort,
		AllZoneOpen:    true,
		ConnectServers: make(map[string]interface{}),
	}
	masterlua.ConnectServers["Log"] = Connect{
		ID:   1,
		IP:   masterM.IP,
		Port: LogPort,
	}
	m.sHead.StartService[0]["nType"] = MasterServer
	m.sHead.LOG_INDEX = "master"

	trans := struct2lua.ToLuaConfig(m.conf.CommonConf, "Master", masterlua, m.sHead, 0)
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
			UserName:       m.conf.MysqlUsr,
			Password:       m.conf.MysqlPwd,
			FlushFrequency: 300,
			DataBase:       "gameAccount",
		},
		Redis: RedisLua{
			IP:       accountDBM.IP,
			Port:     RedisAccountPort,
			Password: m.conf.RedisAccountPWd,
		},
		ConnectServers: make(map[string]interface{}),
	}
	accountDBlua.ConnectServers["Log"] = Connect{
		ID:   1,
		IP:   masterM.IP,
		Port: LogPort,
	}
	m.sHead.StartService[0]["nType"] = DbproxyServer
	m.sHead.LOG_INDEX = "accountdb"

	trans := struct2lua.ToLuaConfig(m.conf.CommonConf, "AccountDB", accountDBlua, m.sHead, 0)
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
	loglua.ConnectServers["DataLog"] = Connect{
		ID:   1,
		IP:   m.conf.DataLogIP,
		Port: m.conf.DataLogPort,
	}
	m.sHead.StartService[0]["nType"] = LogServer
	m.sHead.LOG_INDEX = "masterlog"

	trans := struct2lua.ToLuaConfig(m.conf.CommonConf, "Log", loglua, m.sHead, 0)
	if trans == false {
		return errors.New("masterLog cannt wirte lua file")
	}
	return nil

}
