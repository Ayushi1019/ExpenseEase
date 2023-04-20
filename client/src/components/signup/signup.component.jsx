import { LockOutlined, MailOutlined } from '@ant-design/icons';
import { Button, Form, Input } from 'antd';
import './style.css'
import axios from 'axios';

const Signup = () => {


  const onFinish=(values)=>{
    const body = {
      "name" : values.name,
      "email" : values.email,
      "password" : values.password
    }

    console.log(body)
    axios.post('http://localhost:8080/user',body)
    .then(({data, status}) => {
      console.log(data)
      if (status == 201){
        console.log("User created. We will go to login page in 3, 2, 1...")
        window.location.href='/'
      }
      else{
        console.log("User not created.")
        window.location.href='/'
      }
    }).catch((error)=>{
      console.log(error)});
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
        <Input className='login-input' prefix={<MailOutlined className="site-form-item-icon" />} placeholder="Name" />
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
