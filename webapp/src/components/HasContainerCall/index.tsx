import React from 'react'
import { Call } from '../HasCall'
import brandColors from '../../brandColors'

interface Props {
    call: Call
}

export default (
    {
        call
    }: Props
) => {
    return (
        <div
            style={{
                border: `solid thin ${brandColors.lightGray}`,
                width: '5rem'
            }}
        >
            Container
        </div>
    )
}