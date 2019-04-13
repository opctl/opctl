import React, { Fragment, PureComponent } from 'react'
import Header from './Header'
import Inputs from './Inputs'
import EventStream from '../EventStream'
import { Modal, ModalBody } from 'reactstrap'
import ReactModal from 'react-modal'
import EventStore from '../eventStore'
import { css } from 'emotion'

interface Props {
  args: { [name: string]: any }
  eventStore: EventStore
  isKillable: boolean
  isStartable: boolean
  name?: string | null | undefined
  onConfigure: (cfg: any) => any
  onDelete: () => any
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
      <Fragment>
        <Header
          isFullScreen={this.state.isFullScreen}
          isKillable={this.state.isKillable}
          isStartable={this.props.isStartable}
          onToggleFullScreen={() => this.setState(prevState => ({isFullScreen: !prevState.isFullScreen}))}
          onConfigure={this.toggleConfigurationModal}
          onStart={this.props.onStart}
          onKill={this.props.onKill}
          onDelete={this.props.onDelete}
          name={this.state.name}
        />
        <Modal
          size='lg'
          isOpen={this.state.isConfigurationVisible}
          toggle={this.toggleConfigurationModal}
        >
          <ModalBody>
            <h2>Display</h2>
            <div
              className='form-group'
            >
              <label
                className='form-control-label'
              >
                name
                      <input
                  className='form-control'
                  name='name'
                  value={this.state.name}
                  onChange={({ target }) => this.handleNameChanged(target.value)}
                />
              </label>
            </div>
            <Inputs
              inputs={this.props.op.inputs}
              onInvalid={this.handleInvalidArg}
              onValid={this.handleValidArg}
              opRef={this.props.opRef}
              scope={this.args}
            />
          </ModalBody>
        </Modal>
        {
          this.props.opId
            ? <div className={css({
              height: 'calc(100% - 37px)'
            })}>
              <EventStream
                eventStore={this.props.eventStore}
                key={this.props.opId}
                filter={{ roots: [this.props.opId] }}
              />
            </div>
            : null
        }
      </Fragment>
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
