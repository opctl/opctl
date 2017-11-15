import React from 'react';
import ReactDOM from 'react-dom';
import Number from './Number';

it('renders without crashing', () => {
  /* arrange */
  const div = document.createElement('div');

  /* act/assert */
  ReactDOM.render(
    <Number
      param={{description: ''}}
    />,
    div);
});
