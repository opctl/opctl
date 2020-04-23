import React from 'react'
import { Window } from '../WindowContext'
import ReactEasyPanzoom from 'react-easy-panzoom'
import CallHasOp from '../CallHasOp'
import CallHasName from '../CallHasName'
import path from 'path'

interface Props {
    window: Window
}

/**
 * A window which enables viewing/editing/running ops WYSIWYG (What You See Is What You Get) style
 */
export default (
    {
        window
    }: Props
) => {
    const call = {
        op: {
            ref: path.dirname(window.fsEntry.path)
        }
    }

    return (
        <ReactEasyPanzoom
            style={{
                outline: 'none'
            }}
        >
            <div
                style={{
                    alignItems: 'center',
                    display: 'flex',
                    flexDirection: 'column'
                }}
            >
                <CallHasName
                    call={call}
                />
                <CallHasOp
                    callOp={call.op}
                />
            </div>
        </ReactEasyPanzoom>
    )
}