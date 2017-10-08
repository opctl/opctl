import React, {Component} from 'react';

export default class NumberInput extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: props.number.default,
    };

    this.handleChange = this.handleChange.bind(this);
  }

  handleChange(e) {
    this.setState({value: e.target.value});
  };

  render() {
    return (
      <div className='form-group'>
        <label className='form-control-label' htmlFor={this.props.name}>{this.props.name}</label>
        <p className='custom-control-description'>{this.props.number.description}</p>
        <input
          className='form-control'
          id={this.props.name}
          type={this.props.number.isSecret? 'password': 'number'}
          value={this.state.value}
          onChange={this.handleChange}
        />
      </div>
    );
  }
}
