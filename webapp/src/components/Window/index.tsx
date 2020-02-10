import React from 'react'
import { Window } from '../WindowContext'
import OpWindow from '../OpWindow'
import CodeWindow from '../CodeWindow'
import path from 'path'

interface Props {
    window: Window
}

export default (
    {
        window
    }: Props
) => {
    if (
        path.basename(window.fsEntry.path) === 'op.yml'
    ) {
        return (
            <OpWindow
                window={window}
            />
        )
    }

    return (
        <CodeWindow
            window={window}
        />
    )
}