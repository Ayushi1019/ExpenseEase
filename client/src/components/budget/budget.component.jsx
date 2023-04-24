import axios from 'axios';
import {EditOutlined,DeleteOutlined,CheckOutlined,CloseOutlined} from "@ant-design/icons"

import { Button, Form, Input, Popconfirm, Select } from 'antd';
import './style.css'
import React, { useState } from 'react';
import { Table } from 'antd';
import { useEffect } from 'react';
import { API_URL } from '../../api';
import moment from "moment"
import { Pie } from '@ant-design/plots';

const globalMonths = [
  {
    key: 1,
    value: 'Jan'
  },
  {
    key: 2,
    value: 'Feb'
  },
  {
    key: 3,
    value: 'Mar'
  },
  {
    key: 4,
    value: 'Apr'
  },
  {
    key: 5,
    value: 'May'
  },
  {
    key: 6,
    value: 'Jun'
  },
  {
    key: 7,
    value: 'Jul'
  },
  {
    key: 8,
    value: 'Aug'
  },
  {
    key: 9,
    value: 'Sept'
  },
  {
    key: 10,
    value: 'Oct'
  },
  {
    key: 11,
    value: 'Nov'
  },
  {
    key: 12,
    value: 'Dec'
  },
]

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
  const [graphData,setGraphData] = useState({1:[],2:[],3:[],4:[],5:[],6:[],7:[],8:[],9:[],10:[],11:[],12:[]})
  const [selectedMonth,setMonth] = useState(1)

  const [editingKey, setEditingKey] = useState('');
  const isEditing = (record) => record.ID === editingKey;

  useEffect(()=>{
    getAllBudgets()
    getFilteredData()
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
    axios.post(API_URL+'budget',body)
    .then(({data, status}) => {
      getAllBudgets()
      budgetForm.resetFields()
    }).catch((error)=>{
      console.log(error)});
  }

  const getFilteredData = ()=>{

    axios.defaults.headers.common['Authorization'] = `${localStorage.getItem('token')}`;
    axios.get(API_URL + 'budget_by_month')
    .then(({data}) => {

      let tmp = {}
      let finalData = {}

      for(let [key,val] of Object.entries(data)){

        let newKey = parseInt(key.split("-")[1])
        tmp = {
          ...tmp,
          [newKey] : val
        }
      }

      for(let key of Object.keys(tmp)){

        let newDict = {}

        for(let [k,v] of Object.entries(tmp[key])){

          if(k in newDict){
            newDict[k] += parseInt(v.map(i=> i['Amount']))
          }
          else{
            newDict[k] = parseInt(v.map(i=> i['Amount']))
          }

        }

        tmp[key] = newDict
      }

      

      for(let k of Object.keys(tmp)){
        let d = []

        for(let [key,val] of Object.entries(tmp[k])){
          let obj = {
            'type' : key,
            'value' : val
          }

          d.push(obj)
        }

        finalData[k] = d
      }
      setGraphData(finalData)
      
    }).catch((error)=>{
      console.log(error)
    });
    
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
          getFilteredData()
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
        getFilteredData()
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

  const config = {
    appendPadding: 10,
    angleField: 'value',
    colorField: 'type',
    radius: 0.9,
    label: {
      type: 'inner',
      offset: '-30%',
      content: ({ percent }) => `${(percent * 100).toFixed(0)}%`,
      style: {
        fontSize: 14,
        textAlign: 'center',
      },
    },
    interactions: [
      {
        type: 'element-active',
      },
    ],
  };

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
  <div style={{marginTop:'30px'}}>
    <Select value={selectedMonth} style={{
        width: 80,
      }} 
      onChange={(val)=>setMonth(val)}>
      {globalMonths.map(m=>

      <Select.Option value={m.key}>{m.value}</Select.Option>
        )}
    </Select>
    {
      graphData[selectedMonth] !== undefined ?
        <Pie {...config} data={graphData[selectedMonth]}/>
          :
        <div>
          No data
      </div>
    }
  </div>
    </div>
  );
};

export default Budget;
