jest.mock('@opctl/sdk/lib/api/client', () => ({
  dataGet: async () => ({
    json: async () => ''
  })
}))

import React from 'react'
import ReactDOM from 'react-dom'
import Object from './Object'

describe('Object', () => {
  it('renders without crashing', async () => {
    /* arrange */
    const div = document.createElement('div')
    const dummyYamlObject = JSON.stringify({})

    /* act/assert */
    ReactDOM.render(
      <Object
        name=''
        opRef=''
        param={{ description: '' }}
        value={dummyYamlObject}
      />,
      div)
  })
})
