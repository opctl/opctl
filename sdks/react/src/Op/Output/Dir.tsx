import React from 'react'
import ModelParamDir from '@opctl/sdk/src/model/param/dir'
import Description from '../Param/Description'

interface Props {
  name: string
  param: ModelParamDir
  opRef: string
  value: any
}

export default (
  {
    name,
    param,
    opRef,
    value
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
        value={value || param.default || ''}
      />
    </div>
  )
}
