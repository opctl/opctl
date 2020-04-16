import React, { useState, useEffect } from 'react'
import getFsEntryData from '../../queries/getFsEntryData'
import jsYaml from 'js-yaml'
import HasCall from '../HasCall'
import brandColors from '../../brandColors'
import AddCallPopper from '../AddCallPopper'
import { ReactComponent as PlusIcon } from '../../icons/Plus.svg'
import { toast } from 'react-toastify'

interface Props {
    opRef: string
}

export default (
    {
        opRef
    }: Props
) => {
    const [op, setOp] = useState(null as any)

    useEffect(
        () => {
            const load = async () => {
                try {
                    setOp(
                        jsYaml.safeLoad(
                            await getFsEntryData(opRef)
                        )
                    )
                } catch (err) {
                    if (err.message.includes('authentication')) {
                        toast.warn(`Loading ${opRef} skipped; because it requires authentication.`)
                    } else if (err.message.includes('no such host')) {
                        const ref = opRef.replace('/op.yml', '')
                        toast.warn(`Loading 'ref: ${ref}' skipped because you're using deprecated syntax! To fix, use 'ref: ../${ref.replace('/op.yml', '')}'.`)
                    } else {
                        toast.error(err)
                    }
                }
            }

            load()
        },
        [
            opRef
        ]
    )

    if (!op) {
        return null
    }

    return (
        <div
            style={{
                border: `solid thin ${brandColors.lightGray}`
            }}
        >
            <div
                style={{
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    flexDirection: 'column',
                    width: 'max-content',
                    minWidth: '100%',
                    minHeight: '100%',
                }}
            >
                <div
                    style={{
                        backgroundColor: brandColors.lightGray,
                        minHeight: '2.5rem',
                        height: '100%',
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
                        height: '2.5rem',
                        width: '.1rem'
                    }}
                ></div>
                <HasCall
                    call={op.run}
                    opRef={opRef}
                />
                <div
                    style={{
                        backgroundColor: brandColors.lightGray,
                        height: '2.5rem',
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
                        minHeight: '2.5rem',
                        height: '100%',
                        width: '.1rem'
                    }}
                ></div>
            </div>
        </div>
    )
}