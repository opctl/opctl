import React, {Component} from 'react';
import jsYaml from 'js-yaml';
import Textarea from 'react-textarea-autosize';

export default class ArrayInput extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: props.array.default,
    };

    this.handleArgChange = this.handleArgChange.bind(this);
  }

  handleArgChange(e) {
    const value = jsYaml.safeLoad(e.target.value);
    this.props.onArgChange({array: value});
    this.setState({value});
  };

  render() {
    return (
      <div className='form-group'>
        <label className='form-control-label' htmlFor={this.props.name}>{this.props.name}</label>
        <p className='custom-control-description'>{this.props.array.description}</p>
        <Textarea
          className='form-control'
          id={this.props.name}
          value={jsYaml.safeDump(this.state.value)}
          onChange={this.handleArgChange}
        />
      </div>
    );
  }
}
