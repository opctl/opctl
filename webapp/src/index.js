import React from 'react';
import ReactDOM from 'react-dom';
import {HashRouter as Router, Route} from 'react-router-dom'
import PkgBrowser from './PkgBrowser';
import EventBrowser from './EventBrowser';
import TopMenu from './TopMenu';
import SideMenu from './SideMenu';
import './bootstrap.css';

import {ToastContainer} from 'react-toastify';
import 'react-toastify/dist/ReactToastify.min.css'

ReactDOM.render(
  <Router>
    <div>
      <TopMenu/>
      <SideMenu/>
      <ToastContainer
        autoClose={20000}
        style={{zIndex: 100000}}
      />
      <div style={{marginLeft: '269px', marginTop: '57px'}}>
        <Route exact path="/" component={PkgBrowser}/>
        <Route path="/events" component={EventBrowser}/>
      </div>
    </div>
  </Router>,
  document.getElementById('root')
);
