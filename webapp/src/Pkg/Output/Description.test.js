import React from 'react';
import ReactDOM from 'react-dom';
import Description from './Description';

it('renders without crashing', () => {
  /* arrange */
  const div = document.createElement('div');

  /* act/assert */
  ReactDOM.render(
    <Description
      pkgRef={'dummyPkgRef'}
      value={'dummyValue'}
    />,
    div);
});
