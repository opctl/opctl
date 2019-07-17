import React, { Component } from 'react'
import Description from '../Param/Description'
import ModelParamNumber from '@opctl/sdk/src/model/param/number'
import { dataGet } from '@opctl/sdk/src/api/client'

interface Props {
  apiBaseUrl: string
  name: string
  opRef: string
  param: ModelParamNumber
  value: any
}

interface State {
  value: any
}

export default class Number extends Component<Props, State> {
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
          type={this.props.param.isSecret ? 'password' : 'number'}
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
