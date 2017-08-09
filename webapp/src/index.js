import React from 'react';
import ReactDOM from 'react-dom';
import PkgBrowser from './PkgBrowser';
import 'bootstrap/dist/css/bootstrap.css';

ReactDOM.render(
    <PkgBrowser location={location}/>, // eslint-disable-line no-restricted-globals
    document.getElementById('root')
);
