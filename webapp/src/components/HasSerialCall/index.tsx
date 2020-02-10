import React, { Fragment } from 'react'
import HasCall, { Call } from '../HasCall'
import brandColors from '../../brandColors'
import AddCallPopper from '../AddCallPopper'
import { ReactComponent as PlusIcon } from '../../icons/Plus.svg'

interface Props {
    call: Call
}

export default (
    {
        call
    }: Props
) => {
    return (
        <div>
            {
                call.serial!.map(
                    (childCall, index, array) =>
                        <div
                            key={index}
                            style={{
                                minWidth: 0,
                                flex: '1 0 0',
                                display: 'flex',
                                marginLeft: '2rem',
                                marginRight: '2rem',
                                alignItems: 'center',
                                flexDirection: 'column'
                            }}
                        >
                            <HasCall
                                key={index}
                                call={childCall}
                            />
                            {
                                index + 1 < array.length
                                    ? <Fragment>
                                        <div
                                            style={{
                                                backgroundColor: brandColors.lightGray,
                                                height: '100%',
                                                minHeight: '2.5rem',
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
                                                height: '100%',
                                                minHeight: '2.5rem',
                                                width: '.1rem'
                                            }}
                                        ></div>
                                    </Fragment>
                                    // is leaf
                                    : null
                            }
                        </div>
                )
            }
        </div>
    )
}