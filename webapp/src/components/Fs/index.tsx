import React, { CSSProperties, Fragment, useState, useEffect } from 'react'
import FsHasEntry, { FsEntry } from '../FsHasEntry'
import listFsEntries from '../../queries/listFsEntries'
import LoadingSpinner from '../LoadingSpinner'
import path from 'path'

interface Props {
    style?: CSSProperties
}

export default (
    {
        style
    }: Props
) => {
    const [fsEntries, setFsEntries] = useState([] as FsEntry[])
    const urlSearchParams = new window.URLSearchParams(window.location.search)
    const mount = urlSearchParams.get('mount')

    useEffect(
        () => {
            const load = async () => {
                if (!mount) {
                    return
                }

                const rawFsEntries = await listFsEntries(mount)

                setFsEntries(
                    // reduce (flat) rawFsEntries to tree structure
                    rawFsEntries.reduce(
                        (fsEntries, rawFsEntry) => {
                            const absoluteRawPath = path.join(mount, rawFsEntry.Path)
                            const pathParts = absoluteRawPath.split("/")
                            if ('' === pathParts[0]) {
                                pathParts[0] = '/'
                            }

                            let currentFsEntries = fsEntries

                            // walk the parts
                            pathParts.forEach(
                                (pathPart, pathPartsIndex, pathPartsArray) => {
                                    let currentFsEntry = currentFsEntries.find(fsEntry => pathPart === fsEntry.name)

                                    if (rawFsEntry.Size === 96 || pathPartsIndex + 1 !== pathPartsArray.length) {
                                        // handle dir
                                        // @TODO: replace this method of identifying dirs
                                        if (!currentFsEntry) {
                                            currentFsEntry = {
                                                dir: [],
                                                name: pathPart,
                                                path: '/' + pathPartsArray.slice(1, pathPartsIndex + 1).join('/')
                                            }
                                            currentFsEntries.push(currentFsEntry)
                                        }
                                        currentFsEntries = currentFsEntry.dir!
                                    } else {
                                        // handle file
                                        if (!currentFsEntry) {
                                            currentFsEntry = {
                                                file: '',
                                                name: pathPart,
                                                path: absoluteRawPath
                                            }
                                            currentFsEntries.push(currentFsEntry)
                                        }
                                    }
                                }
                            )

                            return fsEntries
                        },
                        [] as FsEntry[]
                    )
                )
            }

            load()
        },
        [
            mount
        ]
    )

    return (
        <div
            style={{
                ...style,
                overflowY: 'auto'
            }}
        >
            {
                fsEntries.length
                    ? fsEntries.map(
                        fsEntry => <FsHasEntry
                            key={fsEntry.name}
                            fsEntry={fsEntry}
                        />
                    )
                    : <div
                        style={{
                            marginTop: '.5rem',
                            display: 'flex',
                            flexDirection: 'column',
                            alignItems: 'center'
                        }}
                    >
                        {
                            mount
                                ? <Fragment>
                                    <LoadingSpinner />
                                    <h4>Indexing...</h4>
                                </Fragment>
                                : <div
                                    style={{
                                        whiteSpace: 'pre-wrap',
                                        wordBreak: 'break-word',
                                    }}>
                                    <h3>Welcome!</h3>
                                    <p>Add a "mount=/filesystem/path" query parameter to explore</p>
                                    <h4>Example</h4>
                                    <p>http://{window.location.host}/?mount=/Users/source/repos/github.com/opctl</p>
                                </div>
                        }
                    </div>
            }
        </div>
    )
}