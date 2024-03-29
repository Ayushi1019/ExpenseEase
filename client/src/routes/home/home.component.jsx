import { Outlet, useNavigate } from 'react-router-dom';
import './style.css';
import { Tabs, Image } from "antd";
import Login from '../../components/login/login.component';
import Signup from '../../components/signup/signup.component';
import { useEffect } from 'react';


const Home = () => {
  const navigate = useNavigate()

  useEffect(()=>{
    let token = localStorage.getItem("token")
    console.log(token)
    if(token){
      navigate('/home')
    }

  },[navigate])

    const items = [
        {
          key: '1',
          label: `Login`,
          children: <div className='form-section'>
                <Login/>
          </div>,
        },
        {
          key: '2',
          label: `Signup`,
          children: <div className='form-section'>
               <Signup/> 
          </div>,
        },
      ];

  return (

    <div className='container'>
        <div className='setting-segment'>
            <Image preview={false} src='/home-wall.jpg' />
        </div>
        <div className='segment-container'>
            <Tabs size='large' defaultActiveKey="1" items={items}/>
            
        </div>

      <Outlet />
    </div>
  );
};

export default Home;
