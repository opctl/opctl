import React from 'react'
import { Window } from '../WindowContext'
import ReactEasyPanzoom from 'react-easy-panzoom'
import CallHasOp from '../CallHasOp'
import CallHasSummary from '../CallHasSummary'
import path from 'path'
import brandColors from '../../brandColors'

interface Props {
    window: Window
}

/**
 * A window which enables viewing/editing/running ops WYSIWYG (What You See Is What You Get) style
 */
export default function OpWindow(
    {
        window
    }: Props
) {
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
                    //justifyContent: 'center',
                    display: 'flex',
                    flexDirection: 'column',
                    border: `solid .1rem ${brandColors.lightGray}`,
                    marginLeft: '1rem',
                    marginRight: '1rem'
                }}
            >
                <CallHasSummary
                    call={call}
                />
                <CallHasOp
                    opCall={call.op}
                />
            </div>
        </ReactEasyPanzoom>
    )
}