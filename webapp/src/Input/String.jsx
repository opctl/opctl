import React, {Component} from 'react';

export default class StringInput extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: props.string.default,
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
        <p className='custom-control-description'>{this.props.string.description}</p>
        <input
          className='form-control'
          id={this.props.name}
          type={this.props.string.isSecret? 'password': 'text'}
          value={this.state.value}
          onChange={this.handleChange}
        />
      </div>
    );
  }
}
