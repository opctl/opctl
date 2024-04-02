import React from 'react'
import ReactDOM from 'react-dom'
import _Object from './Object'
jest.mock('@opctl/sdk/src/api/client', () => ({
  dataGet: async () => ({
    json: async () => ''
  })
}))

describe('Object', () => {
  it('renders without crashing', () => {
    /* arrange */
    const div = document.createElement('div')
    const dummyYamlObject = JSON.stringify({})

    /* act/assert */
    ReactDOM.render(
      <_Object
        apiBaseUrl=''
        name=''
        opRef=''
        param={{ description: '' }}
        value={dummyYamlObject}
      />,
      div
      )
  })
})
