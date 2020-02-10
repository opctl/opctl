import React, { useState, ReactElement } from 'react'
import { ReactComponent as CloseIcon } from '../../icons/Close.svg'
import CallTypeInput from '../CallTypeInput'
import CallNameInput from '../CallNameInput'
import Panel from '../Panel'
import PanelRow from '../PanelRow'
import Form from '../Form'
import brandColors from '../../brandColors'

interface Props {
    children: ReactElement | ReactElement[]
}

export default (
    {
        children
    }: Props
) => {
    const [callType, setCallType] = useState('')
    const [callName, setCallName] = useState('')
    const [isVisible, setIsVisible] = useState(false)

    return (
        <div
            style={{
                display: 'inline-block'
            }}
        >
            {
                isVisible ?
                    <div
                        style={{
                            position: 'absolute',
                            boxShadow: '0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)',
                            borderRadius: '.3rem',
                            padding: '1rem',
                            margin: '1rem',
                            backgroundColor: brandColors.white,
                            width: '20rem'
                        }}
                    >
                        <div
                            style={{
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'space-between'
                            }}
                        >
                            <h3>Add call</h3>
                            <CloseIcon
                                style={{
                                    cursor: 'pointer',
                                    flexShrink: 0,
                                    position: 'relative',
                                    right: 0,
                                    fill: brandColors.active
                                }}
                                onClick={() => setIsVisible(false)}
                            />
                        </div>
                        <Form
                            isValid={!!(callType && callName)}
                            onSubmit={() => { }}
                        >
                            <Panel>
                                <PanelRow>
                                    <label>Type</label>
                                    <CallTypeInput
                                        onChange={setCallType}
                                    />
                                </PanelRow>
                                <PanelRow>
                                    <label>Name</label>
                                    <CallNameInput
                                        onChange={setCallName}
                                        value={callName}
                                    />
                                </PanelRow>
                            </Panel>
                        </Form>
                    </div>
                    : null
            }
            <div
                onClick={() => setIsVisible(true)}
            >
                {children}
            </div>
        </div>
    )
}