import React, { Component } from 'react'
import { Table, Tag } from 'antd'
import './zone.css'



class ZoneShowTable extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            show: []
        }
        this.columns = [{
            title: '区服名',
            dataIndex: 'zoneName',
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
        console.log(show)
        show.forEach(v => {
            v.key = v.zoneName
        })
        return (
            <Table columns={this.columns} dataSource={show} size="small" />
        )
    }
}
export default ZoneShowTable