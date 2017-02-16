import React, { Component } from 'react';
import { Select, Row, Col, Button } from 'antd';
const Option = Select.Option;

class zoneHead extends Component {
    constructor(props) {
        super(props);
        console.log(props)
        let {zoneData, channelData} = props
        this.state = {
            zonelst: zoneData[channelData[0]],
            selectZone: zoneData[channelData[0]][0],
            curChannel: channelData[0],
        }
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
        console.log("aaa", zone)
        this.setState({ selectZone: zone });
    }
    onbtnAddZone(value){
        value.preventDefault()
        let {addZoneFunc} = this.props
        console.log("headheadhead",value)
        addZoneFunc()
    }

    render() {
        let {zoneData, channelData, addZoneFunc} = this.props
        let {zonelst, selectZone} = this.state
        const channelOptions = channelData.map(channel => <Option key={channel}>{channel}</Option>);
        const zoneOptions = zonelst.map((zone, index) => <Option key={index}>{zone.zoneName}</Option>);
        return (
            <div>
                <Row>
                    <Col span={2}>
                        <Select
                            defaultValue={channelData[0]}
                            onChange={(e) => this.handleChannelChange(e)}>
                            {channelOptions}
                        </Select>
                    </Col>
                    <Col span={5}>
                        <Select
                            value={selectZone.zoneName}
                            onChange={(e) => this.onZoneChange(e)}>
                            {zoneOptions}
                        </Select>
                    </Col>
                    <Col span={5}>
                        <Button type="primary" onClick={addZoneFunc}>添加区服信息</Button>
                    </Col>
                </Row>
            </div>
        );
    }
}
export default zoneHead;