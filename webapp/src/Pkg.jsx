import React, {Component} from 'react';
import PkgRef from './PkgRef'
import Inputs from './Inputs'
import Outputs from './Outputs'
import pkgFetcher from './pkgFetcher'
import queryString from 'query-string';

class Pkg extends Component {
    constructor(props) {
        super(props);
        this.state = {
            inputs: {}
        };
    }

    componentDidMount() {
        const pkgRef = this.getPkgRef();
        if (pkgRef) {
            pkgFetcher.fetch(pkgRef).then(pkg => this.setState(pkg));
        }
    }

    getPkgRef() {
        return queryString.parse(this.props.location.search).pkgRef;
    }

    render() {
        return (
            <div>
                <h1><PkgRef name={this.state.name} version={this.state.version}/></h1>
                <p className="lead">{this.state.description}</p>
                <Inputs value={this.state.inputs}/>
                <Outputs value={this.state.outputs}/>
            </div>
        );
    }
}

export default Pkg;
