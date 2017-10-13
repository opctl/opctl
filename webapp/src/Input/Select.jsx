import React, {Component} from 'react';

export default class Select extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: props.value
    };
  }

  componentWillMount() {
    this.processValue(this.state.value);
  }

  processValue(value) {
    const validationErrs = this.props.validate(value) || [];
    this.setState(prevState => ({validationErrs, value}));

    if (validationErrs.length === 0) {
      this.props.onValid(value);
    } else {
      this.props.onInvalid();
    }

  }

  render() {
    return <div className='form-group'>
      <label className='form-control-label' htmlFor={this.props.name}>{this.props.name}</label>
      <p className='custom-control-description'>{this.props.description}</p>
      <select
        className={`form-control`}
        id={this.props.name}
        value={this.state.value}
        onChange={e => this.processValue(e.target.value)}
      >
        {this.props.options.map((option, index) => (
          <option key={`${this.props.name}_${index}`} value={option.value}>{option.name}</option>))}
      </select>
    </div>;
  }
}
