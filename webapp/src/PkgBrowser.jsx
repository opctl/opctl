import React, {Component} from 'react';
import {Button, Form, FormControl, FormGroup, Grid, InputGroup} from 'react-bootstrap';
import pkgFetcher from './pkgFetcher'

import Pkg from './Pkg'

class PkgBrowser extends Component {
  constructor(props) {
    super(props);

    this.state = {
      pkgRef: 'github.com/opspec-pkgs/git.clean#1.0.0',
    };

    this.openPkg(this.state.pkgRef);

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(e) {
    this.setState({pkgRef: e.target.value});
  }

  openPkg(pkgRef) {
    pkgFetcher.fetch(pkgRef).then(pkg => this.setState({pkgRef, pkg}));
  }

  handleSubmit(e) {
    e.preventDefault();

    this.openPkg(this.state.pkgRef);
  }

  render() {
    const form = (
      <Form onSubmit={this.handleSubmit} style={{'padding-top': '25px'}}>
        <FormGroup controlId="pkgRef">
          <InputGroup bsSize="lg" >
            <FormControl type="text" value={this.state.pkgRef} onChange={this.handleChange}
                         placeholder="/absolute/path or host/path/git-repo#tag"/>
            <InputGroup.Button>
              <Button type="submit">open</Button>
            </InputGroup.Button>
          </InputGroup>
        </FormGroup>
      </Form>
    );

    if (!this.state.pkg) {
      return (
        <Grid>
          <div>
            {form}
          </div>
        </Grid>
      )
    }
    return (
      <Grid>
        <div>
          {form}
          <Pkg
            description={this.state.pkg.description}
            inputs={this.state.pkg.inputs}
            name={this.state.pkg.name}
            outputs={this.state.pkg.outputs}
            version={this.state.pkg.version}
          />
        </div>
      </Grid>
    );
  }

}

export default PkgBrowser;
