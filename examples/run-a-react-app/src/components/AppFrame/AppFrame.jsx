import React, { Component } from 'react';
import './AppFrame.css';

class AppFrame extends Component {
  render() {
    return (
      <div className="app-frame">
        <header className="app-frame-header">
          <img src="/images/logo.svg" className="app-frame-logo" alt="logo" />
          <h1 className="app-frame-title">Application frame!</h1>
        </header>
        <div>
          { this.props.children }
        </div>
      </div>
    );
  }
}

export { AppFrame };
