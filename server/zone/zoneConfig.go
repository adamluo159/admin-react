package zone

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"encoding/json"

	"github.com/adamluo159/admin-react/server/comInterface"
	"github.com/adamluo159/gameAgent/utils"
	"github.com/adamluo159/struct2lua"
	"gopkg.in/mgo.v2/bson"
)

const (
	confDir string = "/gConf/"
)

func WriteZoneConfigLua(zid int, ret *ZoneRsp, hostName string) {
	zone := Zone{}
	query := bson.M{"zid": zid}
	cl.Find(query).One(&zone)

	if zone.ZoneHost != hostName {
		ret.Result = "zid cannt match hostname fail"
		return
	}

	zonequery := bson.M{"zoneHost": zone.ZoneHost}
	zoneCount, zerr := cl.Find(zonequery).Count()
	if zerr != nil {
		ret.Result = "cannt Find zoneCount--" + zerr.Error()
		return
	}
	zonem := zMgr.machineMgr.GetMachineByName(zone.ZoneHost)
	if zonem == nil {
		ret.Result = "cannt Find zoneMachine--"
		return
	}
	hostdir := os.Getenv("HOME") + confDir + zonem.Hostname
	os.Mkdir(hostdir, os.ModePerm)

	dir := hostdir + "/zone" + strconv.Itoa(zone.Zid)
	os.Mkdir(dir, os.ModePerm)
	curDir := dir + "/"
	gerr := GateLua(&zone, zonem, zoneCount, curDir)
	if gerr != nil {
		ret.Result = "gate " + gerr.Error()
		return
	}
	cerr := CenterLua(&zone, zonem, zoneCount, curDir)
	if cerr != nil {
		ret.Result = "center" + cerr.Error()
		return
	}
	lerr := LogLua(&zone, zonem, zoneCount, curDir)
	if lerr != nil {
		ret.Result = "log" + lerr.Error()
		return
	}
	logicerr := LogicLua(&zone, zonem, zoneCount, curDir)
	if logicerr != nil {
		ret.Result = "logic" + logicerr.Error()
		return
	}
	charErr := CharDBLua(&zone, zonem, zoneCount, curDir)
	if charErr != nil {
		ret.Result = "chardb" + charErr.Error()
		return
	}

	commitstr := os.Getenv("HOME") + confDir + "gitCommit"
	_, exeErr := utils.ExeShell("sh", commitstr, "add or update zone"+strconv.Itoa(zone.Zid))
	if exeErr != nil {
		ret.Result = exeErr.Error()
		return
	}
}
func DelZoneConfig(zid int, hostname string) error {
	commitstr := os.Getenv("HOME") + confDir + "gitDelete"
	dir := hostname + "/" + "zone" + strconv.Itoa(zid)
	_, exeErr := utils.ExeShell("sh", commitstr, dir)
	if exeErr != nil {
		return exeErr
	}
	return nil
}

func GateLua(zone *Zone, zonem *comInterface.Machine, zoneCount int, Dir string) error {
	masterm := zMgr.machineMgr.GetMachineByName("master")
	if masterm == nil {
		return errors.New(" GateLua cannt find machine")
	}

	s := make([]int, len(zone.Channels))
	n := 0
	for _, v := range zone.Channels {
		s[n] = Str2IntChannels[v]
		n++
	}

	gateLua := Gate{
		ID:             zone.Zid,
		Zid:            zone.Zid,
		ServerIP:       zonem.IP,
		ServerPort:     comInterface.GatePort + zoneCount,
		ClientIP:       zonem.OutIP,
		ClientPort:     comInterface.ClientPort + zoneCount,
		ChannelIds:     s,
		Open:           zone.Whitelst,
		Name:           zone.ZoneName,
		ConnectServers: make(map[string]interface{}),
	}
	gateLua.ConnectServers["CharDB"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CharDBPort + zoneCount,
	}
	gateLua.ConnectServers["Master"] = Connect{
		ID:   1,
		IP:   masterm.IP,
		Port: comInterface.MasterPort + comInterface.MasterCount,
	}
	gateLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.LogPort + zoneCount,
	}

	srv := make(map[string]int)
	srv["nType"] = comInterface.GateServer
	head := comInterface.ServerConfigHead{
		NET_TIMEOUT_MSEC:  3600000,
		NET_MAX_CONNETION: comInterface.NetMaxConnection,
		StartService:      []comInterface.SRV{srv},
		LOG_INDEX:         "gate",
		LOG_MAXLINE:       comInterface.LogMaxLine,
		OpenGM:            1,
	}

	trans := struct2lua.ToLuaConfig(Dir, "Gate", gateLua, head, 0)
	if trans == false {
		log.Println("gate cannt wirte lua file")
	}
	return nil
}

