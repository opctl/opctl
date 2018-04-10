import React from 'react'
import ReactDOM from 'react-dom'
import Op from './index'

it('renders without crashing', () => {
  /* arrange */
  const div = document.createElement('div')

  /* act/assert */
  ReactDOM.render(
    <Op
      value={{description: ''}}
      opRef={'dummyOpRef'}
    />,
    div
  )
})
