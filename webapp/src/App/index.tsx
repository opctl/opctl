import React, { Component } from 'react'
import Operations from './views/Operations'
import Events from './views/Events'
import Environment from './views/Environment'
import { Route, HashRouter } from 'react-router-dom'
import TopMenu from './TopMenu'
import SideMenu from './SideMenu'
import { ToastContainer } from 'react-toastify'
import OpView from './views/Op'
import './bootstrap.scss'

// shared vendor css
import 'react-toastify/dist/ReactToastify.min.css'
import 'react-virtualized/styles.css'
import 'react-resizable/css/styles.css'

export default class Page extends Component<any,any> {
  constructor (props) {
    super(props)
  }

  state={
    isSideMenuCollapsed: false
  }

  handleSideMenuToggled () {
    this.setState(state => {
      return {
        isSideMenuCollapsed: !state.isSideMenuCollapsed
      }
    })
  }

  render () {
    return (
      <HashRouter>
        <div style={{ height: '100vh' }}>
          <TopMenu onCollapseToggled={() => this.handleSideMenuToggled()} />
          <SideMenu isCollapsed={this.state.isSideMenuCollapsed} />
          <ToastContainer
            autoClose={20000}
            style={{ zIndex: 100000 }}
          />
          <div
            style={{
              height: 'calc(100vh - 58px)',
              marginLeft: this.state.isSideMenuCollapsed ? '0' : '269px'
            }}>
            <Route exact path='/' component={Operations} />
            <Route path='/events' component={Events} />
            <Route path='/op' component={OpView} />
            <Route path='/environment' component={Environment} />
          </div>
        </div>
      </HashRouter>
    )
  }
}
