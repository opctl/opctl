import React from 'react';
import './index.css';
import { NavLink } from 'react-router-dom'

export default () => {

  return (
    <ul className='side-menu'>
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

