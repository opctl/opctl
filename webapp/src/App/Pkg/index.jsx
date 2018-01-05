import React, { Component } from 'react';
import Markdown from './Markdown';
import Icon from './Icon';
import Inputs from './Inputs';
import Outputs from './Outputs';
import EventBrowser from '../EventBrowser';
import opspecNodeApiClient from '../../utils/clients/opspecNodeApi';
import { toast } from 'react-toastify';

export default class Pkg extends Component {
  constructor(props) {
    super(props);

    this.state = {
      isStartable: (props.value.inputs || []).length === 0,
      isKillable: false,
      outputs: {},
    };
  }

  handleInvalid = (name) => {
    this.args = this.args || {};
    delete this.args[name];
    this.setState({ isStartable: false });
  };

  handleValid = (name, value) => {
    this.args = this.args || {};
    this.args[name] = value;
    this.setState({ isStartable: Object.keys(this.props.value.inputs).length === Object.keys(this.args).length });
  };

  kill = () => {
    opspecNodeApiClient.op_kill({
      opId: this.state.opId
    })
      .then(() => this.setState({ isKillable: false }))
      .catch(error => {
        toast.error(error.message);
      });
  };

  start = () => {
    opspecNodeApiClient.op_start({
      args: this.args || {},
      pkg: {
        ref: this.props.pkgRef,
      }
    })
      .then(opId => {
        this.setState({ opId, isKillable: true });

        opspecNodeApiClient.event_stream_get({
          filter: {
            roots: [opId],
          },
          onEvent: event => {
            if (event.opEnded && event.opEnded.opId === opId) {
              this.setState({
                isKillable: false,
                outputs: event.opEnded.outputs,
              });
            }
          },
        })
      })
      .catch(error => {
        toast.error(error.message);
      });
  };

  render() {
    return (
      <div style={{height:'100%'}}>
        <form onSubmit={e => {
          e.preventDefault()
        }}>
          <Icon
            value={this.props.value.icon}
            pkgRef={this.props.pkgRef}
          />
          <Markdown
            value={this.props.value.description}
            pkgRef={this.props.pkgRef}
          />
          <Inputs
            value={this.props.value.inputs}
            onInvalid={this.handleInvalid}
            onValid={this.handleValid}
            pkgRef={this.props.pkgRef}
            values={{}}
          />
          <div className='form-group'>
            {
              this.state.isKillable ?
                <button
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
            pkgRef={this.props.pkgRef}
            values={this.state.outputs || {}}
          />
        </form>
        <br />
        {
          this.state.opId ?
            <div style={{height: '100%'}}>
              <h2>Events Stream:</h2>
              <EventBrowser key={this.state.opId} filter={{ roots: [this.state.opId] }} />
            </div>
            : null
        }
      </div>
    );
  }
}
