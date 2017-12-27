import React from 'react';
import './index.css'
import MdDeHaze from 'react-icons/lib/md/dehaze';

export default ({ onCollapseToggled }) =>
  <div>
    <div className='top-menu'>
      <div className='menu-item logo'>
        <a
          href='/#/'
          onClick={e => {
            onCollapseToggled();
            e.preventDefault();
          }}
        >
          <MdDeHaze size={30} />
        </a>
      </div>
      <div className='menu-item logo'>
        <a href='https://opctl.io'>
          <img src='logo-text_rect.png' alt='opctl logo' height={30} />
        </a>
      </div>
    </div>
    <div style={{ height: '58px' }}>
    </div>
  </div>;

