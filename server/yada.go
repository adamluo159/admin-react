// Package yadaprovides ...

package yada

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/adamluo159/admin-react/utils"
	"github.com/labstack/echo"
	permissions "github.com/xyproto/permissions2"
	mgo "gopkg.in/mgo.v2"
)

type (
	//回复信息
	ZoneRsp struct {
		Result  string
		Item    Zone
		Items   []Zone
		Zstates []ZoneStates
	}

	SaveZoneReq struct {
		OldZoneName string
		OldZid      int
		Item        Zone
	}

	ZoneReq struct {
		Zid  int
		Host string
	}
	InitMachine struct {
		Items []Machine
	}

	SaveMachineReq struct {
		Oldhost string
		Item    Machine
	}

	//回复信息
	MachineRsp struct {
		Result string
		Item   Machine
	}

	MachineAllRsp struct {
		Result string
		Items  []Machine
	}
)

type (
	Yada interface {
		Run()
	}
	yada struct {
		as         Aserver
		e          *echo.Echo
		session    *mgo.Session
		conf       Conf
		zMgr       ZoneMgr
		machineMgr MachineMgr
		perm       *permissions.Permissions
	}
	Conf struct {
		MongoIP         string
		Channels        map[string]int
		LogicCount      int
		GateCount       int
		CommonConf      string
		GitCommit       string
		GitDelete       string
		GConf           string
		PerRedisHost    string
		PerRedisPwd     string
		OpWebIP         string
		MysqlUsr        string
		MysqlPwd        string
		RedisCharPwd    string
		RedisAccountPWd string
	}
)

var (
	localHost string
)

func New(pathFile string) Yada {
	data, err := ioutil.ReadFile(pathFile)
	if err != nil {
		log.Fatal("read file ", err)
	}

	fileConf := Conf{}
	err = json.Unmarshal([]byte(data), &fileConf)
	if err != nil {
		log.Fatal("config json ", err)
	}

	s, merr := mgo.Dial(fileConf.MongoIP)
	if merr != nil {
		log.Fatal("connect mongodb ", err)
	}

	s.SetMode(mgo.Monotonic, true)

	return &yada{
		e:          echo.New(),
		as:         NewAS(":3300"),
		zMgr:       NewZoneMgr(s),
		conf:       fileConf,
		machineMgr: NewMachineMgr(s, fileConf),
		session:    s,
	}
}

func (y *yada) RegisterWeb() {
	y.RegisterPerm()

	y.e.POST("/login", y.WebLogin)

	y.e.GET("/zone", y.GetZones)
	y.e.POST("/zone/add", y.AddZone)
	y.e.POST("/zone/save", y.SaveZone)
	y.e.GET("/zone/synMachine", y.SynMachine)
	y.e.POST("/zone/startZone", y.StartZone)
	y.e.POST("/zone/stopZone", y.StopZone)
	y.e.POST("/zone/del", y.DelZone)
	y.e.POST("/zone/updateZonelogdb", y.UpdateZonelogdb)
	y.e.POST("/zone/startAllZone", y.StartAllZone)
	y.e.POST("/zone/stopAllZone", y.StopAllZone)
	y.e.GET("/zone/zonelist", y.Zonelist)

	y.e.GET("/machine", y.GetMachines)
	y.e.POST("/machine/add", y.AddMachine)
	y.e.POST("/machine/save", y.SaveMachine)
	y.e.POST("/machine/del", y.DelMachine)
	y.e.GET("machine/common", y.CommonConfig)
	y.e.POST("machine/svnUpdate", y.SvnUpdate)
	y.e.GET("machine/svnUpdateAll", y.SvnUpdateAll)
}

func (y *yada) RegisterPerm() {
	userstate, err := permissions.NewUserStateWithPassword2(y.conf.PerRedisHost, y.conf.PerRedisPwd)
	if err != nil {
		log.Fatal(err)
	}
	y.perm = permissions.NewPermissions(userstate)
	y.perm.AddUserPath("/machine")
	y.perm.AddUserPath("/zone")

	y.perm.UserState().AddUser("cgyx", "cgyx!123", "bo@zombo.com")

	f := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if y.perm.Rejected(c.Response(), c.Request()) {
				// Deny the request
				//return echo.NewHTTPError(http.StatusForbidden, denyMessage)
				return c.String(http.StatusOK, "verify")
			}
			// Continue the chain of middleware
			log.Println(c.Request().Method, c.Request().RequestURI, err)
			return next(c)
		}
	}

	y.e.Use(f)
}

