import React, { Component } from 'react';
import { Route, Router, Switch } from 'react-router-dom';
import { createBrowserHistory as createHistory } from 'history';
import { Home } from './pages/Home/Home.jsx';
import { AppFrame } from './components/AppFrame/AppFrame.jsx';

const history = createHistory();

const routes = {
    Home: '/',
};

class Routes extends Component {
    render() {
        return (
            <Router history={history} >
                <AppFrame>
                    <Switch>
                        <Route exact path={routes.Home} component={Home} />
                    </Switch>
                </AppFrame>
            </Router>
        );
    }
}

export { history, Routes, routes };
