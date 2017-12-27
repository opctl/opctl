import React, {Component} from 'react';
import Spaces from './views/Spaces';
import Events from "./views/Events";
import {HashRouter} from 'react-router-dom'
import TopMenu from './TopMenu';
import SideMenu from './SideMenu';
import {Route} from 'react-router-dom'
import {ToastContainer} from 'react-toastify';

// shared vendor css
import './bootstrap.css';
import 'react-toastify/dist/ReactToastify.min.css'
import 'react-virtualized/styles.css'
import 'react-grid-layout/css/styles.css';
import 'react-resizable/css/styles.css';

export default class Page extends Component {
  constructor(props) {
    super(props);

    this.state = {};
  }

  handleSideMenuToggled() {
    this.setState(state => {
      return {
        isSideMenuCollapsed: !state.isSideMenuCollapsed,
      };
    })
  }

  render() {
    return (
      <HashRouter>
        <div style={{height: '100vh'}}>
          <TopMenu onCollapseToggled={() => this.handleSideMenuToggled()}/>
          <SideMenu isCollapsed={this.state.isSideMenuCollapsed}/>
          <ToastContainer
            autoClose={20000}
            style={{zIndex: 100000}}
          />
          <div
            style={{
              height: 'calc(100vh - 58px)',
              marginLeft: this.state.isSideMenuCollapsed ? '0' : '269px',
            }}>
            <Route exact path="/" component={Spaces}/>
            <Route path="/events" component={Events}/>
          </div>
        </div>
      </HashRouter>
    );
  }
}
