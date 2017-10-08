import React, {Component} from 'react';
import PkgRef from './PkgRef';
import Inputs from './Inputs';
import Outputs from './Outputs';
import EventBrowser from './EventBrowser';
import opSpecNodeApiClient from './opspecNodeApiClient'
import {toast} from 'react-toastify';

export default class Pkg extends Component {
  constructor(props) {
    super(props);

    this.state = {};
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleArgChange = (name, value) => {
    this.args = this.args || {};
    this.args[name] = value;
  };

  handleSubmit = (e) => {
    e.preventDefault();

    const req = {
      args: this.args || {},
      pkg: {
        ref: this.props.pkgRef,
      }
    };

    opSpecNodeApiClient.startOp(req)
      .then(opId => (this.setState({opId})))
      .catch(error => {
        toast.error(error.message);
      });
  };

  render() {
    return (
      <div>
        <form onSubmit={this.handleSubmit}>
          <h1><PkgRef name={this.props.value.name} version={this.props.value.version}/></h1>
          <p className="lead">{this.props.value.description}</p>
          <Inputs value={this.props.value.inputs} onArgChange={this.handleArgChange}/>
          <Outputs value={this.props.value.outputs}/>
          <input className='btn btn-primary btn-lg' id='startOp_Submit' type='submit' value='start'/>
        </form>
        <br/>
        {
          this.state.opId ?
            <div>
              <h2>Events Stream:</h2>
              <EventBrowser key={this.state.opId} filter={{root: this.state.opId}}/>
            </div>
            : null
        }
      </div>
    );
  }
}
