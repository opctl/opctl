import React, { useContext } from 'react'
import brandColors from '../../brandColors'
import Window from '../Window'
import { WindowContext } from '../WindowContext'
import { ReactComponent as CloseIcon } from '../../icons/Close.svg'
import path from 'path'

export default () => {
    const {
        closeWindow,
        openWindows,
        setActiveWindow
    } = useContext(WindowContext)

    return (
        <div
            style={{
                flexGrow: 0,
                display: 'flex',
                flexDirection: 'column',
                width: '100%',
                height: '100%',
                overflow: 'hidden'
            }}
        >
            <div
                style={{
                    display: 'flex',
                    height: '2.4rem',
                    flexShrink: 0,
                    borderBottom: `solid .1rem ${brandColors.reallyLightGray}`,
                    backgroundColor: brandColors.reallyReallyLightGray
                }}
            >
                {
                    openWindows.map(
                        window =>
                            <div
                                key={window.fsEntry.path}
                                style={{
                                    cursor: 'pointer',
                                    paddingLeft: '1rem',
                                    paddingRight: '1rem',
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'space-between',
                                    backgroundColor: window.isActive
                                        ? brandColors.white
                                        : brandColors.reallyReallyLightGray,
                                    borderRight: `solid .1rem ${brandColors.reallyLightGray}`,
                                    minWidth: 0
                                }}
                                onClick={() => setActiveWindow(window.fsEntry.path)}
                            >
                                <div
                                    style={{
                                        textOverflow: 'ellipsis',
                                        whiteSpace: 'nowrap',
                                        overflow: 'hidden'
                                    }}
                                >
                                    {window.fsEntry.name} ../{path.basename(path.dirname(window.fsEntry.path))}
                                </div>
                                <CloseIcon
                                    style={{
                                        cursor: 'pointer',
                                        flexShrink: 0,
                                        position: 'relative',
                                        right: 0,
                                        fill: brandColors.active
                                    }}
                                    onClick={e => {
                                        closeWindow(window.fsEntry.path)
                                        e.stopPropagation()
                                    }}
                                />
                            </div>
                    )
                }
            </div>
            { 
                openWindows.map(
                    window =>
                        <div
                            key={window.fsEntry.path}
                            style={{
                                ...!window.isActive
                                    ? { display: 'none' }
                                    : {},
                                width: '100%',
                                height: '100%',
                                overflow: 'auto'
                            }}
                        >
                            <Window
                                window={window}
                            />
                        </div>
                )
            }
        </div>
    )
}