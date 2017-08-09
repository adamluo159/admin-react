package zone

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"encoding/json"

	"github.com/adamluo159/admin-react/server/comInterface"
	"github.com/adamluo159/gameAgent/utils"
	"github.com/adamluo159/struct2lua"
	"gopkg.in/mgo.v2/bson"
)

var srvHead = comInterface.ServerConfigHead{
	NET_TIMEOUT_MSEC:  30000,
	NET_MAX_CONNETION: comInterface.NetMaxConnection,
	//StartService:      []comInterface.SRV{srv},
	//LOG_INDEX:         "gate",
	LOG_MAXLINE: comInterface.LogMaxLine,
	OpenGM:      1,
}

func WriteZoneConfigLua(zid int, ret *ZoneRsp, hostName string) {
	zone := Zone{}
	query := bson.M{"zid": zid}
	cl.Find(query).One(&zone)

	if zone.ZoneHost != hostName {
		ret.Result = "zid cannt match hostname fail"
		return
	}

	//每台机器只允许开
	zonem := zMgr.machineMgr.GetMachineByName(zone.ZoneHost)
	if zonem == nil {
		ret.Result = "cannt Find zoneMachine--"
		return
	}
	hostdir := os.Getenv("HOME") + comInterface.ConfDir + zonem.Hostname
	os.Mkdir(hostdir, os.ModePerm)

	dir := hostdir + "/zone" + strconv.Itoa(zone.Zid)
	os.Mkdir(dir, os.ModePerm)
	curDir := dir + "/"
	gerr := GateLua(&zone, zonem, curDir)
	if gerr != nil {
		ret.Result = "gate " + gerr.Error()
		return
	}
	cerr := CenterLua(&zone, zonem, curDir)
	if cerr != nil {
		ret.Result = "center" + cerr.Error()
		return
	}
	lerr := LogLua(&zone, zonem, curDir)
	if lerr != nil {
		ret.Result = "log" + lerr.Error()
		return
	}
	logicerr := LogicLua(&zone, zonem, curDir)
	if logicerr != nil {
		ret.Result = "logic" + logicerr.Error()
		return
	}
	charErr := CharDBLua(&zone, zonem, curDir)
	if charErr != nil {
		ret.Result = "chardb" + charErr.Error()
		return
	}

	commitstr := os.Getenv("HOME") + comInterface.ConfDir + "gitCommit"
	_, exeErr := utils.ExeShell("sh", commitstr, "add or update zone"+strconv.Itoa(zone.Zid))
	if exeErr != nil {
		ret.Result = exeErr.Error()
		return
	}
}
func DelZoneConfig(zid int, hostname string) error {
	commitstr := os.Getenv("HOME") + comInterface.ConfDir + "gitDelete"
	dir := hostname + "/" + "zone" + strconv.Itoa(zid)
	_, exeErr := utils.ExeShell("sh", commitstr, dir)
	if exeErr != nil {
		return exeErr
	}
	return nil
}

func GateLua(zone *Zone, zonem *comInterface.Machine, Dir string) error {
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

	theTime, _ := time.Parse("2006-01-02 15:04:05", zone.OpenTime) //使用模板在对应时区转化为time.time类型
	gateLua := comInterface.Gate{
		ID:             zone.Zid,
		Zid:            zone.Zid,
		ServerIP:       zonem.IP,
		ServerPort:     comInterface.GatePort + zone.PortNumber,
		ClientIP:       zonem.OutIP,
		ClientPort:     comInterface.ClientPort + zone.PortNumber,
		ChannelIds:     s,
		Open:           zone.Whitelst,
		Name:           zone.ZoneName,
		OpenTime:       theTime.Unix(),
		ConnectServers: make(map[string]interface{}),
	}
	gateLua.ConnectServers["CharDB"] = comInterface.Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CharDBPort + zone.PortNumber,
	}
	gateLua.ConnectServers["Master"] = comInterface.Connect{
		ID:   1,
		IP:   masterm.IP,
		Port: comInterface.MasterPort + comInterface.MasterCount,
	}
	gateLua.ConnectServers["Log"] = comInterface.Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.LogPort + zone.PortNumber,
	}

	srv := make(map[string]int)
	srv["nType"] = comInterface.GateServer
	srvHead.StartService = []comInterface.SRV{srv}
	srvHead.LOG_INDEX = "gate"

	trans := struct2lua.ToLuaConfig(Dir, "Gate", gateLua, srvHead, 0)
	if trans == false {
		log.Println("gate cannt wirte lua file")
	}
	return nil
}

func CenterLua(zone *Zone, zonem *comInterface.Machine, Dir string) error {
	theTime, _ := time.Parse("2006-01-02 15:04:05", zone.OpenTime) //使用模板在对应时区转化为time.time类型
	centerLua := comInterface.Center{
		ID:   zone.Zid,
		Zid:  zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CenterPort + zone.PortNumber,
		OnlineNumberCheckTime: 60 * 5,
		SingleServerLoad:      4000,
		ConnectServers:        make(map[string]interface{}),
		OpenTime:              theTime.Unix(),
	}

	centerLua.ConnectServers["CharDB"] = comInterface.Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CharDBPort + zone.PortNumber,
	}
	centerLua.ConnectServers["Gate"] = comInterface.Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.GatePort + zone.PortNumber,
	}
	centerLua.ConnectServers["Log"] = comInterface.Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.LogPort + zone.PortNumber,
	}
	srv := make(map[string]int)
	srv["nType"] = comInterface.CenterServer
	srvHead.StartService = []comInterface.SRV{srv}
	srvHead.LOG_INDEX = "center" + strconv.Itoa(zone.Zid)

	trans := struct2lua.ToLuaConfig(Dir, "Center", centerLua, srvHead, 0)
	if trans == false {
		log.Println("center cannt wirte lua file")
	}
	return nil
}

