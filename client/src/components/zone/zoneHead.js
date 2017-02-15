import React, {Component} from 'react';
import {Select} from 'antd';
const Option = Select.Option;

const channelData = ['IOS', 'yyb'];
const cityData = {
    IOS: [
        '一区大保健', '二区小保健'
    ],
    yyb: ['1服大宝剑', '2服小宝剑']
};

class zoneHead extends Component {
    constructor(props) {
        super(props);
        this.state = {
            cities: cityData[channelData[0]],
            secondCity: cityData[channelData[0]][0]
        }
    }
   handleChannelChange(value) {
        this.setState({cities: cityData[value], secondCity: cityData[value][0]
        });
    }
    onZoneChange(value) {
        this.setState({secondCity: value});
    }
    render() {
        const channelOptions = channelData.map(province => <Option key={province}>{province}</Option>);
        const zoneOptions = this
            .state
            .cities
            .map(city => <Option key={city}>{city}</Option>);
        return (
            <div>
                <Select
                    defaultValue={channelData[0]}
                    style={{
                    width: 90
                }}
                    onChange={(e)=>this.handleChannelChange(e)}>
                    {channelOptions}
                </Select>
                <Select
                    value={this.state.secondCity}
                    style={{
                    width: 90
                }}
                    onChange={(e)=>this.onZoneChange(e)}>
                    {zoneOptions}
                </Select>
            </div>
        );
    }
}
export default zoneHead;