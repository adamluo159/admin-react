package yada

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

	//
	MachineAllRsp struct {
		Result string
		Items  []Machine
	}
)
