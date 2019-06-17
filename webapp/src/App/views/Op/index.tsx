import React, { Component } from 'react'
import OpSelector from '../../OpSelector'
import Op from '../../Op/index'

interface Props {
  initialOpRef: string
  location
}

interface State {
  op?
  opRef?: string
}

export default class OpView extends Component<Props, State> {
  initialOpRef = new window.URLSearchParams(this.props.location.search).get('op') || 'github.com/opspec-pkgs/uuid.v4.generate#1.1.0'

  state: State = {}

  handleSelect(selection) {
    this.setState(selection)
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
