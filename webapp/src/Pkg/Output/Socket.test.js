import React from 'react';
import ReactDOM from 'react-dom';
import Socket from './Socket';

it('renders without crashing', () => {
  /* arrange */
  const div = document.createElement('div');

  /* act/assert */
  ReactDOM.render(<Socket param={{}}/>, div);
});
