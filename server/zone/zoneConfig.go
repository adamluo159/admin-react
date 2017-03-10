package zone

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/adamluo159/admin-react/server/machine"
	"github.com/adamluo159/struct2lua"
	"gopkg.in/mgo.v2/bson"
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
	zonem, err := machine.GetMachineByName(zone.ZoneHost)
	if err != nil {
		ret.Result = "cannt Find zoneMachine--" + err.Error()
		return
	}
	hostdir := os.Getenv("HOME") + "/GameConfig/" + zonem.Hostname
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

	commitstr := os.Getenv("HOME") + "/GameConfig/gitCommit"
	exeErr := ExeShell(commitstr, "add or update zone"+strconv.Itoa(zone.Zid))
	if exeErr != nil {
		ret.Result = exeErr.Error()
		return
	}
}

func GateLua(zone *Zone, zonem *machine.Machine, zoneCount int, Dir string) error {
	masterm, merr := machine.GetMachineByName("master")
	if merr != nil {
		return merr
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
		ServerPort:     machine.GatePort + zoneCount,
		ClientIP:       zonem.OutIP,
		ClientPort:     machine.ClientPort + zoneCount,
		ChannelIds:     s,
		Open:           zone.Whitelst,
		Name:           zone.ZoneName,
		ConnectServers: make(map[string]interface{}),
	}
	gateLua.ConnectServers["CharDB"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.CharDBPort + zoneCount,
	}
	gateLua.ConnectServers["Master"] = Connect{
		ID:   1,
		IP:   masterm.IP,
		Port: machine.MasterPort + machine.MasterCount,
	}
	gateLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.LogPort + zoneCount,
	}
	srv := make(map[string]int)
	srv["nType"] = machine.GateServer
	head := machine.ServerConfigHead{
		NET_TIMEOUT_MSEC:  machine.NetTimeOut,
		NET_MAX_CONNETION: machine.NetMaxConnection,
		StartService:      []machine.SRV{srv},
	}

	trans := struct2lua.ToLuaConfig(Dir, "Gate", gateLua, head, 0)
	if trans == false {
		fmt.Println("gate cannt wirte lua file")
	}
	return nil
}

func CenterLua(zone *Zone, zonem *machine.Machine, zoneCount int, Dir string) error {
	centerLua := Center{
		ID:   zone.Zid,
		Zid:  zone.Zid,
		IP:   zonem.IP,
		Port: machine.CenterPort + zoneCount,
		OnlineNumberCheckTime: 60 * 5,
		SingleServerLoad:      4000,
		ConnectServers:        make(map[string]interface{}),
	}

	centerLua.ConnectServers["CharDB"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.CharDBPort + zoneCount,
	}
	centerLua.ConnectServers["Gate"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.GatePort + zoneCount,
	}
	centerLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.LogPort + zoneCount,
	}
	srv := make(map[string]int)
	srv["nType"] = machine.CenterServer
	head := machine.ServerConfigHead{
		NET_TIMEOUT_MSEC:  machine.NetTimeOut,
		NET_MAX_CONNETION: machine.NetMaxConnection,
		StartService:      []machine.SRV{srv},
	}

	trans := struct2lua.ToLuaConfig(Dir, "Center", centerLua, head, 0)
	if trans == false {
		fmt.Println("center cannt wirte lua file")
	}
	return nil
}

