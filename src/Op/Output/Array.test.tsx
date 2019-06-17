import React from 'react'
import ReactDOM from 'react-dom'
import Array from './Array'
jest.mock('@opctl/sdk/lib/api/client', () => ({
  dataGet: async () => ({
    json: async () => ''
  })
}))

describe('Array', () => {
  it('renders without crashing', () => {
    /* arrange */
    const div = document.createElement('div')
    const dummyYamlArray = JSON.stringify([])

    /* act/assert */
    ReactDOM.render(
      <Array
        apiBaseUrl=''
        description=''
        name=''
        opRef=''
        param={{ description: '' }}
        value={dummyYamlArray}
      />,
      div)
  })
})
