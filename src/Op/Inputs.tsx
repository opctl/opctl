import React from 'react'
import ModelParam from '@opctl/sdk/lib/model/param'
import Input from './Input/index'

interface Props {
  inputs: { [name: string]: ModelParam }
  onInvalid?: (name: string) => any | null | undefined
  onValid: (name: string, value: any) => any
  opRef: string
  scope: { [name: string]: any }
}

export default (
  {
    inputs,
    onInvalid,
    onValid,
    opRef,
    scope
  }: Props
) => {
  if (!inputs || Object.entries(inputs).length === 0) return (null)

  return (
    <div
      style={{
        marginTop: '1rem'
      }}
    >
      <h4>Inputs</h4>
      {Object.entries(inputs).map(([name, input]) =>
        <Input
          {
          ...onInvalid && {
            onInvalid: () => onInvalid(name)
          }
          }
          onValid={value => (onValid(name, value))}
          name={name}
          opRef={opRef}
          input={input}
          key={name}
          value={scope[name] || null}
        />
      )}
    </div>
  )
}
