import React, {Component} from 'react';

export default class Input extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: props.value || '',
      validationErrs: [],
    };
  }

  handleChange(e) {
    const value = e.target.value;
    this.props.onChange(value);
    this.setState({value});
  };

  handleBlur() {
    this.setState(prevState => ({validationErrs: this.props.validate(prevState.value) || []}));
  }

  render() {
    let invalidFeedback;
    switch (this.state.validationErrs.length) {
      case 0:
        invalidFeedback = null;
        break;
      case 1:
        invalidFeedback = this.state.validationErrs[0].message;
        break;
      default:
        invalidFeedback =
          (<ul>
            {this.state.validationErrs.map(err => <li>{err.message}</li>)}
          </ul>);
    }

    return (
      <div className='form-group'>
        <label className='form-control-label' htmlFor={this.props.name}>{this.props.name}</label>
        <p className='custom-control-description'>{this.props.description}</p>
        <input
          className={`form-control ${this.state.validationErrs.length > 0 ? 'is-invalid' : ''}`}
          id={this.props.name}
          type={this.props.type}
          value={this.state.value}
          onChange={e => this.handleChange(e)}
          onBlur={() => this.handleBlur()}
        />
        <span className='invalid-feedback'>{invalidFeedback}</span>
      </div>
    );
  }
}
