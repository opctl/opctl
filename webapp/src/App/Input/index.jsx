import React, { Component } from 'react'
import Description from '../Op/Param/Description'

export default class Input extends Component {
  state = {
    value: this.props.value || '',
    validationErrs: []
  };

  componentWillMount () {
    this.processValue(this.state.value)
  }

  processValue (value) {
    const validationErrs = this.props.validate(value) || []
    this.setState({ validationErrs, value })

    if (validationErrs.length === 0) {
      this.props.onValid(value)
    } else {
      this.props.onInvalid()
    }
  }

  render () {
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
        {
          this.props.name
            ? <label className='form-control-label' htmlFor={this.props.name}>{this.props.name}</label>
            : null
        }
        {
          this.props.description && this.props.opRef
            ? <Description value={this.props.description} opRef={this.props.opRef} />
            : null
        }
        <input
          className={`form-control ${this.state.validationErrs.length > 0 ? 'is-invalid' : ''}`}
          id={this.props.id || this.props.name}
          type={this.props.type}
          value={this.state.value}
          // if checkbox, process target.checked; otherwise target.value
          onChange={e => this.processValue(e.target[this.props.type === 'checkbox' ? 'checked' : 'value'])}
          checked={this.state.value}
        />
        <span className='invalid-feedback'>{invalidFeedback}</span>
      </div>
    )
  }
}
