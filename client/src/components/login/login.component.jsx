import axios from 'axios';
import { LockOutlined, MailOutlined } from '@ant-design/icons';
import { Button, Form, Input, message } from 'antd';
import './style.css'
import { API_URL } from '../../api';
import { useNavigate } from 'react-router-dom';


const Login = () => {
  const navigate = useNavigate()


  const onFinish=(values)=>{
    const body = {
      "email" : values.email,
      "password" : values.password
    }
    axios.post(API_URL + 'login',body)
    .then(({data, status}) => {

      if(status === 200){
        message.success("Login Successful!")
        localStorage.setItem('token', data.token)
        navigate("/home")

      }
      else{
        message.error("Something went Wrong!")
      }
      
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

export default Login;
