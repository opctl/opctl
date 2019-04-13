import React, { Component } from 'react'
import Description from '../Param/Description'
import ModelParamString from '@opctl/sdk/lib/model/param/string'
import { dataGet } from '@opctl/sdk/lib/api/client'

interface Props {
  apiBaseUrl: string
  name: string
  opRef: string
  param: ModelParamString
  value: any
}

interface State {
  value: any
}

export default class String extends Component<Props, State> {
  state: State = {
    value: ''
  };

  componentDidMount() {
    this._loadValue()
  }

  componentDidUpdate(
    prevProps: Props
  ) {
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
          type={this.props.param.isSecret ? 'password' : 'text'}
          value={this.state.value || ''}
        />
      </div>)
  }

  _loadValue() {
    const { value, param } = this.props
    if (value) {
      dataGet(
        this.props.apiBaseUrl,
        value
      )
        .then(data => data.text())
        .then(value => this.setState({ value }))
    } else {
      this.setState({ value: param.default || '' })
    }
  }
}
