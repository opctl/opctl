import React, { Fragment } from 'react'
import HasCall, { Call } from '../HasCall'
import { ReactComponent as PlusIcon } from '../../icons/Plus.svg'
import brandColors from '../../brandColors'
import AddCallPopper from '../AddCallPopper'


interface Props {
    call: Call
}

export default (
    {
        call
    }: Props
) => {
    // splice a dummy call into the middle
    const middle = Math.round(call.parallel!.length / 2)
    let spliced = [...call.parallel!]
    spliced.splice(middle, 0, {})

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
                    gridTemplateColumns: `repeat(${spliced.length},1fr)`,
                    borderTop: `solid thin ${brandColors.lightGray}`,
                    borderBottom: `solid thin ${brandColors.lightGray}`
                }}
            >
                {
                    spliced.map(
                        (childCall, index) => {
                            const isDummyCall = !(childCall.container || childCall.op || childCall.parallel || childCall.serial)

                            return (
                                <div
                                    key={index}
                                    style={{
                                        display: 'flex',
                                        alignItems: 'center',
                                        flexDirection: 'column'
                                    }}
                                >
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
                                </div>
                            )
                        }
                    )
                }
            </div>
        </div>
    )
}