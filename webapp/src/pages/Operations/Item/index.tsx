import React, { PureComponent } from 'react'
import { AutoSizer } from 'react-virtualized'
import { HotKeys } from 'react-hotkeys'
import Header from './Header'
import Inputs from '../../../components/Op/Inputs'
import EventStream from '../../../components/EventStream'
import { Modal, ModalBody } from 'reactstrap'
import Input from '../../../components/Input'
import eventStore from '../../../core/eventStore'
import reactClickOutside from 'react-click-outside'
import './index.css'

interface Props {
  args
  isAllItemsSelected
  isKillable
  isFullScreen: boolean
  isStartable: boolean
  name: string
  onConfigured
  onDelete
  onKill
  onStart
  onToggleFullScreen
  op
  opId: string
  opRef: string
}

interface State {
  isConfigurationVisible: boolean
  isKillable: boolean
  isSelected: boolean
  name: string
  outputs?
}

class Item extends PureComponent<Props,State> {
  args = this.props.args;
  eventStreamCloser

  state: State = {
    isConfigurationVisible: false,
    isKillable: false,
    isSelected: false,
    name: this.props.name || this.props.opRef
  };

  ensureEventStreamClosed = () => {
    if (this.eventStreamCloser) this.eventStreamCloser()
  };

  toggleConfigurationModal = () => {
    this.setState(prevState => ({ isConfigurationVisible: !prevState.isConfigurationVisible }))
  };

  // this method is required by react-click-outside
  handleClickOutside () {
    this.setState({ isSelected: false })
  }

  handleSelected = () => {
    this.setState({ isSelected: true })
  };

  handleInvalidArg = (name) => {
    delete this.args[name]

    this.props.onConfigured({ args: this.args })
  };

  handleValidArg = (name, value) => {
    this.args[name] = value

    this.props.onConfigured({ args: this.args })
  };

  handleNameChanged = (name) => {
    this.setState({ name })
    this.props.onConfigured({ name })
  };

  processEventStream = ({ opId }) => {
    this.eventStreamCloser = eventStore.getStream({
      filter: {
        roots: [opId]
      },
      onEvent: event => {
        if (event.opStarted) {
          this.setState({ isKillable: true })
        }
        if (event.opEnded && event.opEnded.opId === opId) {
          this.setState({
            isKillable: false,
            outputs: event.opEnded.outputs
          })
        }
      }
    })
  };

  componentWillReceiveProps (nextProps) {
    if (nextProps.opId !== this.props.opId) {
      this.ensureEventStreamClosed()
      this.processEventStream({ opId: nextProps.opId })
    }
  }

  componentWillMount () {
    if (this.props.opId) {
      this.processEventStream({ opId: this.props.opId })
    }
  }

  componentWillUnmount () {
    this.ensureEventStreamClosed()
  }

  render () {
    return (
      <AutoSizer>
        {({ height, width }) => (
          <HotKeys
            keyMap={{
              del: 'del',
              kill: 'ctrl+c',
              start: 'enter'
            }}
            handlers={{
              del: this.props.onDelete,
              kill: this.props.onKill,
              start: this.props.onStart
            }}
          >
            <div
              tabIndex={-1}
              style={{ height, width, border: 'dashed 3px #ececec' }}
              onClick={this.handleSelected}
            >
              <Header
                isFullScreen={this.props.isFullScreen}
                isKillable={this.state.isKillable}
                isStartable={this.props.isStartable}
                onToggleFullScreen={this.props.onToggleFullScreen}
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
                  <Input
                    name='name'
                    value={this.state.name}
                    validate={() => []}
                    onInvalid={() => {
                    }}
                    onValid={this.handleNameChanged}
                  />
                  <Inputs
                    value={this.props.op.inputs}
                    onInvalid={this.handleInvalidArg}
                    onValid={this.handleValidArg}
                    opRef={this.props.opRef}
                    values={this.args}
                  />
                </ModalBody>
              </Modal>
              {
                this.props.isFullScreen || this.state.isSelected || this.props.isAllItemsSelected
                  ? null
                  : <div style={{
                    opacity: 0.2,
                    backgroundColor: '#000',
                    position: 'absolute',
                    width: `${width - 6}px`,
                    top: '37px',
                    height: `${height - 40}px`,
                    zIndex: 1,
                    cursor: 'pointer'
                  }} />
              }
              <div style={{ marginTop: '37px', height: 'calc(100% - 37px)' }}>
                {
                  this.props.opId
                    ? <EventStream key={this.props.opId} filter={{ roots: [this.props.opId] }} />
                    : null
                }
              </div>
            </div>
          </HotKeys>
        )}
      </AutoSizer>
    )
  }
}

export default reactClickOutside(Item)
