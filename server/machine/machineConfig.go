package machine

import (
	"errors"
	"log"
	"strconv"

	"github.com/adamluo159/admin-react/server/comInterface"
	"github.com/adamluo159/struct2lua"
)

var srvHead = comInterface.ServerConfigHead{
	NET_TIMEOUT_MSEC:  comInterface.NetTimeOut,
	NET_MAX_CONNETION: comInterface.NetMaxConnection,
	LOG_MAXLINE:       comInterface.LogMaxLine,
	OpenGM:            1,
}

func LoginLua(loginId int, dir string, loginWebIP string, masterIP string, accountdbIP string) error {
	mName := "login" + strconv.Itoa(loginId)
	mInfo := mhMgr.GetMachineByName(mName)
	if mInfo == nil {
		return errors.New("cannt find machine info," + mName)
	}
	loginLua := comInterface.Login{
		ID:             loginId,
		IP:             mInfo.IP,
		Port:           comInterface.LoginPort + loginId,
		VesionStr:      "3.4.1",
		ConnectServers: make(map[string]interface{}),
	}
	loginLua.ConnectServers["LoginWeb"] = comInterface.Connect{
		ID:   0,
		IP:   loginWebIP,
		Port: comInterface.LoginWebPort,
	}
	loginLua.ConnectServers["Log"] = comInterface.Connect{
		ID:   0,
		IP:   masterIP,
		Port: comInterface.LogPort,
	}
	loginLua.ConnectServers["AccountDB"] = comInterface.Connect{
		ID:   0,
		IP:   accountdbIP,
		Port: comInterface.AccountDBPort,
	}
	loginLua.ConnectServers["Master"] = comInterface.Connect{
		ID:   0,
		IP:   masterIP,
		Port: comInterface.MasterPort,
	}
	srv := make(map[string]int)
	srv["nType"] = comInterface.LoginServer
	srvHead.StartService = []comInterface.SRV{srv}
	srvHead.LOG_INDEX = mName

	trans := struct2lua.ToLuaConfig(dir, "Login", loginLua, srvHead, loginId)
	if trans == false {
		log.Println("log cannt wirte lua file")
	}
	return nil
}

func MasterLua(dir string, masterIP string) error {
	masterlua := comInterface.Master{
		ID:             0,
		IP:             masterIP,
		Port:           comInterface.MasterPort,
		AllZoneOpen:    true,
		ConnectServers: make(map[string]interface{}),
	}

	masterlua.ConnectServers["Log"] = comInterface.Connect{
		ID:   0,
		IP:   masterIP,
		Port: comInterface.LogPort,
	}

	srv := make(map[string]int)
	srv["nType"] = comInterface.MasterServer
	srvHead.StartService = []comInterface.SRV{srv}
	srvHead.LOG_INDEX = "master"

	trans := struct2lua.ToLuaConfig(dir, "Master", masterlua, srvHead, 0)
	if trans == false {
		return errors.New("master cannt wirte lua file")
	}
	return nil
}

func AccountDBLua(dir string, accountIP string, dbIP string, logIP string) error {
	accountDBlua := comInterface.AccountDB{
		ID:   0,
		IP:   accountIP,
		Port: comInterface.AccountDBPort,
		Zid:  0,
		Mysql: comInterface.MysqlLua{
			IP:             dbIP,
			Port:           comInterface.MysqlPort,
			UserName:       comInterface.UserName,
			Password:       comInterface.PassWord,
			FlushFrequency: 300,
			DataBase:       "",
		},
		Redis: comInterface.RedisLua{
			IP:       dbIP,
			Port:     comInterface.RedisAccountPort,
			Password: comInterface.RedisAccountPassWord,
		},
		ConnectServers: make(map[string]interface{}),
	}

	accountDBlua.ConnectServers["Log"] = comInterface.Connect{
		ID:   0,
		IP:   logIP,
		Port: comInterface.LogPort,
	}

	srv := make(map[string]int)
	srv["nType"] = comInterface.DbproxyServer
	srvHead.StartService = []comInterface.SRV{srv}
	srvHead.LOG_INDEX = "accountdb"

	trans := struct2lua.ToLuaConfig(dir, "AccountDB", accountDBlua, srvHead, 0)
	if trans == false {
		return errors.New("accountdblua cannt wirte lua file")
	}
	return nil
}

func MasterLogLua(dir string, masterIP string, errLogIP string) error {
	loglua := comInterface.Log{
		ID:             0,
		IP:             masterIP,
		Port:           comInterface.LogPort,
		ConnectServers: make(map[string]interface{}),
	}
	loglua.ConnectServers["Collect"] = comInterface.Connect{
		ID:   0,
		IP:   errLogIP,
		Port: comInterface.ErrLogPort,
	}

	srv := make(map[string]int)
	srv["nType"] = comInterface.DbproxyServer
	srvHead.StartService = []comInterface.SRV{srv}
	srvHead.LOG_INDEX = "masterlog"

	trans := struct2lua.ToLuaConfig(dir, "MasterLog", loglua, srvHead, 0)
	if trans == false {
		return errors.New("masterLog cannt wirte lua file")
	}
	return nil
}
