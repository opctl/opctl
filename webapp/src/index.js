import React from 'react';
import ReactDOM from 'react-dom';
import { Route, BrowserRouter as Router} from 'react-router-dom'
import PkgBrowser from './PkgBrowser';
import 'bootstrap/dist/css/bootstrap.css';

ReactDOM.render(
  <Router>
    <Route path="/" component={PkgBrowser}/>
  </Router>,
    document.getElementById('root')
);
