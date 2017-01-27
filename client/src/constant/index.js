export const adminMenu = {
  "游戏指令": [
    {
        key:"machineMgr", 
        text:"机器管理  "
    },
    {
        key:"gM",
        text:"GM命令"
    }
    ],
}

export const machineColumns = [{
        title: 'HostName',
        dataIndex: 'hostname',
        key: 'hostname',
        width: '10%',
      },
      {
        title: '类型',
        dataIndex: 'type',
        key: 'type',
        width: '10%',
     },
      {
        title: '内网IP',
        dataIndex: 'IP',
        key: 'IP',
        width: '15%',
     },
      {
        title: '外网IP',
        dataIndex: 'outIP',
        key: 'outIP',
        width: '15%',
     },
     {
       title: '机器用途',
       dataIndex:'describe',
       key :'describe',
       width: '35%'
     },
      {
        title: 'Action',
        key: 'action',
        width: '15%',
}]
export const typeOption = ['login', 'master', 'zone', 'pay', 'db', 'logdb']
