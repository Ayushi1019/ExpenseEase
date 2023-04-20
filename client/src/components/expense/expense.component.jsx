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

  const data = [
    {
      key: '1',
      amount: '10',
      sourceofincome: 'job',
      date: 32,
      
    },
    {
        key: '1',
        amount: '10',
        sourceofincome: 'job',
        date: 32,
     
    },
    {
        key: '1',
        amount: '10',
        sourceofincome: 'job',
        date: 32,
      
    },
  ];
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
      <Input className='login-input' placeholder="Enter amount" />
      </Form.Item>

      <Form.Item
        name="Income category"
        rules={[{ required: true, message: 'Please add category!' }]}
      >
      <Input className='login-input' placeholder="Enter category" />
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
