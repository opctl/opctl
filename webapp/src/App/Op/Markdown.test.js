import React from 'react'
import ReactDOM from 'react-dom'
import Markdown from './Markdown'

describe('value null', () => {
  it('renders without crashing', () => {
    /* arrange */
    const div = document.createElement('div')

    /* act/assert */
    ReactDOM.render(<Markdown value={null} />, div)
  })
})
describe('value empty', () => {
  it('renders without crashing', () => {
    /* arrange */
    const div = document.createElement('div')

    /* act/assert */
    ReactDOM.render(<Markdown value={''} />, div)
  })
})
describe('value not null or empty', () => {
  it('renders without crashing', () => {
    /* arrange */
    const div = document.createElement('div')

    /* act/assert */
    ReactDOM.render(<Markdown value={'dummyMarkdown'} />, div)
  })
})
