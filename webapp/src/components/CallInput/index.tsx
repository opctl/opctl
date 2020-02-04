import React from 'react'
import { css } from 'emotion'
import formElementStyles from '../formElementStyles'

type CallType = 'container' | 'op' | 'parallel' | 'parallelLoop' | 'serial' | 'serialLoop'

interface Props {
    onChange: (callType: CallType) => void
    callType: CallType
}

export default (
    {
        onChange,
        callType
    }: Props
) => <input
        className={css(formElementStyles)}
        onChange={({ target }) => onChange(target.value as CallType)}
        value={callType}
    />