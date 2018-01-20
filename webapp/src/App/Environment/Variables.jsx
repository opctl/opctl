import React, {Component} from 'react';
import Entry from './Variable';

export default class Variables extends Component {
  state = {
    idCounter: 1,
    variables: []
  };

  handleInvalid = () => {
    this.setState({isVariableAddDisabled: true});
  };

  handleValid = (value) => {
    this.setState(prevState => {
      const entryIndex = prevState.variables.findIndex(entry => entry.id === value.id);
      prevState.variables[entryIndex] = Object.assign(prevState.variables[entryIndex], value);
      return prevState;
    });
    this.setState({isVariableAddDisabled: false});
  };

  addVariable = () => {
    this.setState(prevState => {
      prevState.idCounter++;
      prevState.variables.push({
        id: prevState.idCounter,
        name: "",
        value: "",
      });
      return prevState;
    })
  };

  render() {
    return (
      <div>
        <h2>Variables</h2>
        <div className='d-none d-xl-block'>
          <div className='row'>
            <div className='col-xl-4'>name</div>
            <div className='col-xl-8'>value</div>
          </div>
        </div>
        <ul className='list-group'>
          {this.state.variables.map((entry) =>
            <Entry
              key={entry.id}
              value={entry}
              onValid={value => (this.handleValid(value))}
              onInvalid={() => (this.handleInvalid(entry.id))}
            />
          )}
        </ul>
        <br/>
        <div className='form-group'>
          <button
            type="button"
            className='btn btn-secondary btn-sm'
            onClick={this.addVariable}
            disabled={this.state.isVariableAddDisabled}
          >
            add
          </button>
        </div>
      </div>
    );
  }
}