func (y *yada) Run() {
	go y.YadaCheck()
	go y.as.Listen()

	y.RegisterWeb()
	y.e.Static("/", "../client/dist/")
	y.e.File("/", "../client/dist/index.html")
	y.e.Logger.Fatal(y.e.Start(":1323"))
}

func (y *yada) YadaCheck() {
	for {
		if err := y.session.Ping(); err != nil {
			s, merr := mgo.Dial(y.conf.MongoIP)
			if merr != nil {
				log.Printf("reconnect mongodb fail:%v", merr)
			} else {
				y.session = s
			}

		}
		time.Sleep(time.Second * 10)
	}
}
func (y *yada) WebLogin(c echo.Context) error {
	user := c.FormValue("user")
	passwd := c.FormValue("password")
	if y.perm.UserState().CorrectPassword(user, passwd) {
		err := y.perm.UserState().Login(c.Response().Writer, user)
		if err != nil {
			c.String(http.StatusOK, err.Error())
		} else {
			c.String(http.StatusOK, "admin")
		}
	} else {
		c.String(http.StatusInternalServerError, "Login fail")
	}
	return nil
}

//获取区服信息
func (y *yada) GetZones(c echo.Context) error {
	rsp := ZoneRsp{
		Zstates: y.as.OnlineZones(),
		Items:   *y.zMgr.GetAllZoneInfo(),
		Result:  "OK",
	}
	return c.JSON(http.StatusOK, rsp)
}

//添加区服信息
func (y *yada) AddZone(c echo.Context) error {
	zone := Zone{}
	ret := ZoneRsp{}

	if err := c.Bind(&zone); err != nil {
		ret.Result = fmt.Sprintf("add zone bind data err %v", err)
		return c.JSON(http.StatusOK, ret)
	}

	if err := y.zMgr.AddZone(&zone); err != nil {
		ret.Result = fmt.Sprintf("add zone  err %v", err)
		return c.JSON(http.StatusOK, ret)
	} else {
		ret.Item = zone
		//新增的zone用到的机器加入到各自的用途中
		r := RelationZone{
			Zid:           zone.Zid,
			ZoneHost:      zone.ZoneHost,
			ZoneDBHost:    zone.ZoneDBHost,
			ZonelogdbHost: zone.ZonelogdbHost,
		}
		ret.Result = "OK"

		//新增机器用途信息
		y.machineMgr.OpZoneRelation(&r, RelationAdd)

		//通知agent更新配置
		//y.as.AddNewZone(zone.ZoneHost, zone.ZoneName, zone.Zid)
	}
	return c.JSON(http.StatusOK, ret)
}

//保存
func (y *yada) SaveZone(c echo.Context) error {
	m := SaveZoneReq{}
	ret := ZoneRsp{}

	if err := c.Bind(&m); err != nil {
		ret.Result = fmt.Sprintf("save zone bind data err", err)
		return c.JSON(http.StatusOK, ret)
	}

	if err := y.zMgr.SaveZone(m.OldZid, m.OldZoneName, &m.Item); err != nil {
		ret.Result = fmt.Sprintf("save zone info err, %v", err)
	} else {
		ret.Item = m.Item
		ret.Result = "OK"
		newRelation := &RelationZone{
			Zid:           m.Item.Zid,
			ZoneDBHost:    m.Item.ZoneDBHost,
			ZoneHost:      m.Item.ZoneHost,
			ZonelogdbHost: m.Item.ZonelogdbHost,
		}

		//更换机器用途信息
		oldRelation := y.zMgr.GetZoneRelation(m.OldZid)
		y.machineMgr.UpdateZone(oldRelation, newRelation)

	}

	return c.JSON(http.StatusOK, ret)
}

//删除
func (y *yada) DelZone(c echo.Context) error {
	ret := ZoneRsp{}
	m := ZoneReq{}
	if err := c.Bind(&m); err != nil {
		ret.Result = fmt.Sprintf("del zone bind data err", err)
		return c.JSON(http.StatusOK, ret)
	}
	zone := y.zMgr.GetZoneInfoByZid(m.Zid)
	if zone == nil {
		ret.Result = fmt.Sprintf("del zone get zone info err")
		return c.JSON(http.StatusOK, ret)

	}
	if err := y.zMgr.DelZone(m.Zid); err != nil {
		ret.Result = fmt.Sprintf("del zone err", err)
		return c.JSON(http.StatusOK, ret)
	} else {
		//删除相关配置，删除对应机器用途信息
		r := RelationZone{
			ZoneDBHost:    zone.ZoneDBHost,
			ZoneHost:      zone.ZoneHost,
			ZonelogdbHost: zone.ZonelogdbHost,
			Zid:           zone.Zid,
		}
		if err := y.machineMgr.DelZoneConf(&r); err != nil {
			ret.Result = fmt.Sprintf("del zone conf err, %v", err)
			return c.JSON(http.StatusOK, ret)
		}

		ret.Result = "OK"
		ret.Item = *zone
	}

	return c.JSON(http.StatusOK, ret)
}

