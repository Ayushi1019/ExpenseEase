import axios from 'axios';
import { LockOutlined, MailOutlined } from '@ant-design/icons';
import { Button, Form, Input } from 'antd';
import './style.css'
import { useNavigate } from 'react-router-dom';


function validateEmail(email) {
  const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(String(email).toLowerCase());
}
function validatePassword(password) {
  if (password.length < 8) {
      return false;
  }
  if (password.search(/[a-z]/i) < 0) {
      return false;
  }
  if (password.search(/[0-9]/) < 0) {
      return false;
  }
  return true;
}


const Login = () => {
  const navigate = useNavigate();

  const onFinish=(values)=>{
    const body = {
      "email" : values.email,
      "password" : values.password
    }
    console.log(body)
    axios.post('http://localhost:8080/login',body)
    .then(({data, status}) => {
      console.log(data)
      if(data.token !== null){
        localStorage.token = data.token
        navigate('/home');
      }
      else{
        console.log("Error karke kuch aaya tho aaaa login nahi kar paa rahe hai aap log.")
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
