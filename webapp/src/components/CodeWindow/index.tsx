import React, { useEffect, useState } from 'react'
import ReactAce from 'react-ace'
import { Window } from '../WindowContext'
import getFsEntryData from '../../queries/getFsEntryData'
import modelist from 'ace-builds/src-noconflict/ext-modelist'

// editor modes (alphabetical)
import 'ace-builds/src-noconflict/mode-css'
import 'ace-builds/src-noconflict/mode-dockerfile'
import 'ace-builds/src-noconflict/mode-gitignore'
import 'ace-builds/src-noconflict/mode-golang'
import 'ace-builds/src-noconflict/mode-groovy'
import 'ace-builds/src-noconflict/mode-handlebars'
import 'ace-builds/src-noconflict/mode-html'
import 'ace-builds/src-noconflict/mode-java'
import 'ace-builds/src-noconflict/mode-javascript'
import 'ace-builds/src-noconflict/mode-json'
import 'ace-builds/src-noconflict/mode-kotlin'
import 'ace-builds/src-noconflict/mode-markdown'
import 'ace-builds/src-noconflict/mode-mysql'
import 'ace-builds/src-noconflict/mode-nginx'
import 'ace-builds/src-noconflict/mode-pgsql'
import 'ace-builds/src-noconflict/mode-php'
import 'ace-builds/src-noconflict/mode-powershell'
import 'ace-builds/src-noconflict/mode-puppet'
import 'ace-builds/src-noconflict/mode-python'
import 'ace-builds/src-noconflict/mode-r'
import 'ace-builds/src-noconflict/mode-rdoc'
import 'ace-builds/src-noconflict/mode-scala'
import 'ace-builds/src-noconflict/mode-sh'
import 'ace-builds/src-noconflict/mode-sql'
import 'ace-builds/src-noconflict/mode-svg'
import 'ace-builds/src-noconflict/mode-text'
import 'ace-builds/src-noconflict/mode-tsx'
import 'ace-builds/src-noconflict/mode-typescript'
import 'ace-builds/src-noconflict/mode-xml'
import 'ace-builds/src-noconflict/mode-yaml'

import 'ace-builds/src-noconflict/theme-github'

interface Props {
    window: Window
}

/**
 * A window which enables viewing/editing code
 */
export default (
    {
        window
    }: Props
) => {
    const [data, setData] = useState('')
    useEffect(
        () => {
            const load = async () => {
                setData(await getFsEntryData(window.fsEntry.path))
            }

            load()
        },
        [
            window.fsEntry.path
        ]
    )

    return (
        <ReactAce
            style={{
                width: '100%',
                height: '100%'
            }}
            mode={modelist.getModeForPath(window.fsEntry.name).mode.replace('ace/mode/', '')}
            theme='github'
            //onChange={onChange}
            name='txt-window'
            editorProps={{ $blockScrolling: true }}
            value={data}
        />
    )
} 