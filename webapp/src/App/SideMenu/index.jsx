import React from 'react';
import './index.css';
import {NavLink} from 'react-router-dom'

export default ({isCollapsed}) => {
  return (
    <ul className='side-menu' style={{visibility: isCollapsed ? 'hidden' : 'visible'}}>
      <li className='menu-item'>
        <NavLink to='/events'>
          events
        </NavLink>
      </li>
      <li className='menu-item'>
        <NavLink to='/' exact={true}>
          run
        </NavLink>
      </li>
    </ul>
  );
}

