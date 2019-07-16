import React, { Component } from 'react'
import Description from '../Param/Description'
import ModelParamBoolean from '@opctl/sdk/lib/model/param/boolean'
import { dataGet } from '@opctl/sdk/lib/api/client'

interface Props {
  apiBaseUrl: string
  name: string
  opRef: string
  param: ModelParamBoolean
  value: any
}

interface State {
  value: any
}

export default class Boolean extends Component<Props, State> {
  state: State = {
    value: this.props.value
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
          type='checkbox'
          value={this.state.value || false}
        />
      </div>)
  }

  _loadValue() {
    const { value, param } = this.props
    if (typeof value !== 'undefined') {
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
