import React from 'react'
import './index.css'
import { NavLink } from 'react-router-dom'

export default ({ isCollapsed }) => {
  return (
    <ul className='side-menu' style={{ visibility: isCollapsed ? 'hidden' : 'visible' }}>
      <li className='menu-item'>
        <NavLink to='/events'>
          events
        </NavLink>
      </li>
      <li className='menu-item'>
        <NavLink to='/' exact>
          operations
        </NavLink>
      </li>
      <li className='menu-item'>
        <NavLink to='/environment' exact>
          environment
        </NavLink>
      </li>
    </ul>
  )
}
