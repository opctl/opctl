import React, { Fragment, PureComponent } from 'react'
import Header from './Header'
import { HotKeys } from 'react-hotkeys'
import Inputs from './Inputs'
import EventStream from '../EventStream'
import ReactModal from 'react-modal'
import EventStore from '../eventStore'

interface Props {
  args: { [name: string]: any }
  eventStore: EventStore
  isKillable: boolean
  isStartable: boolean
  name?: string | null | undefined
  onConfigure: (cfg: any) => any
  onDelete?: () => any | null | undefined
  onKill: () => any
  onStart: () => any
  op: any
  opId: string
  opRef: string
}

interface State {
  isConfigurationVisible: boolean
  isFullScreen: boolean
  isKillable: boolean
  isSelected: boolean
  name: string
  outputs?: { [name: string]: any }
}

export default class Op extends PureComponent<Props, State> {
  args = this.props.args;
  eventStreamCloser?: any

  state: State = {
    isConfigurationVisible: false,
    isFullScreen: false,
    isKillable: false,
    isSelected: false,
    name: this.props.name || this.props.opRef
  };

  ensureEventStreamClosed = () => {
    this.eventStreamCloser && this.eventStreamCloser()
  };

  toggleConfigurationModal = () => {
    this.setState(prevState => ({ isConfigurationVisible: !prevState.isConfigurationVisible }))
  }

  handleInvalidArg = (name: string) => {
    delete this.args[name]

    this.props.onConfigure({ args: this.args })
  };

  handleValidArg = (
    name: string,
    value: any
  ) => {
    this.args[name] = value

    this.props.onConfigure({ args: this.args })
  };

  handleNameChanged = (
    name: string
  ) => {
    this.setState({ name })
    this.props.onConfigure({ name })
  };

  processEventStream = ({ opId }: { opId: string }) => {
    this.eventStreamCloser = this.props.eventStore.getStream(
      {
        roots: [opId]
      },
      event => {
        if (event.opStarted) {
          this.setState({ isKillable: true })
        }
        if (event.opEnded && event.opEnded.opId === opId) {
          this.setState({
            isKillable: false,
            outputs: event.opEnded.outputs
          })
        }
      })
  };

  componentWillReceiveProps(nextProps: Props) {
    if (nextProps.opId !== this.props.opId) {
      this.ensureEventStreamClosed()
      this.processEventStream({ opId: nextProps.opId })
    }
  }

  componentWillMount() {
    if (this.props.opId) {
      this.processEventStream({ opId: this.props.opId })
    }
  }

  componentWillUnmount() {
    this.ensureEventStreamClosed()
  }

  render() {
    const component = (
      <HotKeys
        keyMap={{
          del: 'del',
          kill: 'ctrl+c',
          start: 'enter'
        }}
        style={{
          height: '100%'
        }}
        handlers={{
          ...this.props.onDelete && {
            del: this.props.onDelete
          },
          kill: this.props.onKill,
          start: this.props.onStart
        }}
      >
        <Header
          isFullScreen={this.state.isFullScreen}
          isKillable={this.state.isKillable}
          isStartable={this.props.isStartable}
          onToggleFullScreen={() => this.setState(prevState => ({ isFullScreen: !prevState.isFullScreen }))}
          onStart={this.props.onStart}
          onKill={this.props.onKill}
          onDelete={this.props.onDelete}
          name={this.state.name}
        />
        <ul
          className='nav nav-tabs'
          style={{
            height: '37px'
          }}
        >
          <li className='nav-item'>
            <span
              style={{
                cursor: 'pointer'
              }}
              onClick={() => this.setState({ isConfigurationVisible: true })}
              className={`nav-link ${this.state.isConfigurationVisible && 'active'}`}
            >
              Configuration
          </span>
          </li>
          <li className='nav-item'>
            <span
              style={{
                cursor: 'pointer'
              }}
              onClick={() => this.setState({ isConfigurationVisible: false })}
              className={`nav-link ${!this.state.isConfigurationVisible && 'active'}`}
            >
              Logs
          </span>
          </li>
        </ul>
        <div
          style={{
            height: 'calc(100% - 74px)'
          }}
        >
          {
            this.state.isConfigurationVisible
              ? <Inputs
                inputs={this.props.op.inputs}
                onInvalid={this.handleInvalidArg}
                onValid={this.handleValidArg}
                opRef={this.props.opRef}
                scope={this.args}
              />
              : <EventStream
                eventStore={this.props.eventStore}
                key={this.props.opId}
                filter={{ roots: [this.props.opId] }}
              />
          }
        </div>
      </HotKeys>
    )

    if (this.state.isFullScreen) {
      // wrap in modal for fullscreen
      return <ReactModal
        isOpen={true}
        style={{
          overlay: {
            zIndex: 1000
          },
          content: {
            height: '100vh',
            left: 0,
            padding: 0,
            top: 0,
            width: '100vw'
          }
        }}
      >
        {
          component
        }
      </ReactModal>
    }
    return component
  }
}
