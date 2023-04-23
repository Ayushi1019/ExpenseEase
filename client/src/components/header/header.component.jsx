import { Fragment } from 'react';
import { Link, Outlet, useNavigate } from 'react-router-dom';
import Footer from '../footer/footer.component';
import {Button, Image} from "antd";
import './style.css'
import { useEffect } from 'react';
import { useState } from 'react';
import axios from 'axios';
import { API_URL } from '../../api';

const Header = () => {

  const [flag,setFlag] = useState(false)
  const navigate = useNavigate()

  useEffect(()=>{
    let token = localStorage.getItem("token")
    
    console.log(flag)
    if(token){
      console.log(token)
      setFlag(true)
    }
  },[])

  const signOut=()=>{
    axios.post(API_URL+'signout')
        .then(({data, status}) => {
          localStorage.removeItem("token")
          navigate("/")
        }).catch((error)=>{
          console.log(error)})
  }

  
  return (
    <Fragment>
      <div className='navigation-container'>
        <Link className='logo-container' to='/'>
          <Image preview={false} style={{borderRadius:'20px'}} width={90} height={60} src='ee-icon.jpeg'/>
        </Link>
        
        { flag &&
          <div className='nav-right-bar'>
            <Button className='nav-signup-button' onClick={signOut}>SIGN OUT</Button>
        </div>
        }
      </div>
      <Outlet />
      <Footer/>
    </Fragment>
  );
};

export default Header;
