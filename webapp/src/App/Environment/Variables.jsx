import React, {Component} from 'react'
import Entry from './Variable'
import contentStore from '../../core/contentStore'
import uuidV4 from 'uuid/v4'
const key = 'environment'

export default class Variables extends Component {
  state =
    {
      variables: contentStore.get({key}) || []
    };

  handleInvalid = () => {
    this.setState({isVariableAddDisabled: true})
  };

  handleValid = value => {
    this.setState(prevState => {
      const nextVariables = [...prevState.variables]
      const entryIndex = prevState.variables.findIndex(entry => entry.id === value.id)
      nextVariables[entryIndex] = value
      contentStore.set({key, value: nextVariables})
      return {
        variables: nextVariables,
        isVariableAddDisabled: false
      }
    })
  };

  addVariable = () => {
    this.setState(prevState => {
      const nextVariables = [
        ...prevState.variables,
        {
          id: uuidV4(),
          name: '',
          value: ''
        }
      ]

      return {
        variables: nextVariables
      }
    })
  };

  render () {
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
        <br />
        <div className='form-group'>
          <button
            type='button'
            className='btn btn-secondary btn-sm'
            onClick={this.addVariable}
            disabled={this.state.isVariableAddDisabled}
          >
            add
          </button>
        </div>
      </div>
    )
  }
}
