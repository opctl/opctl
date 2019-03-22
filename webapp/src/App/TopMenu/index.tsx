import React from 'react'
import './index.css'
import { MdDehaze } from 'react-icons/md'

export default ({ onCollapseToggled }) =>
  <div>
    <div className='top-menu'>
      <div className='menu-item logo'>
        <a
          href='/#/'
          onClick={e => {
            onCollapseToggled()
            e.preventDefault()
          }}
        >
          <MdDehaze size={30} />
        </a>
      </div>
      <div className='menu-item logo'>
        <a href='https://opctl.io'>
          <img src='logo.svg' alt='opctl logo' height={30} />
        </a>
      </div>
    </div>
    <div style={{ height: '58px' }} />
  </div>
