import React from 'react'
import ReactDOM from 'react-dom'
import Array from './Array'

describe('Array', () => {
  it('renders without crashing', () => {
    /* arrange */
    const div = document.createElement('div')
    const dummyYamlArray = JSON.stringify([])

    /* act/assert */
    ReactDOM.render(
      <Array
        param={{ description: '' }}
        value={dummyYamlArray}
      />,
      div)
  })
})
