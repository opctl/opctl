import React from 'react';
import ReactDOM from 'react-dom';
import {HashRouter as Router, Route} from 'react-router-dom'
import PkgBrowser from './PkgBrowser';
import EventBrowser from './EventBrowser';
import 'bootstrap/dist/css/bootstrap.min.css';

import {ToastContainer} from 'react-toastify';
import 'react-toastify/dist/ReactToastify.min.css'

ReactDOM.render(
  <Router>
    <div>
      <ToastContainer
        autoClose={20000}
        style={{zIndex: 100000}}
      />
      <Route exact path="/" component={PkgBrowser}/>
      <Route path="/events" component={EventBrowser}/>
    </div>
  </Router>,
  document.getElementById('root')
);
