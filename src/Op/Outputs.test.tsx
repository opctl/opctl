import React from 'react'
import ReactDOM from 'react-dom'
import Outputs from './Outputs'

it('renders without crashing', () => {
  /* arrange */
  const div = document.createElement('div')

  /* act/assert */
  ReactDOM.render(
    <Outputs
      apiBaseUrl=''
      outputs={{}}
      opRef=''
      scope={{}}
    />,
    div
  )
})
