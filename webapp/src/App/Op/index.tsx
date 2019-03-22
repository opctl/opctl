import React, { Component } from 'react'
import Markdown from './Markdown'
import Icon from './Icon'
import Inputs from './Inputs'
import Outputs from './Outputs'
import EventStream from '../EventStream'
import opspecNodeApiClient from '../../core/clients/opspecNodeApi'
import { toast } from 'react-toastify'

export default class Op extends Component<any,any> {
  constructor (props) {
    super(props)
  }
  args

  state = {
    isStartable: (this.props.value.inputs || []).length === 0,
    isKillable: false,
    opId: '',
    outputs: {}
  }

  handleInvalid = (name) => {
    this.args = this.args || {}
    delete this.args[name]
    this.setState({ isStartable: false })
  };

  handleValid = (name, value) => {
    this.args = this.args || {}
    this.args[name] = value
    this.setState({ isStartable: Object.keys(this.props.value.inputs).length === Object.keys(this.args).length })
  };

  kill = () => {
    opspecNodeApiClient.op_kill({
      opId: this.state.opId
    })
      .then(() => this.setState({ isKillable: false }))
      .catch(error => toast.error(error.message))
  };

  start = () => {
    const args = Object.entries(this.props.value.inputs || [])
      .reduce((args, [name, param]: [string,any]) => {
        if (param.array) args[name] = { array: this.args[name] }
        if (param.boolean) args[name] = { boolean: this.args[name] }
        if (param.dir) args[name] = { dir: this.args[name] }
        if (param.file) args[name] = { file: this.args[name] }
        if (param.number) args[name] = { number: this.args[name] }
        if (param.object) args[name] = { object: this.args[name] }
        if (param.socket) args[name] = { socket: this.args[name] }
        if (param.string) args[name] = { string: this.args[name] }
        return args
      }, {})

    opspecNodeApiClient.op_start({
      args,
      op: {
        ref: this.props.opRef
      }
    })
      .then(opId => {
        this.setState({ opId, isKillable: true })

        opspecNodeApiClient.event_stream_get({
          filter: {
            roots: [opId]
          },
          onEvent: event => {
            if (event.opEnded && event.opEnded.opId === opId) {
              this.setState({
                isKillable: false,
                outputs: event.opEnded.outputs
              })
            }
          }
        })
      })
      .catch(error => toast.error(error.message))
  };

  render () {
    return (
      <div style={{ height: '100%' }}>
        <form onSubmit={e => {
          e.preventDefault()
        }}>
          <Icon
            value={this.props.value.icon}
            opRef={this.props.opRef}
          />
          <Markdown
            value={this.props.value.description}
            opRef={this.props.opRef}
          />
          <Inputs
            value={this.props.value.inputs}
            onInvalid={this.handleInvalid}
            onValid={this.handleValid}
            opRef={this.props.opRef}
            values={{}}
          />
          <div className='form-group'>
            {
              this.state.isKillable
                ? <button
                  className='col-12 btn btn-primary btn-lg'
                  id='opKill'
                  onClick={this.kill}
                >kill</button>
                : <button
                  className='col-12 btn btn-primary btn-lg'
                  id='opStart'
                  onClick={this.start}
                  disabled={!this.state.isStartable}
                >start</button>
            }
          </div>
          <Outputs
            params={this.props.value.outputs}
            opRef={this.props.opRef}
            values={this.state.outputs || {}}
          />
        </form>
        <br />
        {
          this.state.opId
            ? <div style={{ height: '100%' }}>
              <h2>Event Stream:</h2>
              <EventStream key={this.state.opId} filter={{ roots: [this.state.opId] }} />
            </div>
            : null
        }
      </div>
    )
  }
}
