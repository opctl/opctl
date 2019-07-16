import React from 'react'
import ReactDOM from 'react-dom'
import Dir from './Dir'

it('renders without crashing', () => {
  /* arrange */
  const div = document.createElement('div')

  /* act/assert */
  ReactDOM.render(
    <Dir
      name=''
      opRef=''
      param={{ description: '' }}
      value=''
    />,
    div)
})
