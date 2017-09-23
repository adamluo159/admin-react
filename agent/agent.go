package agent

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/adamluo159/admin-react/utils"
)

type (

	//服务信息
	ServiceInfo struct {
		Sname          string //服务名
		Started        bool   //游戏区服是否已启动
		RegularlyCheck bool   //是否开启定时检查进程功能
		ClientPorts    *[]int
		Zid            int
	}

	Agent interface {
		Run() error
	}

	agent struct {
		conn     *net.Conn
		msgMap   map[uint32]func([]byte)
		srvs     map[string]*ServiceInfo //后台配置该机器所有的服务信息
		SvnVer   string                  //svn版本号
		hostName string                  //主机名
	}

	Conf struct {
		ConAddress string

		RemoteIP      string //远端游戏配置仓库地址
		RemoteConfDir string //远端游戏配置仓库目录
		LocalConfDir  string //本地游戏配置地址
		GameShell     string //游戏启动关闭脚本
		HttpStr       string //在线人数上报地址
		HttpKey       string //上报数据时做签名的key
	}

	OnlineReport struct {
		Method    string       `json:"method"`
		Param     []ZoneOnline `json:"param"`
		TimeStamp int64        `json:"timestamp"`
	}
	ZoneOnline struct {
		Zid       int `json:"zid"`
		OnlineNum int `json:"onlineNum"`
	}
)

var (
	conf *Conf //agent配置信息
)

//加载配置信息
func LoadConfig(cfPath string) {

	data, err := ioutil.ReadFile(cfPath)
	if err != nil {
		log.Fatal(err)
	}

	datajson := []byte(data)
	err = json.Unmarshal(datajson, &conf)
	if err != nil {
		log.Fatal(err)
	}

	os.Mkdir(conf.LocalConfDir, os.ModePerm)
}

//起服
func StartZone(zone string) bool {
	run := CheckProcess(zone)
	if run == false {
		utils.ExeShellArgs3("sh", conf.GameShell, "start", zone, "")
		for index := 0; index < 6; index++ {
			if run {
				break
			}
			run = CheckProcess(zone)
			time.Sleep(time.Second * 5)
		}
	}
	return run
}

//停服
func StopZone(zone string) bool {
	utils.ExeShellArgs2("sh", conf.GameShell, "stop", zone)
	run := CheckProcess(zone)
	return !run
}

//获取svn版本
func SvnInfo() string {
	ver, err := utils.ExeShell("sh", "scripts/svnInfo", "")
	if err != nil {
		log.Fatal("update svn ", err)
	}
	return ver
}

//svn更新
func SvnUp() error {
	_, upErr := utils.ExeShell("sh", "scripts/svnUp", "")
	return upErr
}

//更新本地配置
func UpdateGameConf() error {
	_, err := utils.ExeShellArgs3("expect", "scripts/synGameConf_expt", conf.RemoteIP, conf.RemoteConfDir, conf.LocalConfDir)
	return err
}

//检查进程是否存在
func CheckProcess(dstName string) bool {
	ret, _ := utils.ExeShell("sh", "scripts/checkZoneProcess", dstName)
	s := strings.Replace(string(ret), " ", "", -1)
	if s != "" {
		return false
	}
	return true
}

//端口的连接数
func OnlinePlayers(port string) int {
	ret, _ := utils.ExeShell("sh", "scripts/onlinePlayers", port)
	if n, err := strconv.Atoi(ret); err == nil {
		return n
	}
	return 0
}
