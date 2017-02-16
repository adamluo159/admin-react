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

export const machineColumns = [{
  title: 'HostName',
  dataIndex: 'hostname',
  key: 'hostname',
  width: '10%'
},
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
    width: '10%'
  },
  {
    title: '内网IP',
    dataIndex: 'IP',
    key: 'IP',
    width: '15%'
  },
  {
    title: '外网IP',
    dataIndex: 'outIP',
    key: 'outIP',
    width: '15%'
  },
  {
    title: '机器用途',
    dataIndex: 'describe',
    key: 'describe',
    width: '35%'
  },
  {
    title: 'Action',
    key: 'action',
    width: '15%'
  }]
export const typeOption = ['login', 'master', 'zone', 'pay', 'db', 'logdb']

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
  wihtelst: {
    Id: 'whitelst',
    label: '白名单'
  },
  channels: {
    Id: 'channels',
    label: '渠道',
    kinds: ['ios', 'yyb', 'xiaomi']
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
    span: 5
  },
  wrapperCol: {
    span: 12
  }
}

export const zoneData = {
  "1": {
    channels: ['IOS', 'yyb'],
    whitelst: true,
    zid: '1',
    zoneDBHost: '1',
    zoneHost: 'host1',
    zoneName: '一区大保健',
    zonelogdbHost: 'logdb02'
  },
  "2": {
    channels: ['IOS'],
    whitelst: true,
    zid: '2',
    zoneDBHost: 'hostdb1',
    zoneHost: 'host1',
    zoneName: '一区小宝剑',
    zonelogdbHost: 'logdb02'
  },
}
