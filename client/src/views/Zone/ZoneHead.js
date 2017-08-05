import React, { Component } from 'react';
import { Select, Row, Col, Button } from 'antd';
const Option = Select.Option;

class ZoneHead extends Component {
    constructor(props) {
        super(props);
        this.state = {
        }
    }
    handleChannelChange(value) {
        let {curChannel} = this.state
        if (value == curChannel) {
            return
        }
        let {zoneData, channelData} = this.props
        this.ZoneOptions = []
        let lastzName = ""
        for (var key in this.channelLst[value]) {
            lastzName = this.channelLst[value][key].zoneName
            this.ZoneOptions.push(<Option key={key}>{lastzName}</Option>)
        }
        this.setState({
            curChannel: value,
            curZoneName: lastzName,
        });
    }
    onZoneChange(value) {
        let {showFunc} = this.props
        this.setState({
            curZoneName: this.zoneData[value].zoneName,
        })
        showFunc(value)
    }
    Init(zoneData, showZid) {
        this.chanOptions = []
        this.ZoneOptions = []
        this.channelLst = {}
        this.zoneData = zoneData
        let {curChannel, curZoneName} = this.state
        for (var k in zoneData) {
            zoneData[k].channels.forEach(channv => {
                if (this.channelLst[channv] == null) {
                    this.channelLst[channv] = {}
                    this.chanOptions.push(<Option key={channv}>{channv}</Option>)
                }
                curChannel = curChannel || channv
                this.channelLst[channv][k] = zoneData[k]
            })
        }
        let lastzName = ""
        let lastZid = 0
        for (var key in this.channelLst[curChannel]) {
            lastzName = this.channelLst[curChannel][key].zoneName
            lastZid = Number(key)
            this.ZoneOptions.push(<Option key={key}>{lastzName}</Option>)
        }
        if (showZid) {
            lastZid = showZid
            curChannel = zoneData[showZid].channels[0]
        }

        let zoneName = ""
        if (zoneData[lastZid]) {
            zoneName = zoneData[lastZid].zoneName
        }
        this.setState({
            curChannel: curChannel || "",
            curZoneName: zoneName,
        });
        return lastZid
    }
    render() {
        let {curChannel, curZoneName} = this.state
        let {addZoneFunc} = this.props
        let width = {
            width: '100px'
        }
        return (
            <div>
                <Row>
                    <Col span={1}>
                        <Select
                            value={curChannel}
                            onChange={(e) => this.handleChannelChange(e)}
                            style={width}>
                            {this.chanOptions}
                        </Select>
                    </Col>
                    <Col span={1} offset={1}>
                        <Select
                            value={curZoneName}
                            onChange={(e) => this.onZoneChange(e)}
                            style={width}>
                            {this.ZoneOptions}
                        </Select>
                    </Col>
                    <Col span={1} offset={1}>
                        <Button type="primary" onClick={() => addZoneFunc()}>添加区服信息</Button>
                    </Col>
                </Row>
            </div>
        );
    }
}
export default ZoneHead;