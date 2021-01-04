import React, { useState, useEffect } from 'react'
import getFsEntryData from '../../queries/getFsEntryData'
import { Call } from '../HasCall'
import constructDataUrl from '../../constructDataUrl'
import path from 'path'

interface Props {
  call: Call
  parentOpRef?: string
}

export default function CallHasIcon(
  {
    call,
    parentOpRef
  }: Props
) {
  const [url, setUrl] = useState('')

  useEffect(
    () => {
      const load = async () => {
        if (call.op) {
          const opRef = call.op.ref.startsWith('.') && parentOpRef
            ? path.join(parentOpRef, call.op.ref)
            : call.op.ref

          const iconDataRef = path.join(opRef, 'icon.svg')
          try {
            // try to load icon
            await getFsEntryData(iconDataRef)
            setUrl(constructDataUrl(iconDataRef))
          } catch (err) {
            // ignore icon load errors
          }
        }
      }

      load()
    },
    [
      call,
      parentOpRef
    ]
  )

  if (!url) {
    return null
  }

  return (
    <img
      src={url}
      alt={'icon'}
      style={{ height: '2.8rem', width: '2.8rem', marginRight: '.5rem' }}
    />
  )
}