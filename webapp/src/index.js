import React from 'react';
import ReactDOM from 'react-dom';
import Pkg from './Pkg';
import 'bootstrap/dist/css/bootstrap.css';

ReactDOM.render(
    <Pkg location={location}/>, // eslint-disable-line no-restricted-globals
    document.getElementById('root')
);
