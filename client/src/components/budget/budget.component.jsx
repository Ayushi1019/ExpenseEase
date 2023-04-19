import axios from 'axios';
import { LockOutlined, MailOutlined } from '@ant-design/icons';
import { Button, Form, Input } from 'antd';
import './style.css'




const Income = () => {


  const onFinish=(values)=>{
    const body = {
      "amount" : values.amount,
      "source" : values.source
    }
    console.log(body)
//     axios.post('http://localhost:8080/login',body)
//     .then(({data, status}) => {
//       console.log(data)
//       localStorage.token = data.token
//     }).catch((error)=>{
//       console.log(error)});
  }
  return (
    <Form
      name="normal_income"
      className="login-form"
      initialValues={{ remember: true }}
      onFinish={onFinish}
    >
      <Form.Item
        name="Amount"
        rules={[{ required: true, message: 'Please add your income' }]}
      >
        <Input className='login-input' prefix={<MailOutlined className="site-form-item-icon" />} placeholder="Enter amount" />
      </Form.Item>

      <Form.Item
        name="Income source"
        rules={[{ required: true, message: 'Please add income source!' }]}
      >
      </Form.Item>

      <Form.Item>
        <Button className="login-form-button" type="primary" htmlType="submit">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
};

export default Income;