func CharDBLua(zone *Zone, zonem *comInterface.Machine, Dir string) error {
	//这个后面再加，现在假设所有服都用6379作缓存
	//zoneDBquery := bson.M{"zoneDBHost": zone.ZoneDBHost}
	//zoneDBCount, zdberr := cl.Find(zoneDBquery).Count()
	//if zdberr != nil {
	//	return zdberr
	//}
	zonedbm := zMgr.machineMgr.GetMachineByName(zone.ZoneDBHost)
	if zonedbm == nil {
		return errors.New("comInterface.CharDBLua cannt find machine")
	}

	mysqldbName := "cgzone" + strconv.Itoa(zone.Zid)
	charDBLua := comInterface.CharDB{
		ID:   zone.Zid,
		Zid:  zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CharDBPort + zone.PortNumber,
		Mysql: comInterface.MysqlLua{
			IP:             zonedbm.IP,
			Port:           comInterface.MysqlPort,
			UserName:       comInterface.UserName,
			Password:       comInterface.PassWord,
			FlushFrequency: 300,
			DataBase:       mysqldbName,
		},
		Redis: comInterface.RedisLua{
			IP:       zonedbm.IP,
			Port:     comInterface.RedisPort,
			Password: comInterface.RedisPassWord,
		},
		ConnectServers: make(map[string]interface{}),
	}

	charDBLua.ConnectServers["Log"] = comInterface.Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.LogPort + zone.PortNumber,
	}

	srv := make(map[string]int)
	srv["nType"] = comInterface.DbproxyServer
	srvHead.StartService = []comInterface.SRV{srv}
	srvHead.LOG_INDEX = "charDB" + strconv.Itoa(zone.Zid)
	trans := struct2lua.ToLuaConfig(Dir, "CharDB", charDBLua, srvHead, 0)
	if trans == false {
		log.Println("chardb cannt wirte lua file")
	}
	return nil
}

func LogicLua(zone *Zone, zonem *comInterface.Machine, Dir string) error {
	theTime, _ := time.Parse("2006-01-02 15:04:05", zone.OpenTime) //使用模板在对应时区转化为time.time类型
	logicLua := comInterface.Logic{
		Zid:            zone.Zid,
		IP:             zonem.IP,
		ConnectServers: make(map[string]interface{}),
		LoadAllMapIds:  false,
		OpenTime:       theTime.Unix(),
	}
	logicLua.ConnectServers["CharDB"] = comInterface.Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CharDBPort + zone.PortNumber,
	}
	logicLua.ConnectServers["Gate"] = comInterface.Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.GatePort + zone.PortNumber,
	}
	logicLua.ConnectServers["Center"] = comInterface.Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.CenterPort + zone.PortNumber,
	}
	logicLua.ConnectServers["Log"] = comInterface.Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: comInterface.LogPort + zone.PortNumber,
	}
	srv := make(map[string]int)
	srv["nType"] = comInterface.LogicServer
	srvHead.StartService = []comInterface.SRV{srv}

	for k, v := range LogicMap {
		logicLua.ID = k
		logicLua.Port = comInterface.LogicPort + (k-1)*100 + zone.PortNumber
		logicLua.MapIds = v

		s := "logic" + strconv.Itoa(k)
		srvHead.LOG_INDEX = s
		srvHead.LOG_MAXLINE = comInterface.LogMaxLine

		trans := struct2lua.ToLuaConfig(Dir, "Logic", logicLua, srvHead, k)
		if trans == false {
			log.Println("logic cannt wirte lua file")
		}
	}

	return nil
}

func LogLua(zone *Zone, zonem *comInterface.Machine, Dir string) error {
	logm := zMgr.machineMgr.GetMachineByName(zone.ZonelogdbHost)
	if logm == nil {
		return errors.New("LogLua cannt find machine")
	}
	errCollect := zMgr.machineMgr.GetMachineByName("errLog")
	if errCollect == nil {
		return errors.New("LogLua cannt find errCollect")
	}

	logLua := comInterface.Log{
		ID:             zone.Zid,
		IP:             zonem.IP,
		Port:           comInterface.LogPort + zone.PortNumber,
		ConnectServers: make(map[string]interface{}),
	}
	logLua.ConnectServers["Collect"] = comInterface.Connect{
		ID:   0,
		IP:   errCollect.IP,
		Port: comInterface.ErrLogPort,
	}

	srv := make(map[string]int)
	srv["nType"] = comInterface.LogServer
	srvHead.StartService = []comInterface.SRV{srv}
	srvHead.LOG_INDEX = "logserver"

	trans := struct2lua.ToLuaConfig(Dir, "Log", logLua, srvHead, 0)
	if trans == false {
		log.Println("log cannt wirte lua file")
	}

	l := comInterface.LogDBConf{
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
