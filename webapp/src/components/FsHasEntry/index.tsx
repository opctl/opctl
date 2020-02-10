import React, { CSSProperties, Fragment, useContext, useState } from 'react'
import { ReactComponent as NotesIcon } from '../../icons/Notes.svg'
import { ReactComponent as ArrowDownIcon } from '../../icons/ArrowDown.svg'
import { ReactComponent as ArrowRightIcon } from '../../icons/ArrowRight.svg'
import { ReactComponent as MoreHorizIcon } from '../../icons/MoreHoriz.svg'
import brandColors from '../../brandColors'
import { WindowContext } from '../WindowContext'

export interface FsEntry {
    name: string,
    dir?: FsEntry[]
    file?: string
    path: string
}

interface Props {
    fsEntry: FsEntry
    style?: CSSProperties
}

const FsHasEntry = (
    {
        fsEntry,
        style
    }: Props
) => {
    const [isOpen, setIsOpen] = useState(false)
    const { openWindow } = useContext(WindowContext)

    return (
        <div
            style={style}
        >
            <div
                style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center'
                }}
            >
                <div
                    style={{
                        display: 'flex',
                        alignItems: 'center',
                        minWidth: 0
                    }}
                >
                    {
                        fsEntry.dir
                            ? <Fragment>
                                {
                                    isOpen
                                        ? <ArrowDownIcon
                                            onClick={() => setIsOpen(false)}
                                            style={{
                                                flexShrink: 0,
                                                cursor: 'pointer',
                                                fill: brandColors.active
                                            }}
                                        />
                                        : <ArrowRightIcon
                                            onClick={() => setIsOpen(true)}
                                            style={{
                                                flexShrink: 0,
                                                cursor: 'pointer',
                                                fill: brandColors.active
                                            }}
                                        />
                                }
                                <div
                                    style={{
                                        textOverflow: 'ellipsis',
                                        whiteSpace: 'nowrap',
                                        overflow: 'hidden'
                                    }}
                                >
                                    {fsEntry.name}
                                </div>
                            </Fragment>
                            : <Fragment>
                                <NotesIcon />
                                <div
                                    onClick={() => openWindow(
                                        {
                                            fsEntry
                                        }
                                    )}
                                    style={{
                                        textOverflow: 'ellipsis',
                                        whiteSpace: 'nowrap',
                                        overflow: 'hidden',
                                        cursor: 'pointer',
                                        color: brandColors.active
                                    }}
                                >
                                    {fsEntry.name}
                                </div>
                            </Fragment>
                    }
                </div>
                <MoreHorizIcon
                    style={{
                        // @TODO: unhide once we need it
                        visibility: 'hidden',
                        flexShrink: 0,
                        cursor: 'pointer',
                        fill: brandColors.active 
                    }}
                />
            </div>
            {
                fsEntry.dir && isOpen
                    ? fsEntry.dir.map(
                        fsEntry => <FsHasEntry
                            key={fsEntry.name}
                            style={{
                                marginLeft: '.5rem'
                            }}
                            fsEntry={fsEntry} />
                    )
                    : null
            }
        </div>
    )
}

export default FsHasEntry