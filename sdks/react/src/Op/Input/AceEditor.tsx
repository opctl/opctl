import React, { Component } from 'react'
import AceEditor from 'react-ace'
import Description from '../Param/Description'
import 'brace'
import 'brace/mode/yaml'
import 'brace/theme/github'
import { ErrorObject } from 'ajv'

interface Props {
  description?: string | null | undefined
  name: string
  onInvalid?: () => any | null | undefined
  onValid: (value: any) => any
  opRef: string
  validate: (value: any) => ErrorObject[]
  value: any
}

interface State {
  value: any
  validationErrs: ErrorObject[]
}

export default class TextArea extends Component<Props, State> {
  state = {
    value: this.props.value || '',
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

    return (
      <div className='form-group'>
        <label className='form-control-label' htmlFor={this.props.name}>{this.props.name}</label>
        <Description
          opRef={this.props.opRef}
          value={this.props.description}
        />
        <AceEditor
          className={`form-control ${this.state.validationErrs.length > 0 ? 'is-invalid' : ''}`}
          mode='yaml'
          onChange={value => this.processValue(value)}
          theme='github'
          value={this.state.value}
          width='100%'
          tabSize={2}
          minLines={1}
          maxLines={30}
        />
        <span className='invalid-feedback'>{invalidFeedback}</span>
      </div>
    )
  }
}
