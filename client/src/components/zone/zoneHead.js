import React, { Component } from 'react';
import { Select, Row, Col, Button } from 'antd';
const Option = Select.Option;

class zoneHead extends Component {
    constructor(props) {
        super(props);
        let {zoneData, channelData, registerFunc} = props
        if (channelData.length > 0) {
            this.state = {
                zonelst: zoneData[channelData[0]],
                selectZone: zoneData[channelData[0]][0],
                curChannel: channelData[0],
            }
        } else {
            this.state = {
                zonelst: [],
                selectZone: {},
                curChannel: "",
            }
        }
        registerFunc(this.freshData)
    }
    handleChannelChange(value) {
        let {zoneData, channelData} = this.props
        this.setState({
            zonelst: zoneData[value], selectZone: zoneData[value][0], curChannel: value
        });
    }
    onZoneChange(value) {
        let {zoneData, channelData, showFunc} = this.props
        let {curChannel} = this.state
        let zone = zoneData[curChannel][value]
        showFunc(zone.zid)
        //this.setState({ selectZone: zone });
    }
    render() {
        let {zoneData, channelData, addZoneFunc} = this.props
        let {zonelst, selectZone} = this.state
        const channelOptions = channelData.map(channel => <Option key={channel}>{channel}</Option>);
        const zoneOptions = zonelst.map((zone, index) => <Option key={index}>{zone.zoneName}</Option>);
        let channelValue, zoneValue
        if (channelData.length > 0) {
            channelValue = channelData[0]
            zoneValue = selectZone.zoneName
        } else {
            channelValue = ""
            zoneValue = ""

        }
        let width = {
            width: '100px'
        }
        return (
            <div>
                <Row>
                    <Col span={5}>
                        <Select
                            defaultValue={channelData[0]}
                            onChange={(e) => this.handleChannelChange(e)}
                            style={width}>
                            {channelOptions}
                        </Select>
                    </Col>
                    <Col span={5}>
                        <Select
                            defaultValue={selectZone.zoneName}
                            onChange={(e) => this.onZoneChange(e)}
                            style={width}>
                            {zoneOptions}
                        </Select>
                    </Col>
                    <Col span={5}>
                        <Button type="primary" onClick={()=>addZoneFunc()}>添加区服信息</Button>
                    </Col>
                </Row>
            </div>
        );
    }
}
export default zoneHead;