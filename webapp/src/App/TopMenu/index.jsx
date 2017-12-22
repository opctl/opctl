import React from 'react';
import './index.css'
import DeHaze from 'react-icons/lib/md/dehaze';

export default ({onCollapseToggled}) =>
  <div className='top-menu'>
    <div className='menu-item logo'>
      <a
        href='/#/'
        onClick={e => {
          onCollapseToggled();
          e.preventDefault();
        }}
      >
        <DeHaze size={30}/>
      </a>
    </div>
    <div className='menu-item logo'>
      <a href='https://opctl.io'>
        <img src='logo-text_rect.png' alt='opctl logo' height={30}/>
      </a>
    </div>
  </div>;

