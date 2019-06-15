import React from 'react'
import ModelParamSocket from '@opctl/sdk/lib/model/param/socket'
import Description from '../Param/Description'

interface Props {
  name: string
  param: ModelParamSocket
  opRef: string
}

export default (
  {
    name,
    param,
    opRef
  }: Props
) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <Description value={param.description} opRef={opRef} />
      <input
        className='form-control'
        id={name}
        readOnly
        type='text'
      />
    </div>
  )
}
