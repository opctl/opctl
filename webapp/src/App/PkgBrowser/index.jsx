import React, { Component } from 'react';
import PkgSelector from '../PkgSelector';
import Pkg from '../Pkg';


export default class PkgBrowser extends Component {
  constructor(props) {
    super(props);

    this.state = {};
    this.initialPkgRef = new URLSearchParams(props.location.search).get('pkg') || 'github.com/opspec-pkgs/uuid.v4.generate#1.0.0';
  }

  handleSelect(selection) {
    this.setState(selection);
  }

  render() {
    return (
      <div className='container'>
        <div>
          <PkgSelector initialPkgRef={this.initialPkgRef} onSelect={selection => this.handleSelect(selection)} />
          {this.state.pkg ? <Pkg value={this.state.pkg} pkgRef={this.state.pkgRef} /> : null}
        </div>
      </div>
    )
  }

}
