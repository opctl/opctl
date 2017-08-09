import React from 'react';
import ReactDOM from 'react-dom';
import Pkg from './Pkg';

it('renders without crashing', () => {
    const div = document.createElement('div');
    ReactDOM.render(<Pkg />, div);
});
