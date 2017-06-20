export const adminMenu = {
  '游戏指令': [
    {
      key: 'machineMgr',
      text: '机器管理  '
    },
    {
      key: 'gM',
      text: 'GM命令'
    },
    {
      key: 'zone',
      text: '游戏区服配置'
    }
  ]
}

export const Commonhost = {
  "master": true,
  "loginWeb": true,
  "login1": true,
  "errLog": true,
  "accountDB": true,
}

export const machineColumns = [{
  title: 'HostName',
  dataIndex: 'hostname',
  key: 'hostname',
  width: '10%',
},
{
  title: '内网IP',
  dataIndex: 'IP',
  key: 'IP',
  width: '10%'
},
{
  title: '外网IP',
  dataIndex: 'outIP',
  key: 'outIP',
  width: '10%'
},
{
  title: '机器用途',
  dataIndex: 'applications',
  key: 'applications',
  width: '50%',
  filters: [
    { text: 'empty', value: 0 },
    { text: 'zone', value: 1 },
    { text: 'zonelogdb', value: 3 },
    { text: 'zonedb', value: 2 },
  ],
},
{
  title: 'Action',
  key: 'action',
  width: '20%'
}]

export const zoneConfig = {
  zoneInput: [
    {
      Id: 'zid',
      label: 'zid'
    }, {
      Id: 'zoneName',
      label: '区服名'
    }, {
      Id: 'zoneHost',
      label: '区服主机名'
    }, {
      Id: 'zoneDBHost',
      label: 'db主机名'
    }, {
      Id: 'zonelogdbHost',
      label: 'logdb主机名'
    }
  ],
  whitelst: {
    Id: 'whitelst',
    label: '白名单'
  },
  channels: {
    Id: 'channels',
    label: '渠道',
    kinds: ['ios', 'yyb', 'xiaomi']
  },
  switchEdit: {
    Id: 'edit',
    label: '',
    layout: {
      wrapperCol: {
        span: 12
      }
    },
    options: {}
  }


}
export const zoneOptions = {
  rules: [
    {
      required: true,
      message: '不能为空'
    }
  ]
}
export const formItemLayout = {
  labelCol: {
    span: 4,
  },
  wrapperCol: {
    span: 10,
  }
}