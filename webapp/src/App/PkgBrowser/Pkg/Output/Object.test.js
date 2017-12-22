import React from 'react';
import ReactDOM from 'react-dom';
import Object from './Object';

it('renders without crashing', () => {
  /* arrange */
  const div = document.createElement('div');
  const dummyYamlObject = JSON.stringify({});

  /* act/assert */
  ReactDOM.render(
    <Object
      param={{description: ''}}
      value={dummyYamlObject}
    />,
    div);
});
