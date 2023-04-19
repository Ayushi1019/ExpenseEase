import { Outlet } from 'react-router-dom';
import './style.css';
import { Tabs, Image } from "antd";
import Income from '../../components/income/income.component';
import Expense from '../../components/expense/expense.component';
import Budget from '../../components/budget/budget.component';


const Userhome = () => {

    const items = [
        {
          key: '1',
          label: `Income`,
          children: <div className='form-section'>
                <Income/>
          </div>,
        },
        {
          key: '2',
          label: `Expense`,
          children: <div className='form-section'>
               <Expense/> 
          </div>,
        },
        {
            key: '3',
            label: `Budget`,
            children: <div className='form-section'>
                 <Budget/> 
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

export default Userhome;
