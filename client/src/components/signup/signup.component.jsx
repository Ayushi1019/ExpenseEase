import { LockOutlined, MailOutlined } from '@ant-design/icons';
import { Button, Form, Input } from 'antd';
import './style.css'

const Signup = () => {


  const onFinish=(values)=>{
    const body = {
      "name" : values.name,
      "email" : values.email,
      "password" : values.password
    }

    console.log(body)
    
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
