import React from 'react'
import ReactDOM from 'react-dom'
import Number from './Number'

describe('Number', () => {
  it('renders without crashing', () => {
    /* arrange */
    const div = document.createElement('div')

    /* act/assert */
    ReactDOM.render(
      <Number
        name=''
        opRef=''
        param={{ description: '' }}
        value=''
      />,
      div)
  })
})
