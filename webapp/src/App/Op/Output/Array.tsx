import React, { Component } from 'react'
import jsYaml from 'js-yaml'
import Textarea from 'react-textarea-autosize'
import Description from '../Param/Description'
import getApiBaseUrl from '../../../core/getApiBaseUrl'
import { dataGet } from '@opctl/sdk/lib/api/client'

interface Props {
  description?: string
  name: string
  opRef: string
  param
  value
}

export default class Array extends Component<Props> {
  state = {
    value: ''
  };

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
        <Textarea
          className='form-control'
          value={this.state.value || ''}
          id={this.props.name}
          readOnly
        />
      </div>)
  }

  _loadValue () {
    const { value, param } = this.props
    if (value) {
      dataGet(
        getApiBaseUrl(),
        value
      )
        .then(data => data.json())
        .then(value => this.setState({ value: jsYaml.safeDump(value) }))
    } else {
      this.setState({ value: value || param.default ? jsYaml.safeDump(param.default) : '' })
    }
  }
}
