import React, {Component} from 'react';
import {AutoSizer} from 'react-virtualized';
import {HotKeys} from 'react-hotkeys';
import Header from './Header'
import Inputs from '../../../Pkg/Inputs'
import EventStream from '../../../EventStream'
import {Modal, ModalBody} from 'reactstrap';
import Input from '../../../Input';
import opspecNodeApiClient from '../../../../utils/clients/opspecNodeApi';
import reactClickOutside from 'react-click-outside';
import './index.css';

class Item extends Component {
  args = this.props.args;
  state = {
    name: this.props.name || this.props.pkgRef,
  };

  toggleConfigurationModal = () => {
    this.setState(prevState => ({isConfigurationVisible: !prevState.isConfigurationVisible}));
  };

  // this method is required by react-click-outside
  handleClickOutside() {
    this.setState({isSelected: false});
  }

  handleSelected = () => {
    this.setState({isSelected: true});
  };

  handleInvalidArg = (name) => {
    delete this.args[name];

    this.props.onConfigured({args: this.args});
  };

  handleValidArg = (name, value) => {
    this.args[name] = value;

    this.props.onConfigured({args: this.args});
  };

  handleNameChanged = (name) => {
    this.setState({name});
    this.props.onConfigured({name});
  };

  processEventStream = ({opId}) => {
    opspecNodeApiClient.event_stream_get({
      filter: {
        roots: [opId],
      },
      onEvent: event => {
        if (event.opStarted) {
          this.setState({isKillable: true});
        }
        if (event.opEnded && event.opEnded.opId === opId) {
          this.setState({
            isKillable: false,
            outputs: event.opEnded.outputs,
          });
        }
      },
    })
  };

  componentWillReceiveProps(nextProps) {
    this.processEventStream({opId: nextProps.opId});
  }

  componentWillMount() {
    this.processEventStream({opId: this.props.opId});
  }

  render() {
    return (
      <AutoSizer>
        {({height, width}) => (
          <HotKeys
            keyMap={{
              del: 'del',
              kill: 'ctrl+c',
              start: 'enter',
            }}
            handlers={{
              del: this.props.onDelete,
              kill: this.props.onKill,
              start: this.props.onStart,
            }}
          >
            <div
              tabIndex='-1'
              style={{height, width, border: 'dashed 3px #ececec'}}
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
                    value={this.props.pkg.inputs}
                    onInvalid={this.handleInvalidArg}
                    onValid={this.handleValidArg}
                    pkgRef={this.props.pkgRef}
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
                    cursor: 'pointer',
                  }}>
                  </div>
              }
              <div style={{marginTop: '37px', height: 'calc(100% - 37px)'}}>
                {
                  this.props.opId
                    ?
                    <EventStream key={this.props.opId} filter={{roots: [this.props.opId]}}/>
                    :
                    null
                }
              </div>
            </div>
          </HotKeys>
        )}
      </AutoSizer>
    );
  }
}

export default reactClickOutside(Item);