func (y *yada) UpdateZonelogdb(c echo.Context) error {
	zReq := ZoneReq{}
	ret := ZoneRsp{}
	if err := c.Bind(&zReq); err != nil {
		ret.Result = fmt.Sprintf("update zonelogdb  bind data err", err)
		return c.JSON(http.StatusOK, ret)
	}

	m := y.machineMgr.GetMachineByName(zReq.Host)
	if m == nil {
		ret.Result = "FAlse"
		return c.JSON(http.StatusOK, ret)
	}
	logdb := "zonelog" + strconv.Itoa(zReq.Zid)
	s, _ := utils.ExeShellArgs2("sh", "update_zonelogdb", logdb, m.IP)
	if logdb != s {
		ret.Result = fmt.Sprintf("update zonelogdb fail,%s", s)
	} else {
		ret.Result = "OK"
	}
	return c.JSON(http.StatusOK, ret)
}

func (y *yada) SynMachine(c echo.Context) error {
	ret := ZoneRsp{}
	zid, _ := strconv.Atoi(c.QueryParam("zid"))
	hostname := c.QueryParam("hostname")

	zone := y.zMgr.GetZoneInfoByZid(zid)
	if zone == nil || zone.ZoneHost != hostname {
		return c.JSON(http.StatusOK, ret)
	}

	hostdir := y.conf.GConf + zone.ZoneHost
	curDir := hostdir + "/zone" + strconv.Itoa(zone.Zid) + "/"
	os.MkdirAll(curDir, os.ModePerm)

	zerr := y.machineMgr.ZoneLua(zone, curDir)
	gerr := y.machineMgr.GateLua(zone, curDir)
	cerr := y.machineMgr.CenterLua(zone, curDir)
	lerr := y.machineMgr.LogLua(zone, curDir)
	logicerr := y.machineMgr.LogicLua(zone, curDir)
	charErr := y.machineMgr.CharDBLua(zone, curDir)

	if zerr != nil || gerr != nil || cerr != nil || lerr != nil || logicerr != nil || charErr != nil {
		log.Printf("zone:%v,gate:%v,center:%v,log:%v, logic:%v, chardb:%v, logdberr:%v", zerr, gerr, cerr, lerr, logicerr, charErr)
		return c.JSON(http.StatusOK, ret)
	}

	if _, err := utils.ExeShell("sh", y.conf.GitCommit, "add or update zone"+strconv.Itoa(zone.Zid)); err != nil {
		log.Printf("exeshell fail %v", err)
		return err
	}

	ret.Result = y.as.UpdateZone(hostname)
	return c.JSON(http.StatusOK, ret)
}

func (y *yada) Zonelist(c echo.Context) error {
	zones := y.zMgr.GetAllZoneInfo()
	zlst := ZoneYunYingLst{}
	for _, v := range *zones {
		z := ZoneYunYing{
			Zid:      v.Zid,
			ZoneName: v.ZoneName,
			Channels: v.Channels,
		}
		m := y.machineMgr.GetMachineByName(v.ZonelogdbHost)
		if m == nil {
			log.Println("zonelist cannt find machine info machineName:", v.ZonelogdbHost)
			continue
		}
		z.ZonelogdbIP = m.IP
		zlst.Zlist = append(zlst.Zlist, z)

	}
	return c.JSON(http.StatusOK, zlst)
}

func (y *yada) StartZone(c echo.Context) error {
	ret := ZoneRsp{}
	m := ZoneReq{}
	if err := c.Bind(&m); err != nil {
		ret.Result = fmt.Sprintf("start zone  bind data err", err)
		return c.JSON(http.StatusOK, ret)
	}
	ret.Result = y.as.StartZone(m.Host, m.Zid)
	ret.Zstates = y.as.OnlineZones()

	log.Println("start ", ret)
	return c.JSON(http.StatusOK, ret)
}

func (y *yada) StopZone(c echo.Context) error {
	ret := ZoneRsp{}
	m := ZoneReq{}
	if err := c.Bind(&m); err != nil {
		ret.Result = fmt.Sprintf("stop zone  bind data err", err)
		return c.JSON(http.StatusOK, ret)
	}
	ret.Result = y.as.StopZone(m.Host, m.Zid)
	ret.Zstates = y.as.OnlineZones()
	return c.JSON(http.StatusOK, ret)
}

