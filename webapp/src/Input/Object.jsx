import React, {Component} from 'react';
import jsYaml from 'js-yaml';
import Textarea from 'react-textarea-autosize';

export default class ObjectInput extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: props.object.default,
    };

    this.handleChange = this.handleChange.bind(this);
  }

  handleChange(e) {
    const value = jsYaml.safeLoad(e.target.value);
    this.props.onChange({object: value});
    this.setState({value});
  };

  render() {
    return (
      <div className='form-group'>
        <label className='form-control-label' htmlFor={this.props.name}>{this.props.name}</label>
        <p className='custom-control-description'>{this.props.object.description}</p>
        <Textarea
          className='form-control'
          id={this.props.name}
          value={jsYaml.safeDump(this.state.value)}
          onChange={this.handleChange}
        />
      </div>
    );
  }
}