func CharDBLua(zone *Zone, zonem *machine.Machine, zoneCount int, Dir string) error {
	zoneDBquery := bson.M{"zoneDBHost": zone.ZoneDBHost}
	zoneDBCount, zdberr := cl.Find(zoneDBquery).Count()
	if zdberr != nil {
		return zdberr
	}
	zonedbm, dberr := machine.GetMachineByName(zone.ZoneDBHost)
	if dberr != nil {
		return dberr
	}

	mysqldbName := "cgzone" + strconv.Itoa(zone.Zid)
	charDBLua := CharDB{
		ID:   zone.Zid,
		Zid:  zone.Zid,
		IP:   zonem.IP,
		Port: machine.CharDBPort + zoneCount,
		Mysql: MysqlLua{
			IP:             zonedbm.IP,
			Port:           machine.MysqlPort,
			UserName:       machine.UserName,
			Password:       machine.PassWord,
			FlushFrequency: 300,
			DataBase:       mysqldbName,
		},
		Redis: RedisLua{
			IP:       zonedbm.IP,
			Port:     machine.RedisPort + zoneDBCount,
			Password: "",
		},
	}
	srv := make(map[string]int)
	srv["nType"] = machine.DbproxyServer
	head := machine.ServerConfigHead{
		NET_TIMEOUT_MSEC:  machine.NetTimeOut,
		NET_MAX_CONNETION: machine.NetMaxConnection,
		StartService:      []machine.SRV{srv},
	}

	trans := struct2lua.ToLuaConfig(Dir, "CharDB", charDBLua, head, 0)
	if trans == false {
		fmt.Println("chardb cannt wirte lua file")
	}
	return nil
}

func LogicLua(zone *Zone, zonem *machine.Machine, zoneCount int, Dir string) error {
	logicLua := Logic{
		//ID:  1,
		Zid: zone.Zid,
		IP:  zonem.IP,
		//Port:           machine.LogicPort + zoneCount*3 + 1,
		ConnectServers: make(map[string]interface{}),
	}
	logicLua.ConnectServers["CharDB"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.CharDBPort + zoneCount,
	}
	logicLua.ConnectServers["Gate"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.GatePort + zoneCount,
	}
	logicLua.ConnectServers["Center"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.CenterPort + zoneCount,
	}
	logicLua.ConnectServers["Log"] = Connect{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.LogPort + zoneCount,
	}
	srv := make(map[string]int)
	srv["nType"] = machine.LogicServer
	head := machine.ServerConfigHead{
		NET_TIMEOUT_MSEC:  machine.NetTimeOut,
		NET_MAX_CONNETION: machine.NetMaxConnection,
		StartService:      []machine.SRV{srv},
	}

	for k, v := range LogicMap {
		logicLua.ID = k
		logicLua.Port = machine.LogicPort + zoneCount*3 + k
		logicLua.MapIds = v

		trans := struct2lua.ToLuaConfig(Dir, "Logic", logicLua, head, k)
		if trans == false {
			fmt.Println("logic cannt wirte lua file")
		}
	}

	return nil
}

func LogLua(zone *Zone, zonem *machine.Machine, zoneCount int, Dir string) error {
	logm, logerr := machine.GetMachineByName(zone.ZonelogdbHost)
	if logerr != nil {
		return logerr
	}
	logLua := Log{
		ID:   zone.Zid,
		IP:   zonem.IP,
		Port: machine.LogPort + zoneCount,
		ZoneLogMysql: MysqlLua{
			IP:             logm.IP,
			Port:           machine.MysqlPort,
			UserName:       machine.UserName,
			Password:       machine.PassWord,
			FlushFrequency: 300,
			DataBase:       "zonelog" + strconv.Itoa(zone.Zid),
		},
		GlobalLogMysql: GlobalDB,
	}
	srv := make(map[string]int)
	srv["nType"] = machine.LogServer
	head := machine.ServerConfigHead{
		NET_TIMEOUT_MSEC:  machine.NetTimeOut,
		NET_MAX_CONNETION: machine.NetMaxConnection,
		StartService:      []machine.SRV{srv},
	}
	trans := struct2lua.ToLuaConfig(Dir, "Log", logLua, head, 0)
	if trans == false {
		fmt.Println("log cannt wirte lua file")
	}

	return nil
}

func ExeShell(dir string, args string) error {

	fmt.Println("begin execute shell.....", dir, "--", args)
	// 执行系统命令
	// 第一个参数是命令名称
	// 后面参数可以有多个，命令参数
	cmd := exec.Command("sh", dir, args) //"GameConfig/gitCommit", "zoneo")
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		return err
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println(string(opBytes))
	return nil
}
