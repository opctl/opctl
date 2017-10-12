import React, {Component} from 'react';

export default class Select extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: props.value
    };
  }

  handleChange(e) {
    const value = e.target.value;
    this.props.onChange(value);
    this.setState({value});
  }

  render() {
    return <div className='form-group'>
      <label className='form-control-label' htmlFor={this.props.name}>{this.props.name}</label>
      <p className='custom-control-description'>{this.props.description}</p>
      <select
        className={`form-control`}
        id={this.props.name}
        value={this.state.value}
        onChange={e => this.handleChange(e)}
      >
        {this.props.options.map((option, index) => (
          <option key={`${this.props.name}_${index}`} value={option.value}>{option.name}</option>))}
      </select>
    </div>;
  }
}
