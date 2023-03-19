import { Fragment } from 'react';
import { Link, Outlet } from 'react-router-dom';
import Footer from '../footer/footer.component';
import {Button} from "antd";
import './style.css'

const Header = () => {

  return (
    <Fragment>
      <div className='navigation-container'>
        <Link className='logo-container' to='/'>
          Logo
        </Link>
        
        <div className='nav-right-bar'>
            <Button className='nav-signup-button'>SIGN UP</Button>
        </div>
      </div>
      <Outlet />
      <Footer/>
    </Fragment>
  );
};

export default Header;
