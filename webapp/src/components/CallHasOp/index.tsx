import React, { useState, useEffect } from 'react'
import getFsEntryData from '../../queries/getFsEntryData'
import jsYaml from 'js-yaml'
import HasCall, { OpCall } from '../HasCall'
import brandColors from '../../brandColors'
import AddCallPopper from '../AddCallPopper'
import { ReactComponent as PlusIcon } from '../../icons/Plus.svg'
import { toast } from 'react-toastify'
import path from 'path'

interface Props {
  opCall: OpCall
  parentOpRef?: string
}

export default function CallHasOp(
  {
    opCall,
    parentOpRef
  }: Props
) {
  const opRef = opCall.ref.startsWith('.') && parentOpRef
    ? path.join(parentOpRef, opCall.ref)
    : opCall.ref

  const [op, setOp] = useState(null as any)

  useEffect(
    () => {
      const load = async () => {
        try {
          setOp(
            jsYaml.safeLoad(
              await getFsEntryData(path.join(opRef, 'op.yml'))
            )
          )
        } catch (err) {
          const errAsError = err as Error
          if (errAsError.message.includes('authentication')) {
            toast.warn(`Loading ${opRef} skipped; because it requires authentication.`)
          } else if (errAsError.message.includes('service=git-upload-pack')) {
            toast.warn(`Loading 'ref: ${opRef}' skipped because you're using deprecated syntax! To fix, use 'ref: ../${opRef}'.`)
          } else {
            toast.error(errAsError.toString())
          }
        }
      }

      load()
    },
    [
      parentOpRef,
      opRef
    ]
  )

  if (!op) {
    return null
  }

  return (
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
        parentOpRef={opRef}
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
  )
}