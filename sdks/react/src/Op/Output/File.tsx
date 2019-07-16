import React from 'react'
import ModelParamFile from '@opctl/sdk/lib/model/param/file'
import Description from '../Param/Description'

interface Props {
  name: string
  param: ModelParamFile
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
      {
        value
          ? <a className='form-control' href={value || param.default || ''}>view file</a>
          : <input
            className='form-control'
            id={name}
            readOnly
            type='text'
            value={value || param.default || ''}
          />
      }
    </div>
  )
}
