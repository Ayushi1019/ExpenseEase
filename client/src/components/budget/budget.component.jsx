import axios from 'axios';
import {EditOutlined,DeleteOutlined,CheckOutlined,CloseOutlined} from "@ant-design/icons"

import { Button, Form, Input, Popconfirm } from 'antd';
import './style.css'
import React, { useState } from 'react';
import { Table } from 'antd';
import { useEffect } from 'react';
import { API_URL } from '../../api';
import moment from "moment"

const EditableCell = ({
  editing,
  dataIndex,
  title,
  inputType,
  record,
  index,
  children,
  ...restProps
}) => {
  const inputNode = <Input />;
  return (
    <td {...restProps}>
      {editing ? (
        <Form.Item
          name={dataIndex}
          style={{
            margin: 0,
          }}
          rules={[
            {
              required: true,
              message: `Please Input ${title}!`,
            },
          ]}
        >
          {inputNode}
        </Form.Item>
      ) : (
        children
      )}
    </td>
  );
};

const Budget = () => {

  const [form] = Form.useForm();
  const [budgetForm] = Form.useForm();
  const [budgets,setBudgets] = useState([])

  const [editingKey, setEditingKey] = useState('');
  const isEditing = (record) => record.ID === editingKey;

  useEffect(()=>{
    getAllBudgets()
  },[])

  const edit = (record) => {
    form.setFieldsValue({
      amount: '',
      category: '',
      ...record,
    });
    setEditingKey(record.ID);
  };

  const cancel = () => {
    setEditingKey('');
  };

  const getAllBudgets = ()=>{

    axios.defaults.headers.common['Authorization'] = `${localStorage.getItem('token')}`;
    axios.get(API_URL + 'budgets')
    .then(({data}) => {

      setBudgets(data)
      
    }).catch((error)=>{
      console.log(error)
    });
    
  }


  const onFinish=(values)=>{
    const body = {
      "amount" : parseFloat(values.amount),
      "category" : values.category,
      "created_at" : moment(new Date()).format("YYYY-MM-DD")
    }
    console.log(body)
    axios.post(API_URL+'budget',body)
    .then(({data, status}) => {
      getAllBudgets()
      budgetForm.resetFields()
    }).catch((error)=>{
      console.log(error)});
  }

  const onEditBudget= async (id)=>{
      try {
        const row = await form.validateFields();
          const body = {
          "amount" : parseFloat(row['Amount']),
          "category" : row['Category'],
          "created_at" : moment(new Date()).format("YYYY-MM-DD")
        }
        console.log(body)
        axios.put(API_URL+'budget/'+id,body)
        .then(({data, status}) => {
          setEditingKey('');
          getAllBudgets()
        }).catch((error)=>{
          console.log(error)});
      } catch (errInfo) {
        console.log('Validate Failed:', errInfo);
      }
    
  }
  const onDeleteBudget=(id)=>{
      axios.delete(API_URL+'budget/'+id)
      .then(({data, status}) => {
        getAllBudgets()
      }).catch((error)=>{
        console.log(error)});
  }

  const columns = [
    {
      title: 'Amount',
      dataIndex: 'Amount',
      editable: true,
    },
    {
      title: 'Category',
      dataIndex: 'Category',
      editable: true,
    },
    {
      title: 'Date',
      dataIndex: 'Created_at',
      editable: false,
    },
    {
      title: 'Action',
      dataIndex: 'ID',
      render: (_, record) => {
        const editable = isEditing(record);
        return editable ? (
          <span>
            <Button style={{backgroundColor:'green'}} className='action-btn' 
            onClick={()=>onEditBudget(record.ID)}>
              <CheckOutlined style={{fontSize:'18px'}} />
            </Button>

    
            <Popconfirm title="Sure to cancel?" onConfirm={cancel}>
            <Button style={{backgroundColor:'red',marginLeft:'5px'}} className='action-btn'><CloseOutlined style={{fontSize:'18px'}} /></Button>
            </Popconfirm>
          </span>
        ) : (
          <span>
          <Button style={{backgroundColor:'green'}} className='action-btn' onClick={() => edit(record)}><EditOutlined style={{fontSize:'18px'}} /></Button>
          
          <Popconfirm title="Sure to delete?" onConfirm={()=>onDeleteBudget(record.ID)}>         
            <Button style={{backgroundColor:'red',marginLeft:'5px'}} className='action-btn'><DeleteOutlined style={{fontSize:'18px'}} /></Button>
            </Popconfirm>
          </span>
        );
      },
    },
  ];

  const mergedColumns = columns.map((col) => {
    if (!col.editable) {
      return col;
    }
    return {
      ...col,
      onCell: (record) => ({
        record,
        inputType: 'text',
        dataIndex: col.dataIndex,
        title: col.title,
        editing: isEditing(record),
      }),
    };
  });

  return (
    <div>

    <Form
      form={budgetForm}
      name="normal_income"
      className="login-form"
      initialValues={{ remember: true }}
      onFinish={onFinish}
    >
      <Form.Item
        name="amount"
        rules={[{ required: true, message: 'Please add your budget' }]}
      >
        <Input className='login-input' placeholder="Enter amount" />
      </Form.Item>

      <Form.Item
        name="category"
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
    <Form form={form} component={false}>
    <Table dataSource={budgets} components={{
          body: {
            cell: EditableCell,
          },
        }}
        columns={mergedColumns}
        rowClassName="editable-row"
        pagination={{
          position: ['none']
        }}
      />
  </Form>
    </div>
  );
};

export default Budget;
