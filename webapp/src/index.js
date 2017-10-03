import React from 'react';
import ReactDOM from 'react-dom';
import { Route, HashRouter as Router} from 'react-router-dom'
import PkgBrowser from './PkgBrowser';
import EventBrowser from './EventBrowser';
import 'bootstrap/dist/css/bootstrap.css';

ReactDOM.render(
  <Router>
    <div>
      <Route exact path="/" component={PkgBrowser}/>
      <Route path="/events" component={EventBrowser}/>
    </div>
  </Router>,
    document.getElementById('root')
);
