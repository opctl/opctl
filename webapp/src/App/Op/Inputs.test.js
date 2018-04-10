import React from 'react'
import ReactDOM from 'react-dom'
import Inputs from './Inputs'

it('renders without crashing', () => {
  /* arrange */
  const div = document.createElement('div')

  /* act/assert */
  ReactDOM.render(<Inputs value={{}} />, div)
})
