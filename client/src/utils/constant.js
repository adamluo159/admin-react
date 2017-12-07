export const adminMenu = {
	'游戏指令': [
	{
		key: 'machineMgr',
		text: '机器管理  '
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
	title: 'Host',
	dataIndex: 'hostname',
	key: 'hostname',
	width: '5%',
},
{
	title: "SvnVer",
	dataIndex: 'codeVersion',
	key: 'codeVersion',
	width: '5%',
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
	{ text: '空闲机器', value: 0 },
	{ text: 'zone', value: 1 },
	{ text: 'zonelogdb', value: 3 },
	{ text: 'zonedb', value: 2 },
	{ text: '包含公共进程', value: 4},
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
	}], 
	cghostSelect:[
	{
		Id: 'zoneHost',
		label: '区服主机名'
	}, {
		Id: 'zoneDBHost',
		label: 'db主机名'
	}, {
		Id: 'zonelogdbHost',
		label: '日志db主机名'
	}, {
		Id: 'datalogdbHost',
		label: '运营db主机名'
	}],
	whitelst: {
		Id: 'whitelst',
		label: '白名单'
	},
	channels: {
		Id: 'channels',
		label: '渠道',
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
	},
	datePick: {
		Id: 'opentime',
		label: '开服时间'
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
