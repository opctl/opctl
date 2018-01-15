import React, { Component } from 'react';
import { AutoSizer } from 'react-virtualized';
import Header from './Header'
import Inputs from '../../../Pkg/Inputs'
import EventBrowser from '../../../EventBrowser'
import { Modal, ModalBody } from 'reactstrap';
import opspecNodeApiClient from '../../../../utils/clients/opspecNodeApi';
import { toast } from 'react-toastify';

export default class Item extends Component {
    args = this.props.args;
    state = {
        isConfigurationVisible: false,
        isKillable: false,
        opId: this.props.opId
    };

    toggleConfigurationModal = () => {
        this.setState(prevState => ({ isConfigurationVisible: !prevState.isConfigurationVisible }));
    };

    handleInvalid = (name) => {
        delete this.args[name];

        this.props.onConfigured({ args: this.args });
    };

    isStartable = () => Object.keys(this.props.pkg.inputs || []).length === Object.keys(this.args).length

    handleValid = (name, value) => {
        this.args[name] = value;

        this.props.onConfigured({ args: this.args });
    };

    kill = () => {
        opspecNodeApiClient.op_kill({
            opId: this.props.opId
        })
            .then(() => {
                this.props.onConfigured({ opId: null });
                this.setState({ isKillable: false })
            })
            .catch(error => {
                toast.error(error.message);
            });
    };

    processEventStream = ({ opId }) => {
        opspecNodeApiClient.event_stream_get({
            filter: {
                roots: [opId],
            },
            onEvent: event => {
                if (event.opStarted) {
                    this.setState({ isKillable: true });
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

    start = () => {
        const args = Object.entries(this.props.pkg.inputs || [])
            .reduce((args, [name, param]) => {
                if (param.array) args[name] = { array: this.args[name] };
                if (param.dir) args[name] = { dir: this.args[name] };
                if (param.file) args[name] = { file: this.args[name] };
                if (param.number) args[name] = { number: this.args[name] };
                if (param.object) args[name] = { object: this.args[name] };
                if (param.socket) args[name] = { socket: this.args[name] };
                if (param.string) args[name] = { string: this.args[name] };
                return args;
            }, {});

        opspecNodeApiClient.op_start({
            args,
            pkg: {
                ref: this.props.pkgRef,
            }
        })
            .then(opId => {
                this.props.onConfigured({ opId });
                this.processEventStream({ opId: this.props.opId });
            })
            .catch(error => {
                toast.error(error.message);
            });
    };

    componentWillMount() {
        this.processEventStream({ opId: this.props.opId });
    }

    render() {
        return (
            <AutoSizer>
                {({ height, width }) => (
                    <div style={{ height, width, border: 'dashed 3px #ececec' }}>
                        <Header
                            isKillable={this.state.isKillable}
                            isStartable={this.isStartable()}
                            onConfigure={this.toggleConfigurationModal}
                            onStart={this.start}
                            onKill={this.kill}
                            onDelete={this.props.onDelete}
                            pkgRef={this.props.pkgRef}
                        />
                        <Modal
                            isOpen={this.state.isConfigurationVisible}
                            toggle={this.toggleConfigurationModal}
                        >
                            <ModalBody>
                                <Inputs
                                    value={this.props.pkg.inputs}
                                    onInvalid={this.handleInvalid}
                                    onValid={this.handleValid}
                                    pkgRef={this.props.pkgRef}
                                    values={this.args}
                                />
                            </ModalBody>
                        </Modal>
                        <div style={{ marginTop: '37px', height: 'calc(100% - 37px)', overflowY: 'auto' }}>
                            {
                                this.props.opId
                                    ?
                                    <EventBrowser key={this.props.opId} filter={{ roots: [this.props.opId] }} />
                                    :
                                    null
                            }
                        </div>
                    </div>
                )}
            </AutoSizer>
        );
    }
}
