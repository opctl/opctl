import React, { Fragment } from 'react'
import HasCall, { CallParallel } from '../HasCall'
import { ReactComponent as PlusIcon } from '../../icons/Plus.svg'
import brandColors from '../../brandColors'
import AddCallPopper from '../AddCallPopper'


interface Props {
    callParallel: CallParallel
    parentOpRef: string
}

export default (
    {
        callParallel,
        parentOpRef
    }: Props
) => {
    // splice a dummy call into the middle
    // @TODO: re-enable once edit supported
    //const middle = Math.round(call.parallel!.length / 2)
    let spliced = [...callParallel!]
    //spliced.splice(middle, 0, {})

    return (
        <div
            style={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center'
            }}
        >
            <div
                style={{
                    // https://stackoverflow.com/questions/29503227/how-to-make-flexbox-items-the-same-size
                    display: 'grid',
                    gridTemplateColumns: `repeat(${spliced.length},1fr)`
                }}
            >
                {
                    spliced.map(
                        (childCall, index) => {
                            const isDummyCall = !(childCall.container || childCall.op || childCall.parallel || childCall.serial)
                            const isFirst = index === 0
                            const isLast = index + 1 === spliced.length

                            return (
                                <div
                                    key={index}
                                    style={{
                                        display: 'flex',
                                        alignItems: 'center',
                                        flexDirection: 'column',
                                        position: 'relative'
                                    }}
                                >
                                    <div
                                        style={{
                                            ...!(isFirst || isLast)
                                                ? {
                                                    width: '100%'
                                                }
                                                : null,
                                            ...isFirst
                                                ? {
                                                    right: 0,
                                                    width: '50%'
                                                }
                                                : null,
                                            ...isLast
                                                ? {
                                                    left: 0,
                                                    width: '50%'
                                                }
                                                : null,
                                            position: 'absolute',
                                            top: 0,
                                            height: '.1rem',
                                            backgroundColor: brandColors.lightGray
                                        }}
                                    ></div>
                                    <div
                                        style={{
                                            backgroundColor: brandColors.lightGray,
                                            minHeight: '2.5rem',
                                            height: '100%',
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
                                    <div
                                        style={{
                                            backgroundColor: brandColors.lightGray,
                                            minHeight: '2.5rem',
                                            height: '100%',
                                            width: '.1rem'
                                        }}
                                    ></div>
                                    {
                                        isDummyCall
                                            ? null
                                            : <Fragment>
                                                <HasCall
                                                    key={index}
                                                    call={childCall}
                                                    parentOpRef={parentOpRef}
                                                />
                                                <div
                                                    style={{
                                                        backgroundColor: brandColors.lightGray,
                                                        minHeight: '2.5rem',
                                                        height: '100%',
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
                                                <div
                                                    style={{
                                                        backgroundColor: brandColors.lightGray,
                                                        minHeight: '2.5rem',
                                                        height: '100%',
                                                        width: '.1rem'
                                                    }}
                                                ></div>
                                            </Fragment>
                                    }
                                    <div
                                        style={{
                                            ...!(isFirst || isLast)
                                                ? {
                                                    width: '100%'
                                                }
                                                : null,
                                            ...isFirst
                                                ? {
                                                    right: 0,
                                                    width: '50%'
                                                }
                                                : null,
                                            ...isLast
                                                ? {
                                                    left: 0,
                                                    width: '50%'
                                                }
                                                : null,
                                            position: 'absolute',
                                            bottom: 0,
                                            height: '.1rem',
                                            backgroundColor: brandColors.lightGray
                                        }}
                                    ></div>
                                </div>
                            )
                        }
                    )
                }
            </div>
        </div>
    )
}