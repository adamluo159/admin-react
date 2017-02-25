package zone

type Connect struct {
	ID   int
	Port int
	IP   string
}

type MysqlLua struct {
	IP             string
	Port           int
	UserName       string
	PassWord       string
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
	ChannelIds     []string
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
	ID    int
	Zid   int
	IP    string
	Port  int
	Mysql MysqlLua
	Redis RedisLua
}

type Logic struct {
	ID             int
	Zid            int
	IP             string
	Port           int
	ConnectServers map[string]interface{}
	MapIds         []int
}

type Log struct {
	ID             int
	IP             string
	Port           int
	ZoneLogMysql   MysqlLua
	GlobalLogMysql MysqlLua
}