func CenterLua(zone *Zone, zonem *comInterface.Machine, zoneCount int, Dir string) error {
	centerLua := Center{
		ID:   zone.Zid,
		Zid:  zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CenterPort + zoneCount,
		OnlineNumberCheckTime: 60 * 5,
		SingleServerLoad:      4000,
		ConnectServers:        make(map[string]interface{}),
	}

	centerLua.ConnectServers["CharDB"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CharDBPort + zoneCount,
	}
	centerLua.ConnectServers["Gate"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.GatePort + zoneCount,
	}
	centerLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.LogPort + zoneCount,
	}
	srv := make(map[string]int)
	srv["nType"] = comInterface.CenterServer
	head := comInterface.ServerConfigHead{
		NET_TIMEOUT_MSEC:  comInterface.NetTimeOut,
		NET_MAX_CONNETION: comInterface.NetMaxConnection,
		StartService:      []comInterface.SRV{srv},
		LOG_INDEX:         "cener",
		LOG_MAXLINE:       comInterface.LogMaxLine,
		OpenGM:            1,
	}

	trans := struct2lua.ToLuaConfig(Dir, "Center", centerLua, head, 0)
	if trans == false {
		log.Println("center cannt wirte lua file")
	}
	return nil
}

func CharDBLua(zone *Zone, zonem *comInterface.Machine, zoneCount int, Dir string) error {
	//这个后面再加，现在假设所有服都用6379作缓存
	//zoneDBquery := bson.M{"zoneDBHost": zone.ZoneDBHost}
	//zoneDBCount, zdberr := cl.Find(zoneDBquery).Count()
	//if zdberr != nil {
	//	return zdberr
	//}
	zonedbm := zMgr.machineMgr.GetMachineByName(zone.ZoneDBHost)
	if zonedbm == nil {
		return errors.New("CharDBLua cannt find machine")
	}

	mysqldbName := "cgzone" + strconv.Itoa(zone.Zid)
	charDBLua := CharDB{
		ID:   zone.Zid,
		Zid:  zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CharDBPort + zoneCount,
		Mysql: MysqlLua{
			IP:             zonedbm.IP,
			Port:           comInterface.MysqlPort,
			UserName:       comInterface.UserName,
			Password:       comInterface.PassWord,
			FlushFrequency: 300,
			DataBase:       mysqldbName,
		},
		Redis: RedisLua{
			IP: zonedbm.IP,
			//Port:     comInterface.RedisPort + zoneDBCount,
			Port:     comInterface.RedisPort,
			Password: "",
		},
		ConnectServers: make(map[string]interface{}),
	}

	charDBLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.LogPort + zoneCount,
	}

	srv := make(map[string]int)
	srv["nType"] = comInterface.DbproxyServer
	head := comInterface.ServerConfigHead{
		NET_TIMEOUT_MSEC:  comInterface.NetTimeOut,
		NET_MAX_CONNETION: comInterface.NetMaxConnection,
		StartService:      []comInterface.SRV{srv},
		LOG_INDEX:         "chardb",
		LOG_MAXLINE:       comInterface.LogMaxLine,
		OpenGM:            1,
	}

	trans := struct2lua.ToLuaConfig(Dir, "CharDB", charDBLua, head, 0)
	if trans == false {
		log.Println("chardb cannt wirte lua file")
	}
	return nil
}

