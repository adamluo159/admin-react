import React, { Component } from 'react'
import { Table, Tag } from 'antd'

class ZoneShowTable extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            show: []
        }
        this.columns = [{
            title: '区服名',
            dataIndex: 'name',
        }, {
            title: '状态',
            dataIndex: 'online',
            render: (text, record, index) => (this.onlineHandle(text, record, index))
        }, {
            title: '机器名',
            dataIndex: 'host',
        }];
    }
    onlineHandle(text, record, index) {
        let onlineColor, onlineText
        if (record.online) {
            onlineColor = "green"
            onlineText = "已连接"
        } else {
            onlineColor = "pink"
            onlineText = "未连接"
        }
        return (
            <Tag color={onlineColor}>{onlineText}</Tag>
        )
    }
    render() {
        let {show} = this.state
        const {getZoneName} = this.props
        show.forEach(v => {
            let k = v.zoneName.replace(/[^0-9]/ig,"")
            v.key = v.zoneName
            v.name = getZoneName(k)
        })
        return (
            <Table columns={this.columns} dataSource={show} size="small" />
        )
    }
}
export default ZoneShowTable
