import React from 'react'
import { Call } from '../HasCall'
import brandColors from '../../brandColors'

interface Props {
    call: Call
}

function getName(call: Call): string | null {
    if (call.name) {
        return call.name
    } else if (call.container) {
        return call.container.image?.ref || 'Container'
    } else if (call.op) {
        return call.op.ref || 'Op'
    }
    return null
}

export default (
    {
        call
    }: Props
) => {
    const name = getName(call)
    if (!name) {
        return null
    }

    return (
        <div
            style={{
                border: `solid thin ${brandColors.lightGray}`,
                minWidth: '5rem',
                textAlign: 'center',
                padding: '0 .5rem',
                margin: '0 .5rem'
            }}
        >
            {
                name
            }
        </div>
    )
}