func LogicLua(zone *Zone, zonem *comInterface.Machine, zoneCount int, Dir string) error {
	logicLua := Logic{
		//ID:  1,
		Zid: zone.Zid,
		IP:  zonem.IP,
		//Port:           comInterface.LogicPort + zoneCount*3 + 1,
		ConnectServers: make(map[string]interface{}),
	}
	logicLua.ConnectServers["CharDB"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CharDBPort + zoneCount,
	}
	logicLua.ConnectServers["Gate"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.GatePort + zoneCount,
	}
	logicLua.ConnectServers["Center"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CenterPort + zoneCount,
	}
	logicLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.LogPort + zoneCount,
	}
	srv := make(map[string]int)
	srv["nType"] = comInterface.LogicServer
	head := comInterface.ServerConfigHead{
		NET_TIMEOUT_MSEC:  comInterface.NetTimeOut,
		NET_MAX_CONNETION: comInterface.NetMaxConnection,
		StartService:      []comInterface.SRV{srv},
		OpenGM:            1,
	}

	for k, v := range LogicMap {
		logicLua.ID = k
		logicLua.Port = comInterface.LogicPort + zoneCount*3 + k
		logicLua.MapIds = v

		s := "logic" + strconv.Itoa(k)
		head.LOG_INDEX = s
		head.LOG_MAXLINE = comInterface.LogMaxLine

		trans := struct2lua.ToLuaConfig(Dir, "Logic", logicLua, head, k)
		if trans == false {
			log.Println("logic cannt wirte lua file")
		}
	}

	return nil
}

func LogLua(zone *Zone, zonem *comInterface.Machine, zoneCount int, Dir string) error {
	logm := zMgr.machineMgr.GetMachineByName(zone.ZonelogdbHost)
	if logm == nil {
		return errors.New("LogLua cannt find machine")
	}

	masterm := zMgr.machineMgr.GetMachineByName("master")
	if masterm == nil {
		return errors.New("LogLua cannt find machine")
	}
	logLua := Log{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.LogPort + zoneCount,
		ZoneLogMysql: MysqlLua{
			IP:             logm.IP,
			Port:           comInterface.MysqlPort,
			UserName:       comInterface.UserName,
			Password:       comInterface.PassWord,
			FlushFrequency: 300,
			DataBase:       "zonelog" + strconv.Itoa(zone.Zid),
		},
		GlobalLogMysql: GlobalDB,
		ConnectServers: make(map[string]interface{}),
	}
	logLua.ConnectServers["Collect"] = Connect{
		ID:   1,
		IP:   masterm.IP,
		Port: 1237,
	}

	srv := make(map[string]int)
	srv["nType"] = comInterface.LogServer
	head := comInterface.ServerConfigHead{
		NET_TIMEOUT_MSEC:  comInterface.NetTimeOut,
		NET_MAX_CONNETION: comInterface.NetMaxConnection,
		StartService:      []comInterface.SRV{srv},
		LOG_INDEX:         "yylog",
		LOG_MAXLINE:       comInterface.LogMaxLine,
		OpenGM:            1,
	}
	trans := struct2lua.ToLuaConfig(Dir, "Log", logLua, head, 0)
	if trans == false {
		log.Println("log cannt wirte lua file")
	}

	l := LogDBConf{
		DirName: "zonelog" + strconv.Itoa(zone.Zid),
		IP:      logm.IP,
	}

	c, err := json.Marshal(l)
	if err != nil {
		log.Println("Loglua cannt code logdbconf json, ", zone.Zid, err.Error())
		return nil
	}
	f, err := os.Create(Dir + "logdbconf")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	f.Write(c)
	defer f.Close()

	return nil
}
