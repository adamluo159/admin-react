package comInterface

type Connect struct {
	ID   int
	Port int
	IP   string
}

type MysqlLua struct {
	IP             string
	Port           int
	UserName       string
	Password       string
	FlushFrequency int
	DataBase       string
}

type RedisLua struct {
	IP       string
	Port     int
	Password string
}

type Gate struct {
	ID             int
	Zid            int
	ServerIP       string
	ServerPort     int
	ClientIP       string
	ClientPort     int
	ChannelIds     []int
	Open           bool
	Name           string
	ConnectServers map[string]interface{}
}

type Center struct {
	ID                    int
	Zid                   int
	IP                    string
	Port                  int
	OnlineNumberCheckTime int
	SingleServerLoad      int
	ConnectServers        map[string]interface{}
}

type CharDB struct {
	ID             int
	Zid            int
	IP             string
	Port           int
	Mysql          MysqlLua
	Redis          RedisLua
	ConnectServers map[string]interface{}
}

type Logic struct {
	ID             int
	Zid            int
	IP             string
	Port           int
	ConnectServers map[string]interface{}
	MapIds         []int
	LoadAllMapIds  bool
}

type Log struct {
	ID             int
	IP             string
	Port           int
	ConnectServers map[string]interface{}
}

type LogDBConf struct {
	DirName string
	IP      string
}

type Login struct {
	ID             int
	IP             string
	Port           int
	VesionStr      string
	ConnectServers map[string]interface{}
}

type Master struct {
	ID             int
	IP             string
	Port           int
	AllZoneOpen    bool
	ConnectServers map[string]interface{}
}
type AccountDB struct {
	ID             int
	Zid            int
	IP             string
	Port           int
	Mysql          MysqlLua
	Redis          RedisLua
	ConnectServers map[string]interface{}
}
