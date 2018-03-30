import React, { Component } from 'react';
import OpSelector from '../../OpSelector';
import Op from '../../Op/index';


export default class OpView extends Component {
  constructor(props) {
    super(props);

    this.state = {};
    this.initialOpRef = new URLSearchParams(props.location.search).get('op') || 'github.com/opspec-pkgs/uuid.v4.generate#1.1.0';
  }

  handleSelect(selection) {
    this.setState(selection);
  }

  render() {
    return (
      <div className='container'>
        <div>
          <OpSelector initialOpRef={this.initialOpRef} onSelect={selection => this.handleSelect(selection)} />
          {this.state.op ? <Op value={this.state.op} opRef={this.state.opRef} /> : null}
        </div>
      </div>
    )
  }

}