func (y *yada) StartAllZone(c echo.Context) error {
	ret := ZoneRsp{
		Result:  y.as.StartAllZone(),
		Zstates: y.as.OnlineZones(),
	}
	return c.JSON(http.StatusOK, ret)
}

func (y *yada) StopAllZone(c echo.Context) error {
	ret := ZoneRsp{
		Result:  y.as.StopAllZone(),
		Zstates: y.as.OnlineZones(),
	}
	return c.JSON(http.StatusOK, ret)
}

func (y *yada) MachineRspFunc() []Machine {
	ms := y.machineMgr.GetAllMachines()
	for _, v := range ms {
		v.Online, v.CodeVersion = y.as.CheckOnlineMachine(v.Hostname)
	}
	return ms
}

//获取机器信息
func (y *yada) GetMachines(c echo.Context) error {
	rsp := InitMachine{
		Items: y.MachineRspFunc(),
	}
	return c.JSON(http.StatusOK, rsp)
}

//添加机器信息
func (y *yada) AddMachine(c echo.Context) error {
	m := Machine{
		IP:       c.FormValue("IP"),
		Hostname: c.FormValue("hostname"),
		OutIP:    c.FormValue("outIP"),
	}
	ret := MachineAllRsp{Result: "OK"}
	if err := y.machineMgr.AddMachine(&m); err != nil {
		ret.Result = fmt.Sprintf("add machine add  %v", err)
	} else {
		ret.Items = y.MachineRspFunc()
	}
	return c.JSON(http.StatusOK, ret)
}

//保存
func (y *yada) SaveMachine(c echo.Context) error {
	m := Machine{}
	rsp := MachineAllRsp{Result: "OK"}

	jdata := c.FormValue("Item")
	err := json.Unmarshal([]byte(jdata), &m)
	if err != nil {
		rsp.Result = err.Error()
		return c.JSON(http.StatusOK, rsp)
	}

	oldhost := c.FormValue("Oldhost")
	if err := y.machineMgr.SaveMachine(oldhost, &m); err != nil {
		rsp.Result = fmt.Sprintf("save machine err %v", err)
	} else {
		rsp.Items = y.MachineRspFunc()
	}

	return c.JSON(http.StatusOK, rsp)
}

//删除
func (y *yada) DelMachine(c echo.Context) error {
	m := Machine{}
	ret := InitMachine{}
	if err := c.Bind(&m); err != nil {
		log.Printf("delmachine bind data err %v", err)
		return c.JSON(http.StatusOK, ret)
	}
	if err := y.machineMgr.DelMachine(m.Hostname); err != nil {
		log.Printf("delmachine err %v", err)
		return c.JSON(http.StatusOK, ret)
	}

	ret.Items = y.MachineRspFunc()
	return c.JSON(http.StatusOK, ret)
}

//生成登陆服、master/masterLog等上层服务器配置
func (y *yada) CommonConfig(c echo.Context) error {
	rsp := MachineRsp{Result: "OK"}
	os.Mkdir(y.conf.CommonConf, os.ModePerm)

	if err := y.machineMgr.LoginLua(); err != nil {
		rsp.Result = fmt.Sprintf(" common config %v", err)
		return c.JSON(http.StatusOK, rsp)
	}
	if err := y.machineMgr.MasterLua(); err != nil {
		rsp.Result = fmt.Sprintf(" common config  %v", err)
		return c.JSON(http.StatusOK, rsp)
	}
	if err := y.machineMgr.AccountDBLua(); err != nil {
		rsp.Result = fmt.Sprintf(" common config  accountdbLua %v", err)
		return c.JSON(http.StatusOK, rsp)
	}
	if err := y.machineMgr.MasterLogLua(); err != nil {
		rsp.Result = fmt.Sprintf(" common config  masterLua %v", err)
		return c.JSON(http.StatusOK, rsp)
	}

	if _, err := utils.ExeShell("sh", y.conf.GitCommit, "updata common Config"); err != nil {
		rsp.Result = fmt.Sprintf(" common config file excute shell err %v", err)
	}

	return c.JSON(http.StatusOK, rsp)
}

func (y *yada) SvnUpdate(c echo.Context) error {
	hostName := c.FormValue("HostName")
	rsp := MachineAllRsp{
		Result: y.as.UpdateSvn(hostName),
		Items:  y.MachineRspFunc(),
	}
	return c.JSON(http.StatusOK, rsp)
}

func (y *yada) SvnUpdateAll(c echo.Context) error {
	rsp := MachineAllRsp{
		Result: y.as.UpdateSvnAll(),
		Items:  y.MachineRspFunc(),
	}
	return c.JSON(http.StatusOK, rsp)
}
