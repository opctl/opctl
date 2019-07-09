import React, { Component } from 'react';
import { HelloWorld } from '../../components/HelloWorld.jsx'
import './Home.css';

class Home extends Component {
  render() {
    return (
      <div className="Home">
        <h1>This is the / Home page</h1>
        <HelloWorld />
      </div>
    );
  }
}

export { Home };
