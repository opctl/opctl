import React, { Component } from 'react'
import Description from '../Param/Description'
import opspecNodeApiClient from '../../../core/clients/opspecNodeApi'

export default class Boolean extends Component {
  state = {};

  componentDidMount () {
    this._loadValue()
  }

  componentDidUpdate (prevProps) {
    if (this.props.value !== prevProps.value) {
      this._loadValue()
    }
  }

  render () {
    return (
      <div className='form-group'>
        <label className='form-control-label' htmlFor={this.props.name}>{this.props.name}</label>
        <Description value={this.props.param.description} opRef={this.props.opRef} />
        <input
          className='form-control'
          id={this.props.name}
          readOnly
          type='checkbox'
          value={this.state.value || false}
        />
      </div>)
  }

  _loadValue () {
    const { value, param } = this.props
    if (typeof value !== 'undefined') {
      opspecNodeApiClient.data_get({
        dataRef: value
      })
        .then(data => data.text())
        .then(value => this.setState({ value }))
    } else {
      this.setState({ value: param.default || '' })
    }
  }
}
