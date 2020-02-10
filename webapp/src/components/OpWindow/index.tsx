import React, { useState, useEffect } from 'react'
import getFsEntryData from '../../queries/getFsEntryData'
import { Window } from '../WindowContext'
import jsYaml from 'js-yaml'
import HasCall from '../HasCall'
import brandColors from '../../brandColors'
import AddCallPopper from '../AddCallPopper'
import { ReactComponent as PlusIcon } from '../../icons/Plus.svg'
import ReactEasyPanzoom from 'react-easy-panzoom'

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
    const [op, setOp] = useState(null as any)

    useEffect(
        () => {
            const load = async () => {
                setOp(
                    jsYaml.safeLoad(
                        await getFsEntryData(window.fsEntry.path)
                    )
                )
            }

            load()
        },
        [
            window.fsEntry.path
        ]
    )

    if (!op) {
        return null
    }

    return (
        <ReactEasyPanzoom
            style={{
                outline: 'none'
            }}
        >
            <div
                style={{
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    flexDirection: 'column',
                    width: 'max-content',
                    minWidth: '100%',
                    minHeight: '100%',
                }}
            >
                <AddCallPopper>
                    <PlusIcon
                        style={{
                            backgroundColor: brandColors.white,
                            cursor: 'pointer',
                            fill: brandColors.active,
                            display: 'block'
                        }}
                    />
                </AddCallPopper>
                <div
                    style={{
                        backgroundColor: brandColors.lightGray,
                        height: '2.5rem',
                        width: '.1rem'
                    }}
                ></div>
                <HasCall
                    call={op.run}
                />
                <div
                    style={{
                        backgroundColor: brandColors.lightGray,
                        height: '2.5rem',
                        width: '.1rem'
                    }}
                ></div>
                <AddCallPopper>
                    <PlusIcon
                        style={{
                            backgroundColor: brandColors.white,
                            cursor: 'pointer',
                            fill: brandColors.active,
                            display: 'block'
                        }}
                    />
                </AddCallPopper>
            </div>
        </ReactEasyPanzoom>
    )
}