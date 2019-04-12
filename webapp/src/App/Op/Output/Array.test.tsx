jest.mock('@opctl/sdk/lib/api/client', () => ({
  dataGet: async () => ({
    json: async () => ''
  })
}))

import React from 'react'
import ReactDOM from 'react-dom'
import Array from './Array'

describe('Array', () => {
  it('renders without crashing', async () => {
    /* arrange */
    const div = document.createElement('div')
    const dummyYamlArray = JSON.stringify([])

    /* act/assert */
    ReactDOM.render(
      <Array
        description=''
        name=''
        opRef=''
        param={{ description: '' }}
        value={dummyYamlArray}
      />,
      div)
  })
})
