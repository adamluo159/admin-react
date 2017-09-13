package yada

const longForm = "2006-01-02 15:04:05"

type (
	ZoneConf struct {
		ID             int
		Zid            int
		ServerIP       string
		ServerPort     int
		ClientIP       string
		ClientPort     int
		ChannelIds     []int
		Open           bool
		Name           string
		OpenTime       int64
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
		ID                    int
		Zid                   int
		IP                    string
		Port                  int
		OnlineNumberCheckTime int
		SingleServerLoad      int
		ConnectServers        map[string]interface{}
		OpenTime              int64
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
		ConnectServers map[string]interface{}
	}

	LogDBConf struct {
		DirName string
		IP      string
	}

	LoginConf struct {
		ID             int
		IP             string
		Port           int
		VesionStr      string
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
		NET_TIMEOUT_MSEC  int
		NET_MAX_CONNETION int
		StartService      []map[string]int
		LOG_INDEX         string
		LOG_MAXLINE       int
		OpenGM            int
	}

	WirteGame interface {
	}

	wirteGame struct {
		conf Conf
	}
)

const (
	CharDBPort     int = 7000
	CenterPort     int = 7100
	LogPort        int = 7200
	ClientPort     int = 7300
	ZonePort       int = 7400
	ZoneClientPort int = 7500
	GatePort       int = 7600 //gate1 7500起 gate2 7510起
	LogicPort      int = 7700 //logic1 7600起 logic2 7610起

	AccountDBPort    int = 6500
	RedisPort        int = 6379
	MysqlPort        int = 3306
	LoginWebPort     int = 1236
	ErrLogPort       int = 1237
	RedisAccountPort int = 6380
	MasterPort       int = 9500
	LoginPort        int = 9550

	NetTimeOut       int = 1000 * 30
	NetMaxConnection int = 5000
	DbproxyServer    int = 1
	LoginServer      int = 2
	CenterServer     int = 3
	LogicServer      int = 4
	LogServer        int = 5
	MasterServer     int = 6
	GateServer       int = 7
	ZoneServer       int = 8

	MasterCount int = 1
	LogMaxLine  int = 10000

	UserName             string = "root"
	PassWord             string = "cg2016"
	ConfDir              string = "/gConf/"
	RedisPassWord        string = ""
	RedisAccountPassWord string = "cg2016"
)

func NewGameWirter(conf Conf) WirteGame {
	return &wirteGame{
		conf: conf,
	}
}
