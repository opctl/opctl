import React, { Component } from 'react'
import opFetcher from '../core/opFetcher'
import { toast } from 'react-toastify'

export default class OpSelector extends Component {
  constructor (props) {
    super(props)

    this.state = {
      opRef: this.props.initialOpRef || ''
    }

    if (this.props.initialOpRef) {
      this.openOp(this.state.opRef)
    }
  }

  openOp (opRef) {
    opFetcher.fetch(opRef)
      .then(op => this.props.onSelect({ opRef, op }))
      .catch(error => {
        toast.error(error.message)
      })
  }

  handleSubmit (e) {
    e.preventDefault()

    this.openOp(this.state.opRef)
  }

  render () {
    return (
      <form onSubmit={e => this.handleSubmit(e)} style={{ paddingTop: '25px' }}>
        <div className='form-group'>
          <span className='input-group input-group-lg'>
            <input className='form-control'
              id='opRef'
              type='text'
              value={this.state.opRef}
              onChange={e => this.setState({ opRef: e.target.value })}
              placeholder='/absolute/path or host/path/git-repo#tag' />
            <span className='input-group-btn'>
              <button className='btn btn-primary' id='opRef_Submit' type='submit'>select</button>
            </span>
          </span>
        </div>
      </form>
    )
  }
}
