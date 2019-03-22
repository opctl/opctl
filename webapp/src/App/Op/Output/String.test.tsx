import React from 'react'
import ReactDOM from 'react-dom'
import String from './String'

describe('String', () => {
  it('renders without crashing', () => {
    /* arrange */
    const div = document.createElement('div')

    /* act/assert */
    ReactDOM.render(
      <String
        name=''
        opRef=''
        param={{ description: '' }}
        value=''
      />,
      div)
  })
})
