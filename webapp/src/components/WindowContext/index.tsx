import React, { createContext, useState } from 'react'
import { FsEntry } from '../FsHasEntry'

export interface Window {
    fsEntry: FsEntry
    isActive?: boolean
}

interface ContextProps {
    /**
     * currently open windows
     */
    openWindows: Window[]

    /**
     * makes the named window active;
     * note; named window must already be open
     */
    setActiveWindow: (name: string) => void

    /**
     * opens the provided window
     */
    openWindow: (window: Window) => void

    /**
     * closes the window w/ provided name
     */
    closeWindow: (name: string) => void
}

export const WindowContext = createContext<ContextProps>(null as any)

export const WindowProvider = ({ children }) => {
    const [openWindows, setOpenWindows] = useState([] as Window[])

    return (
        <WindowContext.Provider value={{
            openWindows,
            setActiveWindow: path => {
                // ensure window w/ name is open
                if (openWindows.find(openWindow => openWindow.fsEntry.path === path)) {
                    setOpenWindows(
                        openWindows.map(
                            openWindow => ({
                                ...openWindow,
                                isActive:
                                    openWindow.fsEntry.path === path
                                        ? true
                                        : false
                            })
                        )
                    )
                }
            },
            openWindow: window => setOpenWindows([
                ...openWindows.reduce((acc, cur) => {
                    if (cur.fsEntry.path !== window.fsEntry.path) {
                        acc.push({
                            ...cur,
                            isActive: false
                        })
                    }
                    return acc
                },
                    [] as Window[]
                ),
                {
                    ...window,
                    isActive: true
                }
            ]),
            closeWindow: path => {
                let isNextItemActive = false
                setOpenWindows(
                    openWindows.reduce(
                        (acc, curr) => {
                            if (curr.fsEntry.path === path) {
                                // this is the item to be removed

                                if (curr.isActive && acc.length) {
                                    // item is active and not first; make previous item active
                                    acc.push({
                                        ...acc.pop()!,
                                        isActive: true
                                    })
                                } else {
                                    // item is active and first; make next item active ()
                                    isNextItemActive = true
                                }
                            } else {
                                acc.push({
                                    ...curr,
                                    isActive: isNextItemActive
                                })
                                isNextItemActive = false
                            }

                            return acc
                        },
                        [] as Window[]
                    ))
            }
        }}>
            {children}
        </WindowContext.Provider>
    )
}