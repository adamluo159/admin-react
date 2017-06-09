package comInterface

type SRV map[string]int
type ServerConfigHead struct {
	NET_TIMEOUT_MSEC  int
	NET_MAX_CONNETION int
	StartService      []SRV
	LOG_INDEX         string
	LOG_MAXLINE       int
	OpenGM            int
}

const (
	RelationDel int = 1
	RelationAdd int = 2
)
const (
	CharDBPort    int    = 7000
	GatePort      int    = 7100
	CenterPort    int    = 7200
	LogicPort     int    = 7300 //logic1 7300起 logic2 7400起
	ClientPort    int    = 7800
	MasterPort    int    = 9500
	LoginPort     int    = 7550
	LogPort       int    = 7600
	AccountDBPort int    = 6500
	UserName      string = "root"
	PassWord      string = "cg2016"
	RedisPort     int    = 6379
	MysqlPort     int    = 3306

	NetTimeOut       int = 1000 * 30
	NetMaxConnection int = 5000
	DbproxyServer    int = 1
	LoginServer      int = 2
	CenterServer     int = 3
	LogicServer      int = 4
	LogServer        int = 5
	MasterServer     int = 6
	GateServer       int = 7

	MasterCount int = 1
	LogMaxLine  int = 10000
)

type ZoneStates struct {
	Host     string `json:"host"`
	ZoneName string `json:"zoneName"`
	Online   bool   `json:"online"`
}

type RelationZone struct {
	Zid           int
	ZoneHost      string
	ZoneDBHost    string
	ZonelogdbHost string
}

//机器信息
type Machine struct {
	Hostname     string `json:"hostname" bson:"hostname"`
	IP           string
	OutIP        string   `json:"outIP" bson:"outIP"`
	Applications []string `json:"applications" bson:"applications"`
	Online       bool
}

type Aserver interface {
	StartZone(host string, zid int) int
	StopZone(host string, zid int) int
	CheckOnlineMachine(mName string) bool
	UpdateZone(host string) int
	StartAllZone() int
	StopAllZone() int
	OnlineZones() []ZoneStates
	AddNewZone(host string, zone string, zid int)
}

type MachineMgr interface {
	UpdateZone(old *RelationZone, new *RelationZone)
	GetMachineByName(name string) *Machine
	OpZoneRelation(r *RelationZone, op int)
	GetAllMachines() []Machine
}
