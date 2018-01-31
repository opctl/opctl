import React, { Component } from 'react';
import pkgFetcher from '../core/pkgFetcher';
import { toast } from 'react-toastify';

export default class PkgSelector extends Component {
    constructor(props) {
        super(props);

        this.state = {
            pkgRef: this.props.initialPkgRef || ''
        };
        
        if (this.props.initialPkgRef) {
            this.openPkg(this.state.pkgRef);
        }
    }

    openPkg(pkgRef) {
        pkgFetcher.fetch(pkgRef)
            .then(pkg => this.props.onSelect({ pkgRef, pkg }))
            .catch(error => {
                toast.error(error.message);
            });
    }

    handleSubmit(e) {
        e.preventDefault();

        this.openPkg(this.state.pkgRef);
    }

    render() {
        return (
            <form onSubmit={e => this.handleSubmit(e)} style={{ paddingTop: '25px' }}>
                <div className='form-group'>
                    <span className='input-group input-group-lg'>
                        <input className='form-control'
                            id='pkgRef'
                            type='text'
                            value={this.state.pkgRef}
                            onChange={e => this.setState({ pkgRef: e.target.value })}
                            placeholder="/absolute/path or host/path/git-repo#tag" />
                        <span className='input-group-btn'>
                            <button className='btn btn-primary' id='pkgRef_Submit' type='submit'>select</button>
                        </span>
                    </span>
                </div>
            </form>
        )
    }

}
