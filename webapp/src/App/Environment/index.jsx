import React, {Component} from 'react'
import Variables from './Variables'
import './index.css'

export default class Environment extends Component {
  render () {
    return (
      <div className='container'>
        <Variables />
      </div>
    )
  }
}
