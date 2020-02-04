import React, { Component } from 'react'
import Description from '../Param/Description'
import getApiBaseUrl from '../../../core/getApiBaseUrl'
import { dataGet } from '@opctl/sdk/lib/api/client'

interface Props {
  name: string
  opRef: string
  param
  value
}

interface State {
  value?
}

export default class Number extends Component<Props, State> {
  state: State = {};

  componentDidMount() {
    this._loadValue()
  }

  componentDidUpdate(prevProps) {
    if (this.props.value !== prevProps.value) {
      this._loadValue()
    }
  }

  render() {
    return (
      <div className='form-group'>
        <label className='form-control-label' htmlFor={this.props.name}>{this.props.name}</label>
        <Description value={this.props.param.description} opRef={this.props.opRef} />
        <input
          className='form-control'
          id={this.props.name}
          readOnly
          type={this.props.param.isSecret ? 'password' : 'number'}
          value={this.state.value || ''}
        />
      </div>)
  }

  _loadValue() {
    const { value, param } = this.props
    if (value) {
      dataGet(
        getApiBaseUrl(),
        value
      )
        .then(data => data.text())
        .then(value => this.setState({ value }))
    } else {
      this.setState({ value: param.default || '' })
    }
  }
}
