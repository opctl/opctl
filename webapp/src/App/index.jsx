import React, {Component} from 'react';
import PkgBrowser from "./PkgBrowser";
import EventBrowser from "./EventBrowser";
import {HashRouter} from 'react-router-dom'
import TopMenu from './TopMenu';
import SideMenu from './SideMenu';
import {Route} from 'react-router-dom'
import {ToastContainer} from 'react-toastify';
import './bootstrap.css';
import 'react-toastify/dist/ReactToastify.min.css'

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
        <div>
          <TopMenu onCollapseToggled={() => this.handleSideMenuToggled()}/>
          <SideMenu isCollapsed={this.state.isSideMenuCollapsed}/>
          <ToastContainer
            autoClose={20000}
            style={{zIndex: 100000}}
          />
          <div
            style={{
              marginLeft: this.state.isSideMenuCollapsed ? '0' : '269px',
              marginTop: '57px'
            }}>
            <Route exact path="/" component={PkgBrowser}/>
            <Route path="/events" component={EventBrowser}/>
          </div>
        </div>
      </HashRouter>
    );
  }
}
