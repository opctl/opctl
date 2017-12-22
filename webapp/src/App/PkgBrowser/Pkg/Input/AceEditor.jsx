import React, {Component} from 'react';
import AceEditor from 'react-ace';
import Description from './Description';
import 'brace';
import 'brace/mode/yaml';
import 'brace/theme/github';

export default class TextArea extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: props.value || '',
      validationErrs: [],
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
        <Description value={this.props.description}/>
        <AceEditor
          className={`form-control ${this.state.validationErrs.length > 0 ? 'is-invalid' : ''}`}
          mode="yaml"
          onChange={value => this.processValue(value)}
          theme="github"
          value={this.state.value}
          width="100%"
          tabSize={2}
          minLines={1}
          maxLines={30}
        />
        <span className='invalid-feedback'>{invalidFeedback}</span>
      </div>
    );
  }
}
