import React, { Component } from 'react';
import { AutoSizer } from 'react-virtualized';
import Header from './Header'
import Inputs from '../../../Pkg/Inputs'
import EventBrowser from '../../../EventBrowser'
import { Modal, ModalBody } from 'reactstrap';
import opspecNodeApiClient from '../../../../utils/clients/opspecNodeApi';
import { toast } from 'react-toastify';

export default class Item extends Component {
    values = {};
    args = {};
    state = {
        isConfigurationVisible: false,
        isKillable: false
    }

    toggleConfigurationModal = () => {
        this.setState(prevState => ({ isConfigurationVisible: !prevState.isConfigurationVisible }));
    }

    handleInvalid = (name) => {
        const args = Object.assign({}, this.props.args);
        delete args[name];

        const values = Object.assign({}, this.props.values);
        delete values[name];

        this.props.onConfigured({ values, args });
    };

    isStartable = () => Object.keys(this.props.pkg.inputs || []).length === Object.keys(this.props.args).length

    handleValid = (name, value) => {
        const args = Object.assign({}, this.props.args);
        args[name] = value;

        const values = Object.assign({}, this.props.values);
        values[name] = value.value;

        this.props.onConfigured({ values, args });
    };

    kill = () => {
        opspecNodeApiClient.op_kill({
            opId: this.state.opId
        })
            .then(() => this.setState({ isKillable: false }))
            .catch(error => {
                toast.error(error.message);
            });
    };

    start = () => {
        opspecNodeApiClient.op_start({
            args: this.props.args,
            pkg: {
                ref: this.props.pkgRef,
            }
        })
            .then(opId => {
                this.setState({ opId, isKillable: true });

                opspecNodeApiClient.event_stream_get({
                    filter: {
                        roots: [opId],
                    },
                    onEvent: event => {
                        if (event.opEnded && event.opEnded.opId === opId) {
                            this.setState({
                                isKillable: false,
                                outputs: event.opEnded.outputs,
                            });
                        }
                    },
                })
            })
            .catch(error => {
                toast.error(error.message);
            });
    };

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
                                    values={this.props.values}
                                />
                            </ModalBody>
                        </Modal>
                        <div style={{ marginTop: '37px', height: 'calc(100% - 37px)', overflowY: 'auto' }}>
                            {
                                this.state.opId
                                    ?
                                    <EventBrowser key={this.state.opId} filter={{ roots: [this.state.opId] }} />
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