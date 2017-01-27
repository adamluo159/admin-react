import React from 'react';
import ReactDOM from 'react-dom';
import {Form, Input, Button, Checkbox, message,Row, Col} from 'antd';
const FormItem = Form.Item;
export default class login extends React.Component {
  constructor(props){
    super(props);
  }

  handleSubmit(e) {
    e.preventDefault();
    console.log(e.target.userName.value)
    let a = {
      user: e.target.userName.value,
      pwd:e.target.password.value,
    }
    if(a.user === '' || a.pwd === ''){
      message.error("error")
    }
    else{
      message.success('收到表单值~~~ ：' + JSON.stringify(a, function (k, v) {
      if (typeof v === 'undefined') {
        return '';
      }
       return v;
       }));
    }
  }

  render() {
   return (
     <div style={{ width: 800, margin: '300px auto' }}>
     <Row type="flex" justify="start">
     <Col xs={20} sm={16} md={12} lg={14}>
      <Form horizontal onSubmit={this.handleSubmit} >
        <FormItem
          id="userName"
          label="账户："
          required>
    <Input placeholder="请输入账户名" id="userName" name="userName" size="large"/>
        </FormItem>
        <FormItem
          id="password"
          label="密码："
          required>
          <Input type="password" placeholder="请输入密码" id="password" name="password"/>
        </FormItem>
       <Button type="primary" htmlType="submit" className='login-form-button'>登录</Button>
      </Form>
      </Col>
      </Row>
      </div>
    );
  }
}
//ReactDOM.render(<Login/>, document.getElementById('root'));
