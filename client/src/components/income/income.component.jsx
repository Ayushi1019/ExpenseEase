import axios from 'axios';

import { Button, Form, Input } from 'antd';
import './style.css'
import React, { useState } from 'react';
import { Space, Table, Tag } from 'antd';

import { Divider, Radio } from 'antd';
const { Column, ColumnGroup } = Table;


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
    <div>

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
        name="Income source"
        rules={[{ required: true, message: 'Please add income source!' }]}
      >
        <Input className='login-input' placeholder="Enter source of Income" />
      </Form.Item>

      <Form.Item>
        <Button className="login-form-button" type="primary" htmlType="submit">
          Submit
        </Button>
        </Form.Item>  
    </Form>
    <Table dataSource={data}>
   
      <Column title="Amount" dataIndex="amount" key="amount" />
      <Column title="Source of income" dataIndex="sourceofincome" key="sourceofincome" />
   
    <Column title="Date" dataIndex="date" key="date" />
    
  </Table>
    </div>
  );
};

export default Income;
