import { LockOutlined, MailOutlined, UserOutlined } from '@ant-design/icons';
import { Button, Form, Input, message } from 'antd';
import {Axios as axios} from 'axios';
import { API_URL } from '../../api';
import './style.css'

const Signup = () => {


  const onFinish=(values)=>{
    const body = {
      "name" : values.name,
      "email" : values.email,
      "password" : values.password
    }

    axios.post(API_URL + 'user',body)
    .then(({data, status}) => {
      const loginBody = {
        'email': body.email,
        'password': body.password
      }
      axios.post(API_URL + 'login',loginBody)
    .then(({data, status}) => {

      if(status === 200){
        message.success("Login Successful!")
        localStorage.setItem('token', data.token)
      }
      else{
        message.error("Something went Wrong!")
      }
      
    }).catch((error)=>{
      message.error("Error")
    });      
    }).catch((error)=>{
      message.error("Error")
    });
    
  }
  return (
    <Form
      name="normal_login"
      className="login-form"
      initialValues={{ remember: true }}
      onFinish={onFinish}
    >
      <Form.Item
        name="name"
        rules={[{ required: true, message: 'Please add your name' }]}
      >
        <Input className='login-input' prefix={<UserOutlined className="site-form-item-icon" />} placeholder="Name" />
      </Form.Item>
      <Form.Item
        name="email"
        rules={[{ required: true, message: 'Please input your Email!' }]}
      >
        <Input className='login-input' prefix={<MailOutlined className="site-form-item-icon" />} placeholder="Email" />
      </Form.Item>
      <Form.Item
        name="password"
        rules={[{ required: true, message: 'Please input your Password!' }]}
      >
        <Input className='login-input'
          prefix={<LockOutlined className="site-form-item-icon" />}
          type="password"
          placeholder="Password"
        />
      </Form.Item>

      <Form.Item>
        <Button className="login-form-button" type="primary" htmlType="submit">
          Log in
        </Button>
      </Form.Item>
    </Form>
  );
};

export default Signup;
