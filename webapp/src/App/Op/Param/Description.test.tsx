import React from 'react'
import ReactDOM from 'react-dom'
import Description from './Description'

describe('value not null or empty', () => {
  it('renders without crashing', () => {
    /* arrange */
    const div = document.createElement('div')

    /* act/assert */
    ReactDOM.render(<Description value={'dummyDescription'} />, div)
  })
})
