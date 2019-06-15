import React, { Component } from 'react'
import Description from '../Param/Description'
import { ErrorObject } from 'ajv'

interface Props {
  description: string | null | undefined
  name: string
  onInvalid?: () => any | null | undefined
  onValid: (value: any) => any
  opRef: string
  options: { name: string, value: any }[]
  validate: (value: any) => ErrorObject[]
  value: any
}

interface State {
  value: any
  validationErrs: ErrorObject[]
}

export default class Select extends Component<Props, State> {
  state = {
    value: this.props.value,
    validationErrs: [] as ErrorObject[]
  }

  componentWillMount() {
    this.processValue(this.state.value)
  }

  processValue(
    value: any
  ) {
    const validationErrs = this.props.validate(value) || []
    this.setState(prevState => ({ validationErrs, value }))

    if (validationErrs.length === 0) {
      this.props.onValid(value)
    } else {
      this.props.onInvalid && this.props.onInvalid()
    }
  }

  render() {
    let invalidFeedback
    switch (this.state.validationErrs.length) {
      case 0:
        invalidFeedback = null
        break
      case 1:
        invalidFeedback = this.state.validationErrs[0].message
        break
      default:
        invalidFeedback =
          (<ul>
            {this.state.validationErrs.map(err => <li>{err.message}</li>)}
          </ul>)
    }

    return <div className='form-group'>
      <label className='form-control-label' htmlFor={this.props.name}>{this.props.name}</label>
      <Description value={this.props.description} opRef={this.props.opRef} />
      <select
        className={`form-control ${this.state.validationErrs.length > 0 ? 'is-invalid' : ''}`}
        id={this.props.name}
        value={this.state.value}
        onChange={e => this.processValue(e.target.value)}
      >
        {this.props.options.map((option, index) => (
          <option key={`${this.props.name}_${index}`} value={option.value}>{option.name}</option>))}
      </select>
      <span className='invalid-feedback'>{invalidFeedback}</span>
    </div>
  }
}
